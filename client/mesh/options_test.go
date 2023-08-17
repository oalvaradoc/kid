package mesh

import (
	"context"
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/codec"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/config"
	"testing"
	"time"
)

func TestWithServiceConfig(t *testing.T) {
	o := &client.RequestOptions{}
	serviceConfig := &config.Service{}
	opt := WithServiceConfig(serviceConfig)
	opt(o)
	assert.Equal(t, o.ServiceConfig, serviceConfig)
}

func TestWithRequestOptions(t *testing.T) {
	o := &client.RequestOptions{}
	opt := WithTopicType(constant.TopicTypeBusiness)
	opt(o)
	assert.Equal(t, o.TopicType, constant.TopicTypeBusiness)

	opt = WithTopicTypeBusiness()
	opt(o)
	assert.Equal(t, o.TopicType, constant.TopicTypeBusiness)

	opt = WithTopicTypeOps()
	opt(o)
	assert.Equal(t, o.TopicType, constant.TopicTypeOPS)

	opt = WithTopicTypeMetrics()
	opt(o)
	assert.Equal(t, o.TopicType, constant.TopicTypeMetrics)

	opt = WithTopicTypeAlert()
	opt(o)
	assert.Equal(t, o.TopicType, constant.TopicTypeAlert)

	opt = WithTopicTypeLog()
	opt(o)
	assert.Equal(t, o.TopicType, constant.TopicTypeLog)

	opt = WithTopicTypeError()
	opt(o)
	assert.Equal(t, o.TopicType, constant.TopicTypeError)

	opt = WithTopicTypeDxc()
	opt(o)
	assert.Equal(t, o.TopicType, constant.TopicTypeDXC)

	opt = WithTopicTypeDts()
	opt(o)
	assert.Equal(t, o.TopicType, constant.TopicTypeDTS)

	opt = WithTopicTypeHeartbeat()
	opt(o)
	assert.Equal(t, o.TopicType, constant.TopicTypeHeartbeat)

	opt = WithSemiSyncCall()
	opt(o)
	assert.True(t, o.IsSemiSyncCall)

	opt = WithORG("org")
	opt(o)
	assert.Equal(t, o.Org, "org")

	opt = WithOrgIfEmpty("new org")
	opt(o)
	assert.Equal(t, o.Org, "org")

	opt = WithWorkspace("wks")
	opt(o)
	assert.Equal(t, o.Wks, "wks")

	opt = WithWorkspaceIfEmpty("new wks")
	opt(o)
	assert.Equal(t, o.Wks, "wks")

	opt = WithEnvironment("env")
	opt(o)
	assert.Equal(t, o.Env, "env")

	opt = WithEnvironmentIfEmpty("new env")
	opt(o)
	assert.Equal(t, o.Env, "env")

	opt = WithSU("su")
	opt(o)
	assert.Equal(t, o.Env, "env")

	opt = WithNodeID("nodeID")
	opt(o)
	assert.Equal(t, o.NodeID, "nodeID")

	opt = WithInstanceID("instanceID")
	opt(o)
	assert.Equal(t, o.InstanceID, "instanceID")

	opt = WithEventID("eventID")
	opt(o)
	assert.Equal(t, o.EventID, "eventID")

	ctx := context.Background()
	opt = WithContext(ctx)
	opt(o)
	assert.Equal(t, o.Context, ctx)

	header := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}

	opt = WithHeader(header)
	opt(o)
	assert.Equal(t, o.Header, header)

	opt = AddKeyToHeader("k4", "v4")
	opt(o)

	t.Logf("header:%++v", header)
	assert.Equal(t, o.Header, header)

	opt = WithElementType("element type")
	opt(o)
	assert.NotNil(t, o.Header)
	assert.Equal(t, o.Header[constant.GlsElementType], "element type")

	opt = WithElementClass("element class")
	opt(o)
	assert.NotNil(t, o.Header)
	assert.Equal(t, o.Header[constant.GlsElementClass], "element class")

	opt = WithElementID("element ID")
	opt(o)
	assert.NotNil(t, o.Header)
	assert.Equal(t, o.Header[constant.GlsElementID], "element ID")

	jsonCodec := codec.BuildJSONCodec()
	opt = WithCodec(jsonCodec)
	opt(o)
	assert.Equal(t, o.Codec, jsonCodec)

	xmlCodec := codec.BuildXMLCodec()
	opt = WithCodecIfEmpty(xmlCodec)
	opt(o)
	assert.Equal(t, o.Codec, jsonCodec)

	timeout := 100 * time.Second
	opt = WithTimeout(timeout)
	opt(o)
	assert.Equal(t, o.Timeout, timeout)

	retryWaitingTime := 200 * time.Second
	opt = WithRetryWaitingMilliseconds(retryWaitingTime)
	opt(o)
	assert.Equal(t, o.RetryWaitingTime, retryWaitingTime)

	maxRetryTimes := 3
	opt = WithMaxRetryTimes(maxRetryTimes)
	opt(o)
	assert.Equal(t, o.MaxRetryTimes, maxRetryTimes)

	opt = WithDeleteTransactionPropagationInformation(true)
	opt(o)
	assert.True(t, o.DeleteTransactionPropagationInformation)

	opt = SkipResponseAutoParse()
	opt(o)
	assert.True(t, o.SkipResponseAutoParse)

	opt = DisableMacroModel()
	opt(o)
	assert.True(t, o.DisableMacroModel)

	maxWaitingTime := 400 * time.Second
	opt = WithMaxWaitingTime(maxWaitingTime)
	opt(o)
	assert.Equal(t, o.MaxWaitingTime, maxWaitingTime)

	opt = WithTopicType("topic type")
	opt(o)
	assert.Equal(t, o.TopicType, "topic type")
	opt = WithTopicTypeIfEmpty("new topic type")
	opt(o)
	assert.Equal(t, o.TopicType, "topic type")

	opt = MarkLocalCall()
	opt(o)
	assert.True(t, o.IsLocalCall)

	opt = MarkDMQEligible()
	opt(o)
	assert.True(t, o.IsDMQEligible)

	opt = MarkPersistentDeliveryMode()
	opt(o)
	assert.True(t, o.IsPersistentDeliveryMode)

	opt = WithVersion("version")
	opt(o)
	assert.Equal(t, o.Version, "version")
	opt = WithHTTPRequestInfo("http address", "", "")
	opt(o)
	assert.True(t, o.HTTPCall)
	assert.Equal(t, o.Address, "http address")
	assert.Equal(t, o.HTTPMethod, constant.DefaultHTTPMethodPost)
	assert.Equal(t, o.ContentType, constant.DefaultContentTypeJSON)

	opt = WithHTTPRequestInfo("address", "method", "contextType")
	opt(o)
	assert.True(t, o.HTTPCall)
	assert.Equal(t, o.Address, "address")
	assert.Equal(t, o.HTTPMethod, "method")
	assert.Equal(t, o.ContentType, "contextType")
}

func TestWithResponseOptions(t *testing.T) {
	o := &client.ResponseOptions{}
	opt := WithReplyToAddress("reply to address")
	opt(o)
	assert.Equal(t, o.ReplyToAddress, "reply to address")

	textCodec := codec.BuildTextCodec()
	opt = WithResponseCodec(textCodec)
	opt(o)
	assert.Equal(t, o.Codec, textCodec)

	responseHeader := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}

	opt = WithResponseHeader(responseHeader)
	opt(o)
	assert.Equal(t, o.Header, responseHeader)

	opt = AddKeyToResponseHeader("k4", "v4")
	opt(o)
	assert.Equal(t, o.Header, responseHeader)
	assert.NotNil(t, o.Header)
	assert.Equal(t, o.Header["k4"], "v4")
}
