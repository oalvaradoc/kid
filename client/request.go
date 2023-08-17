package client

import "git.multiverse.io/eventkit/kit/codec"

// Request is an interface that contains the necessary request information for the Request instance for an RPC request
type Request interface {
	Body() interface{}
	Codec() codec.Codec
	RequestOptions() *RequestOptions
	WithOptions(reqOptions ...RequestOption)
}
