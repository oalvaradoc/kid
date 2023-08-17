package server_response

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"testing"
)

var interceptor = &Interceptor{}

func TestInterceptor_String(t *testing.T) {
	nameOfInterceptor := fmt.Sprintf("%s", interceptor)
	assert.Equal(t, nameOfInterceptor, constant.InterceptorServerResponse)
}

func TestInterceptor_PreHandle(t *testing.T) {
	ctx := context.Background()
	err := interceptor.PreHandle(ctx, nil)
	assert.Nil(t, err)
}

func TestInterceptor_PostHandle(t *testing.T) {
	ctx := context.Background()
	err := interceptor.PostHandle(ctx, nil, nil)
	assert.Nil(t, err)

	response := &msg.Message{}

	err = interceptor.PostHandle(ctx, nil, response)
	assert.Nil(t, err)

	response.Body = []byte("the response body")
	err = interceptor.PostHandle(ctx, nil, response)
	assert.Nil(t, err)

	handlerContexts := &contexts.HandlerContexts{
		Org:        "org",
		Az:         "az",
		Su:         "su",
		CommonSu:   "commonSu",
		NodeID:     "nodeID",
		ServiceID:  "serviceID",
		InstanceID: "instanceID",
		Wks:        "wks",
		Env:        "env",
		Lang:       "lang",
		SpanContexts: &contexts.SpanContexts{
			TraceID:             "trace ID",
			SpanID:              "spanID",
			ParentSpanID:        "parentSpanID",
			ReplyToAddress:      "replyToAddress",
			TimeoutMilliseconds: 100,
		},
		ResponseTemplate:            `{"code":"{{.errorCode}}","message":"{{.errorMsg}}","data":{{.data}}}`,
		ResponseAutoParseKeyMapping: nil,
		TransactionContexts:         nil,
	}
	newCtx := contexts.BuildContextFromParentWithHandlerContexts(ctx, handlerContexts)
	assert.NotEqual(t, ctx, newCtx)

	err = interceptor.PostHandle(newCtx, nil, response)
	assert.Nil(t, err)

	response.SetAppProperty(constant.ReturnErrorCode, "error code")
	err = interceptor.PostHandle(newCtx, nil, response)
	assert.Nil(t, err)
}
