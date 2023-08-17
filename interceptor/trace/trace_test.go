package trace

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/constant"
	"testing"
)

var interceptor = &Interceptor{}

func TestInterceptor_String(t *testing.T) {
	nameOfInterceptor := fmt.Sprintf("%s", interceptor)
	assert.Equal(t, nameOfInterceptor, constant.InterceptorTrace)
}

func TestInterceptor_PreHandle(t *testing.T) {
	ctx := context.Background()
	err := interceptor.PreHandle(ctx, nil)
	assert.Nil(t, err)
}

func TestInterceptor_PostHandle(t *testing.T) {
	ctx := context.Background()
	err := interceptor.PostHandle(ctx, nil, nil)
	assert.Nil(t, err)
}
