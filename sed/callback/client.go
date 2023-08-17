package callback

import (
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/model"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/common/protocol"
	"git.multiverse.io/eventkit/kit/common/serializer"
	"git.multiverse.io/eventkit/kit/common/status"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
	"github.com/valyala/fasthttp"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// ErrNullPointer is a null pointer error which raising when input parameter is nil
// The caller can use this variable to check if the error is excepted error.
var ErrNullPointer = errors.New(constant.SystemInternalError, "Null pointer error")

const (
	defaultReadTimeoutMilliseconds  = int64(80 * 1000)
	defaultWriteTimeoutMilliseconds = int64(80 * 1000)
)

var (
	maxReadTimeoutMilliseconds  = defaultReadTimeoutMilliseconds
	maxWriteTimeoutMilliseconds = defaultWriteTimeoutMilliseconds
	once                        sync.Once
	clientWithConnectTimeout    *fasthttp.Client
)

func SetMaxReadAndWriteTimeoutMilliseconds(iMaxReadTimeoutMilliseconds, iMaxWriteTimeoutMilliseconds int64) {
	log.Infosf("Setting read timeout(milliseconds)=[%d], write timeout(milliseconds)=[%d]", iMaxReadTimeoutMilliseconds, iMaxWriteTimeoutMilliseconds)
	maxReadTimeoutMilliseconds = iMaxReadTimeoutMilliseconds
	maxWriteTimeoutMilliseconds = iMaxWriteTimeoutMilliseconds
}

func createClientIfNecessary() {
	once.Do(func() {
		clientWithConnectTimeout = &fasthttp.Client{
			// Default value of MaxConnsPerHost is 512, increase to 16384
			MaxConnsPerHost: 16384,
			// Disable idempotent calls attempts when remote call abnormal.
			//
			// default max attempts is 5 times.
			MaxIdemponentCallAttempts: 1,
		}
		log.Infosf("Create http client with read timeout(milliseconds)=[%d], write timeout(milliseconds)=[%d]", maxReadTimeoutMilliseconds, maxWriteTimeoutMilliseconds)
		clientWithConnectTimeout.ReadTimeout = time.Duration(maxReadTimeoutMilliseconds) * time.Millisecond
		clientWithConnectTimeout.WriteTimeout = time.Duration(maxWriteTimeoutMilliseconds) * time.Millisecond
	})
}

// post is a internal function for client post request to server
// and it's returns responses where from server endpoint
// allow the caller to specify a timeout(millisecond)
func post(protoMsg protocol.ProtoMessage, path string, isNeedDeserializerToMessage bool, timeout time.Duration) (protocol.ProtoMessage, error) {
	var retProtoMsg protocol.ProtoMessage
	var err error
	var fullURL string
	fullURL = serverAddr + path

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	var requestBytes []byte
	var version string

	if status.ServerProtocolLevel >= 2 {
		requestBytes = serializer.ProtoMsg2Bytes(&protoMsg, nil)
		version = "2"
	} else {
		requestBytes = []byte(serializer.ProtoMsg2String(&protoMsg, nil))
		version = "1"
	}

	req.Header.DisableNormalizing()
	req.Header.SetMethod(http.MethodPost)
	req.Header.Set("v", version)

	req.SetRequestURI(fullURL)
	req.SetBody(requestBytes)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	createClientIfNecessary()
	if err = clientWithConnectTimeout.DoTimeout(req, resp, timeout); err != nil {
		log.Errorsf("post to server failed err=%v", err)
		if fasthttp.ErrTimeout == err {
			return retProtoMsg, errors.Wrap(constant.SystemRemoteCallTimeout, err, 0)
		} else if fasthttp.ErrConnectionClosed == err {
			return retProtoMsg, errors.Errorf(constant.SystemErrConnectionClosed, "Connection closed, url=[%s]", fullURL)
		} else {
			switch err.(type) {
			case *net.OpError:
				{
					netOpError := err.(*net.OpError)
					switch netOpError.Err.(type) {
					case *os.SyscallError:
						{
							syscallError := netOpError.Err.(*os.SyscallError)
							if errno, ok := syscallError.Err.(syscall.Errno); ok {
								switch errno {
								case syscall.ECONNREFUSED:
									{
										return retProtoMsg, errors.Wrap(constant.SystemErrConnectionRefused, err, 0)
									}
								case syscall.ECONNRESET:
									{
										return retProtoMsg, errors.Wrap(constant.SystemErrConnectionReset, err, 0)
									}
								case syscall.ECONNABORTED:
									{
										return retProtoMsg, errors.Wrap(constant.SystemErrConnectionAborted, err, 0)
									}
								default:
									return retProtoMsg, errors.Wrap(constant.SystemInternalError, err, 0)
								}
							}
						}
					default:
						return retProtoMsg, errors.Wrap(constant.SystemInternalError, err, 0)
					}
				}
			default:
				return retProtoMsg, errors.Wrap(constant.SystemInternalError, err, 0)
			}

		}
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return retProtoMsg, errors.Errorf(constant.SystemInternalError, "postRequest StatusCode != 200,fullURL=%s,statusCode=%v", fullURL, resp.StatusCode())
	}

	version = string(resp.Header.Peek("v"))
	resBody := resp.Body()
	if isNeedDeserializerToMessage {
		if "2" == version {
			retProtoMsg, err = serializer.Bytes2protoMsg(resBody)
		} else {
			retProtoMsg, err = serializer.String2protoMsg(string(resBody))
		}

		if nil != err {
			return retProtoMsg, err
		}
	} else {
		if "2" == version {
			err = serializer.Bytes2Error(resBody)
		} else {
			err = serializer.String2Error(string(resBody))
		}

		if nil != err {
			return retProtoMsg, err
		}
	}

	return retProtoMsg, nil
}

// SyncCall publishes a request event to the event mesh and waiting response until server reply.
// this version is allow callers specify a timout(millisecond)
func SyncCall(message *msg.Message, timeout time.Duration) (*msg.Message, error) {
	if nil == message {
		return nil, ErrNullPointer
	}
	message.SetAppProperty(constant.To1, strconv.FormatInt(int64(timeout.Seconds()*1000), 10))

	protoMsg := model.MsgToProtocolMsg(message)
	retProtoMsg, err := post(protoMsg, SendRequestReplyMsgPath, true, timeout)
	if nil != err {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	retUserMessage := model.ProtocolMsgToMsg(&retProtoMsg)

	retUserMessage.DeleteProperty(constant.To1)
	retUserMessage.DeleteProperty(constant.To2)
	retUserMessage.DeleteProperty(constant.To3)

	return retUserMessage, nil
}

//ReplySemiSyncCall is use for requester initiate a semi synchronized call.
func ReplySemiSyncCall(msg *msg.Message) (err error) {
	if nil == msg {
		log.Errors("ReplySemiSyncCall msg=nil")
		return ErrNullPointer
	}
	log.Debugsf("ReplySemiSyncCall, msg=[%++v]", msg)
	protoMsg := model.MsgToProtocolMsg(msg)
	if _, err = post(protoMsg, ReplySemiSyncCallPath, false, time.Second*30); nil != err {
		return err
	}

	return nil
}

//Publish is use for publish a event to MESH
func Publish(msg *msg.Message) (err error) {
	if nil == msg {
		log.Errors("Publish msg=nil")
		return ErrNullPointer
	}
	log.Debugsf("publish message, msg.TopicAttribute=[%++v],msg.AppProps=[%s],msg.Body=[%s]",
		msg.TopicAttribute, msg.AppPropsToString(), string(msg.Body))

	protoMsg := model.MsgToProtocolMsg(msg)
	if _, err = post(protoMsg, SendTopicMsgPath, false, time.Second*30); nil != err {
		return err
	}

	return nil
}
