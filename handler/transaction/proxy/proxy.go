package proxy

import (
	"context"
	"fmt"
	kitClient "git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/client/mesh"
	"git.multiverse.io/eventkit/kit/client/mesh/wrapper/apm"
	"git.multiverse.io/eventkit/kit/client/mesh/wrapper/logging"
	"git.multiverse.io/eventkit/kit/client/mesh/wrapper/trace"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/compensable"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/base"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/handler/transaction"
	"git.multiverse.io/eventkit/kit/handler/transaction/imports"
	"git.multiverse.io/eventkit/kit/handler/transaction/manager"
	"git.multiverse.io/eventkit/kit/handler/transaction/register"
	"git.multiverse.io/eventkit/kit/log"
	"reflect"
	"strings"
	"sync"
)

var (
	defaultCallWrapperOption = kitClient.DefaultWrapperCall(
		&apm.Wrapper{},
		&logging.Wrapper{},
		&trace.Wrapper{},
	)

	//transactionClient kitClient.Client = mesh.NewMeshClient(defaultCallWrapperOption)
	transactionClient     kitClient.Client = nil
	transactionClientOnce sync.Once
)

// TransactionProxy Transaction Proxy information
type /**/ TransactionProxy struct {
	ctx               context.Context
	txnManager        manager.TxnManager
	txInvocation      *client.TxInvocation
	tryInstance       reflect.Value
	transactionConfig *config.Transaction
}

// NewTransactionProxy creates a new transaction proxy
func NewTransactionProxy(ctx context.Context, txnInstance interface{}, compensable *compensable.Compensable) (*TransactionProxy, error) {
	log.Debugf(ctx, "transaction begin Transaction service[%v]", compensable)
	var serviceName = compensable.ServiceName
	if serviceName == "" {
		var tType reflect.Type
		switch txnInstance.(type) {
		case reflect.Value:
			tType = txnInstance.(reflect.Value).Type()
		default:
			tType = reflect.TypeOf(txnInstance)
		}

		if tType.Kind() == reflect.Ptr {
			serviceName = tType.Elem().Name()
		} else {
			serviceName = tType.Name()
		}
	}

	txInvocation := register.GetCompensableService(serviceName)
	if txInvocation == nil {
		return nil, errors.Errorf(constant.SystemInternalError, "repository no txInvocation, please check this service:%s has registered, context:[%++v]", serviceName, ctx)
	}

	if err := txInvocation.Compensable.DoUntilHasSucceeded(func() error {
		//check whether has executed imports.EnableTransactionSupports
		if false == imports.HasEnabledTransactionSupport {
			return errors.Errorf(constant.SystemInternalError, "Please execute imports.EnableTransactionSupports at startup first, context:[%++v]", ctx)
		}
		return nil
	}); nil != err {
		log.Errorf(ctx, "Transaction method check failed, error=%++v", err)
		return nil, err
	}
	var tryInstance reflect.Value
	switch txnInstance.(type) {
	case reflect.Value:
		tryInstance = txnInstance.(reflect.Value)
	default:
		tryInstance = reflect.ValueOf(txnInstance)
	}
	transactionConfig := ctx.Value(constant.ContextTransactionKey).(*config.Transaction)

	transactionClientOnce.Do(func() {
		transactionClient = mesh.NewMeshClient(defaultCallWrapperOption)
	})

	result := &TransactionProxy{
		txInvocation:      txInvocation,
		ctx:               ctx,
		txnManager:        manager.NewTxnManager(ctx, transactionConfig, transactionClient),
		tryInstance:       tryInstance,
		transactionConfig: transactionConfig,
	}
	return result, nil
}

// Do wraps the logic of service, do some transaction logic
func (p *TransactionProxy) Do(inputParams ...interface{}) []reflect.Value {
	var serverAddress string
	var err error

	handlerContexts := contexts.HandlerContextsFromContext(p.ctx)
	log.Debugf(p.ctx, "Start do Transaction proxy, Transaction context:[%++v], SpanCtx:[%++v]", p.ctx, handlerContexts.SpanContexts.SpanID)
	paramData, err := util.SerialParams(inputParams...)
	if err != nil {
		err := errors.Errorf(constant.SystemInternalError, "Serialization InputParam fail, err:%s, context:[%++v]", err, p.ctx)
		return []reflect.Value{reflect.ValueOf(err)}
	}

	ins := p.tryInstance.Interface().(base.HandlerInterface)
	currentSu := ins.GetCurrentSU()
	headers := ins.GetRequestHeader()
	if nil == headers {
		headers = make(map[string]string)
	}
	headers[constant.CurrentSU] = currentSu
	isRoot := isRoot(handlerContexts.TransactionContexts)
	if isRoot {
		if serverAddress, err = p.rootBegin(handlerContexts, paramData, headers); nil != err {
			log.Errorf(p.ctx, "root begin failed, error: [%s]", errors.ErrorToString(err))
			return []reflect.Value{reflect.ValueOf(errors.New(constant.TransactionBeginError, err))}
		}
	} else {
		if p.transactionConfig.IsPropagator ||
			p.txInvocation.Compensable.IsPropagator ||
			(nil != p.transactionConfig &&
				nil != p.transactionConfig.PropagatorServicesMap &&
				p.transactionConfig.PropagatorServicesMap[p.txInvocation.Compensable.ServiceName]) {
			log.Debug(p.ctx, "Transaction propagator, only invoke try method and propagate Transaction context.")
			cmpTxnCtx := p.ctx
			// set parent XID as current transaction XID(branch xid)
			handlerContexts.With(contexts.WithBranchXID(handlerContexts.SpanContexts.ParentSpanID))
			targetMethodResult := p.txInvocation.InvokeTryMethod(cmpTxnCtx, p.tryInstance, inputParams...)
			return targetMethodResult
		}
		if serverAddress, err = p.branchJoin(handlerContexts, paramData, headers); nil != err {
			log.Errorf(p.ctx, "branch join failed, error: [%s]", errors.ErrorToString(err))
			return []reflect.Value{reflect.ValueOf(errors.New(constant.TransactionJoinError, err))}
		}
	}

	// call the try method
	targetMethodResult := p.txInvocation.InvokeTryMethod(p.ctx, p.tryInstance, inputParams...)

	lastReturnValue := targetMethodResult[len(targetMethodResult)-1]
	// check result
	isOk := (nil == handlerContexts.TransactionContexts || !handlerContexts.TransactionContexts.ForceCancelGlobalTransaction) &&
		valueIsNil(lastReturnValue)
	respHeader := ins.GetResponseHeader()
	log.Debugf(p.ctx, "Response header:%++v", respHeader)
	_, ok := respHeader[constant.MarkAsErrorResponseKey]
	if ok {
		log.Debugf(p.ctx, "find mark error response key:%s, need invoke cancel", constant.MarkAsErrorResponseKey)
		isOk = false
	}

	// check transaction context, avoid business modify
	err = p.checkContext(handlerContexts)
	if err != nil {
		log.Errorf(p.ctx, "check transaction context(business modify error) failed, err: [%s]", err)
		return []reflect.Value{reflect.ValueOf(err)}
	}

	log.Debugf(p.ctx, "Call try method result:[%v]", isOk)

	//when non-root return
	if !isRoot {
		if !isOk && p.transactionConfig.TryFailedIgnoreCallbackCancel {
			// Report branch try failed, ignore current service to cancel
			if errCode, err := p.doEnd(isOk, serverAddress, handlerContexts, nil, nil); err != nil {
				log.Errorf(p.ctx, "doEnd failed, err: [%s], errCode: [%d]", err, errCode)
				return []reflect.Value{reflect.ValueOf(err)}
			}
		}
		return targetMethodResult
	}

	// R/R mode, only ROOT transaction need report try execute result!
	// sync report try execute result
	tryReturnError, _ := lastReturnValue.Interface().(*errors.Error)
	rootKVToSecondStageHeaders := getRootKVToSecondStageHeaders(respHeader)
	if errCode, err := p.doEnd(isOk, serverAddress, handlerContexts, tryReturnError, rootKVToSecondStageHeaders); err != nil {
		log.Errorf(p.ctx, "doEnd failed, err: [%s], errCode: [%d]", err, errCode)
		// The last return type must be error
		var finalErrorCode string
		switch errCode {
		case constant.TxnEndFailedBranchesNotAllCallbackSuccess:
			{
				if isOk {
					finalErrorCode = constant.TransactionEndCallbackConfirmError
				} else {
					finalErrorCode = constant.TransactionEndCallbackCancelError
				}
			}
		case constant.TxnEndFailedTimeOut:
			{
				if isOk {
					finalErrorCode = constant.TransactionEndCallbackConfirmTimeout
				} else {
					finalErrorCode = constant.TransactionEndCallbackCancelTimeout
				}
			}
		default:
			{
				finalErrorCode = constant.TransactionEndOtherError
			}
		}
		if isOk {
			retErr := errors.Wrap(finalErrorCode, err, 0)
			return []reflect.Value{reflect.ValueOf(retErr)}
		}
		er, ok := lastReturnValue.Interface().(*errors.Error)
		if !ok {
			retErr := errors.Wrap(finalErrorCode, err, 0)
			return []reflect.Value{reflect.ValueOf(retErr)}
		}
		finalErrorCode = er.ErrorCode
		err := fmt.Errorf("business fail:[%s] and tcc do end fail:[%s]", er.Err, err)
		retErr := errors.Wrap(finalErrorCode, err, 0)
		return []reflect.Value{reflect.ValueOf(retErr)}
	}

	return targetMethodResult
}

func (p *TransactionProxy) rootBegin(handlerContexts *contexts.HandlerContexts, paramData []byte, headers map[string]string) (serverAddress string, err error) {
	var serverAddressOld string
	if strings.EqualFold(constant.CommDirect, p.transactionConfig.CommType) {
		// server address + "|" + TxnBegin Path + "|" + TxnJoin Path + "|" + Result try report Path
		serverAddress = p.transactionConfig.TransactionServer.AddressURL + "|" +
			p.transactionConfig.TransactionServer.TxnBeginURLPath + "|" +
			p.transactionConfig.TransactionServer.TxnJoinURLPath + "|" +
			p.transactionConfig.TransactionServer.TxnEndURLPath
	} else {
		// "ORG" + "|" + "WKS" + "|" + "ENV"+ "|" + SU + "|" + NODE ID + "|" + Instance ID + "|" + TxnBegin TopicID + "|" + TxnJoin TopicID + "|" + Result try report TopicID
		serverAddress = p.transactionConfig.TransactionServer.Org + "|" +
			p.transactionConfig.TransactionServer.Wks + "|" +
			p.transactionConfig.TransactionServer.Env + "|" +
			p.transactionConfig.TransactionServer.Su + "|" +
			p.transactionConfig.TransactionServer.NodeID + "|" +
			p.transactionConfig.TransactionServer.InstanceID + "|" +
			p.transactionConfig.TransactionServer.TxnBeginEventID + "|" +
			p.transactionConfig.TransactionServer.TxnJoinEventID + "|" +
			p.transactionConfig.TransactionServer.TxnEndEventID

		// "ORG" + "|" + SU + "|" + NODE ID + "|" + Instance ID + "|" + TxnBegin TopicID + "|" + TxnJoin TopicID + "|" + Result try report TopicID
		serverAddressOld = p.transactionConfig.TransactionServer.Org + "|" +
			p.transactionConfig.TransactionServer.Su + "|" +
			p.transactionConfig.TransactionServer.NodeID + "|" +
			p.transactionConfig.TransactionServer.InstanceID + "|" +
			p.transactionConfig.TransactionServer.TxnBeginEventID + "|" +
			p.transactionConfig.TransactionServer.TxnJoinEventID + "|" +
			p.transactionConfig.TransactionServer.TxnEndEventID
	}
	// register ROOT transaction
	rootXid, err := p.txnManager.TxnBegin(handlerContexts, serverAddress, p.txInvocation.Compensable.CompensableFlagSet, paramData, p.txInvocation.Compensable.ServiceName, headers)
	if err != nil {
		return serverAddress, err
	}
	spanContexts := handlerContexts.SpanContexts
	transactionContexts := handlerContexts.TransactionContexts
	transactionContexts.With(
		contexts.RootXID(rootXid),
		contexts.ParentXID(spanContexts.TraceID),
		contexts.BranchXID(rootXid),
		contexts.TransactionAgentAddress(serverAddress),
		contexts.TransactionAgentAddressOld(serverAddressOld),
	)
	return serverAddress, nil
}

func (p *TransactionProxy) branchJoin(handlerContexts *contexts.HandlerContexts, paramData []byte, headers map[string]string) (serverAddress string, err error) {
	transactionContexts := handlerContexts.TransactionContexts
	serverAddress = transactionContexts.TransactionAgentAddress
	branchXid, err := p.txnManager.TxnJoin(handlerContexts, serverAddress,
		transactionContexts.RootXID, transactionContexts.ParentXID,
		p.txInvocation.Compensable.CompensableFlagSet, paramData, p.txInvocation.Compensable.ServiceName, headers)
	if err != nil {
		return serverAddress, err
	}
	transactionContexts.With(contexts.BranchXID(branchXid))
	return serverAddress, nil
}

func (p *TransactionProxy) doEnd(isOk bool, serverAddress string, handlerContexts *contexts.HandlerContexts, tryReturnError *errors.Error, rootKVToSecondStageHeaders map[string]string) (int, error) {
	transactionContexts := handlerContexts.TransactionContexts
	errCode, err := p.txnManager.TxnEnd(handlerContexts, serverAddress, transactionContexts.RootXID, transactionContexts.ParentXID, transactionContexts.BranchXID, isOk, tryReturnError, rootKVToSecondStageHeaders)
	return errCode, err
}

// isRoot Judge the txn is root txn or not
func isRoot(transactionContexts *contexts.TransactionContexts) bool {
	return (nil == transactionContexts) || ("" == transactionContexts.RootXID && "" == transactionContexts.BranchXID && "" == transactionContexts.TransactionAgentAddress)
}

func (p *TransactionProxy) checkContext(handlerContexts *contexts.HandlerContexts) error {
	if handlerContexts == nil {
		return fmt.Errorf("handlerContexts is nil, Context:[%++v]", p.ctx)
	}
	if handlerContexts.SpanContexts == nil {
		return fmt.Errorf("handlerContexts.SpanContexts is nil, handlerContexts:[%++v]", handlerContexts)
	}
	spanContexts := handlerContexts.SpanContexts
	if spanContexts.TraceID == "" {
		return fmt.Errorf("SpanContexts.TraceID is nil, handlerContexts:[%++v]", handlerContexts)
	}
	if spanContexts.SpanID == "" {
		return fmt.Errorf("SpanContexts.SpanID is nil, handlerContexts:[%++v]", handlerContexts)
	}
	if spanContexts.ParentSpanID == "" {
		return fmt.Errorf("SpanContexts.ParentSpanID is nil, handlerContexts:[%++v]", handlerContexts)
	}
	if handlerContexts.TransactionContexts == nil {
		return fmt.Errorf("handlerContexts.TransactionContexts is nil, handlerContexts:[%++v]", handlerContexts)
	}
	transactionContexts := handlerContexts.TransactionContexts
	if transactionContexts.RootXID == "" {
		return fmt.Errorf("TransactionContexts.RootXID is nil, handlerContexts:[%++v]", handlerContexts)
	}
	if transactionContexts.BranchXID == "" {
		return fmt.Errorf("TransactionContexts.BranchXID is nil, handlerContexts:[%++v]", handlerContexts)
	}
	if transactionContexts.ParentXID == "" {
		return fmt.Errorf("TransactionContexts.ParentXID is nil, handlerContexts:[%++v]", handlerContexts)
	}
	return nil
}

func valueIsNil(value reflect.Value) bool {
	kind := value.Kind()
	if reflect.Interface == kind {
		switch value.Interface().(type) {
		case *errors.Error:
			{
				return value.Interface().(*errors.Error) == nil
			}
		default:
			// DO NOTHING
		}
	}

	return value.IsNil()
}

func getRootKVToSecondStageHeaders(responseHeader map[string]string) map[string]string {
	mp := make(map[string]string)
	for k, v := range responseHeader {
		if strings.HasPrefix(k, constant.RootKVToSecondStageKeyPrefix) {
			mp[k] = v
		}
	}
	// remove sensitive to SecondStage KV
	for k := range mp {
		delete(responseHeader, k)
	}
	return mp
}
