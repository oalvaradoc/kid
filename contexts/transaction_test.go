package contexts

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestBuildTransactionContexts(t *testing.T) {
	transactionContexts := BuildTransactionContexts()
	assert.NotNil(t, transactionContexts)
}

func TestSetTransactionOptions(t *testing.T) {
	transactionContexts := BuildTransactionContexts()
	opts := []TransactionContext{}
	opts = append(opts, TransactionAgentAddress("transactionAgentAddress1"))
	opts = append(opts, TransactionAgentAddressOld("transactionAgentAddressOld1"))
	opts = append(opts, RootXID("rootXID1"))
	opts = append(opts, ParentXID("parentXID1"))
	opts = append(opts, BranchXID("branchXID1"))
	opts = append(opts, ForceCancelGlobalTransaction())

	for _, opt := range opts {
		opt(transactionContexts)
	}

	assert.Equal(t, transactionContexts.ParentXID, "parentXID1")
	assert.Equal(t, transactionContexts.RootXID, "rootXID1")
	assert.Equal(t, transactionContexts.BranchXID, "branchXID1")
	assert.Equal(t, transactionContexts.TransactionAgentAddress, "transactionAgentAddress1")
	assert.Equal(t, transactionContexts.TransactionAgentAddressOld, "transactionAgentAddressOld1")
	assert.True(t, transactionContexts.ForceCancelGlobalTransaction)

	transactionContextsCopied := transactionContexts.Copy()
	assert.True(t, nil != transactionContextsCopied)
	assert.Equal(t, transactionContexts.ParentXID, transactionContextsCopied.ParentXID)
	assert.Equal(t, transactionContexts.RootXID, transactionContextsCopied.RootXID)
	assert.Equal(t, transactionContexts.BranchXID, transactionContextsCopied.BranchXID)
	assert.Equal(t, transactionContexts.TransactionAgentAddress, transactionContextsCopied.TransactionAgentAddress)
	assert.Equal(t, transactionContexts.TransactionAgentAddressOld, transactionContextsCopied.TransactionAgentAddressOld)
	assert.Equal(t, transactionContexts.ForceCancelGlobalTransaction, transactionContextsCopied.ForceCancelGlobalTransaction)

}
