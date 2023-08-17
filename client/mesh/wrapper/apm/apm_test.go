package apm

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/constant"
	commConstant "git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/config"
	"testing"
	"time"
)

var wrapper = &Wrapper{}

func TestWrapper_String(t *testing.T) {
	nameOfWrapper := fmt.Sprintf("%s", wrapper)
	assert.Equal(t, nameOfWrapper, constant.WrapperApm)
}

func TestWrapper_Before(t *testing.T) {
	ctx := context.Background()
	retCtx, err := wrapper.Before(ctx, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, ctx, retCtx)

	requestOptions := &client.RequestOptions{
		Codec:                                   nil,
		Context:                                 nil,
		TopicType:                               "",
		EventID:                                 "",
		Org:                                     "",
		Wks:                                     "",
		Env:                                     "",
		Su:                                      "",
		Version:                                 "",
		NodeID:                                  "",
		InstanceID:                              "",
		MaxRetryTimes:                           0,
		MaxWaitingTime:                          0,
		Timeout:                                 0,
		RetryWaitingTime:                        0,
		Header:                                  nil,
		Backoff:                                 nil,
		Retry:                                   nil,
		IsLocalCall:                             true,
		IsDMQEligible:                           false,
		IsPersistentDeliveryMode:                false,
		DeleteTransactionPropagationInformation: false,
		SkipResponseAutoParse:                   false,
		DisableMacroModel:                       false,
		IsSemiSyncCall:                          false,
		HTTPCall:                                false,
		Address:                                 "",
		ContentType:                             "",
		HTTPMethod:                              "",
	}
	retCtx, err = wrapper.Before(ctx, nil, requestOptions)
	assert.Nil(t, err)
	assert.Equal(t, ctx, retCtx)

	config.SetConfigs(&config.ServiceConfigs{
		ServerAddress:       "",
		Port:                0,
		CallbackPort:        0,
		CommType:            "",
		Service:             config.Service{},
		ClientSideStatusFSM: false,
		Log:                 config.Log{},
		Version:             0,
		Db:                  nil,
		Transaction:         config.Transaction{},
		Heartbeat:           config.Heartbeat{},
		Alert:               config.Alert{},
		Apm: config.Apm{
			Enable:                              true,
			PrintEmptyTraceIdRecordAtClientSide: false,
			Version:                             "v1",
			RootPath:                            "",
			FileRows:                            0,
		},
		Deployment:   config.Deployment{},
		Addressing:   config.Addressing{},
		Downstream:   nil,
		EventKey:     nil,
		GetFn:        nil,
		IsExistsFn:   nil,
		GetIntFn:     nil,
		GetInt32Fn:   nil,
		GetInt64Fn:   nil,
		GetFloat64Fn: nil,
		GetStringFn:  nil,
		GetBoolFn:    nil,
		GetUintFn:    nil,
		GetUint32Fn:  nil,
		GetUint64Fn:  nil,
	})
	retCtx, err = wrapper.Before(ctx, nil, requestOptions)
	assert.Nil(t, err)
	assert.NotEqual(t, ctx, retCtx)

	csStartTimestamp := retCtx.Value(commConstant.CsStartTimestamp)
	assert.NotNil(t, csStartTimestamp)
	t.Logf("client start timestamp:%++v", csStartTimestamp)
}

func TestWrapper_After(t *testing.T) {
	ctx := context.Background()
	retCtx, err := wrapper.After(ctx, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, ctx, retCtx)

	requestOptions := &client.RequestOptions{
		Codec:                                   nil,
		Context:                                 nil,
		TopicType:                               "",
		EventID:                                 "",
		Org:                                     "",
		Wks:                                     "",
		Env:                                     "",
		Su:                                      "",
		Version:                                 "",
		NodeID:                                  "",
		InstanceID:                              "",
		MaxRetryTimes:                           0,
		MaxWaitingTime:                          0,
		Timeout:                                 0,
		RetryWaitingTime:                        0,
		Header:                                  nil,
		Backoff:                                 nil,
		Retry:                                   nil,
		IsLocalCall:                             false,
		IsDMQEligible:                           false,
		IsPersistentDeliveryMode:                false,
		DeleteTransactionPropagationInformation: false,
		SkipResponseAutoParse:                   false,
		DisableMacroModel:                       false,
		IsSemiSyncCall:                          false,
		HTTPCall:                                false,
		Address:                                 "",
		ContentType:                             "",
		HTTPMethod:                              "",
	}
	retCtx, err = wrapper.After(ctx, nil, nil, requestOptions)
	assert.Nil(t, err)
	assert.Equal(t, ctx, retCtx)

	requestOptions.IsLocalCall = true
	retCtx, err = wrapper.After(ctx, nil, nil, requestOptions)
	assert.Nil(t, err)
	assert.Equal(t, ctx, retCtx)

	config.SetConfigs(&config.ServiceConfigs{
		ServerAddress:       "",
		Port:                0,
		CallbackPort:        0,
		CommType:            "",
		Service:             config.Service{},
		ClientSideStatusFSM: false,
		Log:                 config.Log{},
		Version:             0,
		Db:                  nil,
		Transaction:         config.Transaction{},
		Heartbeat:           config.Heartbeat{},
		Alert:               config.Alert{},
		Apm: config.Apm{
			Enable:                              true,
			PrintEmptyTraceIdRecordAtClientSide: false,
			Version:                             "v1",
			RootPath:                            "",
			FileRows:                            0,
		},
		Deployment:   config.Deployment{},
		Addressing:   config.Addressing{},
		Downstream:   nil,
		EventKey:     nil,
		GetFn:        nil,
		IsExistsFn:   nil,
		GetIntFn:     nil,
		GetInt32Fn:   nil,
		GetInt64Fn:   nil,
		GetFloat64Fn: nil,
		GetStringFn:  nil,
		GetBoolFn:    nil,
		GetUintFn:    nil,
		GetUint32Fn:  nil,
		GetUint64Fn:  nil,
	})
	retCtx, err = wrapper.After(ctx, nil, nil, requestOptions)
	assert.Equal(t, ctx, retCtx)
	assert.NotNil(t, err)
	t.Logf("error:%++v", err)

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
		ResponseTemplate:            "",
		ResponseAutoParseKeyMapping: nil,
		TransactionContexts:         nil,
	}
	ctx = contexts.BuildContextFromParentWithHandlerContexts(ctx, handlerContexts)
	retCtx, err = wrapper.After(ctx, nil, nil, requestOptions)
	assert.Equal(t, ctx, retCtx)
	assert.Nil(t, err)

	csStartTimestamp := time.Now()
	nCtx := context.WithValue(ctx, commConstant.CsStartTimestamp, csStartTimestamp)
	retCtx, err = wrapper.After(nCtx, nil, nil, requestOptions)
	assert.Equal(t, nCtx, retCtx)
	assert.Nil(t, err)
}
