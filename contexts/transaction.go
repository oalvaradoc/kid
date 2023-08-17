package contexts

// TransactionContexts is a runtime context, used for transaction-related parameter transfer in various methods and modules,
// mainly including transaction manager address, root transaction ID, branch transaction ID, and parent transaction ID
type TransactionContexts struct {
	TransactionAgentAddress      string
	TransactionAgentAddressOld   string
	RootXID                      string
	ParentXID                    string
	BranchXID                    string
	ForceCancelGlobalTransaction bool
}

// TransactionContext sets a parameter into transaction contexts
type TransactionContext func(*TransactionContexts)

// BuildTransactionContexts create a new transaction contexts with one or more setting functions.
func BuildTransactionContexts(otherTransactionContexts ...TransactionContext) *TransactionContexts {
	transactionContexts := &TransactionContexts{}

	for _, otherTransactionContext := range otherTransactionContexts {
		otherTransactionContext(transactionContexts)
	}

	return transactionContexts
}

// TransactionAgentAddress sets the transaction manager address into transaction contexts
func TransactionAgentAddress(transactionAgentAddress string) TransactionContext {
	return func(contexts *TransactionContexts) {
		contexts.TransactionAgentAddress = transactionAgentAddress
	}
}

// TransactionAgentAddressOld sets the old format transaction manager address into transaction contexts for compatible
func TransactionAgentAddressOld(transactionAgentAddressOld string) TransactionContext {
	return func(contexts *TransactionContexts) {
		contexts.TransactionAgentAddressOld = transactionAgentAddressOld
	}
}

// RootXID sets the root xid into transaction contexts
func RootXID(rootXID string) TransactionContext {
	return func(contexts *TransactionContexts) {
		contexts.RootXID = rootXID
	}
}

// ParentXID sets the parent xid into transaction contexts
func ParentXID(parentXID string) TransactionContext {
	return func(contexts *TransactionContexts) {
		contexts.ParentXID = parentXID
	}
}

// BranchXID sets the branch xid into transaction contexts
func BranchXID(branchXID string) TransactionContext {
	return func(contexts *TransactionContexts) {
		contexts.BranchXID = branchXID
	}
}

// ForceCancelGlobalTransaction marks to cancel global transaction
func ForceCancelGlobalTransaction() TransactionContext {
	return func(contexts *TransactionContexts) {
		contexts.ForceCancelGlobalTransaction = true
	}
}

// Copy creates a new transaction contexts that cloned from current TransactionContexts
func (s *TransactionContexts) Copy() *TransactionContexts {
	transactionContexts := &TransactionContexts{}
	transactionContexts.TransactionAgentAddress = s.TransactionAgentAddress
	transactionContexts.TransactionAgentAddressOld = s.TransactionAgentAddressOld
	transactionContexts.RootXID = s.RootXID
	transactionContexts.BranchXID = s.BranchXID
	transactionContexts.ParentXID = s.ParentXID
	transactionContexts.ForceCancelGlobalTransaction = s.ForceCancelGlobalTransaction

	return transactionContexts
}

// With sets one or more configurations into transaction contexts
func (s *TransactionContexts) With(otherTransactionContexts ...TransactionContext) {
	for _, otherTransactionContext := range otherTransactionContexts {
		otherTransactionContext(s)
	}
}
