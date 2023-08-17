package apm

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/config"
	"strconv"
	"testing"
	"time"
)

var interceptor = &Interceptor{}

func TestInterceptor_String(t *testing.T) {
	nameOfInterceptor := fmt.Sprintf("%s", interceptor)
	assert.Equal(t, nameOfInterceptor, constant.InterceptorApm)
}

func TestInterceptor_PreHandle(t *testing.T) {
	ctx := context.Background()
	err := interceptor.PreHandle(ctx, nil)
	assert.Nil(t, err)

	request := &msg.Message{}
	request.SetAppProps(map[string]string{
		constant.TxnIsLocalCall: "1",
	})
	err = interceptor.PreHandle(ctx, request)
	assert.Nil(t, err)

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
			Version:                             "",
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
	err = interceptor.PreHandle(ctx, request)
	assert.Nil(t, err)

	csStartTimestamp := request.GetAppPropertyIgnoreCaseSilence(constant.CsStartTimestamp)
	assert.NotNil(t, csStartTimestamp)
	t.Logf("client start timestamp:%++v", csStartTimestamp)
}

func TestInterceptor_PostHandle(t *testing.T) {
	ctx := context.Background()
	err := interceptor.PostHandle(ctx, nil, nil)
	assert.Nil(t, err)

	request := &msg.Message{}
	request.SetAppProps(map[string]string{
		constant.TxnIsLocalCall: "0",
	})
	err = interceptor.PostHandle(ctx, request, nil)
	assert.Nil(t, err)

	request.SetAppProperty(constant.TxnIsLocalCall, "1")
	err = interceptor.PostHandle(ctx, request, nil)
	assert.Nil(t, err)

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
			Version:                             "",
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
	err = interceptor.PostHandle(ctx, request, nil)
	assert.Nil(t, err)
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
	err = interceptor.PostHandle(ctx, request, nil)
	assert.Nil(t, err)

	csStartTimestampString := strconv.FormatInt(time.Now().UnixNano()/1000, 10)
	request.SetAppProperty(constant.CsStartTimestamp, csStartTimestampString)
	err = interceptor.PostHandle(ctx, request, nil)
	assert.Nil(t, err)
}
