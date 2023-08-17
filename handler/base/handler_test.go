package base

import (
	"context"
	"git.multiverse.io/eventkit/kit/client/mesh"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/mocks/remote"
	"github.com/golang/mock/gomock"
	"testing"
)

type Request struct {
	Field1 string `validate:"required"`
	Field2 string `validate:"required"`
}

func TestHandler_PreHandle(t *testing.T) {
	handler := Handler{}
	handler.SetLang(constant.LangEnUS)
	err := handler.PreHandle("test")
	assert.Nil(t, err)

	r := Request{
		Field1: "",
		Field2: "",
	}
	err = handler.Validation(r)
	assert.NotNil(t, err)
	t.Logf("the result of error:%++v", err)

	handler.SetCombineErrors(true)
	err = handler.Validation(r)
	assert.NotNil(t, err)
	t.Logf("the result of error:%++v", err)

	r = Request{
		Field1: "test field1",
		Field2: "test field2",
	}
	err = handler.Validation(r)
	assert.Nil(t, err)

	serviceConfigs := &config.Service{
		ServiceID:                   "ServiceID",
		Org:                         "Org",
		Az:                          "Az",
		Wks:                         "Wks",
		Env:                         "Env",
		NodeID:                      "NodeID",
		InstanceID:                  "InstanceID",
		Su:                          "Su",
		GroupSu:                     "GroupSu",
		CommonSu:                    "CommonSu",
		ResponseTemplate:            `{"returnCode":"{{.errorCode}}", "returnMsg":"{{.errorMsg}}", "data":{{.data}}, "debugStack":"{{.errorStack}}"}`,
		ResponseAutoParseKeyMapping: nil,
	}
	extConfigs := map[string]interface{}{
		constant.ExtConfigService: serviceConfigs,
	}

	handler.SetExtConfigs(extConfigs)

	assert.Equal(t, serviceConfigs, handler.ServiceConfig)

	// testing topicattributes
	topicAttributes := map[string]string{
		constant.TopicType:              "TopicType",
		constant.TopicID:                "TopicID",
		constant.TopicSourceORG:         "TopicSourceORG",
		constant.TopicSourceWorkspace:   "TopicSourceWorkspace",
		constant.TopicSourceEnvironment: "TopicSourceEnvironment",
		constant.TopicSourceServiceID:   "TopicSourceServiceID",
		constant.TopicSourceSU:          "TopicSourceSU",
		constant.TopicSourceNodeID:      "TopicSourceNodeID",
		constant.TopicSourceAZ:          "TopicSourceAZ",
		constant.TopicSourceInstanceID:  "TopicSourceInstanceID",
		constant.TopicDestinationSU:     "TopicDestinationSU",
		constant.TopicDestinationDCN:    "TopicDestinationDCN",
	}
	handler.SetTopicAttributes(topicAttributes)
	assert.Equal(t, handler.GetTopicAttributes(), topicAttributes)
	assert.Equal(t, handler.GetMsgTopicType(), "TopicType")
	assert.Equal(t, handler.GetMsgTopicID(), "TopicID")
	assert.Equal(t, handler.GetSourceORG(), "TopicSourceORG")
	assert.Equal(t, handler.GetSourceWKS(), "TopicSourceWorkspace")
	assert.Equal(t, handler.GetSourceENV(), "TopicSourceEnvironment")
	assert.Equal(t, handler.GetSourceServiceID(), "TopicSourceServiceID")
	assert.Equal(t, handler.GetSourceSU(), "TopicSourceSU")
	assert.Equal(t, handler.GetSourceNodeID(), "TopicSourceNodeID")
	assert.Equal(t, handler.GetSourceAZ(), "TopicSourceAZ")
	assert.Equal(t, handler.GetSourceInstanceID(), "TopicSourceInstanceID")
	assert.Equal(t, handler.GetCurrentSU(), "TopicDestinationSU")
	assert.Equal(t, topicAttributes, handler.GetTopicAttributes())

	ctx := context.Background()
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
	handler.SetContext(ctx)
	assert.Equal(t, handlerContexts, handler.GetHandlerContexts())

	// testing request header
	requestAppProps := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	}
	handler.SetRequestHeader(requestAppProps)
	assert.Equal(t, "v1", handler.GetRequestHeaderValWithKeyEitherSilence("k1", "key_does_not_exist"))
	assert.Equal(t, "v2", handler.GetRequestHeaderValWithKeyEitherSilence("key_does_not_exist", "k2"))
	assert.Equal(t, requestAppProps, handler.GetRequestHeader())
	v, ok := handler.GetRequestHeaderValWithKey("k3")
	assert.True(t, ok)
	assert.Equal(t, v, "v3")

	v, ok = handler.GetRequestHeaderValWithKey("key_does_not_exist")
	assert.False(t, ok)
	assert.Equal(t, v, "")

	assert.Equal(t, "", handler.GetRequestHeaderValWithKeySilence("key_does_not_exist"))
	assert.Equal(t, "v2", handler.GetRequestHeaderValWithKeySilence("k2"))
	assert.Equal(t, "v1", handler.GetRequestHeaderValWithKeyIgnoreCaseSilence("K1"))
	assert.Equal(t, "", handler.GetRequestHeaderValWithKeyIgnoreCaseSilence("KEY_DOES_NOT_EXIST"))

	v, ok = handler.GetRequestHeaderValWithKeyIgnoreCase("K3")
	assert.True(t, ok)
	assert.Equal(t, v, "v3")

	clonedRequestHeader := handler.CloneRequestHeader()
	assert.Equal(t, handler.GetRequestHeader(), clonedRequestHeader)

	v, ok = handler.GetRequestHeaderValWithKeyIgnoreCase("KEY_DOES_NOT_EXIST")
	assert.False(t, ok)
	assert.Equal(t, v, "")

	handler.RangeRequestHeader(func(s string, s2 string) {
		t.Logf("request header, key:%s-value:%s", s, s2)
	})
	handler.DeleteRequestHeader("k3")
	assert.Equal(t, "", handler.GetRequestHeaderValWithKeySilence("k3"))
	t.Logf("request header:%s", handler.RequestHeaderToString())

	// testing body
	body := []byte("this is a test body")
	handler.SetBody(body)
	assert.Equal(t, body, handler.RequestBody())

	// testing response
	responseAppProps := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
		"k4": "v4",
	}
	handler.SetResponseHeader(responseAppProps)
	assert.Equal(t, responseAppProps, handler.GetResponseHeader())
	v, ok = handler.GetResponseHeaderValWithKey("k3")
	assert.True(t, ok)
	assert.Equal(t, "v3", v)

	v, ok = handler.GetResponseHeaderValWithKey("key_does_not_exist")
	assert.False(t, ok)
	assert.Equal(t, "", v)

	v, ok = handler.GetResponseHeaderValWithKeyIgnoreCase("K3")
	assert.True(t, ok)
	assert.Equal(t, "v3", v)

	v, ok = handler.GetResponseHeaderValWithKeyIgnoreCase("KEY_DOES_NOT_EXIST")
	assert.False(t, ok)
	assert.Equal(t, "", v)

	assert.Equal(t, "v2", handler.GetResponseHeaderValWithKeySilence("k2"))
	assert.Equal(t, "", handler.GetResponseHeaderValWithKeySilence("key_does_not_exist"))

	handler.AddResponseHeader("k5", "v5")
	assert.Equal(t, len(handler.GetResponseHeader()), 5)
	assert.Equal(t, "v5", handler.GetResponseHeaderValWithKeySilence("k5"))

	handler.RangeResponseHeader(func(s string, s2 string) {
		t.Logf("response header, key:%s-value:%s", s, s2)
	})

	assert.Equal(t, handler.GetResponseHeader(), handler.CloneResponseHeader())
	handler.DeleteResponseHeader("k4")
	assert.Equal(t, "", handler.GetResponseHeaderValWithKeySilence("k4"))
	t.Logf("response header:%s", handler.ResponseHeaderToString())

	// testing discard response
	assert.False(t, handler.IsDiscardResponse())
	handler.DiscardResponse()
	assert.True(t, handler.IsDiscardResponse())

	// testing remote call
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	newMockRemoteCallInc := remote.NewMockCallInc(mockCtrl)
	handler.SetRemoteCall(newMockRemoteCallInc)
	newMockRemoteCallInc.EXPECT().SyncCall(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(a, b, c, d, e interface{}) (h interface{}, err error) {
			response := mesh.NewMeshResponseMeta([]byte(`{"errorCode":"0","errorMsg":"","response":true}`), map[string]string{})
			return response, nil
		}).AnyTimes()

	newMockRemoteCallInc.EXPECT().SemiSyncCall(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(a, b, c, d, e interface{}) (h interface{}, err error) {
			response := mesh.NewMeshResponseMeta([]byte(`{"errorCode":"0","errorMsg":"","response":true}`), map[string]string{})
			return response, nil
		}).AnyTimes()

	newMockRemoteCallInc.EXPECT().AsyncCall(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	_, err = handler.SyncCall(context.Background(), "dstSU", "testServiceKey", nil, nil)
	assert.Nil(t, err)

	_, err = handler.SemiSyncCall(context.Background(), "dstSU", "testServiceKey", nil, nil)
	assert.Nil(t, err)

	err = handler.AsyncCall(context.Background(), "dstSU", "testServiceKey", nil)
	assert.Nil(t, err)

	newMockRemoteCallInc.EXPECT().SyncCallw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(a, b, c, d, e, f interface{}) (h interface{}, err error) {
			response := mesh.NewMeshResponseMeta([]byte(`{"errorCode":"0","errorMsg":"","response":true}`), map[string]string{})
			return response, nil
		}).AnyTimes()
	newMockRemoteCallInc.EXPECT().SemiSyncCallw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(a, b, c, d, e, f interface{}) (h interface{}, err error) {
			response := mesh.NewMeshResponseMeta([]byte(`{"errorCode":"0","errorMsg":"","response":true}`), map[string]string{})
			return response, nil
		}).AnyTimes()
	newMockRemoteCallInc.EXPECT().AsyncCallw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	_, err = handler.SyncCallw(context.Background(), "elementType", "elementID", "testServiceKey", nil, nil)
	assert.Nil(t, err)

	_, err = handler.SemiSyncCallw(context.Background(), "elementType", "elementID", "testServiceKey", nil, nil)
	assert.Nil(t, err)

	err = handler.AsyncCallw(context.Background(), "dstSU", "elementID", "testServiceKey", nil)
	assert.Nil(t, err)

	newMockRemoteCallInc.EXPECT().SyncCalls(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(a, b, c interface{}) (h interface{}, err error) {
			response := mesh.NewMeshResponseMeta([]byte(`{"errorCode":"0","errorMsg":"","response":true}`), map[string]string{})
			return response, nil
		}).AnyTimes()
	newMockRemoteCallInc.EXPECT().SemiSyncCalls(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(a, b, c interface{}) (h interface{}, err error) {
			response := mesh.NewMeshResponseMeta([]byte(`{"errorCode":"0","errorMsg":"","response":true}`), map[string]string{})
			return response, nil
		}).AnyTimes()
	newMockRemoteCallInc.EXPECT().AsyncCalls(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	_, err = handler.SyncCalls(context.Background(), nil, nil)
	assert.Nil(t, err)

	_, err = handler.SemiSyncCalls(context.Background(), nil, nil)
	assert.Nil(t, err)

	err = handler.AsyncCalls(context.Background(), nil)
	assert.Nil(t, err)

	newMockRemoteCallInc.EXPECT().ReplyTo(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	err = handler.ReplyTo(context.Background(), nil)
	assert.Nil(t, err)

}
