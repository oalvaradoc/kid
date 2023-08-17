package mesh

import (
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/codec"
	"git.multiverse.io/eventkit/kit/codec/auto"
	"git.multiverse.io/eventkit/kit/common/util"
)

type meshResponseMeta struct {
	body   []byte
	header map[string]string
}

var (
	// DefaultResponseCodec defines the default codec of each response
	DefaultResponseCodec = auto.BuildAutoCodecWithJSONCodec()
)

// NewMeshResponseMeta constructs a client.ResponseMeta for a response meta of mesh request, contains response data and response header
func NewMeshResponseMeta(response []byte, header map[string]string) client.ResponseMeta {
	meshResponseMeta := &meshResponseMeta{
		body:   response,
		header: header,
	}

	return meshResponseMeta
}

// Body returns the body data of a response meta
func (m *meshResponseMeta) Body() []byte {
	return m.body
}

// Header returns the header of a response ,eta
func (m *meshResponseMeta) Header() map[string]string {
	return m.header
}

func (m meshResponseMeta) String() string {
	str := "{header:" + util.MapToString(m.header)
	str += " body:[" + string(m.body)
	str += "]}"

	return str
}

type meshResponse struct {
	response interface{}
	opts     *client.ResponseOptions
}

// Body returns the body data of a response
func (m *meshResponse) Body() interface{} {
	return m.response
}

// Codec returns the codec of a response
func (m *meshResponse) Codec() codec.Codec {
	return m.opts.Codec
}

// ResponseOptions returns the client.ResponseOptions of a response
func (m *meshResponse) ResponseOptions() *client.ResponseOptions {
	return m.opts
}

// WithOptions sets an optional parameter into response
func (m *meshResponse) WithOptions(resOptions ...client.ResponseOption) {
	for _, resOpt := range resOptions {
		resOpt(m.opts)
	}
}

// NewMeshResponse creates a new client.Response with response and one or more client.ResponseOptions
func NewMeshResponse(response interface{}, resOptions ...client.ResponseOption) client.Response {
	opts := &client.ResponseOptions{
		Codec: DefaultResponseCodec,
	}

	for _, reqOpt := range resOptions {
		reqOpt(opts)
	}

	meshResponse := &meshResponse{
		response: response,
		opts:     opts,
	}

	return meshResponse
}
