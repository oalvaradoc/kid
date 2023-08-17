package handler

import (
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/interceptor"
	"git.multiverse.io/eventkit/kit/sed/callback"
	"git.multiverse.io/eventkit/kit/validation"
	"git.multiverse.io/eventkit/kit/wrapper"
	"github.com/go-playground/validator/v10"
)

// CustomOptions is a set of configuration parameters for register a executor
type CustomOptions struct {
	CustomOptionConfigs map[string]interface{}
}

// DefaultCustomValidationOptions is a set of configuration parameters for register the validation options into executor
type DefaultCustomValidationOptions struct {
	CombineErrors                            bool
	DefaultCustomValidationRegisterFunctions []validation.CustomValidationRegisterFunc
}

// CustomOption sets a parameter of customization configuration of executor
type CustomOption func(*CustomOptions)

// DefaultCustomValidationOption sets a parameter of customization validation configurations
type DefaultCustomValidationOption func(*DefaultCustomValidationOptions)

// WithCallbackOptions sets callback options by default when the handler is registered to the router
func WithCallbackOptions(sedClientOptions ...callback.Option) CustomOption {
	return func(configs *CustomOptions) {
		configs.CustomOptionConfigs[constant.ExtConfigCustomSedClientOptions] = sedClientOptions
	}
}

// WithDefaultInterceptors sets one or more interceptors by default when the handler is registered to the router
func WithDefaultInterceptors(interceptors ...interceptor.Interceptor) CustomOption {
	return func(options *CustomOptions) {
		options.CustomOptionConfigs[constant.ExtConfigCustomDefaultInterceptors] = interceptors
	}
}

// WithCallbackHandleWrapper sets the callback handler wrapper by default when the handler is registered to the router
func WithCallbackHandleWrapper(callbackHandleWrapper CallbackHandleWrapper) CustomOption {
	return func(options *CustomOptions) {
		options.CustomOptionConfigs[constant.ExtConfigCustomCallbackHandleWrapper] = callbackHandleWrapper
	}
}

// WithDefaultResponseTemplate sets the response template by default when the handler is registered to the router
func WithDefaultResponseTemplate(responseTemplate string) CustomOption {
	return func(options *CustomOptions) {
		options.CustomOptionConfigs[constant.ExtConfigCustomResponseTemplate] = responseTemplate
	}
}

// AddDefaultRegisterValidationFunc sets the validation functions by default when the handler is registered to the router,
// only used when enabled the validation
func AddDefaultRegisterValidationFunc(tag string, validationFunc validator.Func, callValidationEvenIfNull ...bool) DefaultCustomValidationOption {
	return func(options *DefaultCustomValidationOptions) {
		if nil == options.DefaultCustomValidationRegisterFunctions {
			options.DefaultCustomValidationRegisterFunctions = make([]validation.CustomValidationRegisterFunc, 0)
		}
		var nilCheck bool
		if len(callValidationEvenIfNull) > 0 {
			nilCheck = callValidationEvenIfNull[0]
		}
		options.DefaultCustomValidationRegisterFunctions = append(options.DefaultCustomValidationRegisterFunctions, validation.CustomValidationRegisterFunc{
			Tag:                      tag,
			Func:                     validationFunc,
			CallValidationEvenIfNull: nilCheck,
		})
	}
}

// WithDefaultEnableValidation sets whether to enable validation by default when the handler is registered to the router
func WithDefaultEnableValidation(combineErrors bool, defaultValidationOptions ...DefaultCustomValidationOption) CustomOption {
	return func(options *CustomOptions) {
		defaultCustomValidationOptions := &DefaultCustomValidationOptions{
			CombineErrors: combineErrors,
		}

		for _, opt := range defaultValidationOptions {
			opt(defaultCustomValidationOptions)
		}

		options.CustomOptionConfigs[constant.ExtConfigDefaultEnableValidation] = defaultCustomValidationOptions
	}
}

// WithCallWrappers sets the call wrappers by default when the router is registered to the router
func WithCallWrappers(cas ...wrapper.Wrapper) CustomOption {
	return func(options *CustomOptions) {
		options.CustomOptionConfigs[constant.ExtConfigCustomCallWrappers] = client.WithCallWrappers(cas)
	}
}

// WithCustomUserLangKey sets the user lang key
func WithCustomUserLangKey(customUserLangKey string) CustomOption {
	return func(options *CustomOptions) {
		options.CustomOptionConfigs[constant.ExtConfigCustomUserLangKey] = customUserLangKey
	}
}

// WithDefaultUserLang sets the default user lang for executor
func WithDefaultUserLang(defaultUserLang string) CustomOption {
	return func(options *CustomOptions) {
		options.CustomOptionConfigs[constant.ExtConfigDefaultUserLang] = defaultUserLang
	}
}

// WithResponseAutoParseKeyMapping sets the response auto parse key mapping into the executor
func WithResponseAutoParseKeyMapping(responseAutoParseKeyMapping map[string]string) CustomOption {
	return func(options *CustomOptions) {
		options.CustomOptionConfigs[constant.ExtConfigCustomResponseAutoParseKeyMapping] = responseAutoParseKeyMapping
	}
}

// WithCallInterceptors sets one or more call interceptors into the executor
func WithCallInterceptors(callInterceptor ...interceptor.Interceptor) CustomOption {
	return func(options *CustomOptions) {
		options.CustomOptionConfigs[constant.ExtConfigCustomCallInterceptors] = client.WithCallInterceptors(callInterceptor)
	}
}

// WithClient sets the client instance into the executor for service remote call
func WithClient(client client.Client) CustomOption {
	return func(configs *CustomOptions) {
		configs.CustomOptionConfigs[constant.ExtConfigCustomClient] = client
	}
}
