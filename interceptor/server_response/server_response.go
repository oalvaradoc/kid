package server_response

import (
	"bytes"
	"context"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/log"
	"text/template"
)

// Interceptor is an interceptor for wrap the server response body
type Interceptor struct{}

// PreHandle does noting
func (l *Interceptor) PreHandle(ctx context.Context, request *msg.Message) error {
	// DO NOTHING
	return nil
}

// PostHandle warps the response according the response template
func (l *Interceptor) PostHandle(ctx context.Context, request *msg.Message, response *msg.Message) error {
	if nil == response {
		log.Debug(ctx, "The response is empty, skip the following logic!")
		return nil
	}
	if "" == response.GetAppPropertyIgnoreCaseSilence(constant.ReturnErrorCode) {
		responseTemplate := constant.DefaultResponseTemplate
		handlerContexts := contexts.HandlerContextsFromContext(ctx)
		if nil != handlerContexts && "" != handlerContexts.ResponseTemplate {
			responseTemplate = handlerContexts.ResponseTemplate
		}
		// wrapper response
		buf := new(bytes.Buffer)
		res, e := template.New("response").Parse(responseTemplate)
		if nil != e {
			return e
		}
		wrapperMap := make(map[string]string)
		wrapperMap["errorCode"] = errors.GetFinalErrorCode(constant.Success)
		wrapperMap["errorMsg"] = ""
		wrapperMap["errorStack"] = ""
		if len(response.Body) > 0 {
			wrapperMap["data"] = string(response.Body)
		} else {
			wrapperMap["data"] = "null"
		}
		res.Execute(buf, wrapperMap)
		response.Body = buf.Bytes()
	}

	return nil
}

func (l Interceptor) String() string {
	return constant.InterceptorServerResponse
}
