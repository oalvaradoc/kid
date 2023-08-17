package client

import (
	"context"
	"git.multiverse.io/eventkit/kit/codec"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/interceptor"
	"git.multiverse.io/eventkit/kit/wrapper"
	"git.multiverse.io/eventkit/kit/sed/callback"
	"sync"
	"time"
)

// TopicType alias int as topic type
type TopicType int

// defines all the topic types
const (
	TRN TopicType = iota
	OPS
	DXC
	HBT
	ALT
	ERR
	LOG
	MTR
	P2P
)

// Options is wrapper of CallOptions
type Options struct {
	CallOptions CallOptions
}

// CallOptions contains the interceptor and wrapper contained in each client execution
type CallOptions struct {
	CallWrappers     []wrapper.Wrapper
	CallInterceptors []interceptor.Interceptor
	CallbackExecutor callback.Executor
}

// RequestOptions contains the relevant parameters required for each request
type RequestOptions struct {
	SessionName                             string
	Codec                                   codec.Codec
	Context                                 context.Context
	TopicType                               string
	EventID                                 string
	Org                                     string
	Wks                                     string
	Env                                     string
	Su                                      string
	Version                                 string
	NodeID                                  string
	InstanceID                              string
	MaxRetryTimes                           int
	MaxWaitingTime                          time.Duration
	Timeout                                 time.Duration
	RetryWaitingTime                        time.Duration
	HeaderLock                              sync.RWMutex
	Header                                  map[string]string
	Backoff                                 BackoffFunc
	Retry                                   RetryFunc
	IsLocalCall                             bool
	IsDMQEligible                           bool
	IsPersistentDeliveryMode                bool
	DeleteTransactionPropagationInformation bool
	SkipResponseAutoParse                   bool
	DisableMacroModel                       bool
	IsSemiSyncCall                          bool
	ServiceConfig                           *config.Service
	ServiceKey                              string
	IsEnabledCircuitBreaker                 bool
	FallbackFunc                            func(error) error
	PassThroughHeaderKeyList                []string
	OriginalHeader   						map[string]string
	Masker                           		config.Masker
	EnableLogging           				bool

	HTTPCall    bool
	Address     string
	ContentType string
	HTTPMethod  string
}

// WithCallbackExecutor sets the callback executor, this optional parameter is used for execute the macro service
func WithCallbackExecutor(callbackExecutor callback.Executor) CallOption {
	return func(options *CallOptions) {
		options.CallbackExecutor = callbackExecutor
	}
}

// WithWrapper adds one or more wrapper.Wrappers into CallOptions
func WithWrapper(cas ...wrapper.Wrapper) CallOption {
	return func(options *CallOptions) {
		options.CallWrappers = append(options.CallWrappers, cas...)
	}
}

// DefaultWrapperCall sets the default one or more wrapper.Wrappers the request will invoked
func DefaultWrapperCall(cas ...wrapper.Wrapper) Option {
	return func(options *Options) {
		options.CallOptions.CallWrappers = append(options.CallOptions.CallWrappers, cas...)
	}
}

// DefaultCallInterceptors sets the default one or more interceptor.Interceptors the request will invoked
func DefaultCallInterceptors(callInterceptors []interceptor.Interceptor) Option {
	return func(options *Options) {
		options.CallOptions.CallInterceptors = callInterceptors
	}
}

// WithCallWrappers sets the array of wrapper.Wrapper for the current request
func WithCallWrappers(callWrappers []wrapper.Wrapper) CallOption {
	return func(options *CallOptions) {
		options.CallWrappers = callWrappers
	}
}

// WithCallInterceptors sets the array of interceptor.Interceptor for the current request
func WithCallInterceptors(callInterceptors []interceptor.Interceptor) CallOption {
	return func(options *CallOptions) {
		options.CallInterceptors = callInterceptors
	}
}

// WithOptionFromCallOption sets an optional parameter of Option from CallOption
func WithOptionFromCallOption(callOption CallOption) Option {
	return func(options *Options) {
		callOption(&options.CallOptions)
	}
}

// NewOptions creates an default Options.
func NewOptions(options ...Option) Options {
	opts := Options{
		CallOptions: CallOptions{},
	}

	for _, o := range options {
		o(&opts)
	}

	return opts
}

// ResponseOptions contains the relevant parameters required for each response
type ResponseOptions struct {
	Codec          codec.Codec
	Context        context.Context
	SessionName    string
	Header         map[string]string
	ReplyToAddress string
}
