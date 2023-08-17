package trace

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"testing"
)

var wrapper = &Wrapper{}

func TestWrapper_Before(t *testing.T) {
	ctx := context.Background()
	retCtx, err := wrapper.Before(ctx, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, ctx, retCtx)

	handlerContexts := &contexts.HandlerContexts{
		Org:        "org",
		Az:         "az",
		Su:         "su",
		CommonSu:   "commonSu",
		NodeID:     "nodeID",
		ServiceID:  "serviceID",
		InstanceID: "instanceID",
		Wks:        "wks",
		Env:        "env",
		Lang:       "lang",
		SpanContexts: &contexts.SpanContexts{
			TraceID:             "trace ID",
			SpanID:              "spanID",
			ParentSpanID:        "parentSpanID",
			ReplyToAddress:      "replyToAddress",
			TimeoutMilliseconds: 100,
		},
		ResponseTemplate:            "",
		ResponseAutoParseKeyMapping: nil,
		TransactionContexts:         nil,
	}
	ctx = contexts.BuildContextFromParentWithHandlerContexts(ctx, handlerContexts)
	retCtx, err = wrapper.Before(ctx, nil, nil)
	assert.Nil(t, err)
	assert.NotEqual(t, ctx, retCtx)
	t.Logf("the return contexts is: %++v", retCtx)
}

func TestWrapper_After(t *testing.T) {
	originalCtx := context.Background()
	retCtx, err := wrapper.After(originalCtx, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, originalCtx, retCtx)
}

func TestWrapper_String(t *testing.T) {
	nameWrapper := fmt.Sprintf("%s", wrapper)
	assert.Equal(t, nameWrapper, constant.WrapperTrace)
}
