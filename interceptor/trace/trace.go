package trace

import (
	"context"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"strconv"
)

// Deprecated: Interceptor trace interceptor will remove soon
type Interceptor struct{}

// PreHandle checks whether the upstream service request exists trace contexts
// and inject the trace contexts into handler contexts if the trace contexts exists in the request
func (l *Interceptor) PreHandle(ctx context.Context, request *msg.Message) error {
	return nil
	handlerContexts := contexts.HandlerContextsFromContext(ctx)

	if nil != handlerContexts {
		traceID := ""
		spanID := ""
		parentSpanID := ""
		replyToAddress := ""
		timeoutMilliseconds := 0
		if nil != request {
			traceID = request.GetAppPropertyEitherSilence(constant.KeyTraceID, constant.KeyTraceIDOld)
			spanID = request.GetAppPropertyEitherSilence(constant.KeySpanID, constant.KeySpanIDOld)
			parentSpanID = request.GetAppPropertyEitherSilence(constant.KeyParentSpanID, constant.KeyParentSpanIDOld)
			replyToAddress = request.GetAppPropertySilence(constant.RrReplyTo)
			timeoutMillisecondsStr := request.GetAppPropertySilence(constant.To3)
			if "" != timeoutMillisecondsStr {
				timeout, err := strconv.Atoi(timeoutMillisecondsStr)
				if err == nil {
					timeoutMilliseconds = timeout
				}
			}
		}

		if traceID == "" {
			// get su id,trace id from context
			traceID = util.GenerateSerialNo(
				handlerContexts.Org,
				handlerContexts.Wks,
				handlerContexts.Env,
				handlerContexts.Su,
				handlerContexts.InstanceID, constant.TraceIDType)
		}

		if spanID == "" {
			spanID = util.GenerateSerialNo(
				handlerContexts.Org,
				handlerContexts.Wks,
				handlerContexts.Env,
				handlerContexts.Su,
				handlerContexts.InstanceID, constant.SpanIDType)
		}

		if parentSpanID == "" {
			parentSpanID = spanID
		}

		if 0 == timeoutMilliseconds {
			timeoutMilliseconds = constant.DefaultTimeoutMilliseconds
		}

		handlerContexts.With(
			contexts.Span(
				contexts.BuildSpanContexts(
					contexts.TraceID(traceID),
					contexts.SpanID(spanID),
					contexts.ParentSpanID(parentSpanID),
					contexts.TimeoutMilliseconds(timeoutMilliseconds),
					contexts.ReplyToAddress(replyToAddress),
				),
			),
		)
	}

	return nil
}

// PostHandle does nothing
func (l *Interceptor) PostHandle(ctx context.Context, request *msg.Message, response *msg.Message) error {
	// DO NOTHING
	return nil
}

func (l Interceptor) String() string {
	return constant.InterceptorTrace
}
