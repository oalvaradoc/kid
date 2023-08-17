package callback

import (
	"context"
	"encoding/base64"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/remote"
	"git.multiverse.io/eventkit/kit/handler/transaction/register"
	"git.multiverse.io/eventkit/kit/log"
)

// NewTxnCallback creates new a TxnCallback which contains transaction's repository
func NewTxnCallback() TxnCallback {
	return &DefaultLocalTxnCallback{}
}

// TxnCallback  the interface of TxnCallback, contains two call-back methods
type TxnCallback interface {
	Confirm(ctx context.Context, remoteCall remote.CallInc, serviceName string, paramData []byte, headers map[string]string, topicAttributes map[string]string) (int, error)
	Cancel(ctx context.Context, remoteCall remote.CallInc, serviceName string, paramData []byte, headers map[string]string, topicAttributes map[string]string) (int, error)
}

// DefaultLocalTxnCallback  a struct which has interface of repository.CmpTxnRep
type DefaultLocalTxnCallback struct {
}

// Confirm call back service's confirm method
//
// rootXid the id of root transaction
// branchXid the id of branch transaction
// errorCode the result code of confirm method
// err error
func (d *DefaultLocalTxnCallback) Confirm(ctx context.Context, remoteCall remote.CallInc, serviceName string, paramData []byte, headers map[string]string, topicAttributes map[string]string) (int, error) {
	log.Debugf(ctx, "start confirm local transaction serviceName:[%s]", serviceName)
	txInvocation := register.GetCompensableService(serviceName)
	if txInvocation == nil {
		err := errors.Errorf(constant.SystemInternalError, "Confirm|Cannot found the transaction with serviceName:[%s], context:[%++v]", serviceName, ctx)
		return constant.TxnEndFailedBranchConfirmFailed, err
	}
	params, err := util.DeSerialParams(paramData)
	if err != nil {
		err = errors.Errorf(constant.SystemInternalError, "Confirm|the transaction with serviceName:[%s] Deserialization fail, err:[%++v], parameter data:[%++v], context:[%++v]", serviceName, err, base64.StdEncoding.EncodeToString(paramData), ctx)
		return constant.TxnEndFailedBranchConfirmFailed, err
	}

	targetMethodResult, err := txInvocation.InvokeConfirmMethod(ctx, remoteCall, params, headers, topicAttributes)
	if err != nil {
		err = errors.Errorf(constant.SystemInternalError, "Confirm|the transaction with serviceName:[%s] Invoke Confirm Method fail, err:[%++v], context:[%++v]", serviceName, err, ctx)
		return constant.TxnEndFailedBranchConfirmFailed, err
	}

	lastResultObj := targetMethodResult[len(targetMethodResult)-1]
	if !lastResultObj.IsNil() {
		err = lastResultObj.Interface().(error)
		err = errors.Errorf(constant.SystemInternalError, "Confirm|local transaction serviceName[%s] failed, target method execute failed, err:[%++v], context:[%++v]", serviceName, err, ctx)
		return constant.TxnEndFailedBranchConfirmFailed, err
	}
	log.Debugf(ctx, "confirm local transaction serviceName[%s] successfully!", serviceName)
	return 0, nil
}

// Cancel call back service's cancel method
//
// rootXid the id of root transaction
// branchXid the id of branch transaction
// errorCode the result code of cancel method
// err error
func (d *DefaultLocalTxnCallback) Cancel(ctx context.Context, remoteCall remote.CallInc, serviceName string, paramData []byte, headers map[string]string, topicAttributes map[string]string) (int, error) {
	log.Debugf(ctx,"start cancel local transaction serviceName:[%s]", serviceName)
	txInvocation := register.GetCompensableService(serviceName)
	if txInvocation == nil {
		err := errors.Errorf(constant.SystemInternalError, "Cancel|Cannot found the transaction with serviceName:[%s], context:[%++v]", serviceName, ctx)
		return constant.TxnEndFailedBranchCancelFailed, err
	}
	params, err := util.DeSerialParams(paramData)
	if err != nil {
		err = errors.Errorf(constant.SystemInternalError, "Cancel|the transaction with serviceName:[%s] Deserialization fail, err:[%++v], parameter data:[%++v], context:[%++v]", serviceName, err, base64.StdEncoding.EncodeToString(paramData), ctx)
		return constant.TxnEndFailedBranchCancelFailed, err
	}
	targetMethodResult, err := txInvocation.InvokeCancelMethod(ctx, remoteCall, params, headers, topicAttributes)
	if err != nil {
		err = errors.Errorf(constant.SystemInternalError, "Cancel|the transaction with serviceName:[%s] Invoke Cancel Method fail, err:[%++v], context:[%++v]", serviceName, err, ctx)
		return constant.TxnEndFailedBranchCancelFailed, err
	}

	lastResultObj := targetMethodResult[len(targetMethodResult)-1]
	if !lastResultObj.IsNil() {
		err = lastResultObj.Interface().(error)
		err = errors.Errorf(constant.SystemInternalError, "Cancel|local transaction serviceName[%s] failed, target method execute failed, err:[%++v], context:[%++v]", serviceName, err, ctx)
		return constant.TxnEndFailedBranchCancelFailed, err
	}
	log.Debugf(ctx, "cancel local transaction serviceName[%s] successfully!", serviceName)
	return 0, nil
}
