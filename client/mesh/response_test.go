package mesh

import (
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

var meshResponseMetaBody = []byte("this is a test response")
var meshResponseMetaHeader = map[string]string{
	"k1": "v1",
	"k2": "v2",
}

func TestNewMeshResponseMeta(t *testing.T) {
	responseMeta := NewMeshResponseMeta(meshResponseMetaBody, meshResponseMetaHeader)
	assert.Equal(t, responseMeta.Header(), meshResponseMetaHeader)
	assert.Equal(t, responseMeta.Body(), meshResponseMetaBody)
}

func TestNewMeshResponse(t *testing.T) {
	meshResponse := NewMeshResponse(nil)
	assert.Equal(t, meshResponse.Codec(), DefaultResponseCodec)

	o := client.ResponseOptions{}
	opt := AddKeyToResponseHeader("k1", "v1")

	opt(&o)
	assert.NotNil(t, o.Header)
	assert.Equal(t, len(o.Header), 1)
	assert.Equal(t, o.Header["k1"], "v1")
}
