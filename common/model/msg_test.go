package model

import (
	"encoding/base64"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/common/protocol"
	"testing"
)

func TestMsgToProtocolMsg(t *testing.T) {
	msg := msg.Message{
		ID: 123,
		TopicAttribute: map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		},
		SessionName: "default",
		Body:        []byte("this is a test body"),
	}
	msg.SetAppProps(map[string]string{
		"A-K-1": "V-K-1",
		"A-K-2": "V-K-2",
		"A-K-3": "V-K-3",
	})
	protocolMsg := MsgToProtocolMsg(&msg)
	assert.NotNil(t, protocolMsg)
	assert.Equal(t, protocolMsg.ID, msg.ID)
	assert.Equal(t, protocolMsg.SessionName, msg.SessionName)
	assert.Equal(t, len(protocolMsg.TopicAttribute), len(msg.TopicAttribute))
	assert.Equal(t, len(protocolMsg.AppProps), len(msg.GetAppProps()))
	assert.Equal(t, protocolMsg.Body, base64.StdEncoding.EncodeToString([]byte("this is a test body")))
}

func TestProtocolMsgToMsg(t *testing.T) {
	protocolMsg := protocol.ProtoMessage{
		ID: 123,
		TopicAttribute: map[string]string{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		},
		NeedReply: true,
		NeedAck:   true,
		Body:      base64.StdEncoding.EncodeToString([]byte("this is a test body")),
	}

	msg := ProtocolMsgToMsg(&protocolMsg)
	assert.NotNil(t, msg)
	assert.Equal(t, protocolMsg.ID, msg.ID)
	assert.Equal(t, len(protocolMsg.TopicAttribute), len(msg.TopicAttribute))
	assert.Equal(t, len(protocolMsg.AppProps), len(msg.GetAppProps()))
	assert.Equal(t, base64.StdEncoding.EncodeToString(msg.Body), base64.StdEncoding.EncodeToString([]byte("this is a test body")))
}
