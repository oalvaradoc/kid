package client

import (
	"context"
)

// Client is an interface defines all the function that client instance should implement.
type Client interface {
	SyncCall(ctx context.Context, request Request, response interface{}, opts ...CallOption) (ResponseMeta, error)
	AsyncCall(ctx context.Context, request Request, opts ...CallOption) error
	ReplySemiSyncCall(ctx context.Context, response Response) error
	Options() Options
}

// Option sets an optional parameter for clients.
type Option func(*Options)

// CallOption sets an optional parameter for clients.
type CallOption func(*CallOptions)

// RequestOption sets an optional parameter for request.
type RequestOption func(*RequestOptions)

//ResponseOption sets an optional parameter for response
type ResponseOption func(options *ResponseOptions)
