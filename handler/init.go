package handler

import (
	"git.multiverse.io/eventkit/kit"
	apmWrapper "git.multiverse.io/eventkit/kit/apm"
	"git.multiverse.io/eventkit/kit/cache/v1/repository"
	"git.multiverse.io/eventkit/kit/common/apm"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/handler/router"
	"git.multiverse.io/eventkit/kit/interceptor"
	"git.multiverse.io/eventkit/kit/log"
	"git.multiverse.io/eventkit/kit/sed/callback"
	"reflect"
	"strings"
)

// NewCallbackExecutor is used to create a SED callback executor. When creating it,
// some necessary default values will be set. For example:
//   * the address of SED Server defaults to `http://127.0.0.1:18080`
//   * the port of SED callback defaults to `18082`
//   * the communication type is `mesh`
//   * The default port for opening http endpoint is `6060`
//   * the state machine of the client side is closed by default
//  You can override the default value by setting one or more method parameters
func NewCallbackExecutor(options ...callback.Option) *EventCallback {
	opts := callback.NewHandlerOptions()
	for _, o := range options {
		o(&opts)
	}

	return &EventCallback{opts: opts}
}

// NewCustomOptions creates an empty CustomOptions
// CustomOptions can be used to set some executor startup or execution parameters
func NewCustomOptions() CustomOptions {
	opts := CustomOptions{
		CustomOptionConfigs: make(map[string]interface{}),
	}

	return opts
}

// StartWithCustomOptions function is similar to Start, and different from the Start method is that this method can be modified at startup
// The default client.Client, the default interceptors, the default call interceptors, the default response template,
// whether to enable validation by default, etc.
//
// example:
// var (
//	 myWrapperOption = client.DefaultWrapperCall(
//		&apm.Wrapper{},
//		&logging.Wrapper{},
//		&trace.Wrapper{})
//
//	 myClient = mesh.NewMeshClient(myWrapperOption)
//
//	 defaultInterceptors = []interceptor.Interceptor{
//		&apmInterceptor.Interceptor{},
//		&loggingInterceptor.Interceptor{},
//		&server_response.Interceptor{},
//		&traceInterceptor.Interceptor{},
//		&transaction.Interceptor{},
//	}
//
//	defaultCallInterceptors = []interceptor.Interceptor{
//		&client_receive.Interceptor{},
//	}
// )
//
// type callbackHandleWrapper struct{}
//
// func (c *callbackHandleWrapper) PreHandle(inCtx context.Context, in *msg.Message) (skip bool, outCtx context.Context, out *msg.Message, err error) {
//	 log.Infos("Execute callback pre handle...")
//
//	  return true, inCtx, &msg.Message{
//		  Body: []byte("this is a test body from PreHandle"),
//	  }, nil
// }
//
// func (c *callbackHandleWrapper) PostHandle(ctx context.Context, m *msg.Message) (out *msg.Message, err error) {
//	 log.Infos("Execute callback post handle...")
//
//	 return &msg.Message{
//		 Body: []byte("this is a test body from PostHandle"),
//	 }, nil
// }
//
//
// 	if err := handler.StartWithCustomOptions(
//		InitHandlerRouters,
//		handler.WithCallbackOptions(configs.GenCallbackOptions()...),
//      // override the default client
//		handler.WithClient(myClient),
//      // override the default interceptors
//		handler.WithDefaultInterceptors(defaultInterceptors...),
//      // override the default call interceptors
//		handler.WithCallInterceptors(defaultCallInterceptors...),
//      // override the default response template
//		handler.WithDefaultResponseTemplate("{\"errorCode\":\"{{.errorCode}}\", \"errorMsg\":\"{{.errorMsg}}\", \"response\":{{.data}}}"),
//      // setting callback handle wrapper
//		handler.WithCallbackHandleWrapper(&callbackHandleWrapper{}),
//      // override the default disabled validation
//		handler.WithDefaultEnableValidation(true,
//			handler.AddDefaultRegisterValidationFunc("is-awesome",ValidateMyVal),
//			handler.AddDefaultRegisterValidationFunc("is-awesome1",ValidateMyVal),
//		),
//	); nil != err {
//		return err
//	}
func StartWithCustomOptions(routerFn func(router *router.HandlerRouter) error, options ...CustomOption) error {
	opts := NewCustomOptions()
	callbackOptions := make([]callback.Option, 0)

	for _, o := range options {
		o(&opts)
	}

	if nil != opts.CustomOptionConfigs[constant.ExtConfigCustomSedClientOptions] {
		callbackOptions = opts.CustomOptionConfigs[constant.ExtConfigCustomSedClientOptions].([]callback.Option)
	}

	delete(opts.CustomOptionConfigs, constant.ExtConfigCustomSedClientOptions)
	for k, v := range opts.CustomOptionConfigs {
		callbackOptions = append(callbackOptions, callback.AddExtConfig(k, v))
	}

	return Start(routerFn, callbackOptions...)
}

// Start is used to launch a handler. In this method,
// a default executor will be created (of course, the necessary parameters of the executor will be set by default)
// and the router information will be registered with the executor. When the executor is created,
// it will be based on the configuration parameters and router function Set the corresponding parameters of the registered handler,
// such as default interceptors, default client, default wrapper, etc., through the options parameter,
// you can modify the default callback.Options related parameters, if you want to modify the handler default interceptors,
// default wrapper and other parameters , You need to use the StartWithCustomOptions method to start the handler
// example:
// 1. Start a handler that does not need to process events (no need to register the handler to the router):
// 	if err := handler.Start(nil, configs.GenCallbackOptions()...); nil != err {
//		return err
//	}
//
// 2.Start a handler that needs to process the event (if need to register the handler to the router):
//    // create a handler router function
//    func InitHandlerRouters(handlerRouter *router.HandlerRouter) error {
//			handlerRouter.Router("TRANSFER",
//				&handlers.TransferHandler{},
//				router.Method("TestTransfer"),
//				router.HandlePost("/v1/transfer"),
//				router.EnableValidation(true),
//			)
//     }
//    // configs.GenCallbackOptions() creates executor parameters based on configuration
//    if err := handler.Start(InitHandlerRouters, configs.GenCallbackOptions()...); nil != err {
//		return err
//	  }
//
// 3.Modify the sed server address to http://127.0.0.1:8080 and the callback port to 12345 at startup
//  if err := handler.Start(routers.InitHandlerRouters,
//		append(append(configs.GenCallbackOptions(),
//			callback.WithCallbackPort(12345)),
//			callback.WithServerAddress("http://127.0.0.1:8080"),
//			)...); nil != err {
//		return err
//	}
func Start(routerFn func(router *router.HandlerRouter) error, options ...callback.Option) error {
	callback.SetKitVersion(kit.Version)
	log.Infosf("**********************version***************************")
	log.Infosf("*                           eventkit/kit     => %s", kit.Version)
	log.Infosf("********************************************************")
	callbackExecutor := NewCallbackExecutor(options...)

	if nil != routerFn {
		router := &router.HandlerRouter{
			URLPathHandlers:         make(map[string]*router.Options),
			DefiniteEventHandlers:   make(map[string]*router.Options),
			ExpressionEventHandlers: make(map[string]*router.Options),
		}
		if !util.IsNil(callbackExecutor.opts.ExtConfigs) {
			if subscribedTopics, ok := callbackExecutor.opts.ExtConfigs[constant.ExtEventKeyMap]; ok {
				router.SetEventKeyMap(subscribedTopics.(map[string]string))
			}
			if defaultInterceptors, ok := callbackExecutor.opts.ExtConfigs[constant.ExtConfigCustomDefaultInterceptors]; ok {
				router.SetDefaultInterceptors(defaultInterceptors.([]interceptor.Interceptor))
			}

			if opts, ok := callbackExecutor.opts.ExtConfigs[constant.ExtConfigDefaultEnableValidation]; ok {
				defaultCustomValidationOptions := opts.(*DefaultCustomValidationOptions)
				router.DefaultEnableValidation(defaultCustomValidationOptions.CombineErrors, defaultCustomValidationOptions.DefaultCustomValidationRegisterFunctions)

			}
		}

		err := routerFn(router)
		if nil != err {
			return err
		}

		callbackExecutor.SetRouter(router)
	}

	// register configuration on change hook function.
	config.RegisterConfigOnChangeHookFunc("Log", rotateLogWhenConfigChanged, false)

	// register error code mapping
	config.RegisterConfigOnChangeHookFunc("ErrorCodeMapping", rotateErrorCodeMappingWhenConfigChanged, false)

	// register cache operator if necessary
	callback.RegisterInitHookFunc("addressing", repository.InitCacheOperatorIfNecessary, true)

	callback.RegisterCallbackExecutor(callbackExecutor)

	// init APM logger if necessary
	configs := config.GetConfigs()
	if nil != configs && configs.Apm.Enable {
		apmLoggerFilePath := configs.Apm.RootPath
		if !strings.HasSuffix(apmLoggerFilePath, "/") {
			apmLoggerFilePath += "/"
		}
		apmLoggerFilePath += configs.Service.InstanceID + "/"
		apmLoggerFilePath += "apm.log"

		apmLog, err := apm.NewLogger(apmLoggerFilePath, configs.Apm.FileRows)
		if err != nil {
			return err
		}

		apmWrapper.ApmLogger = apmLog
	}

	return nil
}

func rotateLogWhenConfigChanged(oldConfig *config.ServiceConfigs, newConfig *config.ServiceConfigs) error {
	if nil == newConfig {
		return nil
	}
	currentLogConfig := log.CurrentConfig()
	newLogConfig := newConfig.GenLogConfig()
	//if currentLogConfig.Level != newLogConfig.Level ||
	//	currentLogConfig.Console != newLogConfig.Console ||
	//	currentLogConfig.Rotate.MaxBackups != newLogConfig.Rotate.MaxBackups ||
	//	currentLogConfig.Rotate.MaxSize != newLogConfig.Rotate.MaxSize ||
	//	currentLogConfig.Rotate.Filename != newLogConfig.Rotate.Filename ||
	//	currentLogConfig.Rotate.FilePath != newLogConfig.Rotate.FilePath ||
	//	currentLogConfig.Rotate.MaxAge != newLogConfig.Rotate.MaxAge {
	//	log.Init(newLogConfig)
	//}
	if !currentLogConfig.Equals(&newLogConfig) {
		log.Init(newLogConfig)
	}

	return nil
}


func rotateErrorCodeMappingWhenConfigChanged(oldConfig *config.ServiceConfigs, newConfig *config.ServiceConfigs) error {
	if nil == newConfig {
		errors.SetErrorCodeMapping(map[string]string{})
		return nil
	}

	if !reflect.DeepEqual(oldConfig.Service.ResponseCodeMapping, newConfig.Service.ResponseCodeMapping) {
		errors.SetErrorCodeMapping(newConfig.Service.ResponseCodeMapping)
	}

	return nil
}
