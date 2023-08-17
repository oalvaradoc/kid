package client

import (
	"context"
	"git.multiverse.io/eventkit/kit/codec"
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
	"time"
)

const FixedWaitingTime = 10 * time.Second

type MyRequest struct{}

func (m MyRequest) Body() interface{}  { return nil }
func (m MyRequest) Codec() codec.Codec { return nil }
func (m MyRequest) RequestOptions() *RequestOptions {
	return &RequestOptions{
		RetryWaitingTime: FixedWaitingTime,
	}
}
func (m MyRequest) WithOptions(reqOptions ...RequestOption) {}

func TestFixedTimeBackoff(t *testing.T) {
	fixedTimeBackoffFunc := FixedTimeBackoff
	duration, err := fixedTimeBackoffFunc(context.Background(), nil, 0)
	assert.Nil(t, err)
	assert.True(t, 0 == duration.Microseconds())

	duration, err = fixedTimeBackoffFunc(context.Background(), nil, 1)
	assert.Nil(t, err)
	assert.True(t, 0 == duration.Microseconds())

	duration, err = fixedTimeBackoffFunc(context.Background(), &MyRequest{}, 1)
	assert.Nil(t, err)
	assert.Equal(t, duration, FixedWaitingTime)
}
