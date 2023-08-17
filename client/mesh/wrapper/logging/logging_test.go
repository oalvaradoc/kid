package logging

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/client/mesh"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/constant"
	"testing"
)

var wrapper = &Wrapper{}

func TestWrapper_After(t *testing.T) {
	ctx := context.Background()
	retCtx, err := wrapper.Before(ctx, mesh.NewMeshRequest(nil), nil)
	assert.Nil(t, err)
	assert.Equal(t, ctx, retCtx)
}

func TestWrapper_Before(t *testing.T) {
	ctx := context.Background()
	retCtx, err := wrapper.After(ctx, nil, mesh.NewMeshResponse(nil), nil)
	assert.Nil(t, err)
	assert.Equal(t, ctx, retCtx)
}

func TestWrapper_String(t *testing.T) {
	nameOfWrapper := fmt.Sprintf("%s", wrapper)
	assert.Equal(t, nameOfWrapper, constant.WrapperLogging)
}
