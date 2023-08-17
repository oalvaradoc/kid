package model

import (
	"encoding/base64"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/common/protocol"
	"git.multiverse.io/eventkit/kit/constant"
)

// QueueAckReq structure for ack queue request message
type QueueAckReq struct {
	MsgID uint64
}

// PublishReq structure for publish messagess
type PublishReq struct {
	Msg protocol.ProtoMessage
}

// ReplyMsgReq structure for sending reply messages
type ReplyMsgReq struct {
	Msg      protocol.ProtoMessage
	AckMsgID uint64
}

// SyncCallReq structure for R/R (request/reply) messages
type SyncCallReq struct {
	Msg protocol.ProtoMessage
}

// MsgToProtocolMsg is convert user message to proto message
func MsgToProtocolMsg(msg *msg.Message) protocol.ProtoMessage {
	protoMsg := protocol.ProtoMessage{
		AppProps:       msg.GetAppProps(),
		ID:             msg.ID,
		TopicAttribute: msg.TopicAttribute,
		SessionName:    msg.SessionName,
	}

	// set session name to "default" if not specified
	if msg.SessionName == "" {
		protoMsg.SessionName = "default"
	}

	// 0: means the delivery mode is direct
	// 1: means the delivery mode is persistent
	// set delivery mode to "direct" if not specified

	if "1" == msg.GetAppPropertySilence(constant.DeliveryMode) {
		protoMsg.DeliveryMode = 1
	} else {
		protoMsg.DeliveryMode = 0
	}

	// encode the message format to base64
	if len(msg.Body) > 0 {
		protoMsg.Body = base64.StdEncoding.EncodeToString(msg.Body)
	}

	return protoMsg
}

// ProtocolMsgToMsg is convert proto message to user message
func ProtocolMsgToMsg(protoMsg *protocol.ProtoMessage) *msg.Message {
	msg := &msg.Message{
		TopicAttribute: protoMsg.TopicAttribute,
		ID:             protoMsg.ID,
		NeedReply:      protoMsg.NeedReply,
		NeedAck:        protoMsg.NeedAck,
	}
	msg.SetAppProps(protoMsg.AppProps)

	// decode base64 message
	if len(protoMsg.Body) > 0 {
		msg.Body, _ = base64.StdEncoding.DecodeString(protoMsg.Body)
	}

	return msg
}
