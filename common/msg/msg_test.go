package msg

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/constant"
	"testing"
)

func TestMessage_AppPropsToString(t *testing.T) {
	msg := &Message{}
	msg.SetAppProps(map[string]string{
		"a":   "b",
		"m-k": "value\nthe new line",
	})
	t.Logf("app properties to string:[%s]", msg.AppPropsToString())
}

func TestMessage_RangeAppProps(t *testing.T) {
	msg := &Message{}
	msg.appProps = map[string]string{"key1": "value1"}
	msg.SetAppProperty("key2", "value2")

	msg.RangeAppProps(func(key string, value string) {
		t.Logf("key=[%s], value=[%s]", key, value)
	})
}

func TestMessage_CloneAppProps(t *testing.T) {
	msg := &Message{}
	msg.appProps = map[string]string{"key1": "value1"}
	msg.SetAppProperty("key2", "value2")

	t.Logf("app properties:[%++v]", msg.CloneAppProps())
}

var msg = &Message{
	appProps: map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	},
	TopicAttribute: map[string]string{
		constant.TopicID:                "id1",
		constant.TopicType:              "TRN",
		constant.TopicSourceORG:         "ORG1",
		constant.TopicSourceAZ:          "AZ1",
		constant.TopicSourceDCN:         "SU1",
		constant.TopicSourceSU:          "SU1",
		constant.TopicSourceNodeID:      "NodeID1",
		constant.TopicSourceInstanceID:  "InstanceID1",
		constant.TopicSourceWorkspace:   "WKS1",
		constant.TopicSourceEnvironment: "ENV1",
		constant.TopicSourceServiceID:   "Service1",
	},
}

func TestMessage_GetAppProps(t *testing.T) {
	r, ok := msg.GetAppProperty("k1")
	assert.True(t, ok)
	assert.Equal(t, r, "v1")

	r, ok = msg.GetAppProperty("k0")
	assert.False(t, ok)
	assert.Equal(t, r, "")
}

func TestMessage_GetAppPropertySilence(t *testing.T) {
	r := msg.GetAppPropertySilence("k1")
	assert.Equal(t, r, "v1")

	r = msg.GetAppPropertySilence("k0")
	assert.Equal(t, r, "")
}

func TestMessage_GetMsgAttributes(t *testing.T) {
	assert.Equal(t, msg.GetMsgTopicId(), "id1")
	assert.Equal(t, msg.GetMsgTopicType(), "TRN")
	assert.Equal(t, msg.GetSourceORG(), "ORG1")
	assert.Equal(t, msg.GetSourceAZ(), "AZ1")
	assert.Equal(t, msg.GetSourceDCN(), "SU1")
	assert.Equal(t, msg.GetSourceSU(), "SU1")
	assert.Equal(t, msg.GetSourceWKS(), "WKS1")
	assert.Equal(t, msg.GetSourceENV(), "ENV1")
	assert.Equal(t, msg.GetSourceServiceID(), "Service1")
	assert.Equal(t, msg.GetSourceNodeID(), "NodeID1")
	assert.Equal(t, msg.GetSourceInstanceID(), "InstanceID1")
	assert.True(t, msg.IsValidTopicType())
}

func TestMessage_JudgeUserLang(t *testing.T) {
	assert.Equal(t, msg.JudgeUserLang(), constant.LangEnUS)
}

func TestMessage_SetAppProperty(t *testing.T) {
	msg.SetAppProperty(constant.UserLang, constant.LangZhCN)
	assert.Equal(t, msg.JudgeUserLang(), constant.LangZhCN)
	msg.DeleteProperty(constant.UserLang)
	assert.Equal(t, msg.JudgeUserLang(), constant.LangEnUS)
}

func TestMessage_TopicAttributesToString(t *testing.T) {
	t.Logf("The string value of topic attributes is:%s", msg.TopicAttributesToString())
}

func TestMessage_JsonPrint(t *testing.T) {
	msg.Body = []byte("{\"abc\":\"def\"}");
	msg.ID = 123123

	t.Logf("message to string:%s", msg)
}
