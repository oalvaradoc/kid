package client

import (
	"context"
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestAlways(t *testing.T) {
	retryFunc := Always

	ok, err := retryFunc(context.Background(), nil, 0, nil)
	assert.True(t, ok)
	assert.Nil(t, err)
}
