package callback

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/model"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/common/protocol"
	"git.multiverse.io/eventkit/kit/common/serializer"
	"git.multiverse.io/eventkit/kit/common/status"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
	"github.com/buaazp/fasthttprouter"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"runtime/debug"
)

// Executor is an executor interface after a SED server receives a message.
// Any service that needs to connect to the SED server needs to implement the methods in this interface.
// When the Executor starts, it will first call the Init method to initialize, and the Init method will only be executed once.
// When the Executor receives the message, it will execute the handle method.
// The handle method needs to execute all the business logic and return the processing result.
type Executor interface {
	Init() error
	Handle(ctx context.Context, in *msg.Message) (out *msg.Message, err error)
	CallbackOptions() *Options
	ResponseTemplate() string
	Destroy() error
}

var (
	callbackExecutor Executor
	server           *fasthttp.Server
)

// RegisterCallbackExecutor provides an application registration callback method
func RegisterCallbackExecutor(ce Executor) {
	status.ClientVersion = constant.EventHandlerVersion
	status.ClientProtocolLevel = constant.ProtocolLevel
	status.SetClientStatus(status.ClientStartingInit)
	ce.Init()

	SetSedServerAddr(ce.CallbackOptions().ServerAddress)
	InitSedClient()

	callbackExecutor = ce

	go StartCallbackServer(ce.CallbackOptions().CallbackPort)

	if ce.CallbackOptions().EnableClientSideStatusFSM {
		if err := StartClientSideStatusFSM(); nil != err {
			panic(fmt.Sprintf("start client side FSM failed, error=%++v", err))
		}
	} else {
		if err := RunHookHandleFunc(); nil != err {
			panic(fmt.Sprintf("run hook handle function failed, error=%++v", err))
		}
		log.Infos("Skip client side status FSM, set server status to `STARTED`...")
		status.SetServerStatus(status.ServerStarted)
		status.SetClientStatus(status.ClientStarted)
	}
}

// callbackHandlerForFastHTTP Provides a fasthttp mode callback method
func callbackHandlerForFastHTTP(ctx *fasthttp.RequestCtx) {
	var err error
	var version string
	var resMessage protocol.ProtoMessage

	defer func() {
		var res []byte
		ctx.Response.Header.Set("v", version)

		if e := recover(); e != nil {
			err := errors.Errorf(constant.SystemInternalError, "callbackHandlerForFastHTTP catch painc:%v", e)
			if "2" == version {
				res = serializer.Error2BytesCompatible(err, true)
			} else {
				res = []byte(serializer.Error2StringCompatible(err, true))
			}

			log.Errorsf("callbackHandlerForFastHTTP painc, error=%++v, call stack:[%s]", e, string(debug.Stack()))
		} else {
			if "2" == version {
				res = serializer.ProtoMsg2BytesCompatible(&resMessage, err, true)
			} else {
				res = []byte(serializer.ProtoMsg2StringCompatible(&resMessage, err, true))
			}
		}

		ctx.Write(res)
		return
	}()

	version = string(ctx.Request.Header.Peek("v"))
	requestBody := ctx.PostBody()
	var protoMsg protocol.ProtoMessage

	if "2" == version {
		protoMsg, err = serializer.Bytes2protoMsg(requestBody)
	} else {
		protoMsg, err = serializer.String2protoMsg(string(requestBody))
	}

	if nil != err {
		err = errors.Errorf(constant.SystemInternalError, "callbackHandlerForFastHTTP:proto message Unmarshal failed, error=%v", err)
		return
	}
	message := model.ProtocolMsgToMsg(&protoMsg)

	context := context.Background()
	reply, err := callbackExecutor.Handle(context, message)
	if nil != err {
		lang := message.JudgeUserLang()
		reply = msg.WrapperErrorResponse(err, lang, callbackExecutor.ResponseTemplate(), nil)
		err = nil
		//return
	}

	if nil != reply {
		resMessage = model.MsgToProtocolMsg(reply)
	} else if protoMsg.NeedReply {
		err = errors.Errorf(constant.SystemInternalError, "callbackHandlerForFastHTTP:reply message is nil,please check")
	} else {
		retMsg := &msg.Message{}
		resMessage = model.MsgToProtocolMsg(retMsg)
	}

	return
}

// getClientStatus gets the current client status (used for fast http)
func getClientStatus(ctx *fasthttp.RequestCtx) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	res := status.Response{
		Status:        status.GetClientStatus(),
		Version:       status.ClientVersion,
		ProtocolLevel: status.ClientProtocolLevel,
	}

	resBytes, err := json.Marshal(res)
	if nil != err {
		log.Errorsf("Marshal Response data failed, error=%++v", err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	ctx.Write(resBytes)
}

// StartCallbackServer start callback service
// default callback http port is 18082
// callback url path is "/v1/newmsg"
func StartCallbackServer(port int) {
	//http.HandleFunc("/v1/newmsg", callbackHandlerForHttp)
	if port == 0 {
		port = 18082
	}

	serverAddr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Infosf("Start callback endpoint, listen addr=%s", serverAddr)
	router := fasthttprouter.New()
	router.POST("/v1/newmsg", callbackHandlerForFastHTTP)
	router.GET("/v1/client/status", getClientStatus)

	server = &fasthttp.Server{
		Handler:                       router.Handler,
		MaxRequestBodySize:            1024 * 1024 * 1024,
		DisableHeaderNamesNormalizing: true,
	}
	if err := server.ListenAndServe(serverAddr); err != nil {
		log.Errors("startUserMsgCallbackServer: start fasthttp failed:", err.Error())
		panic(err)
	}
}

// ShutdownClientListen shutdown the client listener
func ShutdownClientListen() {
	if nil != server {
		log.Infos("Start shutdown the client listen...")
		if err := server.Shutdown(); nil != err {
			log.Errorsf("Shutdown the client listen failed, error = %v", err)
		} else {
			log.Infos("Shutdown the client listen successfully!")
		}
	}
}
