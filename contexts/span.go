package contexts

// SpanContexts is a runtime context, used for trace-related parameter transfer in various methods and modules,
// mainly including trace ID, span ID and parent span ID
type SpanContexts struct {
	TraceID             string
	SpanID              string
	ParentSpanID        string
	ReplyToAddress      string
	TimeoutMilliseconds int
	PassThroughHeaderKeyList string
}

// SpanContext sets a parameter into span contexts
type SpanContext func(*SpanContexts)

// BuildSpanContexts create a new span contexts with one or more setting functions.
func BuildSpanContexts(otherSpanContexts ...SpanContext) *SpanContexts {
	spanContexts := &SpanContexts{}

	for _, otherSpanContext := range otherSpanContexts {
		otherSpanContext(spanContexts)
	}
	return spanContexts
}


// TraceID sets the pass through header key list into span contexts
func PassThroughHeaderKeyList(passThroughHeaderKeyList string) SpanContext {
	return func(contexts *SpanContexts) {
		contexts.PassThroughHeaderKeyList = passThroughHeaderKeyList
	}
}

// TraceID sets the trace ID into span contexts
func TraceID(traceID string) SpanContext {
	return func(contexts *SpanContexts) {
		contexts.TraceID = traceID
	}
}

// SpanID sets the span ID into span contexts
func SpanID(spanID string) SpanContext {
	return func(contexts *SpanContexts) {
		contexts.SpanID = spanID
	}
}

// ParentSpanID sets the parent span ID into span contexts
func ParentSpanID(parentSpanID string) SpanContext {
	return func(contexts *SpanContexts) {
		contexts.ParentSpanID = parentSpanID
	}
}

// ReplyToAddress sets the reply to address into span contexts for semi-synchronous call
func ReplyToAddress(replyToAddress string) SpanContext {
	return func(contexts *SpanContexts) {
		contexts.ReplyToAddress = replyToAddress
	}
}

// TimeoutMilliseconds sets the timeout(unit milliseconds) into span contexts
func TimeoutMilliseconds(timeoutMilliseconds int) SpanContext {
	return func(contexts *SpanContexts) {
		contexts.TimeoutMilliseconds = timeoutMilliseconds
	}
}

// JudgeParentSpanID returns the parent span ID from span context
func (s *SpanContexts) JudgeParentSpanID() string {
	if "" == s.SpanID {
		return s.TraceID
	}

	return s.SpanID
}

// Copy creates a new span contexts that cloned from current TransactionContexts
func (s *SpanContexts) Copy() *SpanContexts {
	spanContexts := &SpanContexts{}
	spanContexts.SpanID = s.SpanID
	spanContexts.TraceID = s.TraceID
	spanContexts.ParentSpanID = s.ParentSpanID
	spanContexts.ReplyToAddress = s.ReplyToAddress
	spanContexts.PassThroughHeaderKeyList = s.PassThroughHeaderKeyList
	spanContexts.TimeoutMilliseconds = s.TimeoutMilliseconds

	return spanContexts
}

// With sets one or more configurations into span contexts
func (s *SpanContexts) With(otherSpanContexts ...SpanContext) {
	for _, otherSpanContext := range otherSpanContexts {
		otherSpanContext(s)
	}
}
