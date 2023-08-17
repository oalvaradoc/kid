package transaction

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"testing"
)

var interceptor = &Interceptor{}

func TestInterceptor_String(t *testing.T) {
	nameOfInterceptor := fmt.Sprintf("%s", interceptor)
	assert.Equal(t, nameOfInterceptor, constant.InterceptorTransaction)
}

func TestInterceptor_PreHandle(t *testing.T) {
	ctx := context.Background()
	err := interceptor.PreHandle(ctx, nil)
	assert.Nil(t, err)

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
	newCtx := contexts.BuildContextFromParentWithHandlerContexts(ctx, handlerContexts)
	assert.NotEqual(t, ctx, newCtx)
	err = interceptor.PreHandle(ctx, nil)
	assert.Nil(t, err)

	request := &msg.Message{}

	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)

	request.SetAppProperty(constant.RootXIDKeyOld, "old root xid")
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.RootXID, "old root xid")

	request.SetAppProperty(constant.RootXIDKey, "new root xid")
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.RootXID, "new root xid")

	request.DeleteProperty(constant.RootXIDKeyOld)
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.RootXID, "new root xid")

	request.SetAppProperty(constant.ParentXIDKeyOld, "old parent xid")
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.ParentXID, "old parent xid")

	request.SetAppProperty(constant.ParentXIDKey, "new parent xid")
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.ParentXID, "new parent xid")

	request.DeleteProperty(constant.ParentXIDKeyOld)
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.ParentXID, "new parent xid")

	request.SetAppProperty(constant.BranchXIDKeyOld, "old branch xid")
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.BranchXID, "old branch xid")

	request.SetAppProperty(constant.BranchXIDKey, "new branch xid")
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.BranchXID, "new branch xid")

	request.DeleteProperty(constant.BranchXIDKeyOld)
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.BranchXID, "new branch xid")

	// testing new direct mode
	directAddress := "http://127.0.0.1:7070|/v1/transaction/start|/v1/transaction/join|/v1/transaction/end"
	request.SetAppProperty(constant.TransactionAgentAddress, directAddress)
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.TransactionAgentAddress, directAddress)

	// testing new mesh mode
	meshAddress := `ORG|WKS|ENV|SU|NODE ID|Instance ID|TxnBegin TopicID|TxnJoin TopicID|TxnEnd TopicID`
	request.SetAppProperty(constant.TransactionAgentAddress, meshAddress)
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.TransactionAgentAddress, meshAddress)

	request.DeleteProperty(constant.TransactionAgentAddress)
	// testing old direct mode
	directAddress = "http://127.0.0.1:7070|/v1/transaction/start|/v1/transaction/join|/v1/transaction/end"
	request.SetAppProperty(constant.TransactionAgentAddressOld, directAddress)
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.TransactionAgentAddress, directAddress)

	// testing old mesh mode
	meshAddress = `ORG|SU|NODE ID|Instance ID|TxnRegister TopicID|TxnEnlist TopicID|TxnTryResultReport TopicID`
	request.SetAppProperty(constant.TransactionAgentAddressOld, meshAddress)
	err = interceptor.PreHandle(newCtx, request)
	assert.Nil(t, err)
	t.Logf("contexts of transaction:%++v", handlerContexts.TransactionContexts)
	assert.Equal(t, handlerContexts.TransactionContexts.TransactionAgentAddress, `ORG|wks|env|SU|NODE ID|Instance ID|TxnRegister TopicID|TxnEnlist TopicID|TxnTryResultReport TopicID`)
}

func TestInterceptor_PostHandle(t *testing.T) {
	ctx := context.Background()
	err := interceptor.PostHandle(ctx, nil, nil)
	assert.Nil(t, err)
}
