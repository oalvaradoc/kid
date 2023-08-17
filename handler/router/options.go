package router

import (
	"fmt"
	"git.multiverse.io/eventkit/kit/codec"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/compensable"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/interceptor"
	"git.multiverse.io/eventkit/kit/log"
	"git.multiverse.io/eventkit/kit/validation"
	"github.com/go-playground/validator/v10"
	"reflect"
)

// Options is a set of configuration parameters that contains HandlerOptions and Compensable and event expression.
type Options struct {
	EventExpression string
	Compensable     *compensable.Compensable
	HandlerOptions  HandlerOptions
}

// HandlerOptions is a set of configuration parameters that perform the handler
type HandlerOptions struct {
	RegisterType                             int
	HandlerName                              string
	HandlerReflectType                       reflect.Type
	HandlerMethodName                        string
	HandlerMethodInParams                    []reflect.Type
	HandlerMethodOutParams                   []reflect.Type
	Interceptors                             []interceptor.Interceptor
	Codec                                    codec.Codec
	URLPath                                  string
	HTTPMethod                               int
	InvokePreHandle                          bool
	EnableValidation                         bool
	CustomValidationOptions                  *CustomValidationOptions
	ResponseTemplate                         string
	ResponseDataWhenErrorForResponseTemplate interface{}
	CustomErrorWrapperFn                     msg.CustomErrorWrapperFn
}

// CustomValidationOptions is a set of configuration parameters that perform the validation of the handler
type CustomValidationOptions struct {
	CombineErrors                     bool
	CustomValidationRegisterFunctions []validation.CustomValidationRegisterFunc
}

// Option sets an optional parameter into the router.
type Option func(*Options)

// HandlerOption sets an optional parameter into the router.
type HandlerOption func(*HandlerOptions)

// CustomValidationOption sets an optional parameter into the validation configuration.
type CustomValidationOption func(*CustomValidationOptions)

// AddRegisterValidationFunc sets the validation function of the transaction handler when register the handler into the router.
func AddRegisterValidationFunc(tag string, validationFunc validator.Func, callValidationEvenIfNull ...bool) CustomValidationOption {
	return func(options *CustomValidationOptions) {
		if nil == options.CustomValidationRegisterFunctions {
			options.CustomValidationRegisterFunctions = make([]validation.CustomValidationRegisterFunc, 0)
		}
		var nilCheck bool
		if len(callValidationEvenIfNull) > 0 {
			nilCheck = callValidationEvenIfNull[0]
		}
		options.CustomValidationRegisterFunctions = append(options.CustomValidationRegisterFunctions, validation.CustomValidationRegisterFunc{
			Tag:                      tag,
			Func:                     validationFunc,
			CallValidationEvenIfNull: nilCheck,
		})

	}
}

// Compensable sets the compensable information of the transaction handler when register the handler into the router,
// only called by transaction handler registration
func Compensable(compensable *compensable.Compensable) Option {
	return func(options *Options) {
		options.Compensable = compensable
		options.HandlerOptions.HandlerMethodName = compensable.TryMethod
	}
}

// WithRegisterType sets the register type of the handler when register the handler into the router
func WithRegisterType(registerType int) Option {
	return func(options *Options) {
		options.HandlerOptions.RegisterType = registerType
	}
}

// WithHandlerName sets the handler name of the handler when register the handler into the router, the same as Method
func WithHandlerName(handlerName string) Option {
	return func(options *Options) {
		options.HandlerOptions.HandlerName = handlerName
	}
}

// WithHandlerReflectType sets the reflect type of the handler when register the handler into the router
func WithHandlerReflectType(handlerReflectType reflect.Type) Option {
	return func(options *Options) {
		options.HandlerOptions.HandlerReflectType = handlerReflectType
	}
}

// WithResponseTemplate sets the response template of the handler when register the handler into the router
func WithResponseTemplate(responseTemplate string, responseDataWhenError ...interface{}) Option {
	return func(options *Options) {
		options.HandlerOptions.ResponseTemplate = responseTemplate
		if len(responseDataWhenError) > 0 {
			options.HandlerOptions.ResponseDataWhenErrorForResponseTemplate = responseDataWhenError[0]
		}
	}
}

// WithResponseTemplateAndHeader sets the response template and response header of the handler when register the handler into the router
func WithResponseTemplateAndHeader(responseTemplate string, responseHeader map[string]string, responseDataWhenError ...interface{}) Option {
	return func(options *Options) {
		options.HandlerOptions.ResponseTemplate = responseTemplate
		options.HandlerOptions.CustomErrorWrapperFn = func(errorCode, errorMessage string, response *msg.Message) {
			if nil != response {
				for k, v := range responseHeader {
					response.SetAppProperty(k, v)
				}
			}
		}
		if len(responseDataWhenError) > 0 {
			options.HandlerOptions.ResponseDataWhenErrorForResponseTemplate = responseDataWhenError[0]
		}
	}
}

// WithResponseAndCustomErrorWrapperFn sets the custom error wrapper function of the handler when register the handler into the router
func WithResponseAndCustomErrorWrapperFn(responseTemplate string, customErrorWrapperFn msg.CustomErrorWrapperFn, responseDataWhenError ...interface{}) Option {
	return func(options *Options) {
		options.HandlerOptions.ResponseTemplate = responseTemplate
		options.HandlerOptions.CustomErrorWrapperFn = customErrorWrapperFn
		if len(responseDataWhenError) > 0 {
			options.HandlerOptions.ResponseDataWhenErrorForResponseTemplate = responseDataWhenError[0]
		}
	}
}

// Method sets the method name of the handler when register the handler into the router
func Method(method string) Option {
	return func(options *Options) {
		options.HandlerOptions.HandlerMethodName = method
	}
}

// HandlePost is used to turn on the http POST method. After it is turned on,
// you can directly send a request to the corresponding handler through the http client and URL path.
func HandlePost(urlPath string) Option {
	return func(options *Options) {
		options.HandlerOptions.HTTPMethod = constant.HTTPMethodPost
		options.HandlerOptions.URLPath = urlPath
	}
}

// HandleGet is used to turn on the http GET method. After it is turned on,
// you can directly send a request to the corresponding handler through the http client and URL path.
func HandleGet(urlPath string) Option {
	return func(options *Options) {
		options.HandlerOptions.HTTPMethod = constant.HTTPMethodGet
		options.HandlerOptions.URLPath = urlPath
	}
}

// HandleOptions is used to turn on the http OPTIONS method. After it is turned on,
// you can directly send a request to the corresponding handler through the http client and URL path.
func HandleOptions(urlPath string) Option {
	return func(options *Options) {
		options.HandlerOptions.HTTPMethod = constant.HTTPMethodOptions
		options.HandlerOptions.URLPath = urlPath
	}
}

// HandleHead is used to turn on the http HEAD method. After it is turned on,
// you can directly send a request to the corresponding handler through the http client and URL path.
func HandleHead(urlPath string) Option {
	return func(options *Options) {
		options.HandlerOptions.HTTPMethod = constant.HTTPMethodHead
		options.HandlerOptions.URLPath = urlPath
	}
}

// HandlePatch is used to turn on the http PATCH method. After it is turned on,
// you can directly send a request to the corresponding handler through the http client and URL path.
func HandlePatch(urlPath string) Option {
	return func(options *Options) {
		options.HandlerOptions.HTTPMethod = constant.HTTPMethodPatch
		options.HandlerOptions.URLPath = urlPath
	}
}

// HandleDelete is used to turn on the http DELETE method. After it is turned on,
// you can directly send a request to the corresponding handler through the http client and URL path.
func HandleDelete(urlPath string) Option {
	return func(options *Options) {
		options.HandlerOptions.HTTPMethod = constant.HTTPMethodDelete
		options.HandlerOptions.URLPath = urlPath
	}
}

// HandlePut is used to turn on the http PUT method. After it is turned on,
// you can directly send a request to the corresponding handler through the http client and URL path.
func HandlePut(urlPath string) Option {
	return func(options *Options) {
		options.HandlerOptions.HTTPMethod = constant.HTTPMethodPut
		options.HandlerOptions.URLPath = urlPath
	}
}

// EnableValidation marks to enable the validation logic when execute the handler
func EnableValidation(combineErrors bool, validationOptions ...CustomValidationOption) Option {
	return func(options *Options) {
		options.HandlerOptions.EnableValidation = true
		if nil == options.HandlerOptions.CustomValidationOptions {
			options.HandlerOptions.CustomValidationOptions = &CustomValidationOptions{}
		}
		options.HandlerOptions.CustomValidationOptions.CombineErrors = combineErrors
		for _, opt := range validationOptions {
			opt(options.HandlerOptions.CustomValidationOptions)
		}
	}
}

// RegisterCustomValidationOptions used to register the custom validation options
func RegisterCustomValidationOptions(combineErrors bool, validationOptions ...CustomValidationOption) Option {
	return func(options *Options) {
		if nil == options.HandlerOptions.CustomValidationOptions {
			options.HandlerOptions.CustomValidationOptions = &CustomValidationOptions{}
		}
		for _, opt := range validationOptions {
			opt(options.HandlerOptions.CustomValidationOptions)
		}
	}
}

// DisableValidation marks to disable the validation logic when execute the handler
func DisableValidation() Option {
	return func(options *Options) {
		options.HandlerOptions.EnableValidation = false
		if nil == options.HandlerOptions.CustomValidationOptions {
			options.HandlerOptions.CustomValidationOptions = &CustomValidationOptions{}
		}
		options.HandlerOptions.CustomValidationOptions.CombineErrors = false
	}
}

// WithHandlerMethodInParams adds the method in parameter type into the router config, usually called by framework
func WithHandlerMethodInParams(handlerMethodInParams []reflect.Type) Option {
	return func(options *Options) {
		options.HandlerOptions.HandlerMethodInParams = handlerMethodInParams
	}
}

// WithHandlerMethodOutParams adds the method out parameter type into the router config, usually called by framework
func WithHandlerMethodOutParams(handlerMethodOutParams []reflect.Type) Option {
	return func(options *Options) {
		options.HandlerOptions.HandlerMethodOutParams = handlerMethodOutParams
	}
}

// AddInterceptors adds one or more interceptors into the router config
func AddInterceptors(interceptors ...interceptor.Interceptor) Option {
	return func(options *Options) {
		options.HandlerOptions.Interceptors = append(options.HandlerOptions.Interceptors, interceptors...)
	}
}

// AddInterceptor adds an interceptor into the router config
func AddInterceptor(interceptor interceptor.Interceptor) Option {
	return AddInterceptors(interceptor)
}

// WithInterceptors sets one or more interceptors into the router config
func WithInterceptors(interceptors ...interceptor.Interceptor) Option {
	return func(options *Options) {
		options.HandlerOptions.Interceptors = interceptors
	}
}

// RemoveInterceptor sets a interceptor name that need to remove
func RemoveInterceptor(name string) Option {
	return func(options *Options) {
		if 0 != len(options.HandlerOptions.Interceptors) {
			for i, interceptor := range options.HandlerOptions.Interceptors {
				if fmt.Sprintf("%s", interceptor) == name {
					log.Infosf("HandlerMethodName:%++v removes interceptor:%s", options.HandlerOptions.HandlerMethodName, name)
					options.HandlerOptions.Interceptors = append(options.HandlerOptions.Interceptors[:i], options.HandlerOptions.Interceptors[i+1:]...)
				}
			}
		}
	}
}

// RemoveInterceptors sets one or more interceptors name that need to remove
func RemoveInterceptors(names ...string) Option {
	return func(options *Options) {
		if 0 != len(options.HandlerOptions.Interceptors) {
			for _, name := range names {
				for i, interceptor := range options.HandlerOptions.Interceptors {
					if fmt.Sprintf("%s", interceptor) == name {
						log.Infosf("HandlerMethodName:%++v removes interceptor:%s", options.HandlerOptions.HandlerMethodName, name)
						options.HandlerOptions.Interceptors = append(options.HandlerOptions.Interceptors[:i], options.HandlerOptions.Interceptors[i+1:]...)
					}
				}
			}
		}
	}
}

// WithCodec sets a codec into the router config
func WithCodec(codec codec.Codec) Option {
	return func(options *Options) {
		options.HandlerOptions.Codec = codec
	}
}
