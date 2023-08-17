package contexts

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestBuildSpanContexts(t *testing.T) {
	spanContexts := BuildSpanContexts()
	assert.NotNil(t, spanContexts)
}

func TestSetSpanOptions(t *testing.T) {
	spanContexts := BuildSpanContexts()
	opts := []SpanContext{}

	opts = append(opts, TraceID("traceID1"))
	opts = append(opts, SpanID("spanID1"))
	opts = append(opts, ParentSpanID("parentSpanID1"))
	opts = append(opts, ReplyToAddress("replyToAddress1"))
	opts = append(opts, TimeoutMilliseconds(30))

	for _, opt := range opts {
		opt(spanContexts)
	}
	assert.Equal(t, spanContexts.TraceID, "traceID1")
	assert.Equal(t, spanContexts.SpanID, "spanID1")
	assert.Equal(t, spanContexts.ParentSpanID, "parentSpanID1")
	assert.Equal(t, spanContexts.ReplyToAddress, "replyToAddress1")
	assert.Equal(t, spanContexts.TimeoutMilliseconds, 30)
	assert.Equal(t, spanContexts.JudgeParentSpanID(), "spanID1")

	spanContextsCopied := spanContexts.Copy()
	assert.True(t, spanContextsCopied != nil)
	assert.Equal(t, spanContexts.TraceID, spanContextsCopied.TraceID)
	assert.Equal(t, spanContexts.SpanID, spanContextsCopied.SpanID)
	assert.Equal(t, spanContexts.ParentSpanID, spanContextsCopied.ParentSpanID)
	assert.Equal(t, spanContexts.ReplyToAddress, spanContextsCopied.ReplyToAddress)
	assert.Equal(t, spanContexts.TimeoutMilliseconds, spanContextsCopied.TimeoutMilliseconds)

	spanContexts.SpanID = ""
	assert.Equal(t, spanContexts.JudgeParentSpanID(), spanContextsCopied.TraceID)
}
