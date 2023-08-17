package audit

import (
	"fmt"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/constant"
	"testing"
)

func TestInterceptor_String(t *testing.T) {
	interceptor := NewAuditInterceptor(nil, nil)
	nameOfInterceptor := fmt.Sprintf("%s", interceptor)
	assert.Equal(t, nameOfInterceptor, constant.InterceptorAudit)
}

//
//func TestWrapper_Before(t *testing.T) {
//	ctx := context.Background()
//	err := interceptor.PreHandle(ctx, nil)
//	assert.Nil(t, err)
//
//	err = interceptor.PreHandle(ctx, nil)
//	assert.Nil(t, err)
//
//	config.SetConfigs(&config.ServiceConfigs{
//		ServerAddress:       "",
//		Port:                0,
//		CallbackPort:        0,
//		CommType:            "",
//		Service:             config.Service{},
//		ClientSideStatusFSM: false,
//		Log:                 config.Log{},
//		Version:             0,
//		Db:                  nil,
//		Transaction:         config.Transaction{},
//		Heartbeat:           config.Heartbeat{},
//		Alert:               config.Alert{},
//		Apm: config.Apm{
//			Enable:                              true,
//			PrintEmptyTraceIdRecordAtClientSide: false,
//			Version:                             false,
//			RootPath:                            "",
//			FileRows:                            0,
//		},
//		Deployment:   config.Deployment{},
//		Addressing:   config.Addressing{},
//		Downstream:   nil,
//		EventKey:     nil,
//		GetFn:        nil,
//		IsExistsFn:   nil,
//		GetIntFn:     nil,
//		GetInt32Fn:   nil,
//		GetInt64Fn:   nil,
//		GetFloat64Fn: nil,
//		GetStringFn:  nil,
//		GetBoolFn:    nil,
//		GetUintFn:    nil,
//		GetUint32Fn:  nil,
//		GetUint64Fn:  nil,
//	})
//	err = interceptor.PreHandle(ctx, nil)
//	assert.Nil(t, err)
//}
//
//func TestWrapper_After(t *testing.T) {
//	ctx := context.Background()
//	err := interceptor.PostHandle(ctx, nil, nil)
//	assert.Nil(t, err)
//
//	requestOptions := &client.RequestOptions{
//		Codec:                                   nil,
//		Context:                                 nil,
//		TopicType:                               "",
//		EventID:                                 "",
//		Org:                                     "",
//		Wks:                                     "",
//		Env:                                     "",
//		Su:                                      "",
//		Version:                                 "",
//		NodeID:                                  "",
//		InstanceID:                              "",
//		MaxRetryTimes:                           0,
//		MaxWaitingTime:                          0,
//		Timeout:                                 0,
//		RetryWaitingTime:                        0,
//		Header:                                  nil,
//		Backoff:                                 nil,
//		Retry:                                   nil,
//		IsLocalCall:                             false,
//		IsDMQEligible:                           false,
//		IsPersistentDeliveryMode:                false,
//		DeleteTransactionPropagationInformation: false,
//		SkipResponseAutoParse:                   false,
//		DisableMacroModel:                       false,
//		IsSemiSyncCall:                          false,
//		HTTPCall:                                false,
//		Address:                                 "",
//		ContentType:                             "",
//		HTTPMethod:                              "",
//	}
//	err = interceptor.PostHandle(ctx, nil, nil)
//	assert.Nil(t, err)
//
//	requestOptions.IsLocalCall = true
//	err = interceptor.PostHandle(ctx, nil, nil)
//	assert.Nil(t, err)
//
//	config.SetConfigs(&config.ServiceConfigs{
//		ServerAddress:       "",
//		Port:                0,
//		CallbackPort:        0,
//		CommType:            "",
//		Service:             config.Service{},
//		ClientSideStatusFSM: false,
//		Log:                 config.Log{},
//		Version:             0,
//		Db:                  nil,
//		Transaction:         config.Transaction{},
//		Heartbeat:           config.Heartbeat{},
//		Alert:               config.Alert{},
//		Apm: config.Apm{
//			Enable:                              true,
//			PrintEmptyTraceIdRecordAtClientSide: false,
//			Version:                             false,
//			RootPath:                            "",
//			FileRows:                            0,
//		},
//		Deployment:   config.Deployment{},
//		Addressing:   config.Addressing{},
//		Downstream:   nil,
//		EventKey:     nil,
//		GetFn:        nil,
//		IsExistsFn:   nil,
//		GetIntFn:     nil,
//		GetInt32Fn:   nil,
//		GetInt64Fn:   nil,
//		GetFloat64Fn: nil,
//		GetStringFn:  nil,
//		GetBoolFn:    nil,
//		GetUintFn:    nil,
//		GetUint32Fn:  nil,
//		GetUint64Fn:  nil,
//	})
//	err = interceptor.PostHandle(ctx, nil, nil)
//	assert.NotNil(t, err)
//	t.Logf("error:%++v", err)
//
//	handlerContexts := &contexts.HandlerContexts{
//		Org:        "org",
//		Az:         "az",
//		Su:         "su",
//		CommonSu:   "commonSu",
//		NodeID:     "nodeID",
//		ServiceID:  "serviceID",
//		InstanceID: "instanceID",
//		Wks:        "wks",
//		Env:        "env",
//		Lang:       "lang",
//		SpanContexts: &contexts.SpanContexts{
//			TraceID:             "trace ID",
//			SpanID:              "spanID",
//			ParentSpanID:        "parentSpanID",
//			ReplyToAddress:      "replyToAddress",
//			TimeoutMilliseconds: 100,
//		},
//		ResponseTemplate:            "",
//		ResponseAutoParseKeyMapping: nil,
//		TransactionContexts:         nil,
//	}
//	ctx = contexts.BuildContextFromParentWithHandlerContexts(ctx, handlerContexts)
//	err = interceptor.PostHandle(ctx, nil, nil)
//	assert.Nil(t, err)
//
//	csStartTimestamp := time.Now()
//	nCtx := context.WithValue(ctx, commConstant.CsStartTimestamp, csStartTimestamp)
//	err = interceptor.PostHandle(nCtx, nil, nil)
//	assert.Nil(t, err)
//}
