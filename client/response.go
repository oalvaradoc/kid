package client

import "git.multiverse.io/eventkit/kit/codec"

// ResponseMeta is the body and header contained in the response instance returned after an RPC request
type ResponseMeta interface {
	Body() []byte
	Header() map[string]string
}

// Response is used to set the relevant content required by the response instance when returning semi-synchronously call
type Response interface {
	Body() interface{}
	Codec() codec.Codec
	ResponseOptions() *ResponseOptions
	WithOptions(reqOptions ...ResponseOption)
}
