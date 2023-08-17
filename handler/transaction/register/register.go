package register

import (
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/compensable"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/transaction"
	"reflect"
	"strings"
	"sync"
)

var (
	cache  = make(map[string]*client.TxInvocation)
	locker sync.RWMutex

	// nil is a literal nil. sometype(nil) is a type conversion of nil to a nil sometype.
	// (*error)(nil) produces a nil *error from which we can take the type via TypeOf.
	// Going back to the error interface is done via Elem
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

// GetCompensableService finds the transaction invocation from transaction manager by key
func GetCompensableService(serviceName string) *client.TxInvocation {
	locker.RLock()
	defer locker.RUnlock()

	return cache[serviceName]
}

// CompensableService used to register transaction invocation to transaction manager
func CompensableService(txInstance interface{}, compensable compensable.Compensable) error {
	locker.Lock()
	defer locker.Unlock()

	if compensable.ServiceName == "" {
		tType := reflect.TypeOf(txInstance)
		if tType.Kind() == reflect.Ptr {
			compensable.ServiceName = tType.Elem().Name()
		} else {
			compensable.ServiceName = tType.Name()
		}
	}

	if err := checkTccMethod(txInstance, &compensable); nil != err {
		return err
	}
	compensable.CompensableFlagSet = genCompensableFlag(&compensable)

	method := reflect.ValueOf(txInstance).MethodByName(compensable.TryMethod)

	mTypes := make([]reflect.Type, method.Type().NumIn())
	for i := 0; i < method.Type().NumIn(); i++ {
		mTypes[i] = method.Type().In(i)
	}

	txInvocation := &client.TxInvocation{
		Compensable:  &compensable,
		InstanceType: reflect.TypeOf(txInstance),
		MethodParams: mTypes,
	}

	_, ok := cache[compensable.ServiceName]
	if ok {
		return errors.Errorf(constant.SystemInternalError, "Service `%s` already exsits, please set a different `ServiceName` in compensable.Compensable", compensable.ServiceName)
	}
	cache[compensable.ServiceName] = txInvocation

	return nil
}

func checkTccMethod(txnInstance interface{}, compensable *compensable.Compensable) error {
	isExistConfirmMethod := "" != compensable.ConfirmMethod && "" != strings.TrimSpace(compensable.ConfirmMethod)
	isExistCancelMethod := "" != compensable.CancelMethod && "" != strings.TrimSpace(compensable.CancelMethod)

	// 1. check try、confirm、cancel method whether exists
	tType := reflect.TypeOf(txnInstance)
	tValue := reflect.ValueOf(txnInstance)

	if _, ok := tType.MethodByName(compensable.TryMethod); !ok {
		return errors.Errorf(constant.SystemInternalError, "[Service:%s]Cannot found try method:[%s]", compensable.ServiceName, compensable.TryMethod)
	}

	if isExistConfirmMethod {
		if _, ok := tType.MethodByName(compensable.ConfirmMethod); !ok {
			return errors.Errorf(constant.SystemInternalError, "[Service:%s]Cannot found confirm method:[%s]", compensable.ServiceName, compensable.ConfirmMethod)
		}
	}
	if isExistCancelMethod {
		if _, ok := tType.MethodByName(compensable.CancelMethod); !ok {
			return errors.Errorf(constant.SystemInternalError, "[Service:%s]Cannot found cancel method:[%s]", compensable.ServiceName, compensable.CancelMethod)
		}
	}
	tryMethod := tValue.MethodByName(compensable.TryMethod)

	// 2.1 if number of return value is less than 1,then raise error
	if tryMethod.Type().NumOut() < 1 {
		return errors.Errorf(constant.SystemInternalError, "[Service:%s]The number of return value must be greater than 1", compensable.ServiceName)
	}

	// 2.2 if type of the last return value is not error, then raise error
	if !tryMethod.Type().Out(tryMethod.Type().NumOut() - 1).Implements(errorInterface) {
		return errors.Errorf(constant.SystemInternalError, "[Service:%s]The type of the last return value must be error", compensable.ServiceName)
	}

	// 3. check try、confirm method whether has the same parameters
	if isExistConfirmMethod {
		confirmMethod := tValue.MethodByName(compensable.ConfirmMethod)
		if tryMethod.Type().NumIn() != confirmMethod.Type().NumIn() {
			return errors.Errorf(constant.SystemInternalError, "[Service:%s]The input number of parameters of the try and confirm methods must be the same!", compensable.ServiceName)
		}

		if tryMethod.Type().NumOut() != confirmMethod.Type().NumOut() {
			return errors.Errorf(constant.SystemInternalError, "[Service:%s]The output number of parameters of the try and confirm methods must be the same!", compensable.ServiceName)
		}

		for i := 0; i < tryMethod.Type().NumIn(); i++ {
			if reflect.ValueOf(tryMethod.Type().In(i)) != reflect.ValueOf(confirmMethod.Type().In(i)) {
				return errors.Errorf(constant.SystemInternalError, "[Service:%s]The type of parameters of the try and confirm methods must be the same[NumIn]!", compensable.ServiceName)
			}
		}

		for i := 0; i < tryMethod.Type().NumOut(); i++ {
			if reflect.ValueOf(tryMethod.Type().Out(i)) != reflect.ValueOf(confirmMethod.Type().Out(i)) {
				return errors.Errorf(constant.SystemInternalError, "[Service:%s]The type of parameters of the try and confirm methods must be the same[NumOut]!", compensable.ServiceName)
			}
		}
	}

	// 3. check try、cancel method whether has the same parameters
	if isExistCancelMethod {
		cancelMethod := tValue.MethodByName(compensable.CancelMethod)
		if tryMethod.Type().NumIn() != cancelMethod.Type().NumIn() {
			return errors.Errorf(constant.SystemInternalError, "[Service:%s]The input number of parameters of the try and cancel methods must be the same!", compensable.ServiceName)
		}

		if tryMethod.Type().NumOut() != cancelMethod.Type().NumOut() {
			return errors.Errorf(constant.SystemInternalError, "[Service:%s]The output number of parameters of the try and cancel methods must be the same!", compensable.ServiceName)
		}

		for i := 0; i < tryMethod.Type().NumIn(); i++ {
			if reflect.ValueOf(tryMethod.Type().In(i)) != reflect.ValueOf(cancelMethod.Type().In(i)) {
				return errors.Errorf(constant.SystemInternalError, "[Service:%s]The type of parameters of the try and cancel methods must be the same[NumIn]!", compensable.ServiceName)
			}
		}

		for i := 0; i < tryMethod.Type().NumOut(); i++ {
			if reflect.ValueOf(tryMethod.Type().Out(i)) != reflect.ValueOf(cancelMethod.Type().Out(i)) {
				return errors.Errorf(constant.SystemInternalError, "[Service:%s]The type of parameters of the try and cancel methods must be the same[NumOut]!", compensable.ServiceName)
			}
		}
	}
	return nil
}

func genCompensableFlag(compensableTxn *compensable.Compensable) int {
	compensableFlagSet := 0

	if "" != compensableTxn.ConfirmMethod && "" != strings.TrimSpace(compensableTxn.ConfirmMethod) {
		compensableFlagSet = compensableFlagSet | constant.ConfirmFlag
	}

	if "" != compensableTxn.CancelMethod && "" != strings.TrimSpace(compensableTxn.CancelMethod) {
		compensableFlagSet = compensableFlagSet | constant.CancelFlag
	}

	return compensableFlagSet
}
