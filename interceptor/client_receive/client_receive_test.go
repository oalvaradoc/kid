package client_receive

import (
	"context"
	"encoding/xml"
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
	assert.Equal(t, nameOfInterceptor, constant.InterceptorClientReceive)
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
	response.SetAppProperty(constant.ReturnStatus, "F")
	err = interceptor.PostHandle(ctx, nil, response)
	assert.NotNil(t, err)

	response.DeleteProperty(constant.ReturnStatus)
	response.SetAppProperty(constant.ReturnErrorCode, constant.SystemInternalError)
	response.SetAppProperty(constant.ReturnErrorMsg, "Failed to execute target service")
	err = interceptor.PostHandle(ctx, nil, response)
	assert.NotNil(t, err)

	response.DeleteProperty(constant.ReturnErrorCode)
	response.DeleteProperty(constant.ReturnErrorMsg)

	nCtx := context.WithValue(ctx, constant.SkipResponseAutoParseKeyMappingFlagKey, true)
	err = interceptor.PostHandle(nCtx, nil, response)
	assert.Nil(t, err)

	err = interceptor.PostHandle(ctx, nil, response)
	assert.Nil(t, err)

	response.Body = []byte(`{"errorCode":"0","errorMsg":"", "response":{"key1":"value1", "key2":"value2"}}`)
	err = interceptor.PostHandle(ctx, nil, response)
	assert.Nil(t, err)
	t.Logf("response data:%++v", string(response.Body))

	response.Body = []byte(`{"code":"0","errorMsg":"", "response":{"key1":"value1", "key2":"value2"}}`)
	err = interceptor.PostHandle(ctx, nil, response)
	assert.NotNil(t, err)
	t.Logf("the error of return is:%++v", err)

	response.Body = []byte(`{"errorCode":"0","message":"", "response":{"key1":"value1", "key2":"value2"}}`)
	err = interceptor.PostHandle(ctx, nil, response)
	assert.NotNil(t, err)
	t.Logf("the error of return is:%++v", err)

	response.Body = []byte(`{"errorCode":"0","errorMsg":"", "data":{"key1":"value1", "key2":"value2"}}`)
	err = interceptor.PostHandle(ctx, nil, response)
	assert.NotNil(t, err)
	t.Logf("the error of return is:%++v", err)

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
		ResponseTemplate: `{"code":"{{.errorCode}}","message":"{{.errorMsg}}","data":{{.data}}}`,
		ResponseAutoParseKeyMapping: map[string]string{
			constant.ResponseAutoParseTypeMappingKey: constant.ResponseAutoDefaultParseType,
			constant.ErrorCodeMappingKey:             constant.DefaultErrorCodeKey,
			constant.ErrorMsgMappingKey:              constant.DefaultErrorMsgKey,
			constant.ResponseDataBodyMappingKey:      constant.DefaultDataBodyKey,
		},
		TransactionContexts: nil,
	}
	newCtx := contexts.BuildContextFromParentWithHandlerContexts(ctx, handlerContexts)
	assert.NotEqual(t, ctx, newCtx)
	response.Body = []byte(`{"errorCode":"0","errorMsg":"", "response":{"key1":"value1", "key2":"value2"}}`)
	err = interceptor.PostHandle(newCtx, nil, response)
	assert.Nil(t, err)
	t.Logf("response data:%++v", string(response.Body))

	handlerContexts.ResponseAutoParseKeyMapping = map[string]string{
		constant.ResponseAutoParseTypeMappingKey: constant.ResponseAutoDefaultParseType,
		constant.ErrorCodeMappingKey:             "code",
		constant.ErrorMsgMappingKey:              "message",
		constant.ResponseDataBodyMappingKey:      "data",
	}
	response.Body = []byte(`{"code":"0","message":"", "data":{"key1":"value1", "key2":"value2"}}`)
	err = interceptor.PostHandle(newCtx, nil, response)
	assert.Nil(t, err)
	t.Logf("response data:%++v", string(response.Body))

	newCtx = context.WithValue(ctx, constant.ResponseAutoParseKeyMappingKey, map[string]string{
		constant.ResponseAutoParseTypeMappingKey: constant.ResponseAutoDefaultParseType,
		constant.ErrorCodeMappingKey:             "code",
		constant.ErrorMsgMappingKey:              "message",
		constant.ResponseDataBodyMappingKey:      "data",
	})
	response.Body = []byte(`{"code":"0","message":"", "data":{"key1":"value1", "key2":"value2"}}`)
	err = interceptor.PostHandle(newCtx, nil, response)
	assert.Nil(t, err)
	t.Logf("response data:%++v", string(response.Body))

	// test XML
	newCtx = context.WithValue(ctx, constant.ResponseAutoParseKeyMappingKey, map[string]string{
		constant.ResponseAutoParseTypeMappingKey: constant.ResponseAutoParseTypeXML,
		constant.ErrorCodeMappingKey:             "code",
		constant.ErrorMsgMappingKey:              "message",
		constant.ResponseDataBodyMappingKey:      "data",
	})
	response.Body = []byte(`<root><code>0</code><message></message><data><field1>v1</field1><field2>v2</field2></data></root>`)
	err = interceptor.PostHandle(newCtx, nil, response)
	assert.Nil(t, err)
	t.Logf("response data:%++v", string(response.Body))

	res := &Response{}
	xml.Unmarshal(response.Body, res)
	t.Logf("response:%++v", res)

	handlerContexts.ResponseAutoParseKeyMapping = map[string]string{
		constant.ResponseAutoParseTypeMappingKey: constant.ResponseAutoParseTypeXML,
		constant.ErrorCodeMappingKey:             "code",
		constant.ErrorMsgMappingKey:              "message",
		constant.ResponseDataBodyMappingKey:      "data",
	}
	newCtx = contexts.BuildContextFromParentWithHandlerContexts(ctx, handlerContexts)
	assert.NotEqual(t, ctx, newCtx)
	response.Body = []byte(`<root><code>0</code><message></message><data><field1>v1</field1><field2>v2</field2></data></root>`)
	err = interceptor.PostHandle(newCtx, nil, response)
	assert.Nil(t, err)
	t.Logf("response data:%++v", string(response.Body))
	res = &Response{}
	xml.Unmarshal(response.Body, res)
	t.Logf("response:%++v", res)

	response.Body = []byte(`<root><errorCode>0</errorCode><message></message><data><field1>v1</field1><field2>v2</field2></data></root>`)
	err = interceptor.PostHandle(newCtx, nil, response)
	assert.NotNil(t, err)
	t.Logf("the error of result:%++v", err)

	response.Body = []byte(`<root><code>0</code><errorMessage></errorMessage><data><field1>v1</field1><field2>v2</field2></data></root>`)
	err = interceptor.PostHandle(newCtx, nil, response)
	assert.NotNil(t, err)
	t.Logf("the error of result:%++v", err)

	response.Body = []byte(`<root><code>0</code><message></message><response><field1>v1</field1><field2>v2</field2></response></root>`)
	err = interceptor.PostHandle(newCtx, nil, response)
	assert.NotNil(t, err)
	t.Logf("the error of result:%++v", err)
}

type Response struct {
	Field1 string `xml:"field1"`
	Field2 string `xml:"field2"`
}
