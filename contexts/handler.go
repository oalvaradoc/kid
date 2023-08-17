package contexts

import (
	"context"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/config"
	"sync"
	"time"
)

// HandlerContexts is a controller that contains the context of request.
type HandlerContexts struct {
	Org                         string
	Az                          string
	Su                          string
	CommonSu                    string
	NodeID                      string
	ServiceID                   string
	InstanceID                  string
	Wks                         string
	Env                         string
	Lang                        string
	SpanContexts                *SpanContexts
	ResponseTemplate            string
	ResponseAutoParseKeyMapping map[string]string
	TransactionContexts         *TransactionContexts
	attach                      map[string]interface{}
	attachRWLock                sync.RWMutex
	StartHandleTime             time.Time
	DownstreamConfigs           map[string]config.Downstream
}

// HandlerContext sets an optional parameter for HandlerContexts
type HandlerContext func(*HandlerContexts)

// BuildHandlerContexts creates a new HandlerContexts with default parameters.
func BuildHandlerContexts(contexts ...HandlerContext) *HandlerContexts {
	handlerContexts := &HandlerContexts{
		SpanContexts:        nil,
		TransactionContexts: nil,
		attach:              make(map[string]interface{}),
	}

	for _, context := range contexts {
		context(handlerContexts)
	}

	return handlerContexts
}

// AddAttach adds a pair of key-value into attach of context
// the same as AppendAttach
func AddAttach(key, value string) HandlerContext {
	return func(contexts *HandlerContexts) {
		if nil == contexts.attach {
			contexts.attach = make(map[string]interface{})
		}
		contexts.attach[key] = value
	}
}

// Lang sets the user lang into HandlerContexts
func Lang(lang string) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.Lang = lang
	}
}

// ServiceID sets the service ID into HandlerContexts
func ServiceID(serviceID string) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.ServiceID = serviceID
	}
}

// NodeID sets the node ID into HandlerContexts
func NodeID(nodeID string) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.NodeID = nodeID
	}
}

// InstanceID sets the instance ID into HandlerContexts
func InstanceID(instanceID string) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.InstanceID = instanceID
	}
}

func DownstreamConfigs(downstreamConfigs map[string]config.Downstream) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.DownstreamConfigs = downstreamConfigs
	}
}

// WKS sets the workspace into HandlerContexts
func WKS(wks string) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.Wks = wks
	}
}

// ENV sets the environment into HandlerContexts
func ENV(env string) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.Env = env
	}
}

// AZ sets the available zone into HandlerContexts
func AZ(az string) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.Az = az
	}
}

// ORG sets the organization into HandlerContexts
func ORG(org string) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.Org = org
	}
}

// ResponseTemplate sets the response template into HandlerContexts
func ResponseTemplate(responseTemplate string) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.ResponseTemplate = responseTemplate
	}
}

// ResponseAutoParseKeyMapping sets the response auto parse key mapping into HandlerContexts
func ResponseAutoParseKeyMapping(responseAutoParseKeyMapping map[string]string) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.ResponseAutoParseKeyMapping = responseAutoParseKeyMapping
	}
}

// SU sets the SU into HandlerContexts
func SU(su string) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.Su = su
	}
}

// CommonSu sets the commonSu into HandlerContexts
func CommonSu(commonSu string) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.CommonSu = commonSu
	}
}

// Attach sets the attach into HandlerContexts
func Attach(attach map[string]interface{}) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.attach = attach
	}
}

// AppendAttach adds a pair of key-value into attach of context
func AppendAttach(key string, value interface{}) HandlerContext {
	return func(contexts *HandlerContexts) {
		if nil == contexts.attach {
			contexts.attach = make(map[string]interface{})
		}
		contexts.attach[key] = value
	}
}

// Span sets the SpanContexts into HandlerContext
func Span(spanContext *SpanContexts) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.SpanContexts = spanContext
	}
}

// Transaction sets the TransactionContexts into HandlerContext
func Transaction(transactionContext *TransactionContexts) HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.TransactionContexts = transactionContext
	}
}

// WithTraceID sets the trace ID into SpanContexts
func WithTraceID(traceID string) HandlerContext {
	return func(contexts *HandlerContexts) {
		if nil == contexts.SpanContexts {
			contexts.SpanContexts = &SpanContexts{}
		}
		contexts.SpanContexts.TraceID = traceID
	}
}

// WithSpanID sets the span ID into SpanContexts
func WithSpanID(spanID string) HandlerContext {
	return func(contexts *HandlerContexts) {
		if nil == contexts.SpanContexts {
			contexts.SpanContexts = &SpanContexts{}
		}
		contexts.SpanContexts.SpanID = spanID
	}
}

// WithParentSpanID sets the parent span ID into SpanContexts
func WithParentSpanID(parentSpanID string) HandlerContext {
	return func(contexts *HandlerContexts) {
		if nil == contexts.SpanContexts {
			contexts.SpanContexts = &SpanContexts{}
		}
		contexts.SpanContexts.ParentSpanID = parentSpanID
	}
}

// WithRootXID sets the root XID into TransactionContexts
func WithRootXID(rootXID string) HandlerContext {
	return func(contexts *HandlerContexts) {
		if nil == contexts.TransactionContexts {
			contexts.TransactionContexts = &TransactionContexts{}
		}
		contexts.TransactionContexts.RootXID = rootXID
	}
}

// WithBranchXID sets the branch XID into TransactionContexts
func WithBranchXID(branchXID string) HandlerContext {
	return func(contexts *HandlerContexts) {
		if nil == contexts.TransactionContexts {
			contexts.TransactionContexts = &TransactionContexts{}
		}
		contexts.TransactionContexts.BranchXID = branchXID
	}
}

// WithParentXID sets the parent span XID into TransactionContexts
func WithParentXID(parentXID string) HandlerContext {
	return func(contexts *HandlerContexts) {
		if nil == contexts.TransactionContexts {
			contexts.TransactionContexts = &TransactionContexts{}
		}
		contexts.TransactionContexts.ParentXID = parentXID
	}
}

// DeleteTransactionPropagationInformation marks delete transaction propagation information when request downstream service
func DeleteTransactionPropagationInformation() HandlerContext {
	return func(contexts *HandlerContexts) {
		contexts.TransactionContexts = nil
	}
}

// Copy creates a new HandlerContexts that cloned from current HandlerContexts
func (h *HandlerContexts) Copy() *HandlerContexts {
	handlerContexts := &HandlerContexts{}
	if nil != h.SpanContexts {
		handlerContexts.SpanContexts = h.SpanContexts.Copy()
	}
	if nil != h.TransactionContexts {
		handlerContexts.TransactionContexts = h.TransactionContexts.Copy()
	}
	handlerContexts.Org = h.Org
	handlerContexts.Az = h.Az
	handlerContexts.Su = h.Su
	handlerContexts.CommonSu = h.CommonSu
	handlerContexts.ServiceID = h.ServiceID
	handlerContexts.NodeID = h.NodeID
	handlerContexts.InstanceID = h.InstanceID
	handlerContexts.Wks = h.Wks
	handlerContexts.Env = h.Env
	handlerContexts.Lang = h.Lang
	handlerContexts.attach = make(map[string]interface{})
	handlerContexts.ResponseTemplate = h.ResponseTemplate

	handlerContexts.attach = h.CloneAttach()
	//handlerContexts.attachRWLock.Lock()
	//for k, v := range h.attach {
	//	handlerContexts.attach[k] = v
	//}
	//handlerContexts.attachRWLock.Unlock()

	handlerContexts.ResponseAutoParseKeyMapping = make(map[string]string)
	for k, v := range h.ResponseAutoParseKeyMapping {
		handlerContexts.ResponseAutoParseKeyMapping[k] = v
	}
	return handlerContexts
}

// With sets one or more parameters into HandlerContext
func (h *HandlerContexts) With(otherHandlerContexts ...HandlerContext) {
	for _, otherHandlerContext := range otherHandlerContexts {
		otherHandlerContext(h)
	}
}

// Inject injects HandlerContexts into app properties of message
func (h *HandlerContexts) Inject(isDeleteTransactionPropagationInformation bool, message *msg.Message) {
	if !isDeleteTransactionPropagationInformation && nil != h.TransactionContexts {
		message.SetAppProperty(constant.RootXIDKey, h.TransactionContexts.RootXID)
		// solve service1(tcc) -> service2(no tcc) -> service3(tcc) problem
		if len(h.TransactionContexts.BranchXID) == 0 {
			message.SetAppProperty(constant.ParentXIDKey, h.TransactionContexts.ParentXID)
		} else {
			message.SetAppProperty(constant.ParentXIDKey, h.TransactionContexts.BranchXID)
		}
		message.SetAppProperty(constant.TransactionAgentAddress, h.TransactionContexts.TransactionAgentAddress)

		message.SetAppProperty(constant.RootXIDKeyOld, h.TransactionContexts.RootXID)
		if len(h.TransactionContexts.BranchXID) == 0 {
			message.SetAppProperty(constant.ParentXIDKeyOld, h.TransactionContexts.ParentXID)
		} else {
			message.SetAppProperty(constant.ParentXIDKeyOld, h.TransactionContexts.BranchXID)
		}
		message.SetAppProperty(constant.TransactionAgentAddressOld, h.TransactionContexts.TransactionAgentAddressOld)
	}

	if nil != h.SpanContexts {
		message.SetAppProperty(constant.KeyTraceID, h.SpanContexts.TraceID)
		message.SetAppProperty(constant.KeySpanID, h.SpanContexts.SpanID)
		message.SetAppProperty(constant.KeyParentSpanID, h.SpanContexts.ParentSpanID)

		message.SetAppProperty(constant.KeyTraceIDOld, h.SpanContexts.TraceID)
		message.SetAppProperty(constant.KeySpanIDOld, h.SpanContexts.SpanID)
		message.SetAppProperty(constant.KeyParentSpanIDOld, h.SpanContexts.ParentSpanID)

		if len(h.SpanContexts.ReplyToAddress) > 0 {
			message.SetAppProperty(constant.RrReplyTo, h.SpanContexts.ReplyToAddress)
		}
	}
}

// RangeAttach is used to range attach in HandlerContexts using customize closure function.
func (h HandlerContexts) RangeAttach(f func(key string, value interface{})) {
	h.attachRWLock.RLock()
	defer func() { h.attachRWLock.RUnlock() }()

	for k, v := range h.attach {
		f(k, v)
	}
}

// DeleteAttach deletes attach value in HandlerContexts by key.
func (h *HandlerContexts) DeleteAttach(key string) {
	h.attachRWLock.Lock()
	defer func() { h.attachRWLock.Unlock() }()

	delete(h.attach, key)
}

// AppendAttach adds a pair of key-value into attach in HandlerContexts
func (h *HandlerContexts) AppendAttach(key string, value interface{}) {
	h.attachRWLock.Lock()
	defer func() { h.attachRWLock.Unlock() }()

	h.attach[key] = value
}

// CloneAttach creates a new map[string]string that copy from attach of HandlerContexts
func (h HandlerContexts) CloneAttach() map[string]interface{} {
	h.attachRWLock.RLock()
	defer func() { h.attachRWLock.RUnlock() }()

	if nil == h.attach {
		return nil
	}
	cloneMaps := make(map[string]interface{})
	for k, v := range h.attach {
		cloneMaps[k] = v
	}
	return cloneMaps

}

// GetAttach gets the value of attach by key
func (h HandlerContexts) GetAttach(key string) (res interface{}, ok bool) {
	h.attachRWLock.RLock()
	defer func() { h.attachRWLock.RUnlock() }()

	r, b := h.attach[key]

	return r, b
}

// GetAttachEitherSilence finds the attach in HandlerContexts value according to the two keys in order
// and return immediately after finding the first value
// returns nil if the both key doesn't exist in the attach
func (h HandlerContexts) GetAttachEitherSilence(key1, key2 string) interface{} {
	h.attachRWLock.RLock()
	defer func() { h.attachRWLock.RUnlock() }()

	if nil == h.attach {
		return ""
	}

	if v, ok := h.attach[key1]; ok {
		return v
	}

	return h.attach[key2]

}

// GetAttachSilence finds the attach value according to the key
// returns nil if the key doesn't exists in the attach
func (h HandlerContexts) GetAttachSilence(key string) interface{} {
	h.attachRWLock.RLock()
	defer func() { h.attachRWLock.RUnlock() }()

	if nil == h.attach {
		return ""
	}
	return h.attach[key]
}

// IsRootTransaction returns whether the current request is root transaction.
func (h *HandlerContexts) IsRootTransaction() bool {
	return nil == h.TransactionContexts || "" == h.TransactionContexts.RootXID
}

// BuildContextFromParent creates a new context.Context from parent context.Context with one or more additional parameters.
func BuildContextFromParent(ctx context.Context, handlerContexts ...HandlerContext) (context.Context, *HandlerContexts) {
	handlerContext := BuildHandlerContexts(handlerContexts...)
	return context.WithValue(ctx, constant.HandlerContextsKey, handlerContext), handlerContext
}

// BuildContextFromParentWithHandlerContexts creates a new context.Context from parent context.Context with HandlerContexts
func BuildContextFromParentWithHandlerContexts(ctx context.Context, handlerContexts *HandlerContexts) context.Context {
	return context.WithValue(ctx, constant.HandlerContextsKey, handlerContexts)
}

// HandlerContextsFromContext returns HandlerContexts from context.Context,
// returns nil if the HandlerContexts key doesn't exist in the context.Context
func HandlerContextsFromContext(ctx context.Context) *HandlerContexts {
	handlerContexts := ctx.Value(constant.HandlerContextsKey)
	if nil == handlerContexts {
		return nil
	}
	return handlerContexts.(*HandlerContexts)
}
