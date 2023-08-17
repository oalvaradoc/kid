package handler

import (
	"context"
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/client/mesh"
	"git.multiverse.io/eventkit/kit/client/mesh/wrapper/trace"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/interceptor"
	"git.multiverse.io/eventkit/kit/interceptor/logging"
	"git.multiverse.io/eventkit/kit/sed/callback"
	"git.multiverse.io/eventkit/kit/wrapper"
	"testing"
)

func TestWithCallbackOptions(t *testing.T) {
	customOptions := CustomOptions{
		CustomOptionConfigs: make(map[string]interface{}, 0),
	}
	opt := WithCallbackOptions([]callback.Option{
		func(options *callback.Options) {
			options.EnableClientSideStatusFSM = true
			options.CallbackPort = 18082
			options.ServerAddress = "http://127.0.0.1:18080"
			options.CommType = "MESH"
			options.Port = 6060
		},
	}...)

	opt(&customOptions)
	assert.Equal(t, len(customOptions.CustomOptionConfigs), 1)
	assert.True(t, customOptions.CustomOptionConfigs[constant.ExtConfigCustomSedClientOptions] != nil)
	callbackOptions := customOptions.CustomOptionConfigs[constant.ExtConfigCustomSedClientOptions].([]callback.Option)

	t.Logf("callback options:%++v", callbackOptions)
	assert.Equal(t, len(callbackOptions), 1)
}

func TestWithDefaultInterceptors(t *testing.T) {
	customOptions := CustomOptions{
		CustomOptionConfigs: make(map[string]interface{}, 0),
	}
	opt := WithDefaultInterceptors([]interceptor.Interceptor{&logging.Interceptor{}}...)

	opt(&customOptions)
	assert.Equal(t, len(customOptions.CustomOptionConfigs), 1)
	assert.True(t, customOptions.CustomOptionConfigs[constant.ExtConfigCustomDefaultInterceptors] != nil)
	interceptors := customOptions.CustomOptionConfigs[constant.ExtConfigCustomDefaultInterceptors].([]interceptor.Interceptor)

	t.Logf("interceptors:%++v", interceptors)
	assert.Equal(t, len(interceptors), 1)
}

type MyCallbackHandleWrapper struct{}

func (m *MyCallbackHandleWrapper) PreHandle(ctx context.Context, in *msg.Message) (skip bool, outCtx context.Context, out *msg.Message, err error) {
	return false, ctx, in, nil
}

func (m *MyCallbackHandleWrapper) PostHandle(ctx context.Context, in *msg.Message) (out *msg.Message, err error) {
	return in, nil
}

func TestWithCallbackHandleWrapper(t *testing.T) {
	customOptions := CustomOptions{
		CustomOptionConfigs: make(map[string]interface{}, 0),
	}
	opt := WithCallbackHandleWrapper(&MyCallbackHandleWrapper{})

	opt(&customOptions)
	assert.Equal(t, len(customOptions.CustomOptionConfigs), 1)
	assert.True(t, customOptions.CustomOptionConfigs[constant.ExtConfigCustomCallbackHandleWrapper] != nil)
	callbackHandleWrapper := customOptions.CustomOptionConfigs[constant.ExtConfigCustomCallbackHandleWrapper].(CallbackHandleWrapper)

	t.Logf("callback handle wrapper:%++v", callbackHandleWrapper)
	assert.False(t, util.IsNil(callbackHandleWrapper))

}

func TestWithDefaultResponseTemplate(t *testing.T) {
	customOptions := CustomOptions{
		CustomOptionConfigs: make(map[string]interface{}, 0),
	}
	opt := WithDefaultResponseTemplate("this is a test response templage")
	opt(&customOptions)
	assert.Equal(t, len(customOptions.CustomOptionConfigs), 1)
	assert.True(t, customOptions.CustomOptionConfigs[constant.ExtConfigCustomResponseTemplate] != nil)
	customResponseTemplate := customOptions.CustomOptionConfigs[constant.ExtConfigCustomResponseTemplate].(string)

	t.Logf("custom response template:%s", customResponseTemplate)
	assert.Equal(t, customResponseTemplate, "this is a test response templage")
}

func TestWithDefaultEnableValidation(t *testing.T) {
	customOptions := CustomOptions{
		CustomOptionConfigs: make(map[string]interface{}, 0),
	}

	opt := WithDefaultEnableValidation(true)

	opt(&customOptions)

	assert.Equal(t, len(customOptions.CustomOptionConfigs), 1)
	assert.True(t, customOptions.CustomOptionConfigs[constant.ExtConfigDefaultEnableValidation] != nil)

	customValidationOptions := customOptions.CustomOptionConfigs[constant.ExtConfigDefaultEnableValidation].(*DefaultCustomValidationOptions)

	assert.True(t, customValidationOptions.CombineErrors)

}

func TestWithCallWrappers(t *testing.T) {
	customOptions := CustomOptions{
		CustomOptionConfigs: make(map[string]interface{}, 0),
	}
	opt := WithCallWrappers([]wrapper.Wrapper{&trace.Wrapper{}}...)

	opt(&customOptions)

	assert.Equal(t, len(customOptions.CustomOptionConfigs), 1)
	assert.True(t, customOptions.CustomOptionConfigs[constant.ExtConfigCustomCallWrappers] != nil)

	callWrappersOption := customOptions.CustomOptionConfigs[constant.ExtConfigCustomCallWrappers].(client.CallOption)

	assert.NotNil(t, callWrappersOption)
}

func TestWithResponseAutoParseKeyMapping(t *testing.T) {
	customOptions := CustomOptions{
		CustomOptionConfigs: make(map[string]interface{}, 0),
	}

	opt := WithResponseAutoParseKeyMapping(map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	})

	opt(&customOptions)

	assert.Equal(t, len(customOptions.CustomOptionConfigs), 1)
	assert.True(t, customOptions.CustomOptionConfigs[constant.ExtConfigCustomResponseAutoParseKeyMapping] != nil)
	responseAutoParseKeyMapping := customOptions.CustomOptionConfigs[constant.ExtConfigCustomResponseAutoParseKeyMapping].(map[string]string)

	t.Logf("response auto parse key mapping:%++v", responseAutoParseKeyMapping)
	assert.Equal(t, len(responseAutoParseKeyMapping), 3)
	assert.Equal(t, responseAutoParseKeyMapping["k1"], "v1")
	assert.Equal(t, responseAutoParseKeyMapping["k2"], "v2")
	assert.Equal(t, responseAutoParseKeyMapping["k3"], "v3")
	assert.Equal(t, responseAutoParseKeyMapping["k4"], "")
}

func TestWithCallInterceptors(t *testing.T) {
	customOptions := CustomOptions{
		CustomOptionConfigs: make(map[string]interface{}, 0),
	}
	opt := WithCallInterceptors([]interceptor.Interceptor{&logging.Interceptor{}}...)
	opt(&customOptions)

	assert.Equal(t, len(customOptions.CustomOptionConfigs), 1)
	assert.True(t, customOptions.CustomOptionConfigs[constant.ExtConfigCustomCallInterceptors] != nil)
	callOption := customOptions.CustomOptionConfigs[constant.ExtConfigCustomCallInterceptors].(client.CallOption)

	t.Logf("call option:%++v", callOption)
	assert.NotNil(t, callOption)
}

func TestWithClient(t *testing.T) {
	customOptions := CustomOptions{
		CustomOptionConfigs: make(map[string]interface{}, 0),
	}
	opt := WithClient(mesh.NewMeshClient())
	opt(&customOptions)

	assert.Equal(t, len(customOptions.CustomOptionConfigs), 1)
	assert.True(t, customOptions.CustomOptionConfigs[constant.ExtConfigCustomClient] != nil)
	client := customOptions.CustomOptionConfigs[constant.ExtConfigCustomClient].(client.Client)

	t.Logf("client:%++v", client)
	assert.NotNil(t, client)
}
