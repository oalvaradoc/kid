package contexts

import (
	"context"
	kitCommon "git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/msg"
	"testing"
)

func TestBuildHandlerContexts(t *testing.T) {
	handlerContexts := BuildHandlerContexts()
	assert.NotNil(t, handlerContexts)
}

func TestOptionsSet(t *testing.T) {
	handlerContexts := HandlerContexts{}
	opts := []HandlerContext{AddAttach("key1", "value1")}
	opts = append(opts, Lang("en-US"))
	opts = append(opts, NodeID("node-id-0001"))
	opts = append(opts, InstanceID("instance-001"))
	opts = append(opts, WKS("workspace1"))
	opts = append(opts, ENV("environment1"))
	opts = append(opts, AZ("az1"))
	opts = append(opts, ORG("org1"))
	opts = append(opts, SU("su1"))
	opts = append(opts, ResponseTemplate("response template"))

	for _, opt := range opts {
		opt(&handlerContexts)
	}

	assert.NotNil(t, handlerContexts.attach)
	assert.Equal(t, handlerContexts.Lang, "en-US")
	assert.Equal(t, handlerContexts.NodeID, "node-id-0001")
	assert.Equal(t, handlerContexts.InstanceID, "instance-001")
	assert.Equal(t, handlerContexts.Wks, "workspace1")
	assert.Equal(t, handlerContexts.Env, "environment1")
	assert.Equal(t, handlerContexts.Az, "az1")
	assert.Equal(t, handlerContexts.Org, "org1")
	assert.Equal(t, handlerContexts.Su, "su1")
	assert.Equal(t, handlerContexts.ResponseTemplate, "response template")
}

func TestResponseAutoParseKeyMapping(t *testing.T) {
	handlerContexts := HandlerContexts{}
	opt := ResponseAutoParseKeyMapping(map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	})
	opt(&handlerContexts)

	assert.Equal(t, len(handlerContexts.ResponseAutoParseKeyMapping), 3)
	assert.Equal(t, handlerContexts.ResponseAutoParseKeyMapping["key2"], "value2")
	assert.Equal(t, handlerContexts.ResponseAutoParseKeyMapping["key1"], "value1")
	assert.Equal(t, handlerContexts.ResponseAutoParseKeyMapping["key3"], "value3")
}

func TestAttach(t *testing.T) {
	handlerContexts := HandlerContexts{}
	opt := Attach(map[string]interface{}{
		"k-A": "v-A",
		"k-B": "v-B",
		"k-C": "v-C",
		"k-D": "v-D",
	})
	opt(&handlerContexts)
	assert.Equal(t, len(handlerContexts.attach), 4)
	assert.Equal(t, handlerContexts.attach["k-A"], "v-A")
	assert.Equal(t, handlerContexts.attach["k-C"], "v-C")
}

func TestAppendAttach(t *testing.T) {
	handlerContexts := HandlerContexts{}
	opt1 := AppendAttach("k-A", "v-A")
	opt2 := AppendAttach("k-B", "v-B")
	opt1(&handlerContexts)
	opt2(&handlerContexts)

	assert.Equal(t, len(handlerContexts.attach), 2)
	assert.Equal(t, handlerContexts.attach["k-A"], "v-A")
	assert.Equal(t, handlerContexts.attach["k-B"], "v-B")
}

func TestSpan(t *testing.T) {
	handlerContexts := HandlerContexts{}
	opt := Span(&SpanContexts{
		TraceID:             "traceID1",
		SpanID:              "spanID1",
		ParentSpanID:        "parentSpanID1",
		ReplyToAddress:      "replyToAddress1",
		TimeoutMilliseconds: 30,
	})

	opt(&handlerContexts)
	assert.True(t, handlerContexts.SpanContexts != nil)
	assert.Equal(t, handlerContexts.SpanContexts.TraceID, "traceID1")
	assert.Equal(t, handlerContexts.SpanContexts.SpanID, "spanID1")
	assert.Equal(t, handlerContexts.SpanContexts.ParentSpanID, "parentSpanID1")
	assert.Equal(t, handlerContexts.SpanContexts.ReplyToAddress, "replyToAddress1")
	assert.Equal(t, handlerContexts.SpanContexts.TimeoutMilliseconds, 30)
}

func TestTransaction(t *testing.T) {
}

func TestWithTraceID(t *testing.T) {
	handlerContexts := HandlerContexts{}
	opt := WithTraceID("traceID1")
	opt(&handlerContexts)
	assert.True(t, handlerContexts.SpanContexts != nil)
	assert.Equal(t, handlerContexts.SpanContexts.TraceID, "traceID1")
}

func TestWithSpanID(t *testing.T) {
	handlerContexts := HandlerContexts{}
	opt := WithSpanID("spanID1")
	opt(&handlerContexts)
	assert.True(t, handlerContexts.SpanContexts != nil)
	assert.Equal(t, handlerContexts.SpanContexts.SpanID, "spanID1")
}

func TestWithParentSpanID(t *testing.T) {
	handlerContexts := HandlerContexts{}
	opt := WithParentSpanID("parentSpanID1")
	opt(&handlerContexts)
	assert.True(t, handlerContexts.SpanContexts != nil)
	assert.Equal(t, handlerContexts.SpanContexts.ParentSpanID, "parentSpanID1")
}

func TestWithRootXID(t *testing.T) {
	handlerContexts := HandlerContexts{}
	opt := WithRootXID("rootXID1")
	opt(&handlerContexts)

	assert.True(t, handlerContexts.TransactionContexts != nil)
	assert.Equal(t, handlerContexts.TransactionContexts.RootXID, "rootXID1")
}

func TestWithBranchXID(t *testing.T) {
	handlerContexts := HandlerContexts{}
	opt := WithBranchXID("branchXID1")
	opt(&handlerContexts)

	assert.True(t, handlerContexts.TransactionContexts != nil)
	assert.Equal(t, handlerContexts.TransactionContexts.BranchXID, "branchXID1")
}

func TestWithParentXID(t *testing.T) {
	handlerContexts := HandlerContexts{}
	opt := WithParentXID("parentXID")
	opt(&handlerContexts)

	assert.True(t, handlerContexts.TransactionContexts != nil)
	assert.Equal(t, handlerContexts.TransactionContexts.ParentXID, "parentXID")
}

func TestDeleteTransactionPropagationInformation(t *testing.T) {
	handlerContexts := HandlerContexts{
		TransactionContexts: &TransactionContexts{},
	}
	opt := DeleteTransactionPropagationInformation()
	opt(&handlerContexts)

	assert.True(t, nil == handlerContexts.TransactionContexts)
}

func TestHandlerContextsCopy(t *testing.T) {
	handlerContexts := HandlerContexts{
		Org:        "org1",
		Az:         "az1",
		Su:         "su1",
		NodeID:     "nodeID1",
		InstanceID: "instanceID1",
		Wks:        "wks1",
		Env:        "env1",
		Lang:       "lang1",
		TransactionContexts: &TransactionContexts{
			TransactionAgentAddress:      "transactionAgentAddress1",
			TransactionAgentAddressOld:   "transactionAgentAddressOld1",
			RootXID:                      "rootXID1",
			ParentXID:                    "parentXID1",
			BranchXID:                    "branchXID1",
			ForceCancelGlobalTransaction: true,
		},
		SpanContexts: &SpanContexts{
			TraceID:             "traceID1",
			SpanID:              "spanID1",
			ParentSpanID:        "parentSpanID1",
			ReplyToAddress:      "replyToAddress1",
			TimeoutMilliseconds: 10,
		},
	}

	handlerContextsCopied := handlerContexts.Copy()
	assert.NotNil(t, handlerContextsCopied)
	assert.Equal(t, handlerContextsCopied.Org, handlerContexts.Org)
	assert.Equal(t, handlerContextsCopied.Az, handlerContexts.Az)
	assert.Equal(t, handlerContextsCopied.Su, handlerContexts.Su)
	assert.Equal(t, handlerContextsCopied.NodeID, handlerContexts.NodeID)
	assert.Equal(t, handlerContextsCopied.InstanceID, handlerContexts.InstanceID)
	assert.Equal(t, handlerContextsCopied.Wks, handlerContexts.Wks)
	assert.Equal(t, handlerContextsCopied.Env, handlerContexts.Env)
	assert.Equal(t, handlerContextsCopied.Lang, handlerContexts.Lang)
	assert.NotNil(t, handlerContextsCopied.TransactionContexts)
	assert.NotNil(t, handlerContextsCopied.SpanContexts)

	assert.Equal(t, handlerContextsCopied.TransactionContexts.ParentXID, handlerContexts.TransactionContexts.ParentXID)
	assert.Equal(t, handlerContextsCopied.TransactionContexts.BranchXID, handlerContexts.TransactionContexts.BranchXID)
	assert.Equal(t, handlerContextsCopied.TransactionContexts.RootXID, handlerContexts.TransactionContexts.RootXID)
	assert.Equal(t, handlerContextsCopied.TransactionContexts.ForceCancelGlobalTransaction, handlerContexts.TransactionContexts.ForceCancelGlobalTransaction)
	assert.Equal(t, handlerContextsCopied.TransactionContexts.TransactionAgentAddressOld, handlerContexts.TransactionContexts.TransactionAgentAddressOld)
	assert.Equal(t, handlerContextsCopied.TransactionContexts.TransactionAgentAddress, handlerContexts.TransactionContexts.TransactionAgentAddress)

	assert.Equal(t, handlerContextsCopied.SpanContexts.ParentSpanID, handlerContexts.SpanContexts.ParentSpanID)
	assert.Equal(t, handlerContextsCopied.SpanContexts.SpanID, handlerContexts.SpanContexts.SpanID)
	assert.Equal(t, handlerContextsCopied.SpanContexts.TraceID, handlerContexts.SpanContexts.TraceID)
	assert.Equal(t, handlerContextsCopied.SpanContexts.ReplyToAddress, handlerContexts.SpanContexts.ReplyToAddress)
	assert.Equal(t, handlerContextsCopied.SpanContexts.TimeoutMilliseconds, handlerContexts.SpanContexts.TimeoutMilliseconds)

	handlerContexts.TransactionContexts = nil
	handlerContexts.SpanContexts = nil

	assert.True(t, handlerContexts.Copy().TransactionContexts == nil)
	assert.True(t, handlerContexts.Copy().SpanContexts == nil)

}

func TestInject(t *testing.T) {
	handlerContexts := HandlerContexts{
		Org:        "org1",
		Az:         "az1",
		Su:         "su1",
		NodeID:     "nodeID1",
		InstanceID: "instanceID1",
		Wks:        "wks1",
		Env:        "env1",
		Lang:       "lang1",
		TransactionContexts: &TransactionContexts{
			TransactionAgentAddress:      "transactionAgentAddress1",
			TransactionAgentAddressOld:   "transactionAgentAddressOld1",
			RootXID:                      "rootXID1",
			ParentXID:                    "parentXID1",
			BranchXID:                    "branchXID1",
			ForceCancelGlobalTransaction: true,
		},
		SpanContexts: &SpanContexts{
			TraceID:             "traceID1",
			SpanID:              "spanID1",
			ParentSpanID:        "parentSpanID1",
			ReplyToAddress:      "replyToAddress1",
			TimeoutMilliseconds: 10,
		},
		attach: map[string]interface{}{
			"k1": "v1",
		},
	}

	msg := msg.Message{}

	handlerContexts.Inject(false, &msg)
	t.Logf("message:[%++v]", msg)
	assert.NotNil(t, msg.TopicAttribute)
}

func TestGetAttach(t *testing.T) {
	handlerContexts := HandlerContexts{
		Org:        "org1",
		Az:         "az1",
		Su:         "su1",
		NodeID:     "nodeID1",
		InstanceID: "instanceID1",
		Wks:        "wks1",
		Env:        "env1",
		Lang:       "lang1",
		TransactionContexts: &TransactionContexts{
			TransactionAgentAddress:      "transactionAgentAddress1",
			TransactionAgentAddressOld:   "transactionAgentAddressOld1",
			RootXID:                      "rootXID1",
			ParentXID:                    "parentXID1",
			BranchXID:                    "branchXID1",
			ForceCancelGlobalTransaction: true,
		},
		SpanContexts: &SpanContexts{
			TraceID:             "traceID1",
			SpanID:              "spanID1",
			ParentSpanID:        "parentSpanID1",
			ReplyToAddress:      "replyToAddress1",
			TimeoutMilliseconds: 10,
		},
		attach: map[string]interface{}{
			"k1": "v1",
		},
	}

	v, ok := handlerContexts.GetAttach("k1")
	assert.True(t, ok)
	assert.Equal(t, v, "v1")

	v = handlerContexts.GetAttachSilence("k1")
	assert.Equal(t, v, "v1")

	v, ok = handlerContexts.GetAttach("k2")
	assert.False(t, ok)
	assert.Equal(t, v, nil)

	v = handlerContexts.GetAttachSilence("k2")
	assert.Equal(t, v, nil)

	v = handlerContexts.GetAttachEitherSilence("k3", "k1")
	assert.Equal(t, v, "v1")

}

func TestIsRootTransaction(t *testing.T) {
	handlerContexts := HandlerContexts{
		TransactionContexts: &TransactionContexts{
			TransactionAgentAddress:      "transactionAgentAddress1",
			TransactionAgentAddressOld:   "transactionAgentAddressOld1",
			RootXID:                      "rootXID1",
			ParentXID:                    "parentXID1",
			BranchXID:                    "branchXID1",
			ForceCancelGlobalTransaction: true,
		},
	}

	assert.False(t, handlerContexts.IsRootTransaction())

	handlerContexts.TransactionContexts.RootXID = ""
	assert.True(t, handlerContexts.IsRootTransaction())
}

func TestBuildContextFromParent(t *testing.T) {
	ctx, handlerContexts := BuildContextFromParent(context.Background())
	assert.NotNil(t, handlerContexts)
	assert.NotNil(t, ctx.Value(kitCommon.HandlerContextsKey))
}

func TestBuildContextFromParentWithHandlerContexts(t *testing.T) {
	handlerContexts := BuildContextFromParentWithHandlerContexts(context.Background(), &HandlerContexts{})
	assert.NotNil(t, handlerContexts)
	assert.NotNil(t, handlerContexts.Value(kitCommon.HandlerContextsKey))
}

func TestHandlerContextsFromContext(t *testing.T) {
	handlerContexts := HandlerContextsFromContext(context.Background())
	assert.True(t, handlerContexts == nil)
}
