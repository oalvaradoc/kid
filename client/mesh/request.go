package mesh

import (
	"fmt"
	"time"

	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/codec"
	"git.multiverse.io/eventkit/kit/codec/auto"
	jsoniter "github.com/json-iterator/go"
)

type meshRequest struct {
	request interface{}
	opts    *client.RequestOptions
}

var (
	// DefaultRetries defines the default retries of the request
	DefaultRetries = 0

	// DefaultTimeout defines the default timeout of the request
	DefaultTimeout = 30 * time.Second

	// DefaultRequestCodec defines the default request codec of the request
	DefaultRequestCodec = auto.BuildAutoCodecWithJSONCodec()

	// DefaultBackoff defines the default backoff strategy of the request
	DefaultBackoff = client.FixedTimeBackoff

	// DefaultRetry defines the default retry strategy of the request
	DefaultRetry = client.Never
)

// NewMeshRequest creates a client.Request with default client.RequestOptions,
// It's easy to set one or more optional parameters via variable args ...client.RequestOption
func NewMeshRequest(request interface{}, reqOptions ...client.RequestOption) client.Request {
	opts := &client.RequestOptions{
		MaxRetryTimes:  DefaultRetries,
		MaxWaitingTime: DefaultTimeout,
		Timeout:        DefaultTimeout,
		Codec:          DefaultRequestCodec,
		Backoff:        DefaultBackoff,
		Retry:          DefaultRetry,
		FallbackFunc: func(e error) error {
			return e
		},
	}

	for _, reqOpt := range reqOptions {
		reqOpt(opts)
	}

	meshRequest := &meshRequest{
		request: request,
		opts:    opts,
	}

	return meshRequest
}

// Body returns the body of the client.Request
func (m *meshRequest) Body() interface{} {
	return m.request
}

// Codec returns the codec of the client.Request
func (m *meshRequest) Codec() codec.Codec {
	return m.opts.Codec
}

// RequestOptions returns the client.RequestOptions of the client.Request
func (m *meshRequest) RequestOptions() *client.RequestOptions {
	return m.opts
}

// WithOptions sets one or more optional parameters into client.Request
func (m *meshRequest) WithOptions(reqOptions ...client.RequestOption) {
	for _, reqOpt := range reqOptions {
		reqOpt(m.opts)
	}
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (m meshRequest) String() string {
	requestStr, _ := json.Marshal(m.request)
	str := "{request:" + string(requestStr)
	str += fmt.Sprintf(" opts::%++v", m.opts)
	str += "}"

	return str
}
