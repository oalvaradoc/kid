package router

import (
	"git.multiverse.io/eventkit/kit/codec"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/compensable"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/base"
	"git.multiverse.io/eventkit/kit/interceptor"
	"git.multiverse.io/eventkit/kit/interceptor/apm"
	"git.multiverse.io/eventkit/kit/interceptor/client_receive"
	"git.multiverse.io/eventkit/kit/interceptor/logging"
	"git.multiverse.io/eventkit/kit/interceptor/server_response"
	"git.multiverse.io/eventkit/kit/interceptor/transaction"
	"reflect"
	"testing"
)

type TestEventHandler struct {
	base.Handler
}

func TestOptions(t *testing.T) {
	o := &Options{}
	opt := WithRegisterType(1)
	opt(o)
	assert.Equal(t, o.HandlerOptions.RegisterType, 1)

	opt = Compensable(&compensable.Compensable{
		TryMethod:     "try method",
		ConfirmMethod: "test confirm method",
		CancelMethod:  "test cancel method",
		ServiceName:   "test service name",
		IsPropagator:  true,
	})

	opt(o)
	assert.Equal(t, o.Compensable.ServiceName, "test service name")
	assert.Equal(t, o.HandlerOptions.HandlerMethodName, "try method")
	assert.Equal(t, o.Compensable.ConfirmMethod, "test confirm method")
	assert.Equal(t, o.Compensable.CancelMethod, "test cancel method")
	assert.True(t, o.Compensable.IsPropagator)

	opt = WithHandlerName("test handler name")
	opt(o)
	assert.Equal(t, o.HandlerOptions.HandlerName, "test handler name")

	reflectType := reflect.TypeOf(&TestEventHandler{})
	opt = WithHandlerReflectType(reflectType)
	opt(o)
	assert.Equal(t, o.HandlerOptions.HandlerReflectType, reflectType)

	responseTemplate := "this is a test response template"
	responseDataWhenError := `{"code":"12345", "message":"failed to run..."}`
	opt = WithResponseTemplate(responseTemplate, responseDataWhenError)
	opt(o)

	assert.Equal(t, o.HandlerOptions.ResponseTemplate, responseTemplate)
	assert.Equal(t, o.HandlerOptions.ResponseDataWhenErrorForResponseTemplate, responseDataWhenError)

	responseTemplate = "new test response template"
	responseDataWhenError = `{"code":"0001", "message":"error response"}`
	responseHeader := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}
	opt = WithResponseTemplateAndHeader(responseTemplate, responseHeader, responseDataWhenError)
	opt(o)
	assert.Equal(t, o.HandlerOptions.ResponseTemplate, responseTemplate)
	assert.Equal(t, o.HandlerOptions.ResponseDataWhenErrorForResponseTemplate, responseDataWhenError)
	m := &msg.Message{}
	o.HandlerOptions.CustomErrorWrapperFn("", "", m)
	assert.True(t, len(m.GetAppProps()) > 0)

	body := []byte("test body")
	f := func(errorCode, errorMessage string, response *msg.Message) {
		if nil != response {
			response.Body = body
		}
	}

	responseTemplate = "new test response template2"
	responseDataWhenError = `{"code":"0001", "message":"error response2"}`
	opt = WithResponseAndCustomErrorWrapperFn(responseTemplate, f, responseDataWhenError)
	opt(o)
	assert.Equal(t, o.HandlerOptions.ResponseTemplate, responseTemplate)
	assert.Equal(t, o.HandlerOptions.ResponseDataWhenErrorForResponseTemplate, responseDataWhenError)
	o.HandlerOptions.CustomErrorWrapperFn("", "", m)
	assert.Equal(t, m.Body, body)

	testMethodName := "test method"
	opt = Method(testMethodName)
	opt(o)
	assert.Equal(t, o.HandlerOptions.HandlerMethodName, testMethodName)

	urlPath := "http://127.0.0.1:8080/test"
	opt = HandlePost(urlPath)
	opt(o)
	assert.Equal(t, o.HandlerOptions.HTTPMethod, constant.HTTPMethodPost)
	assert.Equal(t, o.HandlerOptions.URLPath, urlPath)

	urlPath = "http://127.0.0.1:8080/test1"
	opt = HandleDelete(urlPath)
	opt(o)
	assert.Equal(t, o.HandlerOptions.HTTPMethod, constant.HTTPMethodDelete)
	opt(o)
	assert.Equal(t, o.HandlerOptions.URLPath, urlPath)

	urlPath = "http://127.0.0.1:8080/test2"
	opt = HandleGet(urlPath)
	opt(o)
	assert.Equal(t, o.HandlerOptions.HTTPMethod, constant.HTTPMethodGet)
	opt(o)
	assert.Equal(t, o.HandlerOptions.URLPath, urlPath)

	urlPath = "http://127.0.0.1:8080/test3"
	opt = HandleOptions(urlPath)
	opt(o)
	assert.Equal(t, o.HandlerOptions.HTTPMethod, constant.HTTPMethodOptions)
	opt(o)
	assert.Equal(t, o.HandlerOptions.URLPath, urlPath)

	urlPath = "http://127.0.0.1:8080/test4"
	opt = HandleHead(urlPath)
	opt(o)
	assert.Equal(t, o.HandlerOptions.HTTPMethod, constant.HTTPMethodHead)
	opt(o)
	assert.Equal(t, o.HandlerOptions.URLPath, urlPath)

	urlPath = "http://127.0.0.1:8080/test5"
	opt = HandlePatch(urlPath)
	opt(o)
	assert.Equal(t, o.HandlerOptions.HTTPMethod, constant.HTTPMethodPatch)
	opt(o)
	assert.Equal(t, o.HandlerOptions.URLPath, urlPath)

	urlPath = "http://127.0.0.1:8080/test6"
	opt = HandlePut(urlPath)
	opt(o)
	assert.Equal(t, o.HandlerOptions.HTTPMethod, constant.HTTPMethodPut)
	opt(o)
	assert.Equal(t, o.HandlerOptions.URLPath, urlPath)

	opt = EnableValidation(true)
	opt(o)
	assert.True(t, o.HandlerOptions.EnableValidation)
	assert.NotNil(t, o.HandlerOptions.CustomValidationOptions)
	assert.True(t, o.HandlerOptions.CustomValidationOptions.CombineErrors)

	opt = DisableValidation()
	opt(o)
	assert.False(t, o.HandlerOptions.EnableValidation)
	assert.NotNil(t, o.HandlerOptions.CustomValidationOptions)
	assert.False(t, o.HandlerOptions.CustomValidationOptions.CombineErrors)

	opt = WithHandlerMethodInParams([]reflect.Type{reflectType})
	opt(o)
	assert.Equal(t, o.HandlerOptions.HandlerMethodInParams, []reflect.Type{reflectType})

	opt = WithHandlerMethodOutParams([]reflect.Type{reflectType})
	opt(o)
	assert.Equal(t, o.HandlerOptions.HandlerMethodOutParams, []reflect.Type{reflectType})

	interceptors := []interceptor.Interceptor{&logging.Interceptor{}, &apm.Interceptor{}}
	opt = WithInterceptors(interceptors...)
	opt(o)
	assert.Equal(t, o.HandlerOptions.Interceptors, interceptors)

	newInterceptors := []interceptor.Interceptor{&logging.Interceptor{}, &apm.Interceptor{}, &transaction.Interceptor{}}
	opt = AddInterceptor(&transaction.Interceptor{})
	opt(o)
	assert.Equal(t, o.HandlerOptions.Interceptors, newInterceptors)

	newInterceptors = []interceptor.Interceptor{
		&logging.Interceptor{},
		&apm.Interceptor{},
		&transaction.Interceptor{},
		&server_response.Interceptor{},
		&client_receive.Interceptor{},
	}
	opt = AddInterceptors(&server_response.Interceptor{}, &client_receive.Interceptor{})
	opt(o)
	assert.Equal(t, o.HandlerOptions.Interceptors, newInterceptors)

	newInterceptors = []interceptor.Interceptor{
		&logging.Interceptor{},
		&transaction.Interceptor{},
		&server_response.Interceptor{},
		&client_receive.Interceptor{},
	}

	opt = RemoveInterceptor(constant.InterceptorApm)
	opt(o)
	assert.Equal(t, o.HandlerOptions.Interceptors, newInterceptors)

	newInterceptors = []interceptor.Interceptor{
		&logging.Interceptor{},
		&client_receive.Interceptor{},
	}
	opt = RemoveInterceptors(constant.InterceptorTransaction, constant.InterceptorServerResponse)
	opt(o)
	assert.Equal(t, o.HandlerOptions.Interceptors, newInterceptors)

	myCodec := codec.BuildXMLCodec()
	opt = WithCodec(myCodec)
	opt(o)
	assert.Equal(t, o.HandlerOptions.Codec, myCodec)
}
