package serializer

import (
	"encoding/base64"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/protocol"
	"git.multiverse.io/eventkit/kit/constant"
	"testing"
)

func TestProtoMsg2BytesAndBytes2ProtoMsg(t *testing.T) {
	msg := &protocol.ProtoMessage{}
	msg.ID = 100
	msg.TopicAttribute = map[string]string{
		"k-A": "v-A",
		"k-B": "v-B",
		"k-C": "v-C",
	}
	msg.AppProps = map[string]string{
		"AP-k-A": "AP-v-A",
		"AP-k-B": "AP-v-B",
		"AP-k-C": "AP-v-C",
	}
	msg.NeedAck = true
	msg.NeedReply = true
	msg.Body = base64.StdEncoding.EncodeToString([]byte("test body"))
	bytes := ProtoMsg2Bytes(msg, nil)
	assert.NotNil(t, bytes)

	resMsg, err := Bytes2protoMsg(bytes)
	t.Logf("The result of bytes to protocol message:%++v", resMsg)
	assert.Nil(t, err)
	assert.Equal(t, resMsg.Body, base64.StdEncoding.EncodeToString([]byte("test body")))
	assert.True(t, resMsg.NeedAck)
	assert.True(t, resMsg.NeedReply)
	assert.Equal(t, len(resMsg.AppProps), 3)
	assert.Equal(t, len(resMsg.TopicAttribute), 3)
	assert.Equal(t, resMsg.AppProps["AP-k-B"], "AP-v-B")
	assert.Equal(t, resMsg.TopicAttribute["k-C"], "v-C")
	assert.Equal(t, resMsg.TopicAttribute["the_key_does_not_exist"], "")

	bytes = ProtoMsg2Bytes(nil, errors.New(constant.SystemInternalError, "mocked error from unit test"))
	assert.NotNil(t, bytes)

	resMsg, err = Bytes2protoMsg(bytes)
	assert.NotNil(t, err)
	t.Logf("The result of error:%++v", err)
}

func TestProtoMsg2BytesAndBytes2ProtoMsgCompatible(t *testing.T) {
	msg := &protocol.ProtoMessage{}
	msg.ID = 100
	msg.TopicAttribute = map[string]string{
		"k-A": "v-A",
		"k-B": "v-B",
		"k-C": "v-C",
	}
	msg.AppProps = map[string]string{
		"AP-k-A": "AP-v-A",
		"AP-k-B": "AP-v-B",
		"AP-k-C": "AP-v-C",
	}
	msg.NeedAck = true
	msg.NeedReply = true
	msg.Body = base64.StdEncoding.EncodeToString([]byte("test body"))
	bytes := ProtoMsg2BytesCompatible(msg, nil, true)
	assert.NotNil(t, bytes)

	resMsg, err := Bytes2protoMsg(bytes)
	t.Logf("The result of bytes to protocol message:%++v", resMsg)
	assert.Nil(t, err)
	assert.Equal(t, resMsg.Body, base64.StdEncoding.EncodeToString([]byte("test body")))
	assert.True(t, resMsg.NeedAck)
	assert.True(t, resMsg.NeedReply)
	assert.Equal(t, len(resMsg.AppProps), 3)
	assert.Equal(t, len(resMsg.TopicAttribute), 3)
	assert.Equal(t, resMsg.AppProps["AP-k-B"], "AP-v-B")
	assert.Equal(t, resMsg.TopicAttribute["k-C"], "v-C")
	assert.Equal(t, resMsg.TopicAttribute["the_key_does_not_exist"], "")

	bytes = ProtoMsg2BytesCompatible(nil, errors.New(constant.SystemInternalError, "mocked error from unit test"), true)
	assert.NotNil(t, bytes)

	resMsg, err = Bytes2protoMsg(bytes)
	assert.NotNil(t, err)
	t.Logf("The result of error:%++v", err)
}

func TestProtoMsg2StringAndString2ProtoMsg(t *testing.T) {
	msg := &protocol.ProtoMessage{}
	msg.ID = 100
	msg.TopicAttribute = map[string]string{
		"T-k-A": "T-v-A",
		"T-k-B": "T-v-B",
		"T-k-C": "T-v-C",
	}
	msg.AppProps = map[string]string{
		"S-AP-k-A": "S-AP-v-A",
		"S-AP-k-B": "S-AP-v-B",
		"S-AP-k-C": "S-AP-v-C",
	}
	msg.NeedAck = true
	msg.NeedReply = true
	msg.Body = base64.StdEncoding.EncodeToString([]byte("test body from unit test"))
	str := ProtoMsg2String(msg, nil)
	assert.NotNil(t, str)

	resMsg, err := String2protoMsg(str)
	t.Logf("The result of bytes to protocol message:%++v", resMsg)
	assert.Nil(t, err)
	assert.Equal(t, resMsg.Body, base64.StdEncoding.EncodeToString([]byte("test body from unit test")))
	assert.True(t, resMsg.NeedAck)
	assert.True(t, resMsg.NeedReply)
	assert.Equal(t, len(resMsg.AppProps), 3)
	assert.Equal(t, len(resMsg.TopicAttribute), 3)
	assert.Equal(t, resMsg.AppProps["S-AP-k-B"], "S-AP-v-B")
	assert.Equal(t, resMsg.TopicAttribute["T-k-C"], "T-v-C")
	assert.Equal(t, resMsg.TopicAttribute["the_key_does_not_exist"], "")

	str = ProtoMsg2String(nil, errors.New(constant.SystemInternalError, "mocked error from unit test"))
	assert.NotNil(t, str)

	resMsg, err = String2protoMsg(str)
	assert.NotNil(t, err)
	t.Logf("The result of error:%++v", err)
}

func TestProtoMsg2StringAndString2ProtoMsgCompatible(t *testing.T) {
	msg := &protocol.ProtoMessage{}
	msg.ID = 100
	msg.TopicAttribute = map[string]string{
		"T-k-A": "T-v-A",
		"T-k-B": "T-v-B",
		"T-k-C": "T-v-C",
	}
	msg.AppProps = map[string]string{
		"S-AP-k-A": "S-AP-v-A",
		"S-AP-k-B": "S-AP-v-B",
		"S-AP-k-C": "S-AP-v-C",
	}
	msg.NeedAck = true
	msg.NeedReply = true
	msg.Body = base64.StdEncoding.EncodeToString([]byte("test body from unit test"))
	str := ProtoMsg2StringCompatible(msg, nil, true)
	assert.NotNil(t, str)

	resMsg, err := String2protoMsg(str)
	t.Logf("The result of bytes to protocol message:%++v", resMsg)
	assert.Nil(t, err)
	assert.Equal(t, resMsg.Body, base64.StdEncoding.EncodeToString([]byte("test body from unit test")))
	assert.True(t, resMsg.NeedAck)
	assert.True(t, resMsg.NeedReply)
	assert.Equal(t, len(resMsg.AppProps), 3)
	assert.Equal(t, len(resMsg.TopicAttribute), 3)
	assert.Equal(t, resMsg.AppProps["S-AP-k-B"], "S-AP-v-B")
	assert.Equal(t, resMsg.TopicAttribute["T-k-C"], "T-v-C")
	assert.Equal(t, resMsg.TopicAttribute["the_key_does_not_exist"], "")

	str = ProtoMsg2StringCompatible(nil, errors.New(constant.SystemInternalError, "mocked error from unit test"), true)
	assert.NotNil(t, str)

	resMsg, err = String2protoMsg(str)
	assert.NotNil(t, err)
	t.Logf("The result of error:%++v", err)
}
