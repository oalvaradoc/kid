package router

import (
	"fmt"
	"git.multiverse.io/eventkit/kit/codec/auto"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/handler/transaction/register"
	"git.multiverse.io/eventkit/kit/interceptor"
	"git.multiverse.io/eventkit/kit/interceptor/apm"
	"git.multiverse.io/eventkit/kit/interceptor/logging"
	"git.multiverse.io/eventkit/kit/interceptor/server_response"
	"git.multiverse.io/eventkit/kit/interceptor/transaction"
	"git.multiverse.io/eventkit/kit/validation"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/log"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

// define the callback executor to match and schedule the handler
// according to the eventID and router type after receiving the message
const (
	EventID = iota
	EventPrefix
	EventSuffix
	EventExpression
)

var registerTypeNameMapping = map[int]string{
	EventID:         "Event ID",
	EventPrefix:     "Prefix of Event ID",
	EventSuffix:     "Suffix of Event ID",
	EventExpression: "Matcher expression of Event ID",
}

// HandlerRouter is a router for handler that defines the route information of handler.
type HandlerRouter struct {
	EventKeyMap map[string]string
	sync.RWMutex
	IsExistsCompensableTransaction bool
	URLPathHandlers                map[string]*Options
	DefiniteEventHandlers          map[string]*Options
	ExpressionEventHandlers        map[string]*Options
	//DefaultInterceptors            []interceptor.Interceptor
	DefaultInterceptorsOption                      Option
	DefaultEnableValidationOption                  Option
	DefaultCustomValidationRegisterFunctionsOption Option
}

var defaultCodec = auto.BuildAutoCodecWithJSONCodec()
var defaultInterceptors = []interceptor.Interceptor{
	&logging.Interceptor{},
	&server_response.Interceptor{},
	&apm.Interceptor{},
	&transaction.Interceptor{},
}

// IsMatchedEventID is used to check whether the eventID is matched in the specified expression.
func (h *HandlerOptions) IsMatchedEventID(eventID string, matcherExpression string) bool {
	switch h.RegisterType {
	case EventPrefix:
		{
			if strings.HasPrefix(eventID, matcherExpression) {
				return true
			}
		}
	case EventSuffix:
		{
			if strings.HasSuffix(eventID, matcherExpression) {
				return true
			}
		}
	case EventExpression:
		{
			if isOk, _ := regexp.MatchString(matcherExpression, eventID); isOk {
				return true
			}
		}
	default:
		return eventID == matcherExpression
	}

	return false
}

// SetEventKeyMap sets the event key map into router.
func (e *HandlerRouter) SetEventKeyMap(eventKeyMap map[string]string) {
	e.EventKeyMap = eventKeyMap
}

// SetDefaultInterceptors sets the default interceptors into router.
func (e *HandlerRouter) SetDefaultInterceptors(defaultInterceptors []interceptor.Interceptor) {
	e.DefaultInterceptorsOption = WithInterceptors(defaultInterceptors...)
}

// DefaultEnableValidation sets the default enable validation into router.
func (e *HandlerRouter) DefaultEnableValidation(combineErrors bool, defaultCustomValidationRegisterFunctions []validation.CustomValidationRegisterFunc) {
	validationOptions := make([]CustomValidationOption, 0)
	for _, f := range defaultCustomValidationRegisterFunctions {
		validationOptions = append(validationOptions, AddRegisterValidationFunc(f.Tag, f.Func, f.CallValidationEvenIfNull))
	}
	e.DefaultEnableValidationOption = EnableValidation(combineErrors, validationOptions...)
}

func (e *HandlerRouter) generateHandlerProperties(instance interface{}, registerOptions *Options) *Options {
	handlerReflectType := reflect.TypeOf(instance)
	handlerName := ""
	if _, ok := handlerReflectType.MethodByName(constant.FunctionForGetHandlerName); ok {
		handlerIns := reflect.ValueOf(instance)
		targetMethod := handlerIns.MethodByName(constant.FunctionForGetHandlerName)
		targetMethodResult := targetMethod.Call(nil)
		handlerName = targetMethodResult[0].String()
	} else {
		if handlerReflectType.Kind() == reflect.Ptr {
			handlerName = handlerReflectType.Elem().Name()
		} //else {
		//handlerName = handlerReflectType.Name()
		//}
	}
	registerOptions.HandlerOptions.HandlerReflectType = handlerReflectType

	if "" == registerOptions.EventExpression {
		panic(errors.Errorf(constant.SystemInternalError, "%s is empty,please check, handler[%v] - Method[%s]",
			registerTypeNameMapping[registerOptions.HandlerOptions.RegisterType],
			handlerName, registerOptions.HandlerOptions.HandlerMethodName))
	}
	var combineErrors bool
	if nil != registerOptions.HandlerOptions.CustomValidationOptions {
		combineErrors = registerOptions.HandlerOptions.CustomValidationOptions.CombineErrors
	}

	log.Infosf("Event handler register: %s[%s] - Handler[%v] - Method[%s] - Interceptors:%++v - Enabled validation:[%v - (combine errors:%v)]",
		registerTypeNameMapping[registerOptions.HandlerOptions.RegisterType],
		registerOptions.EventExpression,
		handlerName,
		registerOptions.HandlerOptions.HandlerMethodName,
		registerOptions.HandlerOptions.Interceptors,
		registerOptions.HandlerOptions.EnableValidation,
		combineErrors,
	)
	if _, ok := handlerReflectType.MethodByName(registerOptions.HandlerOptions.HandlerMethodName); !ok {
		panic(fmt.Sprintf("failed to register: %s[%s] - handler[%v], cannot found method name:[%s]",
			registerTypeNameMapping[registerOptions.HandlerOptions.RegisterType],
			registerOptions.EventExpression, handlerName, registerOptions.HandlerOptions.HandlerMethodName))
	}

	if _, ok := handlerReflectType.MethodByName("PreHandle"); ok {
		registerOptions.HandlerOptions.InvokePreHandle = true
	}

	method := reflect.ValueOf(instance).MethodByName(registerOptions.HandlerOptions.HandlerMethodName)
	if method.Type().NumIn() > 1 {
		panic(fmt.Sprintf("failed to register: %s[%s] - handler[%v] - method[%s], number of in parameters is %d, cannot be great than 1. ",
			registerTypeNameMapping[registerOptions.HandlerOptions.RegisterType],
			registerOptions.EventExpression, handlerName, registerOptions.HandlerOptions.HandlerMethodName, method.Type().NumIn()))
	}
	handlerMethodInParams := make([]reflect.Type, method.Type().NumIn())
	for i := 0; i < method.Type().NumIn(); i++ {
		handlerMethodInParams[i] = method.Type().In(i)
	}
	registerOptions.HandlerOptions.HandlerMethodInParams = handlerMethodInParams

	handlerMethodOutParams := make([]reflect.Type, method.Type().NumOut())
	for i := 0; i < method.Type().NumOut(); i++ {
		handlerMethodOutParams[i] = method.Type().Out(i)
	}
	registerOptions.HandlerOptions.HandlerMethodOutParams = handlerMethodOutParams
	// register transaction compensable
	if nil != registerOptions.Compensable {
		log.Infosf("Register compensable service to transaction manager, %++v", registerOptions.Compensable)
		if err := register.CompensableService(instance, *registerOptions.Compensable); nil != err {
			panic(fmt.Sprintf("failed to register compensable to transaction manager, error=%++v", err))
		}
	}

	// register URL path with handler
	if "" != registerOptions.HandlerOptions.URLPath {
		if nil == e.URLPathHandlers {
			e.URLPathHandlers = make(map[string]*Options)
		}
		e.URLPathHandlers[registerOptions.HandlerOptions.URLPath] = registerOptions
	}
	return registerOptions
}

func newRegisterOptions(opts ...Option) *Options {
	tmpInterceptors := make([]interceptor.Interceptor, len(defaultInterceptors))
	copy(tmpInterceptors, defaultInterceptors)
	registerOptions := &Options{
		Compensable: nil,
		HandlerOptions: HandlerOptions{
			Codec:        defaultCodec,
			Interceptors: tmpInterceptors,
		},
	}
	for _, registerOption := range opts {
		registerOption(registerOptions)
	}

	return registerOptions
}

// RouterWithEventKey sets the event key and handler instance into router with one or more parameters.
// will finds the real event ID in config by eventKey.
func (e *HandlerRouter) RouterWithEventKey(eventKey string, instance interface{}, opts ...Option) *Options {
	eventID := config.GetEventID(e.EventKeyMap, eventKey)
	return e.Router(eventID, instance, opts...)
}

// Router sets the event ID and handler instance into router with one or more parameters.
func (e *HandlerRouter) Router(eventID string, instance interface{}, opts ...Option) *Options {
	e.Lock()
	defer e.Unlock()
	if reflect.TypeOf(instance).Kind() != reflect.Ptr {
		panic(errors.Errorf(constant.SystemInternalError, "Router eventID[%s] failed, instance must be a pointer", eventID))
	}
	var finalOpts = make([]Option, 0)
	if nil != e.DefaultInterceptorsOption {
		optsOut := append(finalOpts, e.DefaultInterceptorsOption)
		finalOpts = optsOut
	}

	if nil != e.DefaultEnableValidationOption {
		optsOut := append(finalOpts, e.DefaultEnableValidationOption)
		finalOpts = optsOut
	}

	optsOut := append(finalOpts, opts...)
	finalOpts = optsOut
	registerOptions := newRegisterOptions(finalOpts...)

	if nil != registerOptions.Compensable {
		e.IsExistsCompensableTransaction = true
	}
	registerOptions.EventExpression = eventID
	registerOptions.HandlerOptions.RegisterType = EventID

	e.generateHandlerProperties(instance, registerOptions)
	if nil == e.DefiniteEventHandlers {
		e.DefiniteEventHandlers = make(map[string]*Options)
	}
	if _, ok := e.DefiniteEventHandlers[eventID]; ok {
		panic(errors.Errorf(constant.SystemInternalError, "Duplicate event ID:%s", eventID))
	}
	e.DefiniteEventHandlers[eventID] = registerOptions

	return registerOptions
}

// RouterPrefix sets the prefix of event ID expression and handler instance into router with one or more parameters.
func (e *HandlerRouter) RouterPrefix(prefixOfEventID string, instance interface{}, opts ...Option) *Options {
	e.Lock()
	defer e.Unlock()

	if reflect.TypeOf(instance).Kind() != reflect.Ptr {
		panic(errors.Errorf(constant.SystemInternalError, "Router prefixOfEventID[%s] failed, instance must be a pointer", prefixOfEventID))
	}

	var finalOpts = make([]Option, 0)
	if nil != e.DefaultInterceptorsOption {
		optsOut := append(finalOpts, e.DefaultInterceptorsOption)
		finalOpts = optsOut
	}

	if nil != e.DefaultEnableValidationOption {
		optsOut := append(finalOpts, e.DefaultEnableValidationOption)
		finalOpts = optsOut
	}

	optsOut := append(finalOpts, opts...)
	finalOpts = optsOut
	registerOptions := newRegisterOptions(finalOpts...)

	if nil != registerOptions.Compensable {
		e.IsExistsCompensableTransaction = true
	}
	registerOptions.EventExpression = prefixOfEventID
	registerOptions.HandlerOptions.RegisterType = EventPrefix

	e.generateHandlerProperties(instance, registerOptions)
	if nil == e.ExpressionEventHandlers {
		e.ExpressionEventHandlers = make(map[string]*Options)
	}
	e.ExpressionEventHandlers[prefixOfEventID] = registerOptions

	return registerOptions
}

// RouterSuffix sets the suffix of event ID expression and handler instance into router with one or more parameters.
func (e *HandlerRouter) RouterSuffix(suffixOfEventID string, instance interface{}, opts ...Option) *Options {
	e.Lock()
	defer e.Unlock()

	if reflect.TypeOf(instance).Kind() != reflect.Ptr {
		panic(errors.Errorf(constant.SystemInternalError, "Router RouterSuffix[%s] failed, instance must be a pointer", suffixOfEventID))
	}

	var finalOpts = make([]Option, 0)
	if nil != e.DefaultInterceptorsOption {
		optsOut := append(finalOpts, e.DefaultInterceptorsOption)
		finalOpts = optsOut
	}

	if nil != e.DefaultEnableValidationOption {
		optsOut := append(finalOpts, e.DefaultEnableValidationOption)
		finalOpts = optsOut
	}

	optsOut := append(finalOpts, opts...)
	finalOpts = optsOut
	registerOptions := newRegisterOptions(finalOpts...)

	if nil != registerOptions.Compensable {
		e.IsExistsCompensableTransaction = true
	}
	registerOptions.EventExpression = suffixOfEventID
	registerOptions.HandlerOptions.RegisterType = EventSuffix

	e.generateHandlerProperties(instance, registerOptions)
	if nil == e.ExpressionEventHandlers {
		e.ExpressionEventHandlers = make(map[string]*Options)
	}
	e.ExpressionEventHandlers[suffixOfEventID] = registerOptions

	return registerOptions
}

// RouterExpression sets the expression and handler instance into router with one or more parameters.
func (e *HandlerRouter) RouterExpression(matcherExpressionOfEventID string, instance interface{}, opts ...Option) *Options {
	e.Lock()
	defer e.Unlock()

	if reflect.TypeOf(instance).Kind() != reflect.Ptr {
		panic(errors.Errorf(constant.SystemInternalError, "Router matcherExpressionOfEventID[%s] failed, instance must be a pointer", matcherExpressionOfEventID))
	}

	var finalOpts = make([]Option, 0)
	if nil != e.DefaultInterceptorsOption {
		optsOut := append(finalOpts, e.DefaultInterceptorsOption)
		finalOpts = optsOut
	}

	if nil != e.DefaultEnableValidationOption {
		optsOut := append(finalOpts, e.DefaultEnableValidationOption)
		finalOpts = optsOut
	}

	optsOut := append(finalOpts, opts...)
	finalOpts = optsOut
	registerOptions := newRegisterOptions(finalOpts...)

	if nil != registerOptions.Compensable {
		e.IsExistsCompensableTransaction = true
	}
	registerOptions.EventExpression = matcherExpressionOfEventID
	registerOptions.HandlerOptions.RegisterType = EventExpression

	e.generateHandlerProperties(instance, registerOptions)
	if nil == e.ExpressionEventHandlers {
		e.ExpressionEventHandlers = make(map[string]*Options)
	}
	e.ExpressionEventHandlers[matcherExpressionOfEventID] = registerOptions

	return registerOptions
}

// MatchHandler finds the handler by eventID, return nil if the eventID cannot match any router.
func (e *HandlerRouter) MatchHandler(eventID string) *Options {
	e.RLock()
	defer e.RUnlock()

	// 1. DefiniteEventHandlers
	if hp, ok := e.DefiniteEventHandlers[eventID]; ok {
		return hp
	}

	// 2. ExpressionEventHandlers
	for k, v := range e.ExpressionEventHandlers {
		if v.HandlerOptions.IsMatchedEventID(eventID, k) {
			return v
		}
	}

	return nil
}

// MatchHandlerWithURLPath finds the handler by URL path, return nil if the eventID cannot match any router.
func (e *HandlerRouter) MatchHandlerWithURLPath(urlPath string) *Options {
	e.RLock()
	defer e.RUnlock()

	return e.URLPathHandlers[urlPath]
}
