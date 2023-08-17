package transaction

import (
	"context"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"strings"
)

// Interceptor is an interceptor for transaction parameters propagation
type Interceptor struct{}

// PreHandle checks whether the upstream service request exists transaction contexts
// and inject the transaction contexts into handler contexts if the transaction contexts exists in the request
func (l *Interceptor) PreHandle(ctx context.Context, request *msg.Message) error {
	handlerContexts := contexts.HandlerContextsFromContext(ctx)
	if nil != handlerContexts {
		rootXID := ""
		parentXID := ""
		branchXID := ""
		transactionAgentAddress := ""
		if nil != request {
			rootXID = request.GetAppPropertyEitherSilence(constant.RootXIDKey, constant.RootXIDKeyOld)
			parentXID = request.GetAppPropertyEitherSilence(constant.ParentXIDKey, constant.ParentXIDKeyOld)
			branchXID = request.GetAppPropertyEitherSilence(constant.BranchXIDKey, constant.BranchXIDKeyOld)

			if address, ok := request.GetAppProperty(constant.TransactionAgentAddress); ok {
				// "ORG" + | + "WKS" + "|" + "ENV"+ "|" + "SU" + "|" + NODE ID + "|" + Instance ID + "|" + TxnBegin TopicID + "|" + TxnJoin TopicID + "|" + TxnEnd TopicID
				// or direct mode:DXC server address + "|" + Transaction Begin URL + "|" + Transaction Join URL + "|" + Transaction End URL
				transactionAgentAddress = address
			} else {
				// "ORG"  "|" + "SU" + "|" + NODE ID + "|" + Instance ID + "|" + TxnBegin TopicID + "|" + TxnJoin TopicID + "|" + TxnEnd TopicID
				if address, ok = request.GetAppProperty(constant.TransactionAgentAddressOld); ok {
					serverAddressElem := strings.Split(address, constant.ParticipantAddressSplitChar)
					if len(serverAddressElem) == 7 {
						// mesh mode
						transactionAgentAddress = serverAddressElem[0] + "|" +
							handlerContexts.Wks + "|" +
							handlerContexts.Env + "|" +
							serverAddressElem[1] + "|" +
							serverAddressElem[2] + "|" +
							serverAddressElem[3] + "|" +
							serverAddressElem[4] + "|" +
							serverAddressElem[5] + "|" +
							serverAddressElem[6]
					} else {
						// direct mode
						// DTS server address + "|" + Transaction register URL + "|" + Transaction enlist URL + "|" + Transaction try result report URL
						transactionAgentAddress = address
					}
				}
			}
		}

		handlerContexts.With(
			contexts.Transaction(
				contexts.BuildTransactionContexts(
					contexts.RootXID(rootXID),
					contexts.ParentXID(parentXID),
					contexts.BranchXID(branchXID),
					contexts.TransactionAgentAddress(transactionAgentAddress)),
			),
		)
	}
	return nil
}

// PostHandle does nothing
func (l *Interceptor) PostHandle(ctx context.Context, request *msg.Message, response *msg.Message) error {
	// clean
	return nil
}

func (l Interceptor) String() string {
	return constant.InterceptorTransaction
}
