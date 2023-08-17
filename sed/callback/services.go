package callback

import (
	"encoding/json"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/model"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
	"github.com/valyala/fasthttp"
	"net/http"
	"reflect"
	"time"
)

const (
	// SendRequestReplyMsgPath is an url path for request/reply message
	SendRequestReplyMsgPath = "/v1/send-request-reply-msg"

	// SendReplyMsgPath is an url path for reply message
	SendReplyMsgPath = "/v1/send-reply-msg"

	// SendTopicMsgPath is an url path for publish message
	SendTopicMsgPath = "/v1/send-topic-msg"

	// SendQueueAckPath is an url path for ack a queue message
	SendQueueAckPath = "/v1/send-queue-ack"

	// ReplySemiSyncCallPath is an url path for semi-sync call message
	ReplySemiSyncCallPath = "/v1/reply-semi-sync-call"
	// ServerStatusGetPath is an URL path for get sed server status
	ServerStatusGetPath = "/v2/status"
)

// define global variable for storing sed server address(http address)
// default address is "http://127.0.0.1:18080"
var serverAddr = "http://127.0.0.1:18080"

// define global variable for storing sed server address(unix socket address)
// default address is "http://unix"
//var unixSocketHttpAddr = "http://unix"

// use for tcp/ip socket
var httpClient *http.Client

// use for UNIX socket
//var httpClient2 *http.Client

// InitSedClient init httpClient & httpClient2, create handle for http transport
func InitSedClient() {
	transport := http.Transport{
		// set disable keep alive to false, request can reuse the connection
		// default is true, if true, disables HTTP keep-alives and
		// will only use the connection to the server for a single
		// HTTP request.
		DisableKeepAlives: false,
	}
	httpClient = &http.Client{
		Transport: &transport,
	}
}

// SetSedServerAddr is for application to specify sed server address
// if addr is empty, default value "http://127.0.0.1:18080" will be set.
func SetSedServerAddr(addr string) {
	if addr != "" {
		log.Infosf("Sed server address:%s", addr)
		serverAddr = addr
	}
}

// isValidPointer determine if a response pointer is valid
func isValidPointer(response interface{}) error {
	if nil == response {
		return errors.New(constant.SystemInternalError, "input param should not be nil")
	}

	if kind := reflect.ValueOf(response).Type().Kind(); kind != reflect.Ptr {
		return errors.New(constant.SystemInternalError, "input param should be pointer type")
	}

	return nil
}

// postRequest is a internal function for client post request to server
// and it's returns responses where from server endpoint
// fixed http request timeout is 30 seconds
func postRequest(path string, request interface{}, response interface{}) error {
	commonResp := model.CommonResponse{}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(constant.SystemInternalError, err, 0)
	}
	var fullURL string
	fullURL = serverAddr + path

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(http.MethodPost)
	req.SetRequestURI(fullURL)
	req.SetBody(requestBytes)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	createClientIfNecessary()
	if err = clientWithConnectTimeout.DoTimeout(req, resp, 30*time.Second); err != nil {
		log.Errorsf("post to server failed err=%v", err)
		if fasthttp.ErrTimeout == err {
			return errors.Wrap(constant.SystemRemoteCallTimeout, err, 0)
		} else if fasthttp.ErrConnectionClosed == err {
			return errors.Errorf(constant.SystemErrConnectionClosed, "Connection closed, url=[%s]", fullURL)
		} else {
			return errors.Wrap(constant.SystemInternalError, err, 0)
		}
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return errors.Errorf(constant.SystemInternalError, "postRequest StatusCode != 200,fullURL=%s,statusCode=%v", fullURL, resp.StatusCode())
	}

	resBody := resp.Body()

	err = json.Unmarshal(resBody, &commonResp)
	if err != nil {
		return errors.Wrap(constant.SystemInternalError, err, 0)
	}

	if commonResp.ErrorCode != 0 {
		return errors.Errorf("%v", commonResp.ErrorMsg)
	}

	//do not need response
	if err := isValidPointer(resp); err != nil {
		return nil
	}

	datas, err := json.Marshal(commonResp.Data)
	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "postRequest marsh data err=%v", err)
	}

	err = json.Unmarshal([]byte(datas), response)

	return errors.Wrap(constant.SystemInternalError, err, 0)

}

// glsPostRequest is a internal function for gls client post request to server
// and it's returns responses where from server endpoint
// fixed http request timeout is 30 seconds
func glsPostRequest(path string, request interface{}, response interface{}) error {
	commonResp := model.CommonResponse{}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(constant.SystemInternalError, err, 0)
	}
	var fullURL string
	fullURL = serverAddr + path

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(http.MethodPost)
	req.SetRequestURI(fullURL)
	req.SetBody(requestBytes)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	log.Debugsf("request:[%++v]", request)
	createClientIfNecessary()
	if err = clientWithConnectTimeout.DoTimeout(req, resp, 30*time.Second); err != nil {
		log.Errorsf("post to server failed err=%v", err)
		if fasthttp.ErrTimeout == err {
			return errors.Wrap(constant.SystemRemoteCallTimeout, err, 0)
		} else if fasthttp.ErrConnectionClosed == err {
			return errors.Errorf(constant.SystemErrConnectionClosed, "Connection closed, url=[%s]", fullURL)
		} else {
			return errors.Wrap(constant.SystemInternalError, err, 0)
		}
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return errors.Errorf(constant.SystemInternalError, "postRequest StatusCode != 200,fullURL=%s,statusCode=%v", fullURL, resp.StatusCode())
	}

	resBody := resp.Body()
	log.Debugsf("response body:[%++v]", string(resBody))
	err = json.Unmarshal(resBody, &commonResp)
	if err != nil {
		return errors.Wrap(constant.SystemInternalError, err, 0)
	}

	if commonResp.ErrorCode != 0 {
		return errors.Errorf("%v", commonResp.ErrorMsg)
	}

	//do not need response
	if err := isValidPointer(resp); err != nil {
		return nil
	}

	datas, err := json.Marshal(commonResp.Data)
	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "postRequest marsh data err=%v", err)
	}

	err = json.Unmarshal([]byte(datas), response)

	return err

}
