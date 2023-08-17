package base

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/client/mesh"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/handler/remote"
	"git.multiverse.io/eventkit/kit/log"
	"git.multiverse.io/eventkit/kit/validation"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// Tup2 represents a pair of validator and translator
type Tup2 struct {
	Validate *validator.Validate
	Trans    ut.Translator
}

type cacheValidation struct {
	sync.RWMutex
	Cache map[string]Tup2
}

var cv = &cacheValidation{
	Cache: make(map[string]Tup2),
}

func (c *cacheValidation) get(name string) (Tup2, bool) {
	c.RLock()
	defer c.RUnlock()

	if len(c.Cache) > 0 {
		v, ok := c.Cache[name]
		return v, ok
	}

	return Tup2{}, false
}

func (c *cacheValidation) generateIfNotExist(name, lang string, customValidationRegisterFunctions []validation.CustomValidationRegisterFunc) Tup2 {
	c.Lock()
	defer c.Unlock()

	validate, trans := validation.NewValidator(lang)
	if len(customValidationRegisterFunctions) > 0 {
		for _, f := range customValidationRegisterFunctions {
			validate.RegisterValidation(f.Tag, f.Func, f.CallValidationEvenIfNull)
		}
	}
	nTp := Tup2{
		Validate: validate,
		Trans:    trans,
	}

	if nil == c.Cache {
		c.Cache = make(map[string]Tup2)
	}

	if _, ok := c.Cache[name]; !ok {
		c.Cache[name] = nTp
	}

	return nTp
}

// HandlerInterface is base handler interface
type HandlerInterface interface {
	PreHandle(request interface{}) error
	Validation(request interface{}) error
	SetCombineErrors(combineErrors bool)
	SetCustomValidationRegisterFunctions(customValidationRegisterFunctions []validation.CustomValidationRegisterFunc)
	SetLang(lang string)
	SetRemoteCall(callInc remote.CallInc)
	SetContext(ctx context.Context)
	SetTopicAttributes(topicAttributes map[string]string)
	SetRequestHeader(requestHeader map[string]string)
	SetServiceConfig(serviceConfig *config.Service)
	GetRequestHeader() map[string]string
	GetResponseHeader() map[string]string
	SetExtConfigs(map[string]interface{})
	SetBody(body []byte)
	IsDiscardResponse() bool
	GetCurrentSU() string
}

// Handler is a implement of HandlerInterface and provides common methods for service handler
type Handler struct {
	Lang string
	//Client                            client.Client
	Ctx                               context.Context
	CombineErrors                     bool
	DiscardResponseFlag               bool
	topicAttributes                   map[string]string
	requestHeader                     map[string]string
	requestHeaderRWLock               sync.RWMutex
	responseHeader                    map[string]string
	responseHeaderRWLock              sync.RWMutex
	ServiceConfig                     *config.Service
	Body                              []byte
	RemoteCallInc                     remote.CallInc
	CustomValidationRegisterFunctions []validation.CustomValidationRegisterFunc
}

// PreHandle is execute before handler logic
// service can modify the request in the PreHandle
func (h *Handler) PreHandle(request interface{}) error {
	return nil
}

// Validation is used to request validation if necessary
func (h Handler) Validation(request interface{}) error {
	if nil == request {
		return nil
	}
	if nil == h.ServiceConfig {
		log.Debugf(h.Ctx, "Service config is empty, skip execute the validation")
		return nil
	}

	cvName := h.ServiceConfig.ServiceID + h.Lang

	var validate *validator.Validate
	var trans ut.Translator

	if tp, ok := cv.get(cvName); !ok {
		nTp := cv.generateIfNotExist(cvName, h.Lang, h.CustomValidationRegisterFunctions)
		validate, trans = nTp.Validate, nTp.Trans
	} else {
		validate, trans = tp.Validate, tp.Trans
	}

	errorsStr := ""
	if reflect.TypeOf(request).Kind() == reflect.Slice {
		s := reflect.ValueOf(request)
		for i := 0; i < s.Len(); i++ {
			r := s.Index(i)
			isExistErrors := false
			if err := validate.Struct(r.Interface()); err != nil {
				isExistErrors = true
				for _, err := range err.(validator.ValidationErrors) {
					if !h.CombineErrors {
						errorsStr = err.Translate(trans)
						break
					} else {
						errorsStr += err.Translate(trans) + ", "
					}
				}
				if strings.HasSuffix(errorsStr, ", ") {
					errorsStr = errorsStr[0 : len(errorsStr)-2]
				}
			}
			if isExistErrors {
				errorsStr += fmt.Sprintf(" (slice index:%d)", i)
				if i != s.Len()-1 {
					errorsStr += " | "
				}
			}
		}

		if len(errorsStr) > 0 {
			return errors.New(constant.ValidationError, errorsStr)
		}
	} else {
		if err := validate.Struct(request); err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				if !h.CombineErrors {
					errorsStr = err.Translate(trans)
					break
				} else {
					errorsStr += err.Translate(trans) + ", "
				}
			}
			if strings.HasSuffix(errorsStr, ", ") {
				errorsStr = errorsStr[0 : len(errorsStr)-2]
			}

			return errors.New(constant.ValidationError, errorsStr)
		}
	}

	return nil
}

// SetExtConfigs is used to set additional config
func (h *Handler) SetExtConfigs(extConfigs map[string]interface{}) {
	// init service info config
	if v, ok := extConfigs[constant.ExtConfigService]; ok {
		h.ServiceConfig = v.(*config.Service)
	}
}

// SetTopicAttributes is used to set topic attributes
func (h *Handler) SetTopicAttributes(topicAttributes map[string]string) {
	h.topicAttributes = topicAttributes
}

// SetLang is used to set user lang
func (h *Handler) SetLang(lang string) {
	h.Lang = lang
}

// SetCustomValidationRegisterFunctions is used to register customize validation function
func (h *Handler) SetCustomValidationRegisterFunctions(customValidationRegisterFunctions []validation.CustomValidationRegisterFunc) {
	h.CustomValidationRegisterFunctions = customValidationRegisterFunctions
}

// SetRemoteCall is used to set remote call instance
func (h *Handler) SetRemoteCall(callInc remote.CallInc) {
	h.RemoteCallInc = callInc
}

// SetContext is used to set context
func (h *Handler) SetContext(ctx context.Context) {
	h.Ctx = ctx
}

// SetBody is used to set request body
func (h *Handler) SetBody(body []byte) {
	h.Body = body
}

// SetCombineErrors is used to mark whether combine error or not when validation failed.
func (h *Handler) SetCombineErrors(combineErrors bool) {
	h.CombineErrors = combineErrors
}

// SetRequestHeader is used to set request header
func (h *Handler) SetRequestHeader(requestHeader map[string]string) {
	h.requestHeader = requestHeader
}

// SetServiceConfig is used to set service config
func (h *Handler) SetServiceConfig(serviceConfig *config.Service) {
	h.ServiceConfig = serviceConfig
}

// GetRequestHeader is used to get request header
func (h *Handler) GetRequestHeader() map[string]string {
	return h.requestHeader
}

// IsDiscardResponse is used to check whether discard response
func (h *Handler) IsDiscardResponse() bool {
	return h.DiscardResponseFlag
}

// SetResponseHeader is used to set response header
func (h *Handler) SetResponseHeader(responseHeader map[string]string) {
	h.responseHeader = responseHeader
}

// GetResponseHeader is used to get response header
func (h *Handler) GetResponseHeader() map[string]string {
	return h.responseHeader
}

//AddResponseHeader is used to add pair of key-value to the response header
func (h *Handler) AddResponseHeader(key, value string) {
	h.responseHeaderRWLock.Lock()
	defer func() { h.responseHeaderRWLock.Unlock() }()

	if nil == h.responseHeader {
		h.responseHeader = make(map[string]string)
	}

	h.responseHeader[key] = value
}

// DiscardResponse is used to mark to discard the response
func (h *Handler) DiscardResponse() {
	h.DiscardResponseFlag = true
}

// GetTopicAttributes returns the topic attributes of message
func (h Handler) GetTopicAttributes() map[string]string {
	return h.topicAttributes
}

// GetMsgTopicType returns the topic type of message
func (h Handler) GetMsgTopicType() string {
	if nil == h.topicAttributes {
		return ""
	}

	return h.topicAttributes[constant.TopicType]
}

// GetMsgTopicID returns the topic id of message
func (h Handler) GetMsgTopicID() string {
	if nil == h.topicAttributes {
		return ""
	}

	return h.topicAttributes[constant.TopicID]
}

// GetSourceORG returns the source organization of message
func (h Handler) GetSourceORG() string {
	if nil == h.topicAttributes {
		return ""
	}

	return h.topicAttributes[constant.TopicSourceORG]
}

// GetSourceWKS returns the source workspace of message
func (h Handler) GetSourceWKS() string {
	if nil == h.topicAttributes {
		return ""
	}

	return h.topicAttributes[constant.TopicSourceWorkspace]
}

// GetSourceENV returns the source environment of message
func (h Handler) GetSourceENV() string {
	if nil == h.topicAttributes {
		return ""
	}

	return h.topicAttributes[constant.TopicSourceEnvironment]
}

// GetSourceAZ returns the source available zone of message
func (h Handler) GetSourceAZ() string {
	if nil == h.topicAttributes {
		return ""
	}

	return h.topicAttributes[constant.TopicSourceAZ]
}

// GetSourceServiceID returns the source service ID of message
func (h Handler) GetSourceServiceID() string {
	if nil == h.topicAttributes {
		return ""
	}

	return h.topicAttributes[constant.TopicSourceServiceID]
}

// GetSourceSU returns the source SU of message
func (h Handler) GetSourceSU() string {
	if nil == h.topicAttributes {
		return ""
	}

	return h.topicAttributes[constant.TopicSourceSU]
}

// GetSourceNodeID returns the source node ID of message
func (h Handler) GetSourceNodeID() string {
	if nil == h.topicAttributes {
		return ""
	}

	return h.topicAttributes[constant.TopicSourceNodeID]
}

// GetSourceInstanceID returns the source instance ID of message
func (h Handler) GetSourceInstanceID() string {
	if nil == h.topicAttributes {
		return ""
	}

	return h.topicAttributes[constant.TopicSourceInstanceID]
}

// GetCurrentSU returns the SU of current request
func (h Handler) GetCurrentSU() string {
	if len(h.requestHeader) != 0 {
		if su, ok := h.requestHeader[constant.CurrentSU]; ok {
			return su
		}
	}

	if len(h.topicAttributes) == 0 {
		return h.ServiceConfig.Su
	}

	su := util.GetEither(h.topicAttributes, constant.TopicDestinationSU, constant.TopicDestinationDCN)
	if "" == su && nil != h.ServiceConfig {
		su = h.ServiceConfig.Su
	}
	return su
}

// GetHandlerContexts is used to get handler context
func (h Handler) GetHandlerContexts() *contexts.HandlerContexts {
	return contexts.HandlerContextsFromContext(h.Ctx)
}

// RequestBody returns the request body of message
func (h Handler) RequestBody() []byte {
	return h.Body
}

// RangeRequestHeader is used to range request with a closure function
func (h Handler) RangeRequestHeader(f func(string, string)) {
	h.requestHeaderRWLock.Lock()
	defer func() { h.requestHeaderRWLock.Unlock() }()

	if nil != h.requestHeader {
		for k, v := range h.requestHeader {
			f(k, v)
		}
	}
}

// CloneRequestHeader returns a new map[string]string that clone from request header
func (h Handler) CloneRequestHeader() map[string]string {
	h.requestHeaderRWLock.RLock()
	defer func() { h.requestHeaderRWLock.RUnlock() }()

	if nil == h.requestHeader {
		return nil
	}
	cloneMaps := make(map[string]string)
	for k, v := range h.requestHeader {
		cloneMaps[k] = v
	}
	return cloneMaps

}

// GetRequestHeaderValWithKey is used to get request header value by key
// the first value (res) is assigned the value. If that key doesn't exist, res is the value type's empty string ("").
// The second value (ok) is a bool that is true if the key exists in the map, and false if not
func (h Handler) GetRequestHeaderValWithKey(key string) (res string, ok bool) {
	h.requestHeaderRWLock.RLock()
	defer func() { h.requestHeaderRWLock.RUnlock() }()

	if nil == h.requestHeader {
		return "", false
	}

	r, b := h.requestHeader[key]

	return r, b
}

// GetRequestHeaderValWithKeyEitherSilence finds the request header value according to the two keys in order
// and return immediately after finding the first value
// returns empty string ("") if both key doesn't exist in the request header
func (h Handler) GetRequestHeaderValWithKeyEitherSilence(key1, key2 string) string {
	h.requestHeaderRWLock.RLock()
	defer func() { h.requestHeaderRWLock.RUnlock() }()

	if nil == h.requestHeader {
		return ""
	}

	v, ok := h.requestHeader[key1]
	if ok {
		return v
	}

	return h.requestHeader[key2]
}

// GetTraceID returns the trace ID from header
func (h Handler) GetTraceID() string {
	if nil != h.Ctx {
		handlerContexts := contexts.HandlerContextsFromContext(h.Ctx)
		if nil != handlerContexts && nil != handlerContexts.SpanContexts {
			return handlerContexts.SpanContexts.TraceID
		}
	}

	return ""
}

// GetRequestHeaderValWithKeySilence is used to get request header value by key
// returns empty string ("") if the key doesn't exist
func (h Handler) GetRequestHeaderValWithKeySilence(key string) string {
	h.requestHeaderRWLock.RLock()
	defer func() { h.requestHeaderRWLock.RUnlock() }()

	if nil == h.requestHeader {
		return ""
	}
	return h.requestHeader[key]
}

// GetRequestHeaderValWithKeyIgnoreCase is used to get request header value and ignore case the key
// the first value (res) is assigned the value. If that key doesn't exist, res is the value type's empty string ("").
// The second value (ok) is a bool that is true if the key exists in the map, and false if not
func (h Handler) GetRequestHeaderValWithKeyIgnoreCase(key string) (res string, ok bool) {
	h.requestHeaderRWLock.RLock()
	defer func() { h.requestHeaderRWLock.RUnlock() }()

	if nil != h.requestHeader {
		for k, v := range h.requestHeader {
			if strings.EqualFold(k, key) {
				return v, true
			}
		}
	}

	return "", false
}

// GetRequestHeaderValWithKeyIgnoreCaseSilence is used to get request header value and ignore case the key
func (h Handler) GetRequestHeaderValWithKeyIgnoreCaseSilence(key string) string {
	h.requestHeaderRWLock.RLock()
	defer func() { h.requestHeaderRWLock.RUnlock() }()

	if nil != h.requestHeader {
		for k, v := range h.requestHeader {
			if strings.EqualFold(k, key) {
				return v
			}
		}
	}

	return ""
}

// DeleteRequestHeader is used to delete request header by key
func (h *Handler) DeleteRequestHeader(key string) {
	h.requestHeaderRWLock.Lock()
	defer func() { h.requestHeaderRWLock.Unlock() }()

	if nil == h.requestHeader {
		return
	}

	delete(h.requestHeader, key)
}

// IsRequestHeaderExistsKey is used to check whether the key exists in the request header
func (h Handler) IsRequestHeaderExistsKey(key string) bool {
	h.requestHeaderRWLock.RLock()
	defer func() { h.requestHeaderRWLock.RUnlock() }()
	if nil == h.requestHeader {
		return false
	}
	_, ok := h.requestHeader[key]
	return ok
}

// IsRequestHeaderExistsKeyIgnoreCase is used to check whether the key exists in the request header
// and ignore case the key
func (h Handler) IsRequestHeaderExistsKeyIgnoreCase(key string) bool {
	h.requestHeaderRWLock.RLock()
	defer func() { h.requestHeaderRWLock.RUnlock() }()
	if nil == h.requestHeader {
		return false
	}

	for k := range h.requestHeader {
		if strings.EqualFold(k, key) {
			return true
		}
	}

	return false
}

// RequestHeaderToString is used to format the map to a string
// returns empty string ("") if the key doesn't exist
func (h Handler) RequestHeaderToString() string {
	h.requestHeaderRWLock.RLock()
	defer func() { h.requestHeaderRWLock.RUnlock() }()

	if nil == h.requestHeader {
		return ""
	}

	return util.MapToString(h.requestHeader)
}

// RangeResponseHeader is used to range the response header with a closure
func (h Handler) RangeResponseHeader(f func(string, string)) {
	h.responseHeaderRWLock.Lock()
	defer func() { h.responseHeaderRWLock.Unlock() }()

	if nil != h.responseHeader {
		for k, v := range h.responseHeader {
			f(k, v)
		}
	}
}

// CloneResponseHeader returns a new map[string]string that clone from response header
// returns empty string ("") if the key doesn't exist
func (h Handler) CloneResponseHeader() map[string]string {
	h.responseHeaderRWLock.RLock()
	defer func() { h.responseHeaderRWLock.RUnlock() }()

	if nil == h.responseHeader {
		return nil
	}
	cloneMaps := make(map[string]string)
	for k, v := range h.responseHeader {
		cloneMaps[k] = v
	}
	return cloneMaps

}

// GetResponseHeaderValWithKey is used to get value from the response header
// the first value (res) is assigned the value. If that key doesn't exist, res is the value type's empty string ("").
// The second value (ok) is a bool that is true if the key exists in the map, and false if not
func (h Handler) GetResponseHeaderValWithKey(key string) (res string, ok bool) {
	h.responseHeaderRWLock.RLock()
	defer func() { h.responseHeaderRWLock.RUnlock() }()

	if nil == h.responseHeader {
		return "", false
	}

	r, b := h.responseHeader[key]

	return r, b
}

// GetResponseHeaderValWithKeyIgnoreCase is used to get value from the response header and ignore case the key
// the first value (res) is assigned the value. If that key doesn't exist, res is the value type's empty string ("").
// The second value (ok) is a bool that is true if the key exists in the map, and false if not
func (h Handler) GetResponseHeaderValWithKeyIgnoreCase(key string) (res string, ok bool) {
	h.responseHeaderRWLock.RLock()
	defer func() { h.responseHeaderRWLock.RUnlock() }()

	if nil != h.responseHeader {
		for k, v := range h.responseHeader {
			if strings.EqualFold(k, key) {
				return v, true
			}
		}
	}

	return "", false
}

// GetResponseHeaderValWithKeyEitherSilence finds the response header value according to the two keys in order
// and return immediately after finding the first value
// returns empty string ("") if the both key doesn't exist in the response header
func (h Handler) GetResponseHeaderValWithKeyEitherSilence(key1, key2 string) string {
	h.responseHeaderRWLock.RLock()
	defer func() { h.responseHeaderRWLock.RUnlock() }()

	if nil == h.responseHeader {
		return ""
	}

	if v, ok := h.responseHeader[key1]; ok {
		return v
	}

	return h.responseHeader[key2]
}

// GetResponseHeaderValWithKeySilence finds the response header value according to the key
// returns empty string("") if the key doesn't exists in the response header
func (h Handler) GetResponseHeaderValWithKeySilence(key string) string {
	h.responseHeaderRWLock.RLock()
	defer func() { h.responseHeaderRWLock.RUnlock() }()

	if nil == h.responseHeader {
		return ""
	}
	return h.responseHeader[key]
}

// DeleteResponseHeader is used to delete response header value according to the key
func (h *Handler) DeleteResponseHeader(key string) {
	h.responseHeaderRWLock.Lock()
	defer func() { h.responseHeaderRWLock.Unlock() }()

	if nil == h.responseHeader {
		return
	}

	delete(h.responseHeader, key)
}

// ResponseHeaderToString returns the formatted response header string
func (h Handler) ResponseHeaderToString() string {
	h.responseHeaderRWLock.RLock()
	defer func() { h.responseHeaderRWLock.RUnlock() }()
	if nil == h.responseHeader {
		return ""
	}
	return util.MapToString(h.responseHeader)
}

func deleteCurrentSU(options *client.RequestOptions) {
	if nil != options && nil != options.Header {
		options.HeaderLock.Lock()
		defer options.HeaderLock.Unlock()

		delete(options.Header, constant.CurrentSU)
	}
}

func removeDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (h *Handler) fillingPassThroughHeaderKeyListIfNecessary(requestOptions *client.RequestOptions, serviceKeySlice ...string) {
	handlerContexts := h.GetHandlerContexts()
	if nil != handlerContexts {
		if len(handlerContexts.SpanContexts.PassThroughHeaderKeyList) > 0 {
			requestOptions.PassThroughHeaderKeyList = append(requestOptions.PassThroughHeaderKeyList, strings.Split(handlerContexts.SpanContexts.PassThroughHeaderKeyList, ",")...)
		}
		if len(serviceKeySlice) > 0 {
			serviceKey := serviceKeySlice[0]
			requestOptions.PassThroughHeaderKeyList = append(requestOptions.PassThroughHeaderKeyList, h.getPassThroughKeyList(handlerContexts.DownstreamConfigs, serviceKey)...)
		}
	}

	if len(requestOptions.PassThroughHeaderKeyList) > 0 {
		requestOptions.PassThroughHeaderKeyList = removeDuplicateValues(requestOptions.PassThroughHeaderKeyList)
		opt := mesh.WithOriginalHeader(h.CloneRequestHeader())
		opt(requestOptions)
	}
}

// SyncCallw Synchronous call, before requesting downstream services,
// it will lookup to the target SU in GLS according to the shard type and shard ID,
// which is suitable for service calls that require GLS lookup
//
// ctx carries deadlines, cancellation signals,and other request-scoped values across API boundaries and between processes
// elementType shard type
// elementID shard ID
// serviceKey downstream service key, which needs to correspond to the key in [downstream] in the configuration
// request request for downstream services, the type is client.Request
// response the response type when the downstream successfully returns. It must be a pointer
// opts optional parameters, one or more request parameters can be set or overridden
func (h *Handler) SyncCallw(ctx context.Context, elementType, elementID, serviceKey string, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error) {
	requestOptions := request.RequestOptions()
	deleteCurrentSU(requestOptions)

	opt := mesh.WithServiceConfig(h.ServiceConfig)
	opt(requestOptions)

	h.fillingPassThroughHeaderKeyListIfNecessary(requestOptions, serviceKey)

	if nil != opts {
		return h.RemoteCallInc.SyncCallw(ctx, elementType, elementID, serviceKey, request, response, opts...)
	} else {
		return h.RemoteCallInc.SyncCallw(ctx, elementType, elementID, serviceKey, request, response)
	}
}

// AsyncCallw Asynchronous call, before requesting downstream services,
// it will lookup to the target SU in GLS according to the shard type and shard ID,
// which is suitable for service calls that require GLS lookup
//
// ctx carries deadlines, cancellation signals,and other request-scoped values across API boundaries and between processes
// elementType shard type
// elementID shard ID
// serviceKey downstream service key, which needs to correspond to the key in [downstream] in the configuration
// request request for downstream services, the type is client.Request
// opts optional parameters, one or more request parameters can be set or overridden
func (h *Handler) AsyncCallw(ctx context.Context, elementType, elementID, serviceKey string, request client.Request, opts ...client.CallOption) *errors.Error {
	requestOptions := request.RequestOptions()
	deleteCurrentSU(requestOptions)

	opt := mesh.WithServiceConfig(h.ServiceConfig)
	opt(requestOptions)

	h.fillingPassThroughHeaderKeyListIfNecessary(requestOptions, serviceKey)

	if nil != opts {
		return h.RemoteCallInc.AsyncCallw(ctx, elementType, elementID, serviceKey, request, opts...)
	} else {
		return h.RemoteCallInc.AsyncCallw(ctx, elementType, elementID, serviceKey, request)
	}
}

// SemiSyncCallw Semi-synchronous call, before requesting downstream services,
// it will lookup to the target SU in GLS according to the shard type and shard ID,
// which is suitable for service calls that require GLS lookup
//
// ctx carries deadlines, cancellation signals,and other request-scoped values across API boundaries and between processes
// elementType shard type
// elementID shard ID
// serviceKey downstream service key, which needs to correspond to the key in [downstream] in the configuration
// request request for downstream services, the type is client.Request
// response the response type when the downstream successfully returns. It must be a pointer
// opts optional parameters, one or more request parameters can be set or overridden
func (h *Handler) SemiSyncCallw(ctx context.Context, elementType, elementID, serviceKey string, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error) {
	requestOptions := request.RequestOptions()
	deleteCurrentSU(requestOptions)

	opt := mesh.WithServiceConfig(h.ServiceConfig)
	opt(requestOptions)

	h.fillingPassThroughHeaderKeyListIfNecessary(requestOptions, serviceKey)

	if nil != opts {
		return h.RemoteCallInc.SemiSyncCallw(ctx, elementType, elementID, serviceKey, request, response, opts...)
	} else {
		return h.RemoteCallInc.SemiSyncCallw(ctx, elementType, elementID, serviceKey, request, response)
	}
}

// SyncCalls Synchronous call, before requesting downstream services,
//
// ctx carries deadlines, cancellation signals,and other request-scoped values across API boundaries and between processes
// request request for downstream services, the type is client.Request
// response the response type when the downstream successfully returns. It must be a pointer
// opts optional parameters, one or more request parameters can be set or overridden
func (h *Handler) SyncCalls(ctx context.Context, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error) {
	requestOptions := request.RequestOptions()
	deleteCurrentSU(requestOptions)

	opt := mesh.WithServiceConfig(h.ServiceConfig)
	opt(requestOptions)

	h.fillingPassThroughHeaderKeyListIfNecessary(requestOptions)

	if nil != opts {
		return h.RemoteCallInc.SyncCalls(ctx, request, response, opts...)
	} else {
		return h.RemoteCallInc.SyncCalls(ctx, request, response)
	}
}

// AsyncCalls Asynchronous call, before requesting downstream services,
//
// ctx carries deadlines, cancellation signals,and other request-scoped values across API boundaries and between processes
// request request for downstream services, the type is client.Request
// opts optional parameters, one or more request parameters can be set or overridden
func (h *Handler) AsyncCalls(ctx context.Context, request client.Request, opts ...client.CallOption) *errors.Error {
	requestOptions := request.RequestOptions()
	deleteCurrentSU(requestOptions)

	opt := mesh.WithServiceConfig(h.ServiceConfig)
	opt(requestOptions)

	h.fillingPassThroughHeaderKeyListIfNecessary(requestOptions)

	if nil != opts {
		return h.RemoteCallInc.AsyncCalls(ctx, request, opts...)
	} else {
		return h.RemoteCallInc.AsyncCalls(ctx, request)
	}
}

// SemiSyncCalls Semi-synchronous call, before requesting downstream services,
//
// ctx carries deadlines, cancellation signals,and other request-scoped values across API boundaries and between processes
// request request for downstream services, the type is client.Request
// response the response type when the downstream successfully returns. It must be a pointer
// opts optional parameters, one or more request parameters can be set or overridden
func (h *Handler) SemiSyncCalls(ctx context.Context, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error) {
	requestOptions := request.RequestOptions()
	deleteCurrentSU(requestOptions)

	opt := mesh.WithServiceConfig(h.ServiceConfig)
	opt(requestOptions)

	h.fillingPassThroughHeaderKeyListIfNecessary(requestOptions)

	if nil != opts {
		return h.RemoteCallInc.SemiSyncCalls(ctx, request, response, opts...)
	} else {
		return h.RemoteCallInc.SemiSyncCalls(ctx, request, response)
	}
}

// SyncCall Asynchronous call, before requesting downstream services,
// which is suitable for service calls with destination SU direct.
//
// ctx carries deadlines, cancellation signals,and other request-scoped values across API boundaries and between processes
// dstSu destination SU
// serviceKey downstream service key, which needs to correspond to the key in [downstream] in the configuration
// request request for downstream services, the type is client.Request
// response the response type when the downstream successfully returns. It must be a pointer
// opts optional parameters, one or more request parameters can be set or overridden
func (h *Handler) SyncCall(ctx context.Context, dstSu, serviceKey string, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error) {
	requestOptions := request.RequestOptions()
	deleteCurrentSU(requestOptions)

	opt := mesh.WithServiceConfig(h.ServiceConfig)
	opt(requestOptions)

	h.fillingPassThroughHeaderKeyListIfNecessary(requestOptions, serviceKey)

	if nil != opts {
		return h.RemoteCallInc.SyncCall(ctx, dstSu, serviceKey, request, response, opts...)
	} else {
		return h.RemoteCallInc.SyncCall(ctx, dstSu, serviceKey, request, response)
	}
}

// AsyncCall Asynchronous call, before requesting downstream services,
// which is suitable for service calls with destination SU direct.
//
// ctx carries deadlines, cancellation signals,and other request-scoped values across API boundaries and between processes
// dstSu destination SU
// serviceKey downstream service key, which needs to correspond to the key in [downstream] in the configuration
// request request for downstream services, the type is client.Request
// opts optional parameters, one or more request parameters can be set or overridden
func (h *Handler) AsyncCall(ctx context.Context, dstSu, serviceKey string, request client.Request, opts ...client.CallOption) *errors.Error {
	requestOptions := request.RequestOptions()
	deleteCurrentSU(requestOptions)

	opt := mesh.WithServiceConfig(h.ServiceConfig)
	opt(requestOptions)

	h.fillingPassThroughHeaderKeyListIfNecessary(requestOptions, serviceKey)

	if nil != opts {
		return h.RemoteCallInc.AsyncCall(ctx, dstSu, serviceKey, request, opts...)
	} else {
		return h.RemoteCallInc.AsyncCall(ctx, dstSu, serviceKey, request)
	}
}

// SemiSyncCall Semi-synchronous call, before requesting downstream services,
// which is suitable for service calls with destination SU direct.
//
// ctx carries deadlines, cancellation signals,and other request-scoped values across API boundaries and between processes
// dstSu destination SU
// serviceKey downstream service key, which needs to correspond to the key in [downstream] in the configuration
// request request for downstream services, the type is client.Request
// response the response type when the downstream successfully returns. It must be a pointer
// opts optional parameters, one or more request parameters can be set or overridden
func (h *Handler) SemiSyncCall(ctx context.Context, dstSu, serviceKey string, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error) {
	requestOptions := request.RequestOptions()
	deleteCurrentSU(requestOptions)

	opt := mesh.WithServiceConfig(h.ServiceConfig)
	opt(requestOptions)

	h.fillingPassThroughHeaderKeyListIfNecessary(requestOptions, serviceKey)

	if nil != opts {
		return h.RemoteCallInc.SemiSyncCall(ctx, dstSu, serviceKey, request, response, opts...)
	} else {
		return h.RemoteCallInc.SemiSyncCall(ctx, dstSu, serviceKey, request, response)
	}
}

// ReplyTo is used to reply response semi-synchronous caller
//
// ctx carries deadlines, cancellation signals,and other request-scoped values across API boundaries and between processes
// response response for reply message, the type is client.Response
func (h *Handler) ReplyTo(ctx context.Context, response client.Response) *errors.Error {
	return h.RemoteCallInc.ReplyTo(ctx, response)
}

func (h *Handler) AddRootKVToSecondStageMethod(key, value string) {
	newKey := constant.RootKVToSecondStageKeyPrefix + key
	h.AddResponseHeader(newKey, value)
}

func (h *Handler) getPassThroughKeyList(downstreamConfigs map[string]config.Downstream, serviceKey string) []string {
	if nil != downstreamConfigs {
		serviceConfig := getDownstreamServiceConfig(downstreamConfigs, serviceKey)
		if nil == serviceConfig {
			return nil
		}

		return serviceConfig.PassThroughHeaderKey.List
	}

	return nil
}

func getDownstreamServiceConfig(m map[string]config.Downstream, key string) *config.Downstream {
	if v1, ok := m[key]; ok {
		return &v1
	} else if v2, ok := m[strings.ToLower(key)]; ok {
		return &v2
	}

	return nil
}
