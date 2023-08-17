package handler

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/client/mesh"
	"git.multiverse.io/eventkit/kit/client/mesh/wrapper/addressing"
	"git.multiverse.io/eventkit/kit/client/mesh/wrapper/apm"
	"git.multiverse.io/eventkit/kit/client/mesh/wrapper/logging"
	"git.multiverse.io/eventkit/kit/client/mesh/wrapper/trace"
	"git.multiverse.io/eventkit/kit/codec"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/base"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/handler/remote"
	"git.multiverse.io/eventkit/kit/handler/router"
	"git.multiverse.io/eventkit/kit/handler/transaction/imports"
	"git.multiverse.io/eventkit/kit/handler/transaction/proxy"
	"git.multiverse.io/eventkit/kit/interceptor"
	"git.multiverse.io/eventkit/kit/interceptor/client_receive"
	"git.multiverse.io/eventkit/kit/log"
	"git.multiverse.io/eventkit/kit/sed/callback"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// CallbackHandleWrapper is a wrapper that contains PreHandle performed before execute the handler
// and PostHandle performed after execute the handler
type CallbackHandleWrapper interface {
	PreHandle(ctx context.Context, in *msg.Message) (skip bool, outCtx context.Context, out *msg.Message, err error)
	PostHandle(ctx context.Context, in *msg.Message) (out *msg.Message, err error)
}

// EventCallback is a default implement of callback executor
type EventCallback struct {
	opts                           callback.Options
	responseTemplate               string
	responseAutoParseKeyMapping    map[string]string
	transactionConfig              *config.Transaction
	handlerRouter                  *router.HandlerRouter
	extConfigs                     map[string]interface{}
	serviceConfig                  *config.Service
	client                         client.Client
	callbackHandleWrapper          CallbackHandleWrapper
	isEnabledCircuitBreakerMonitor bool
	defaultUserLang                string
	userLangKey                    string
	duplicateErrorCodeTo           string
	enableExecutorLogging          bool
}

var (
	defaultCallWrapperOption = client.DefaultWrapperCall(
		&apm.Wrapper{},
		&logging.Wrapper{},
		&trace.Wrapper{},
		&addressing.Wrapper{},
	)
	defaultCallInterceptorOption = client.DefaultCallInterceptors([]interceptor.Interceptor{
		&client_receive.Interceptor{},
	})

	//NewMeshClient client.Client                               = mesh.NewMeshClient(defaultCallInterceptorOption, defaultCallWrapperOption)
	NewMeshClient func(option ...client.Option) client.Client = mesh.NewMeshClient

	serviceCofnigCh = make(chan config.ServiceConfigs, 1024)
)

func rotateServiceConfigWhenConfigChanged(oldConfig *config.ServiceConfigs, newConfig *config.ServiceConfigs) error {
	if !oldConfig.Equals(newConfig) {
		serviceCofnigCh <- *newConfig
	}

	return nil
}

// Logger is for adapting the platform logger
type Logger struct{}

// Printf is used to print the log of hystrix
func (l Logger) Printf(format string, items ...interface{}) {
	log.Infosf(format, items...)
}

// Init initializes the executor when launch the executor.
func (e *EventCallback) Init() error {
	e.responseTemplate = constant.DefaultResponseTemplate

	hystrix.SetLogger(Logger{})

	if nil != e.opts.ExtConfigs {
		extConfig := e.opts.ExtConfigs
		e.extConfigs = extConfig

		if v, ok := extConfig[constant.ExtConfigDefaultUserLang]; ok {
			e.defaultUserLang = v.(string)
		} else {
			e.defaultUserLang = constant.LangEnUS
		}

		if v, ok := extConfig[constant.ExtConfigCustomUserLangKey]; ok {
			e.userLangKey = v.(string)
		} else {
			e.userLangKey = constant.UserLang
		}

		// init transaction config
		if v, ok := extConfig[constant.ExtConfigTransaction]; ok {
			e.transactionConfig = v.(*config.Transaction)
		}

		// init service info config
		if v, ok := extConfig[constant.ExtConfigService]; ok {
			e.serviceConfig = v.(*config.Service)
			if "" != e.serviceConfig.CustomResponseTemplate.Value {
				e.responseTemplate = e.serviceConfig.CustomResponseTemplate.Value
			} else if "" != e.serviceConfig.ResponseTemplate {
				e.responseTemplate = e.serviceConfig.ResponseTemplate
			}

			e.responseAutoParseKeyMapping = e.serviceConfig.ResponseAutoParseKeyMapping
			e.duplicateErrorCodeTo = e.serviceConfig.DuplicateErrorCodeTo
		}

		// replace the response template if exists custom client
		if v, ok := extConfig[constant.ExtConfigCustomResponseTemplate]; ok {
			e.responseTemplate = v.(string)
		}

		// replace the response auto parse key mapping if exists custom client
		if v, ok := extConfig[constant.ExtConfigCustomResponseAutoParseKeyMapping]; ok {
			e.responseAutoParseKeyMapping = v.(map[string]string)
		}

		// replace the default client if exists custom client
		if v, ok := extConfig[constant.ExtConfigCustomClient]; ok {
			e.client = v.(client.Client)
		}

		if v, ok := extConfig[constant.ExtConfigCustomCallbackHandleWrapper]; ok {
			e.callbackHandleWrapper = v.(CallbackHandleWrapper)
		}

		if v, ok := extConfig[constant.ExtEnableExecutorLogging]; ok {
			e.enableExecutorLogging = v.(bool)
		}

		clientOptions := make([]client.Option, 0)
		// Adding external client call wrappers(callOption)
		if v, ok := extConfig[constant.ExtConfigCustomCallWrappers]; ok {
			clientOptions = append(clientOptions, client.WithOptionFromCallOption(v.(client.CallOption)))
		} else {
			clientOptions = append(clientOptions, defaultCallWrapperOption)
		}

		// Adding external client call interceptors(callOption)
		if v, ok := extConfig[constant.ExtConfigCustomCallInterceptors]; ok {
			clientOptions = append(clientOptions, client.WithOptionFromCallOption(v.(client.CallOption)))
		} else {
			clientOptions = append(clientOptions, defaultCallInterceptorOption)
		}

		e.client = NewMeshClient(clientOptions...)

		if err := e.enableTransactionProxyIfNecessary(); nil != err {
			return err
		}
	} else {
		e.client = mesh.NewMeshClient(defaultCallInterceptorOption, defaultCallWrapperOption)
	}

	if err := e.enableHTTPRouterIfNecessary(); nil != err {
		return err
	}

	if err := e.initCircuitBreakerIfNecessary(); nil != err {
		return err
	}

	// watch service config channel
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Errorsf("failed to watch service config, error:%++v\n", e)
			}
		}()

		for {
			select {
			case c := <-serviceCofnigCh:
				{
					e.enableExecutorLogging = c.EnableExecutorLogging

					if !c.Service.Equals(e.serviceConfig) {
						isExistsCustom := false
						if nil != e.extConfigs {
							_, isExistsCustom = e.extConfigs[constant.ExtConfigCustomResponseTemplate]
						}
						if !isExistsCustom &&
							"" != c.Service.CustomResponseTemplate.Value &&
							e.responseTemplate != c.Service.CustomResponseTemplate.Value {
							e.responseTemplate = c.Service.CustomResponseTemplate.Value
						} else if !isExistsCustom &&
							"" == c.Service.CustomResponseTemplate.Value &&
							e.responseTemplate != c.Service.ResponseTemplate {
							e.responseTemplate = c.Service.ResponseTemplate
						}

						isExistsCustom = false
						if nil != e.extConfigs {
							_, isExistsCustom = e.extConfigs[constant.ExtConfigCustomResponseAutoParseKeyMapping]
						}

						e.duplicateErrorCodeTo = c.Service.DuplicateErrorCodeTo

						if !isExistsCustom && !reflect.DeepEqual(e.responseAutoParseKeyMapping, c.Service.ResponseAutoParseKeyMapping) {
							e.responseAutoParseKeyMapping = c.Service.ResponseAutoParseKeyMapping
						}

						e.serviceConfig = &c.Service

						if !c.Transaction.Equals(e.transactionConfig) {
							e.transactionConfig = &c.Transaction
						}

						var copiedExtConfigs = make(map[string]interface{})
						for k, v := range e.extConfigs {
							copiedExtConfigs[k] = v
						}
						copiedExtConfigs[constant.ExtConfigService] = e.serviceConfig
						e.extConfigs = copiedExtConfigs
					}

					if !reflect.DeepEqual(e.extConfigs[constant.ExtConfigDownstreamService], c.Downstream) {
						existDownstreamServiceConfigs := e.extConfigs[constant.ExtConfigDownstreamService].(map[string]config.Downstream)

						if len(c.Downstream) > 0 {
							isExistCircuitBreaker := false
							for k, v := range c.Downstream {
								if !reflect.DeepEqual(existDownstreamServiceConfigs[k].CircuitBreaker, v.CircuitBreaker) ||
									existDownstreamServiceConfigs[k].MaxWaitingTimeMilliseconds != v.MaxWaitingTimeMilliseconds {
									timeout := constant.DefaultTimeoutMilliseconds * 1000
									if v.MaxWaitingTimeMilliseconds > 0 {
										timeout = v.MaxWaitingTimeMilliseconds
									}
									if v.CircuitBreaker.Enable {
										config := hystrix.CommandConfig{
											Timeout:                timeout,                                  // request timeout
											MaxConcurrentRequests:  v.CircuitBreaker.MaxConcurrentRequests,   // Maximum concurrency
											SleepWindow:            v.CircuitBreaker.SleepWindowMilliseconds, // How long does it take to try to see if the service is available after the circuit breaker
											RequestVolumeThreshold: v.CircuitBreaker.RequestVolumeThreshold,  // Verify the number of requests for fusing, sampling within 10 seconds
											ErrorPercentThreshold:  v.CircuitBreaker.ErrorPercentThreshold,   // Verification of the percentage of fusing errors
										}
										serviceKey := strings.ToLower(k)
										log.Infosf("Set circuit breaker[%++v] for [%s]", config, serviceKey)
										hystrix.ConfigureCommand(serviceKey, config)

										isExistCircuitBreaker = true
									}
								}
							}

							if isExistCircuitBreaker {
								e.enableCircuitBreakerIfNecessary()
							}
						}

						var copiedExtConfigs = make(map[string]interface{})
						for k, v := range e.extConfigs {
							copiedExtConfigs[k] = v
						}
						copiedExtConfigs[constant.ExtConfigDownstreamService] = c.Downstream
						e.extConfigs = copiedExtConfigs
					}
				}
			}
		}
	}()

	// register configuration on change hook function.
	config.RegisterConfigOnChangeHookFunc("executor", rotateServiceConfigWhenConfigChanged, false)
	return nil
}

// SetRouter sets the router into executor
func (e *EventCallback) SetRouter(handlerRouter *router.HandlerRouter) {
	e.handlerRouter = handlerRouter
}

// Destroy This method will be executed before the executor is destroyed
func (e *EventCallback) Destroy() error {
	return nil
}

// ResponseTemplate returns the response template from the current executor
func (e *EventCallback) ResponseTemplate() string {
	return e.responseTemplate
}

// CallbackOptions returns the callback.Options from the current executor
func (e *EventCallback) CallbackOptions() *callback.Options {
	return &e.opts
}

func (e *EventCallback) enableTransactionProxyIfNecessary() error {
	// Auto enable transaction support if necessary
	if nil != e.handlerRouter && e.handlerRouter.IsExistsCompensableTransaction {
		log.Infosf("start enable transaction support...")
		if err := imports.EnableTransactionSupports(e.handlerRouter, &e.transactionConfig.TransactionClient); nil != err {
			log.Errorsf("failed to enable transaction support, error=%++v", err)
			return err
		}
		log.Infosf("enabled transaction support successfully!")
	}

	return nil
}

func (e *EventCallback) initCircuitBreakerIfNecessary() error {
	if downstreamServiceInc, ok := e.extConfigs[constant.ExtConfigDownstreamService]; ok {
		downstreamService := downstreamServiceInc.(map[string]config.Downstream)
		isExistCircuitBreaker := false
		if len(downstreamService) > 0 {
			for k, v := range downstreamService {
				timeout := constant.DefaultTimeoutMilliseconds
				if v.MaxWaitingTimeMilliseconds > 0 {
					timeout = v.MaxWaitingTimeMilliseconds
				}
				if v.CircuitBreaker.Enable {
					config := hystrix.CommandConfig{
						Timeout:                timeout,                                  // request timeout
						MaxConcurrentRequests:  v.CircuitBreaker.MaxConcurrentRequests,   // Maximum concurrency
						SleepWindow:            v.CircuitBreaker.SleepWindowMilliseconds, // How long does it take to try to see if the service is available after the circuit breaker
						RequestVolumeThreshold: v.CircuitBreaker.RequestVolumeThreshold,  // Verify the number of requests for fusing, sampling within 10 seconds
						ErrorPercentThreshold:  v.CircuitBreaker.ErrorPercentThreshold,   // Verification of the percentage of fusing errors
					}
					serviceKey := strings.ToLower(k)
					log.Infosf("Set circuit breaker[%++v] for [%s]", config, serviceKey)
					hystrix.ConfigureCommand(serviceKey, config)

					isExistCircuitBreaker = true
				}
			}

			if isExistCircuitBreaker {
				e.enableCircuitBreakerIfNecessary()
			}
		}
	}
	return nil
}

func (e *EventCallback) enableCircuitBreakerIfNecessary() {
	if !e.isEnabledCircuitBreakerMonitor {
		hystrixStreamHandler := hystrix.NewStreamHandler()
		hystrixStreamHandler.Start()
		go func() {
			endpointAddr := fmt.Sprintf("0.0.0.0:%d", e.serviceConfig.CircuitBreakerMonitorPort)
			log.Infosf("start enable circuit breaker monitor, endpoint address:%s", endpointAddr)
			err := http.ListenAndServe(endpointAddr, hystrixStreamHandler)
			log.Errorsf("failed to start hystrix monitor, error:%++v", err)
		}()

		e.isEnabledCircuitBreakerMonitor = true
	}
}

func (e *EventCallback) enableHTTPRouterIfNecessary() error {
	// Auto enable register if necessary
	if nil != e.handlerRouter && len(e.handlerRouter.URLPathHandlers) > 0 {
		httpEndpointController := func(ctx *fasthttp.RequestCtx) {
			topicAttributes := make(map[string]string)
			appProps := make(map[string]string)
			ctx.Request.Header.VisitAll(func(key, value []byte) {
				keyStr := string(key)
				if strings.HasPrefix(keyStr, constant.AttributesPrefix) {
					topicAttributes[strings.TrimPrefix(keyStr, constant.AttributesPrefix)] = string(value)
				} else {
					appProps[keyStr] = string(value)
				}
			})

			ctx.QueryArgs().VisitAll(func(key, value []byte) {
				appProps[string(key)] = string(value)
			})

			request := &msg.Message{
				TopicAttribute: topicAttributes,
				RequestURL:     string(ctx.URI().Path()),
				Body:           ctx.Request.Body(),
			}
			request.SetAppProps(appProps)
			log.Debugsf("request URL:%s, URI:%++v", string(ctx.URI().FullURI()), string(ctx.URI().Path()))
			context := context.Background()
			response, err := e.Handle(context, request)
			if nil != err {
				// never happen when we call Eventkit Handler
				log.Errorsf("failed to execute ExecuteHandler, error=%++v", err)
				return
			}
			if nil != response {
				response.RangeAppProps(func(k string, v string) {
					ctx.Response.Header.Set(k, v)
				})
				if len(response.TopicAttribute) != 0 {
					for key, value := range response.TopicAttribute {
						ctx.Response.Header.Set(constant.AttributesPrefix+key, value)
					}
				}
				ctx.Write(response.Body)
			}
		}

		// enable http server, if necessary
		endpointAddr := fmt.Sprintf("0.0.0.0:%d", e.opts.Port)
		router := fasthttprouter.New()
		server := &fasthttp.Server{
			Handler:                       router.Handler,
			MaxRequestBodySize:            1024 * 1024 * 1024,
			DisableHeaderNamesNormalizing: true,
		}
		log.Infosf(`Configured for Endpoint. => address: %s => port:%d`, "0.0.0.0", e.opts.Port)
		for path, opts := range e.handlerRouter.URLPathHandlers {
			switch opts.HandlerOptions.HTTPMethod {
			case constant.HTTPMethodHead:
				{
					log.Infosf("*                         => HEAD %s", path)
					router.HEAD(path, httpEndpointController)
				}
			case constant.HTTPMethodPost:
				{
					log.Infosf("*                         => POST %s", path)
					router.POST(path, httpEndpointController)
				}
			case constant.HTTPMethodOptions:
				{
					log.Infosf("*                         => OPTIONS %s", path)
					router.OPTIONS(path, httpEndpointController)
				}
			case constant.HTTPMethodPut:
				{
					log.Infosf("*                         => PUT %s", path)
					router.PUT(path, httpEndpointController)
				}
			case constant.HTTPMethodPatch:
				{
					log.Infosf("*                         => PATH %s", path)
					router.PATCH(path, httpEndpointController)
				}
			case constant.HTTPMethodDelete:
				{
					log.Infosf("*                         => DELETE %s", path)
					router.DELETE(path, httpEndpointController)
				}
			default:
				log.Infosf("*                         => GET %s", path)
				router.GET(path, httpEndpointController)
			}

		}

		go func(srv *fasthttp.Server, addr string) {
			if err := srv.ListenAndServe(addr); err != nil {
				log.Errorsf("failed to Start: start http server error:%++v", err.Error())
				panic(err)
			}
			log.Infosf("successfully started http server!")
		}(server, endpointAddr)

	}

	return nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (e EventCallback) judgeUserLangWithDefaultValue(request *msg.Message) string {
	lang := request.GetAppPropertySilence(e.userLangKey)

	if "" == lang {
		lang = e.defaultUserLang
	}

	return lang
}

// Handle executes the main process of each callback, matches the corresponding handler according to the eventID,
// and decodes the message and executes the corresponding method through reflection
func (e *EventCallback) Handle(parentCtx context.Context, request *msg.Message) (response *msg.Message, err error) {
	var gerr error
	var ctx context.Context
	var handlerContexts *contexts.HandlerContexts
	var preCtx context.Context
	var responseTemplate = e.responseTemplate
	var duplicateErrorCodeTo = e.duplicateErrorCodeTo
	var responseDataWhenError interface{}
	var customErrorWrapperFn msg.CustomErrorWrapperFn
	var ins base.HandlerInterface
	var st = time.Now()
	if nil != request && e.enableExecutorLogging {
		log.Debugf(parentCtx, "server receive message, header:[%s] topic attributes:[%s] body:[%s]", request.AppPropsToString(), request.TopicAttributesToString(), string(request.Body))
	}

	isLocalCallCheckFunc := func(eventId string) bool {
		hp := e.handlerRouter.MatchHandler(eventId)
		return nil != hp
	}
	var downstreamServiceConfigs map[string]config.Downstream
	if v, ok := e.extConfigs[constant.ExtConfigDownstreamService]; ok {
		downstreamServiceConfigs = v.(map[string]config.Downstream)
	}
	remoteCall := remote.NewDefaultRemoteCall(e.client, e, isLocalCallCheckFunc, downstreamServiceConfigs)

	lang := e.judgeUserLangWithDefaultValue(request)
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf(constant.SystemInternalError, "panic, error=%++v, stack=%s", r, debug.Stack())
		}
		if nil != err {
			response = msg.WrapperErrorResponse(err, lang, responseTemplate, customErrorWrapperFn, responseDataWhenError)
			if len(duplicateErrorCodeTo) > 0 {
				response.SetAppProperty(duplicateErrorCodeTo, response.GetAppPropertySilence(constant.ReturnErrorCode))
			}
			if nil != ins {
				for k, v := range ins.GetResponseHeader() {
					response.SetAppProperty(k, v)
				}
			}
			if nil != handlerContexts {
				handlerContexts.Inject(true, response)
			}
			if e.isSemiSyncCall(request.GetAppProps()) {
				// when semi-sync call, need reply to error response
				replyToAddress := request.GetAppPropertySilence(constant.RrReplyTo)
				err2 := e.replySemiSyncCall(ctx, replyToAddress, response.GetAppProps())
				if err2 != nil {
					log.Errorsf("replySemiSyncCall fail, err:%s", err2)
				}
			}
			err = nil
		}
		if nil != response && e.enableExecutorLogging {
			log.Debugf(ctx, "server response message, header:[%s] topic attributes:[%s] body:[%s]", response.AppPropsToString(), response.TopicAttributesToString(), string(response.Body))
		}
	}()

	defer func() {
		if nil != e.callbackHandleWrapper {
			postResponse, postErr := e.callbackHandleWrapper.PostHandle(ctx, response)
			if nil != postErr {
				err = postErr
				return
			}
			response = postResponse
		}
	}()

	if nil != e.callbackHandleWrapper {
		skip, tmpCtx, preRequest, preErr := e.callbackHandleWrapper.PreHandle(parentCtx, request)
		if skip || nil != preErr {
			ctx = tmpCtx
			return preRequest, preErr
		}
		preCtx = tmpCtx
	} else {
		preCtx = parentCtx
	}

	// new handler contexts
	ctx, handlerContexts = contexts.BuildContextFromParent(preCtx, contexts.Lang(lang))

	if nil != request {
		traceID := request.GetAppPropertyEitherSilence(constant.KeyTraceID, constant.KeyTraceIDOld)
		spanID := request.GetAppPropertyEitherSilence(constant.KeySpanID, constant.KeySpanIDOld)
		parentSpanID := request.GetAppPropertyEitherSilence(constant.KeyParentSpanID, constant.KeyParentSpanIDOld)
		replyToAddress := request.GetAppPropertySilence(constant.RrReplyTo)
		passThroughHeaderKeyList := request.GetAppPropertySilence(constant.PassThroughHeaderKeyListKey)
		timeoutMilliseconds := 0

		timeoutMillisecondsStr := request.GetAppPropertySilence(constant.To3)
		if "" != timeoutMillisecondsStr {
			timeout, err := strconv.Atoi(timeoutMillisecondsStr)
			if err == nil {
				timeoutMilliseconds = timeout
			}
		}

		if 0 == timeoutMilliseconds {
			timeoutMilliseconds = constant.DefaultTimeoutMilliseconds
		}

		if traceID == "" && nil != e.serviceConfig {
			// get su id,trace id from context
			traceID = util.GenerateSerialNo(
				e.serviceConfig.Org,
				e.serviceConfig.Wks,
				e.serviceConfig.Env,
				e.serviceConfig.Su,
				e.serviceConfig.InstanceID, constant.TraceIDType)
		}

		if spanID == "" && nil != e.serviceConfig {
			spanID = util.GenerateSerialNo(
				e.serviceConfig.Org,
				e.serviceConfig.Wks,
				e.serviceConfig.Env,
				e.serviceConfig.Su,
				e.serviceConfig.InstanceID, constant.SpanIDType)
		}

		if "" == parentSpanID {
			parentSpanID = spanID
		}

		handlerContexts.With(
			contexts.Span(
				contexts.BuildSpanContexts(
					contexts.TraceID(traceID),
					contexts.SpanID(spanID),
					contexts.ParentSpanID(parentSpanID),
					contexts.TimeoutMilliseconds(timeoutMilliseconds),
					contexts.ReplyToAddress(replyToAddress),
					contexts.PassThroughHeaderKeyList(passThroughHeaderKeyList),
				),
			))
		ctx = context.WithValue(ctx, constant.KeyST, st)
		ctx = context.WithValue(ctx, constant.To3, timeoutMilliseconds)
		ctx = context.WithValue(ctx, constant.KeyTraceID, traceID)
		ctx = context.WithValue(ctx, constant.KeySpanID, spanID)
		ctx = context.WithValue(ctx, constant.KeyParentSpanID, parentSpanID)
		ctx = context.WithValue(ctx, constant.KeyParentSpanIDForLog, parentSpanID)
		ctx = context.WithValue(ctx, constant.KeyTopicID, request.GetMsgTopicId())
	}

	if nil != e.serviceConfig {
		ctx = context.WithValue(ctx, constant.KeyServiceID, e.serviceConfig.ServiceID)
		handlerContexts.With(
			contexts.ServiceID(e.serviceConfig.ServiceID),
			contexts.NodeID(e.serviceConfig.NodeID),
			contexts.InstanceID(e.serviceConfig.InstanceID),
			contexts.AZ(e.serviceConfig.Az),
			contexts.WKS(e.serviceConfig.Wks),
			contexts.ENV(e.serviceConfig.Env),
			contexts.SU(e.serviceConfig.Su),
			contexts.ORG(e.serviceConfig.Org),
			contexts.CommonSu(e.serviceConfig.CommonSu),
			contexts.ResponseTemplate(e.responseTemplate),
			contexts.ResponseAutoParseKeyMapping(e.responseAutoParseKeyMapping),
			contexts.DownstreamConfigs(downstreamServiceConfigs),
		)
	}

	var hp *router.Options
	indexOfInterceptorsExecuted := 0
	defer func() {
		if nil == response {
			response = &msg.Message{}
		}
		// Do postHandle of interceptors
		for i := indexOfInterceptorsExecuted - 1; i >= 0; i-- {
			interceptor := hp.HandlerOptions.Interceptors[i]
			interceptorName := fmt.Sprintf("%s", interceptor)
			if nil != e.serviceConfig && e.serviceConfig.IsMarkedAsSkippedHandleInterceptor(interceptorName) {
				log.Debugf(ctx, "Skipping the `PostHandle` of handle interceptor:%s", interceptorName)
			} else {
				if gerr = interceptor.PostHandle(ctx, request, response); nil != gerr {
					log.Errorf(ctx, "post handle of interceptor execute failed, error=%s", errors.ErrorToString(gerr))
					break
				}
			}
		}
	}()

	defer func(lang string) {
		if r := recover(); r != nil {
			err = errors.Errorf(constant.SystemInternalError, "Service ID: %s - failed to execute handler, error=%++v, stack=%s", e.serviceConfig.ServiceID, r, debug.Stack())
		}

		if nil != err {
			log.Errorf(ctx, "Service ID: %s - failed to execute handler, error=%s", e.serviceConfig.ServiceID, errors.ErrorToString(err))
			response = msg.WrapperErrorResponse(err, lang, responseTemplate, customErrorWrapperFn, responseDataWhenError)
			if len(duplicateErrorCodeTo) > 0 {
				response.SetAppProperty(duplicateErrorCodeTo, response.GetAppPropertySilence(constant.ReturnErrorCode))
			}
			if nil != ins {
				for k, v := range ins.GetResponseHeader() {
					response.SetAppProperty(k, v)
				}
			}
			if nil != handlerContexts {
				handlerContexts.Inject(true, response)
			}
			if e.isSemiSyncCall(request.GetAppProps()) {
				// when semi-sync call, need reply to error response
				replyToAddress := request.GetAppPropertySilence(constant.RrReplyTo)
				err2 := e.replySemiSyncCall(ctx, replyToAddress, response.GetAppProps())
				if err2 != nil {
					log.Errorsf("replySemiSyncCall fail, err:%s", err2)
				}
			}
			err = nil
		}
	}(lang)

	if "" != request.RequestURL {
		hp = e.handlerRouter.MatchHandlerWithURLPath(request.RequestURL)
		if nil == hp {
			return nil, errors.Errorf(constant.CannotFoundHandlerWithURLError, "Service ID: %s - Cannot found handler with URL: %s", e.serviceConfig.ServiceID, request.RequestURL)
		}
	} else {
		if !request.IsValidTopicType() {
			return nil, errors.Errorf(constant.InvalidEventTypeError, "Service ID: %s - Invalid event type: %++v", e.serviceConfig.ServiceID, request)
		}

		eventID := request.GetMsgTopicId()
		hp = e.handlerRouter.MatchHandler(eventID)
		if nil == hp {
			return nil, errors.Errorf(constant.CannotFoundHandlerWithEventIDError, "Service ID: %s - Cannot found handler with event id: %s", e.serviceConfig.ServiceID, eventID)
		}
	}

	if len(hp.HandlerOptions.ResponseTemplate) > 0 {
		responseTemplate = hp.HandlerOptions.ResponseTemplate
		handlerContexts.With(contexts.ResponseTemplate(responseTemplate))
		responseDataWhenError = hp.HandlerOptions.ResponseDataWhenErrorForResponseTemplate
		customErrorWrapperFn = hp.HandlerOptions.CustomErrorWrapperFn
	}

	// Do preHandle of interceptors
	for i := 0; i < len(hp.HandlerOptions.Interceptors); i++ {
		interceptor := hp.HandlerOptions.Interceptors[i]
		interceptorName := fmt.Sprintf("%s", interceptor)
		if nil != e.serviceConfig && e.serviceConfig.IsMarkedAsSkippedHandleInterceptor(interceptorName) {
			log.Debugf(ctx, "Skipping the `PreHandle` of handle interceptor:%s", interceptorName)
		} else {
			if gerr = interceptor.PreHandle(ctx, request); nil != gerr {
				return nil, gerr
			}
		}
		indexOfInterceptorsExecuted++
	}

	decoder := hp.HandlerOptions.Codec.Decoder()
	lengthOfInParams := len(hp.HandlerOptions.HandlerMethodInParams)
	parameterInValues := make([]reflect.Value, lengthOfInParams)

	for i := 0; i < min(len(hp.HandlerOptions.HandlerMethodInParams), 1); i++ {
		var pValue reflect.Value
		pType := hp.HandlerOptions.HandlerMethodInParams[i]
		if pType.Kind() == reflect.Ptr {
			pValue = reflect.New(pType.Elem())
		} else {
			pValue = reflect.New(pType)
		}

		if gerr = decoder.Decode(request.Body, pValue.Interface()); gerr != nil {
			return nil, errors.New(constant.UpstreamServiceMessageDecodeError, gerr)
		}

		if pType.Kind() == reflect.Ptr {
			parameterInValues[i] = pValue
		} else {
			parameterInValues[i] = pValue.Elem()
		}
	}

	var instance reflect.Value
	if hp.HandlerOptions.HandlerReflectType.Kind() == reflect.Ptr {
		instance = reflect.New(hp.HandlerOptions.HandlerReflectType.Elem())
	} else {
		instance = reflect.New(hp.HandlerOptions.HandlerReflectType)
	}

	ins = instance.Interface().(base.HandlerInterface)
	ins.SetRemoteCall(remoteCall)
	ins.SetLang(handlerContexts.Lang)
	ins.SetContext(ctx)
	ins.SetBody(request.Body)
	ins.SetExtConfigs(e.extConfigs)

	if nil != hp.HandlerOptions.CustomValidationOptions {
		ins.SetCombineErrors(hp.HandlerOptions.CustomValidationOptions.CombineErrors)
		ins.SetCustomValidationRegisterFunctions(hp.HandlerOptions.CustomValidationOptions.CustomValidationRegisterFunctions)
	}

	ins.SetTopicAttributes(request.TopicAttribute)
	ins.SetRequestHeader(request.CloneAppProps())

	var responseBody []byte
	var returnValues []reflect.Value

	if hp.HandlerOptions.EnableValidation && 0 != lengthOfInParams {
		// DO validation
		validationMethod := instance.MethodByName(constant.ValidationMethodName)
		validationReturnValues := validationMethod.Call(parameterInValues)
		validationReturnValue := validationReturnValues[0]
		if !valueIsNil(validationReturnValue) {
			return nil, validationReturnValue.Interface().(error)
		}
	}

	if hp.HandlerOptions.InvokePreHandle && 0 != lengthOfInParams {
		// DO pre handle
		preHandleMethod := instance.MethodByName(constant.PreHandleMethodName)
		preHandleReturnValues := preHandleMethod.Call(parameterInValues)
		preHandleReturnValue := preHandleReturnValues[0]
		if !valueIsNil(preHandleReturnValue) {
			return nil, preHandleReturnValue.Interface().(error)
		}
	}

	if nil == hp.Compensable {
		// invoke non-TCC handler
		targetMethod := instance.MethodByName(hp.HandlerOptions.HandlerMethodName)
		returnValues = targetMethod.Call(parameterInValues)
	} else {
		transactionContext := context.WithValue(ctx, constant.ContextTransactionKey, e.transactionConfig)
		proxy, r := proxy.NewTransactionProxy(transactionContext, instance, hp.Compensable)
		if r != nil {
			gerr = errors.Errorf(constant.NewTransactionProxyError, "Service ID: %s - new transaction proxy error: %v", e.serviceConfig.ServiceID, r)
			return nil, gerr
		}

		convertFunc := func(inputValues []reflect.Value) []interface{} {
			r := make([]interface{}, 0)
			for _, i := range inputValues {
				r = append(r, i.Interface())
			}

			return r
		}
		returnValues = proxy.Do(convertFunc(parameterInValues)...)
	}

	lastReturnValue := returnValues[len(returnValues)-1]
	if !valueIsNil(lastReturnValue) {
		gerr = lastReturnValue.Interface().(error)
		return nil, gerr
	} else if len(returnValues) == 1 {
		responseBody = nil
	} else {
		encoder := hp.HandlerOptions.Codec.Encoder()

		if responseBody, gerr = encoder.Encode(returnValues[0].Interface()); nil != gerr {
			return nil, errors.New(constant.UpstreamServiceMessageEncodeError, gerr)
		}
	}

	response = &msg.Message{
		ID:          request.ID,
		SessionName: request.SessionName,
		Body:        responseBody,
	}
	responseHeader := ins.GetResponseHeader()
	if ins.IsDiscardResponse() {
		if nil == responseHeader {
			responseHeader = make(map[string]string)
		}

		responseHeader[constant.DiscardResponse] = "1"
	}
	response.SetAppProps(responseHeader)
	return response, nil
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

func (e *EventCallback) isSemiSyncCall(headers map[string]string) bool {
	_, ok := headers[constant.RrReplyTo]
	return ok
}

func (e *EventCallback) replySemiSyncCall(ctx context.Context, replyToAddress string, headers map[string]string) error {
	semiResponse := mesh.NewMeshResponse(nil, mesh.WithResponseHeader(headers), mesh.WithResponseCodec(codec.BuildTextCodec()))
	semiResponse.WithOptions(mesh.WithReplyToAddress(replyToAddress))
	err := e.client.ReplySemiSyncCall(ctx, semiResponse)
	if err != nil {
		return err
	}
	return nil
}
