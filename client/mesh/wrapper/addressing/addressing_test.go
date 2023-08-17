package addressing

import (
	"context"
	"fmt"
	v1 "git.multiverse.io/eventkit/kit/cache/v1"
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/config"
	mockcache "git.multiverse.io/eventkit/kit/mocks/cache"
	"github.com/golang/mock/gomock"
	"testing"
)

var wrapper = &Wrapper{}

func TestWrapper_String(t *testing.T) {
	nameOfWrapper := fmt.Sprintf("%s", wrapper)
	assert.Equal(t, nameOfWrapper, constant.WrapperAddressing)
}

func TestWrapper_After(t *testing.T) {
	originalCtx := context.Background()
	retCtx, err := wrapper.After(originalCtx, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, originalCtx, retCtx)
}

func TestWrapper_Before(t *testing.T) {
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
	ctx := context.Background()
	ctx = contexts.BuildContextFromParentWithHandlerContexts(ctx, handlerContexts)
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

	retCtx, err = wrapper.Before(ctx, nil, requestOptions)
	assert.Nil(t, err)
	assert.Equal(t, ctx, retCtx)

	t.Logf("mark the request as local call.")
	requestOptions.IsLocalCall = true
	retCtx, err = wrapper.Before(ctx, nil, requestOptions)
	assert.NotNil(t, err)
	assert.Equal(t, ctx, retCtx)
	t.Logf("request options:%++v", requestOptions)
	assert.Equal(t, requestOptions.Su, handlerContexts.CommonSu)

	requestOptions.Su = ""
	requestOptions.EventID = "event ID"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	newMockOperatorCallInc := mockcache.NewMockOperator(mockCtrl)
	newMockOperatorCallInc.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	newMockOperatorCallInc.EXPECT().Get(gomock.Any(), gomock.Any()).Return("V1", nil).AnyTimes()
	newMockOperatorCallInc.EXPECT().HGet(gomock.Any(), gomock.Any(), gomock.Any()).Return("V2", nil).AnyTimes()
	v1.AddressingCacheOperator = newMockOperatorCallInc

	requestOptions.Header = map[string]string{
		constant.GlsElementType:  "123",
		constant.GlsElementClass: "",
		constant.GlsElementID:    "element ID",
	}

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
	assert.Equal(t, ctx, retCtx)
	assert.Equal(t, requestOptions.Su, "V2")
}
