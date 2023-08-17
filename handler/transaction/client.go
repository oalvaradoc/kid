package client

import (
	"context"
	handlerClient "git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/client/mesh"
	"git.multiverse.io/eventkit/kit/client/mesh/wrapper/apm"
	"git.multiverse.io/eventkit/kit/client/mesh/wrapper/logging"
	"git.multiverse.io/eventkit/kit/client/mesh/wrapper/trace"
	"git.multiverse.io/eventkit/kit/compensable"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/base"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/handler/remote"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"sync"
)

//// ErrCallbackConfirm some branch transactions confirm failed
//var ErrCallbackConfirm = errors.New(constant.TransactionEndCallbackConfirmError, "some branch transactions confirm failed")
//
//// ErrCallbackCancel some branch transactions cancel failed
//var ErrCallbackCancel = errors.New(constant.TransactionEndCallbackCancelError, "some branch transactions cancel failed")
//
//// ErrDoEndConfirmTimeOut do transaction end[confirm] timeout
//var ErrDoEndConfirmTimeOut = errors.New(constant.TransactionEndCallbackConfirmTimeout, "do transaction end[confirm] timeout")
//
//// ErrDoEndCancelTimeOut do transaction end[cancel] timeout
//var ErrDoEndCancelTimeOut = errors.New(constant.TransactionEndCallbackCancelTimeout, "do transaction end[cancel] timeout")

var (
	defaultWrapperOption = handlerClient.DefaultWrapperCall(
		&apm.Wrapper{},
		&logging.Wrapper{},
		&trace.Wrapper{})

	defaultClient     handlerClient.Client = nil
	defaultClientOnce sync.Once
)

// TxInvocation defines the invocation of transaction
type TxInvocation struct {
	Compensable  *compensable.Compensable
	InstanceType reflect.Type
	MethodParams []reflect.Type
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// InvokeTryMethod proxy execution try method
func (t *TxInvocation) InvokeTryMethod(ctx context.Context, tryInstance reflect.Value, params ...interface{}) []reflect.Value {
	pValues := make([]reflect.Value, len(params))
	for idx, param := range params {
		pValues[idx] = reflect.ValueOf(param)
	}
	targetMethod := tryInstance.MethodByName(t.Compensable.TryMethod)
	return targetMethod.Call(pValues)
}

// InvokeConfirmMethod proxy execution confirm method
func (t *TxInvocation) InvokeConfirmMethod(ctx context.Context, remoteCall remote.CallInc, params []string, headers map[string]string, topicAttributes map[string]string) ([]reflect.Value, error) {
	var instance reflect.Value
	if t.InstanceType.Kind() == reflect.Ptr {
		instance = reflect.New(t.InstanceType.Elem())
	} else {
		instance = reflect.New(t.InstanceType)
	}

	ins := instance.Interface().(base.HandlerInterface)
	handlerContexts := contexts.HandlerContextsFromContext(ctx)
	if nil != handlerContexts {
		ins.SetLang(handlerContexts.Lang)
	}
	defaultClientOnce.Do(func() {
		defaultClient = mesh.NewMeshClient(defaultWrapperOption)
	})

	ins.SetRemoteCall(remoteCall)
	ins.SetContext(ctx)
	ins.SetRequestHeader(headers)
	ins.SetTopicAttributes(topicAttributes)
	ins.SetServiceConfig(&config.GetConfigs().Service)
	//instance.Elem().FieldByName(constant.DXC_CTX_FIELD_NAME).Set(reflect.ValueOf(ctx))
	pValues := make([]reflect.Value, len(t.MethodParams))
	for i := 0; i < len(t.MethodParams); i++ {
		var pValue reflect.Value
		pType := t.MethodParams[i]
		if pType.Kind() == reflect.Ptr {
			pValue = reflect.New(pType.Elem())
		} else {
			pValue = reflect.New(pType)
		}
		if err := json.Unmarshal([]byte(params[i]), pValue.Interface()); err != nil {
			return nil, err
		}
		if pType.Kind() == reflect.Ptr {
			pValues[i] = pValue
		} else {
			pValues[i] = pValue.Elem()
		}
	}
	targetMethod := instance.MethodByName(t.Compensable.ConfirmMethod)
	return targetMethod.Call(pValues), nil
}

// InvokeCancelMethod proxy execution cancel method
func (t *TxInvocation) InvokeCancelMethod(ctx context.Context, remoteCall remote.CallInc, params []string, headers map[string]string, topicAttributes map[string]string) ([]reflect.Value, error) {
	var instance reflect.Value
	if t.InstanceType.Kind() == reflect.Ptr {
		instance = reflect.New(t.InstanceType.Elem())
	} else {
		instance = reflect.New(t.InstanceType)
	}
	ins := instance.Interface().(base.HandlerInterface)
	handlerContexts := contexts.HandlerContextsFromContext(ctx)
	if nil != handlerContexts {
		ins.SetLang(handlerContexts.Lang)
	}

	defaultClientOnce.Do(func() {
		defaultClient = mesh.NewMeshClient(defaultWrapperOption)
	})

	ins.SetRemoteCall(remoteCall)
	ins.SetContext(ctx)
	ins.SetRequestHeader(headers)
	ins.SetTopicAttributes(topicAttributes)
	ins.SetServiceConfig(&config.GetConfigs().Service)
	pValues := make([]reflect.Value, len(t.MethodParams))
	for i := 0; i < len(t.MethodParams); i++ {
		var pValue reflect.Value
		pType := t.MethodParams[i]
		if pType.Kind() == reflect.Ptr {
			pValue = reflect.New(pType.Elem())
		} else {
			pValue = reflect.New(pType)
		}
		if err := json.Unmarshal([]byte(params[i]), pValue.Interface()); err != nil {
			return nil, err
		}
		if pType.Kind() == reflect.Ptr {
			pValues[i] = pValue
		} else {
			pValues[i] = pValue.Elem()
		}
	}
	targetMethod := instance.MethodByName(t.Compensable.CancelMethod)
	return targetMethod.Call(pValues), nil
}
