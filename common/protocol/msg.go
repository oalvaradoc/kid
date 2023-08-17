package protocol

// ProtoMessage is used to store the data structure of data conversion between go/c inside sed server.
type ProtoMessage struct {
	// Each interaction has a session unique identifier
	ID uint64

	// Used to save topic-related key-values pairs
	TopicAttribute map[string]string

	// True means this is a synchronous call message
	// False means this is an asynchronous call message
	NeedReply bool

	// If it's a synchronous call message, this field value is true
	// If it's a asynchronous call message, this field value is false
	NeedAck bool

	// SessionName is used to distinguish which session the message was sent to solace
	SessionName string

	// 0: means the delivery mode is direct
	// 1: means the delivery mode is persistent
	DeliveryMode int

	// app properties
	// service can set this attributes pass to target service
	AppProps map[string]string

	// message payload
	Body string
}
