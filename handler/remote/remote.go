package remote

import (
	"context"
	"reflect"
	"strings"
	"time"

	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/client/mesh"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/log"
	"git.multiverse.io/eventkit/kit/sed/callback"
	"github.com/afex/hystrix-go/hystrix"
)

// CallInc defines the interface of the remote call instance
type CallInc interface {
	SyncCalls(ctx context.Context, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error)
	AsyncCalls(ctx context.Context, request client.Request, opts ...client.CallOption) *errors.Error
	SemiSyncCalls(ctx context.Context, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error)

	SyncCallw(ctx context.Context, elementType, elementID, serviceKey string, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error)
	AsyncCallw(ctx context.Context, elementType, elementID, serviceKey string, request client.Request, opts ...client.CallOption) *errors.Error
	SemiSyncCallw(ctx context.Context, elementType, elementID, serviceKey string, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error)

	SyncCall(ctx context.Context, dstSU, serviceKey string, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error)
	AsyncCall(ctx context.Context, dstSU, serviceKey string, request client.Request, opts ...client.CallOption) *errors.Error
	SemiSyncCall(ctx context.Context, dstSU, serviceKey string, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error)

	ReplyTo(ctx context.Context, response client.Response) *errors.Error

	//GetDownstreamServiceConfig() map[string]config.Downstream
}

type defaultRemoteCallImpl struct {
	Client                  client.Client
	CallbackExecutor        callback.Executor
	IsLocalCallCheckFunc    func(eventId string) bool
	DownstreamServiceConfig map[string]config.Downstream
}

//func (h *defaultRemoteCallImpl) GetDownstreamServiceConfig() map[string]config.Downstream {
//	return h.DownstreamServiceConfig
//}

func (h *defaultRemoteCallImpl) ReplyTo(ctx context.Context, response client.Response) *errors.Error {
	handlerContexts := contexts.HandlerContextsFromContext(ctx)
	if util.IsNil(handlerContexts) || util.IsNil(handlerContexts.SpanContexts) || len(handlerContexts.SpanContexts.ReplyToAddress) == 0 {
		return errors.Errorf(constant.SystemInternalError,
			"Failed to reply semi-sync call, cannot get reply to address, handlerContexts:[%++v], response:[%++v]",
			handlerContexts, response)
	}
	response.WithOptions(mesh.WithReplyToAddress(handlerContexts.SpanContexts.ReplyToAddress))

	e := h.Client.ReplySemiSyncCall(ctx, response)
	if e != nil {
		var err *errors.Error

		switch e.(type) {
		case *errors.Error:
			err = e.(*errors.Error)
		case error:
			err = errors.Wrap(constant.SystemInternalError, e.(error), 0)
		default:
			err = errors.Errorf(constant.SystemInternalError, "%++v", e)
		}
		return err
	}

	return nil
}

func (h *defaultRemoteCallImpl) SemiSyncCalls(ctx context.Context, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error) {
	request.WithOptions(mesh.WithSemiSyncCall())

	return h.SyncCalls(ctx, request, response, opts...)
}

func (h *defaultRemoteCallImpl) SemiSyncCallw(ctx context.Context, elementType, elementID, serviceKey string, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error) {
	request.WithOptions(mesh.WithSemiSyncCall())

	return h.SyncCallw(ctx, elementType, elementID, serviceKey, request, response, opts...)
}

func (h *defaultRemoteCallImpl) SemiSyncCall(ctx context.Context, dstSU, serviceKey string, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error) {
	request.WithOptions(mesh.WithSemiSyncCall())

	return h.SyncCall(ctx, dstSU, serviceKey, request, response, opts...)
}

// NewDefaultRemoteCall creates a new remote call instance with default configuration
func NewDefaultRemoteCall(client client.Client,
	callbackExecutor callback.Executor,
	isLocalCallCheckFunc func(eventId string) bool,
	downstreamServiceConfig map[string]config.Downstream) CallInc {
	return &defaultRemoteCallImpl{
		client,
		callbackExecutor,
		isLocalCallCheckFunc,
		downstreamServiceConfig,
	}
}

func getDownstreamServiceConfig(m map[string]config.Downstream, key string) *config.Downstream {
	if v1, ok := m[key]; ok {
		return &v1
	} else if v2, ok := m[strings.ToLower(key)]; ok {
		return &v2
	}

	return nil
}

type OutputPair struct {
	ResponseMeta  client.ResponseMeta
	ResponseError error
}

func (h *defaultRemoteCallImpl) SyncCalls(ctx context.Context, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error) {
	eventID := request.RequestOptions().EventID

	var finalCallCtx context.Context

	if request.RequestOptions().SkipResponseAutoParse {
		finalCallCtx = context.WithValue(ctx, constant.SkipResponseAutoParseKeyMappingFlagKey, true)
	} else {
		finalCallCtx = ctx
	}

	if !request.RequestOptions().DisableMacroModel && nil != h.IsLocalCallCheckFunc && h.IsLocalCallCheckFunc(eventID) {
		log.Debugf(finalCallCtx, "target event id[%s] is local call", eventID)
		request.WithOptions(mesh.MarkLocalCall())
	}

	if nil != response {
		vi := reflect.ValueOf(response)
		if vi.Kind() != reflect.Ptr {
			return nil, errors.New(constant.SystemInternalError, "The parameter of `response` must be a pointer!")
		}
	}

	var responseMeta client.ResponseMeta
	var responseError error

	if request.RequestOptions().IsEnabledCircuitBreaker {
		output := make(chan OutputPair, 1)
		doError := hystrix.Do(strings.ToLower(request.RequestOptions().ServiceKey), func() error {
			responseMetaInner, responseErrorInner := h.Client.SyncCall(finalCallCtx, request, response, append(opts, client.WithCallbackExecutor(h.CallbackExecutor))...)
			if nil != responseErrorInner {
				switch responseErrorInner.(type) {
				case *errors.Error:
					err := responseErrorInner.(*errors.Error)
					if constant.SystemCallbackAppTimeout == err.ErrorCode ||
						constant.SystemMeshRequestReplyTimeout == err.ErrorCode ||
						constant.SystemRemoteCallTimeout == err.ErrorCode {

						return hystrix.ErrTimeout
					}
				}
				return responseErrorInner
			}

			output <- OutputPair{
				responseMetaInner,
				responseErrorInner,
			}
			return nil
		}, request.RequestOptions().FallbackFunc)

		if nil != doError {
			responseError = doError
		} else {
			select {
			case out := <-output:
				{
					responseMeta = out.ResponseMeta
					responseError = out.ResponseError
				}
			}
		}
	} else {
		responseMeta, responseError = h.Client.SyncCall(finalCallCtx, request, response, append(opts, client.WithCallbackExecutor(h.CallbackExecutor))...)
	}

	if responseError != nil {
		var err *errors.Error

		switch responseError.(type) {
		case *errors.Error:
			err = responseError.(*errors.Error)
		case error:
			err = errors.Wrap(constant.SystemInternalError, responseError.(error), 0)
		default:
			err = errors.Errorf(constant.SystemInternalError, "%++v", responseError)
		}
		return nil, err
	}

	return responseMeta, nil
}

func (h *defaultRemoteCallImpl) AsyncCalls(ctx context.Context, request client.Request, opts ...client.CallOption) *errors.Error {
	e := h.Client.AsyncCall(ctx, request, opts...)
	if e != nil {
		var err *errors.Error

		switch e.(type) {
		case *errors.Error:
			err = e.(*errors.Error)
		case error:
			err = errors.Wrap(constant.SystemInternalError, e.(error), 0)
		default:
			err = errors.Errorf(constant.SystemInternalError, "%++v", e)
		}
		return err
	}

	return nil
}

func getCommunicateConfigs(downstreamConfigs *config.Downstream) (timeoutMilliseconds int, retryWaitingMilliseconds int,
	maxWaitingTimeMilliseconds int, maxRetryTimes int, deleteTransactionPropagationInfo bool, protoType string,
	httpAddress string, httpMethod string, httpContextType string) {

	if nil == downstreamConfigs {
		return
	}

	timeoutMilliseconds = downstreamConfigs.TimeoutMilliseconds
	retryWaitingMilliseconds = downstreamConfigs.RetryWaitingMilliseconds
	maxWaitingTimeMilliseconds = downstreamConfigs.MaxWaitingTimeMilliseconds
	maxRetryTimes = downstreamConfigs.MaxRetryTimes
	deleteTransactionPropagationInfo = downstreamConfigs.DeleteTransactionPropagationInfo
	protoType = downstreamConfigs.ProtoType
	httpAddress = downstreamConfigs.HTTPAddress
	httpMethod = downstreamConfigs.HTTPMethod
	httpContextType = downstreamConfigs.HTTPContextType

	cc := downstreamConfigs.CustomConfigurations

	if cc.TimeoutMilliseconds > 0 {
		timeoutMilliseconds = cc.TimeoutMilliseconds
	}

	if cc.RetryWaitingMilliseconds > 0 {
		retryWaitingMilliseconds = cc.RetryWaitingMilliseconds
	}

	if cc.MaxWaitingTimeMilliseconds > 0 {
		maxWaitingTimeMilliseconds = cc.MaxWaitingTimeMilliseconds
	}

	if cc.MaxRetryTimes > 0 {
		maxRetryTimes = cc.MaxRetryTimes
	}

	if cc.DeleteTransactionPropagationInfo {
		deleteTransactionPropagationInfo = cc.DeleteTransactionPropagationInfo
	}

	if len(cc.ProtoType) > 0 {
		protoType = cc.ProtoType
	}

	if len(cc.HTTPAddress) > 0 {
		httpAddress = cc.HTTPAddress
	}

	if len(cc.HTTPMethod) > 0 {
		httpMethod = cc.HTTPMethod
	}

	if len(cc.HTTPContextType) > 0 {
		httpContextType = cc.HTTPContextType
	}

	return
}

func (h *defaultRemoteCallImpl) SyncCall(ctx context.Context, dstSU, serviceKey string, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error) {
	serviceConfig := getDownstreamServiceConfig(h.DownstreamServiceConfig, serviceKey)
	if nil == serviceConfig {
		return nil, errors.Errorf(constant.SystemInternalError, "cannot found downstream service config with service key:[%s], please check", serviceKey)
	}
	var responseAutoParseKeyMapping map[string]string
	if len(serviceConfig.ResponseAutoParseKeyMapping) > 0 {
		responseAutoParseKeyMapping = serviceConfig.ResponseAutoParseKeyMapping
	} else {
		handlerContexts := contexts.HandlerContextsFromContext(ctx)
		responseAutoParseKeyMapping = handlerContexts.ResponseAutoParseKeyMapping
	}
	callCtx := context.WithValue(ctx, constant.ResponseAutoParseKeyMappingKey, responseAutoParseKeyMapping)

	log.Debugf(callCtx, "SyncCall, serviceKey=[%s], serviceConfig=[%++v]", serviceKey, serviceConfig)

	timeoutMilliseconds, retryWaitingMilliseconds, maxWaitingTimeMilliseconds, maxRetryTimes, deleteTransactionPropagationInfo,
		protoType, httpAddress, httpMethod, httpContextType := getCommunicateConfigs(serviceConfig)

	if strings.EqualFold(protoType, constant.RequestTypeHTTP) {
		request.WithOptions(
			mesh.WithHTTPRequestInfo(
				httpAddress,
				httpMethod,
				httpContextType,
			),
			//mesh.WithTimeout(time.Duration(timeoutMilliseconds)*time.Millisecond),
			//mesh.WithMaxWaitingTime(time.Duration(maxWaitingTimeMilliseconds)*time.Millisecond),
			//mesh.WithRetryWaitingMilliseconds(time.Duration(retryWaitingMilliseconds)*time.Millisecond),
			mesh.WithVersion(serviceConfig.Version),
			mesh.WithMaxRetryTimes(maxRetryTimes),
			mesh.WithServiceKey(serviceKey),
			mesh.MarkIsEnableCircuitBreaker(serviceConfig.CircuitBreaker.Enable),
			mesh.WithEnableLogging(serviceConfig.EnableLogging),
			mesh.WithDeleteTransactionPropagationInformation(deleteTransactionPropagationInfo),
		)
	} else {
		request.WithOptions(
			mesh.WithTopicType(serviceConfig.EventType),
			mesh.WithSU(dstSU),
			mesh.WithEventID(serviceConfig.EventID),
			//mesh.WithTimeout(time.Duration(timeoutMilliseconds)*time.Millisecond),
			//mesh.WithMaxWaitingTime(time.Duration(maxWaitingTimeMilliseconds)*time.Millisecond),
			//mesh.WithRetryWaitingMilliseconds(time.Duration(retryWaitingMilliseconds)*time.Millisecond),
			mesh.WithVersion(serviceConfig.Version),
			mesh.WithMaxRetryTimes(maxRetryTimes),
			mesh.WithServiceKey(serviceKey),
			mesh.MarkIsEnableCircuitBreaker(serviceConfig.CircuitBreaker.Enable),
			mesh.WithEnableLogging(serviceConfig.EnableLogging),
			mesh.WithDeleteTransactionPropagationInformation(deleteTransactionPropagationInfo),
		)
	}
	if timeoutMilliseconds > 0 {
		request.WithOptions(mesh.WithTimeout(time.Duration(timeoutMilliseconds) * time.Millisecond))
	} else {
		st := ctx.Value(constant.KeyST)
		to3 := ctx.Value(constant.To3)
		if nil != st && nil != to3 {
			startTime := st.(time.Time)
			to3Time := to3.(int)
			currentTime := time.Now()
			tc := currentTime.Sub(startTime).Milliseconds()
			to := int64(to3Time) - tc
			if to > 0 {
				request.WithOptions(mesh.WithMaxWaitingTime(time.Duration(to) * time.Millisecond))
				request.WithOptions(mesh.WithTimeout(time.Duration(to) * time.Millisecond))
			}
		}
	}

	if maxWaitingTimeMilliseconds > 0 {
		request.WithOptions(mesh.WithMaxWaitingTime(time.Duration(maxWaitingTimeMilliseconds) * time.Millisecond))
	}

	if retryWaitingMilliseconds > 0 {
		request.WithOptions(mesh.WithRetryWaitingMilliseconds(time.Duration(retryWaitingMilliseconds) * time.Millisecond))
	}

	return h.SyncCalls(callCtx, request, response, opts...)
}

func (h *defaultRemoteCallImpl) AsyncCall(ctx context.Context, dstSU, serviceKey string, request client.Request, opts ...client.CallOption) *errors.Error {
	serviceConfig := getDownstreamServiceConfig(h.DownstreamServiceConfig, serviceKey)
	if nil == serviceConfig {
		return errors.Errorf(constant.SystemInternalError, "cannot found downstream service config with service key:[%s], please check", serviceKey)
	}

	var responseAutoParseKeyMapping map[string]string
	if len(serviceConfig.ResponseAutoParseKeyMapping) > 0 {
		responseAutoParseKeyMapping = serviceConfig.ResponseAutoParseKeyMapping
	} else {
		handlerContexts := contexts.HandlerContextsFromContext(ctx)
		responseAutoParseKeyMapping = handlerContexts.ResponseAutoParseKeyMapping
	}
	callCtx := context.WithValue(ctx, constant.ResponseAutoParseKeyMappingKey, responseAutoParseKeyMapping)

	log.Debugf(callCtx, "AsyncCall, serviceKey=[%s], serviceConfig=[%++v]", serviceKey, serviceConfig)
	_, _, _, _, deleteTransactionPropagationInfo, protoType, httpAddress, httpMethod, httpContextType := getCommunicateConfigs(serviceConfig)
	if strings.EqualFold(protoType, constant.RequestTypeHTTP) {
		request.WithOptions(
			mesh.WithHTTPRequestInfo(
				httpAddress,
				httpMethod,
				httpContextType,
			),
			mesh.WithVersion(serviceConfig.Version),
			mesh.WithServiceKey(serviceKey),
			mesh.MarkIsEnableCircuitBreaker(serviceConfig.CircuitBreaker.Enable),
			mesh.WithEnableLogging(serviceConfig.EnableLogging),
			mesh.WithDeleteTransactionPropagationInformation(deleteTransactionPropagationInfo),
		)
	} else {
		request.WithOptions(
			mesh.WithTopicType(serviceConfig.EventType),
			mesh.WithSU(dstSU),
			mesh.WithVersion(serviceConfig.Version),
			mesh.WithEventID(serviceConfig.EventID),
			mesh.WithServiceKey(serviceKey),
			mesh.MarkIsEnableCircuitBreaker(serviceConfig.CircuitBreaker.Enable),
			mesh.WithEnableLogging(serviceConfig.EnableLogging),
			mesh.WithDeleteTransactionPropagationInformation(deleteTransactionPropagationInfo),
		)
	}

	return h.AsyncCalls(callCtx, request, opts...)
}

func (h *defaultRemoteCallImpl) SyncCallw(ctx context.Context, elementType, elementID, serviceKey string, request client.Request, response interface{}, opts ...client.CallOption) (client.ResponseMeta, *errors.Error) {
	serviceConfig := getDownstreamServiceConfig(h.DownstreamServiceConfig, serviceKey)
	if nil == serviceConfig {
		return nil, errors.Errorf(constant.SystemInternalError, "cannot found downstream service config with service key:[%s], please check", serviceKey)
	}

	var responseAutoParseKeyMapping map[string]string
	if len(serviceConfig.ResponseAutoParseKeyMapping) > 0 {
		responseAutoParseKeyMapping = serviceConfig.ResponseAutoParseKeyMapping
	} else {
		handlerContexts := contexts.HandlerContextsFromContext(ctx)
		responseAutoParseKeyMapping = handlerContexts.ResponseAutoParseKeyMapping
	}
	callCtx := context.WithValue(ctx, constant.ResponseAutoParseKeyMappingKey, responseAutoParseKeyMapping)

	log.Debugf(callCtx, "SyncCallw, serviceKey=[%s], serviceConfig=[%++v]", serviceKey, serviceConfig)
	timeoutMilliseconds, retryWaitingMilliseconds, maxWaitingTimeMilliseconds, maxRetryTimes, deleteTransactionPropagationInfo,
		protoType, httpAddress, httpMethod, httpContextType := getCommunicateConfigs(serviceConfig)
	if strings.EqualFold(protoType, constant.RequestTypeHTTP) {
		request.WithOptions(
			mesh.WithHTTPRequestInfo(
				httpAddress,
				httpMethod,
				httpContextType,
			),
			//mesh.WithTimeout(time.Duration(timeoutMilliseconds)*time.Millisecond),
			//mesh.WithMaxWaitingTime(time.Duration(maxWaitingTimeMilliseconds)*time.Millisecond),
			//mesh.WithRetryWaitingMilliseconds(time.Duration(retryWaitingMilliseconds)*time.Millisecond),
			mesh.WithVersion(serviceConfig.Version),
			mesh.WithMaxRetryTimes(maxRetryTimes),
			mesh.WithServiceKey(serviceKey),
			mesh.MarkIsEnableCircuitBreaker(serviceConfig.CircuitBreaker.Enable),
			mesh.WithEnableLogging(serviceConfig.EnableLogging),
			mesh.WithDeleteTransactionPropagationInformation(deleteTransactionPropagationInfo),
		)
	} else {
		request.WithOptions(
			mesh.WithTopicType(serviceConfig.EventType),
			mesh.WithElementType(elementType),
			mesh.WithElementID(elementID),
			mesh.WithEventID(serviceConfig.EventID),
			//mesh.WithTimeout(time.Duration(timeoutMilliseconds)*time.Millisecond),
			//mesh.WithMaxWaitingTime(time.Duration(maxWaitingTimeMilliseconds)*time.Millisecond),
			//mesh.WithRetryWaitingMilliseconds(time.Duration(retryWaitingMilliseconds)*time.Millisecond),
			mesh.WithVersion(serviceConfig.Version),
			mesh.WithMaxRetryTimes(maxRetryTimes),
			mesh.WithServiceKey(serviceKey),
			mesh.MarkIsEnableCircuitBreaker(serviceConfig.CircuitBreaker.Enable),
			mesh.WithEnableLogging(serviceConfig.EnableLogging),
			mesh.WithDeleteTransactionPropagationInformation(deleteTransactionPropagationInfo),
		)
	}

	if timeoutMilliseconds > 0 {
		request.WithOptions(mesh.WithTimeout(time.Duration(timeoutMilliseconds) * time.Millisecond))
	} else {
		st := ctx.Value(constant.KeyST)
		to3 := ctx.Value(constant.To3)
		if nil != st && nil != to3 {
			startTime := st.(time.Time)
			to3Time := to3.(int)
			currentTime := time.Now()
			tc := currentTime.Sub(startTime).Milliseconds()
			to := int64(to3Time) - tc
			if to > 0 {
				request.WithOptions(mesh.WithMaxWaitingTime(time.Duration(to) * time.Millisecond))
				request.WithOptions(mesh.WithTimeout(time.Duration(to) * time.Millisecond))
			}
		}
	}

	if maxWaitingTimeMilliseconds > 0 {
		request.WithOptions(mesh.WithMaxWaitingTime(time.Duration(maxWaitingTimeMilliseconds) * time.Millisecond))
	}

	if retryWaitingMilliseconds > 0 {
		request.WithOptions(mesh.WithRetryWaitingMilliseconds(time.Duration(retryWaitingMilliseconds) * time.Millisecond))
	}

	return h.SyncCalls(callCtx, request, response, opts...)
}

func (h *defaultRemoteCallImpl) AsyncCallw(ctx context.Context, elementType, elementID, serviceKey string, request client.Request, opts ...client.CallOption) *errors.Error {
	serviceConfig := getDownstreamServiceConfig(h.DownstreamServiceConfig, serviceKey)
	if nil == serviceConfig {
		return errors.Errorf(constant.SystemInternalError, "cannot found downstream service config with service key:[%s], please check", serviceKey)
	}

	var responseAutoParseKeyMapping map[string]string
	if len(serviceConfig.ResponseAutoParseKeyMapping) > 0 {
		responseAutoParseKeyMapping = serviceConfig.ResponseAutoParseKeyMapping
	} else {
		handlerContexts := contexts.HandlerContextsFromContext(ctx)
		responseAutoParseKeyMapping = handlerContexts.ResponseAutoParseKeyMapping
	}
	callCtx := context.WithValue(ctx, constant.ResponseAutoParseKeyMappingKey, responseAutoParseKeyMapping)

	log.Debugf(callCtx, "AsyncCallw, serviceKey=[%s], serviceConfig=[%++v]", serviceKey, serviceConfig)
	_, _, _, _, deleteTransactionPropagationInfo, protoType, httpAddress, httpMethod, httpContextType := getCommunicateConfigs(serviceConfig)
	if strings.EqualFold(protoType, constant.RequestTypeHTTP) {
		request.WithOptions(
			mesh.WithHTTPRequestInfo(
				httpAddress,
				httpMethod,
				httpContextType,
			),
			mesh.WithVersion(serviceConfig.Version),
			mesh.WithServiceKey(serviceKey),
			mesh.MarkIsEnableCircuitBreaker(serviceConfig.CircuitBreaker.Enable),
			mesh.WithEnableLogging(serviceConfig.EnableLogging),
			mesh.WithDeleteTransactionPropagationInformation(deleteTransactionPropagationInfo),
		)
	} else {
		request.WithOptions(
			mesh.WithTopicType(serviceConfig.EventType),
			mesh.WithElementType(elementType),
			mesh.WithElementID(elementID),
			mesh.WithVersion(serviceConfig.Version),
			mesh.WithEventID(serviceConfig.EventID),
			mesh.WithServiceKey(serviceKey),
			mesh.MarkIsEnableCircuitBreaker(serviceConfig.CircuitBreaker.Enable),
			mesh.WithEnableLogging(serviceConfig.EnableLogging),
			mesh.WithDeleteTransactionPropagationInformation(deleteTransactionPropagationInfo),
		)
	}

	return h.AsyncCalls(callCtx, request, opts...)
}
