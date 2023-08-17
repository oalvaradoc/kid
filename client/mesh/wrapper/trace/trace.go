package trace

import (
	"context"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
)

// Wrapper is an wrapper for trace
type Wrapper struct{}

// Before wraps the requests, will generate a new span ID each time.
func (t *Wrapper) Before(ctx context.Context, request interface{}, opts interface{}) (context.Context, error) {
	// generate a new span ID
	originalHandlerContexts := contexts.HandlerContextsFromContext(ctx)
	if nil == originalHandlerContexts {
		return ctx, nil
	}

	handlerContexts := originalHandlerContexts.Copy()
	var parentSpanID string
	if nil != handlerContexts.SpanContexts {
		parentSpanID = handlerContexts.SpanContexts.SpanID
	}
	spanID := util.GenerateSerialNo(
		handlerContexts.Org,
		handlerContexts.Wks,
		handlerContexts.Env,
		handlerContexts.Su,
		handlerContexts.InstanceID,
		constant.SpanIDType)

	if len(parentSpanID) == 0 {
		parentSpanID = spanID
	}

	handlerContexts.With(
		contexts.WithSpanID(spanID),
		contexts.WithParentSpanID(parentSpanID),
	)

	return contexts.BuildContextFromParentWithHandlerContexts(ctx, handlerContexts), nil
}

// After do nothing
func (t *Wrapper) After(ctx context.Context, request interface{}, responseMeta interface{}, opts interface{}) (context.Context, error) {
	return ctx, nil
}

func (t Wrapper) String() string {
	return constant.WrapperTrace
}
