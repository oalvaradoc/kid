package mesh

import (
	"context"
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestNewMeshRequest(t *testing.T) {
	request := NewMeshRequest(nil)
	assert.NotNil(t, request.RequestOptions())
	assert.Equal(t, request.RequestOptions().MaxRetryTimes, DefaultRetries)
	assert.Equal(t, request.RequestOptions().Timeout, DefaultTimeout)
	assert.Equal(t, request.RequestOptions().Codec, DefaultRequestCodec)

	dr, err := request.RequestOptions().Backoff(context.Background(), nil, 0)
	assert.Nil(t, err)
	assert.True(t, 0 == dr.Seconds())

	r, err := request.RequestOptions().Retry(context.Background(), nil, 0, nil)
	assert.Nil(t, err)
	assert.True(t, r)
}
