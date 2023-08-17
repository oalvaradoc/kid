package mesh

import (
	"context"
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/codec"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/config"
	"time"
)

// WithSessionName sets the session name of request option
func WithSessionName(sessionName string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.SessionName = sessionName
	}
}

// WithOriginalHeader sets the original http header
func WithOriginalHeader(originalHeader map[string]string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.OriginalHeader = originalHeader
	}
}

// WithServiceConfig sets the service config of the client.Request
func WithServiceConfig(serviceConfig *config.Service) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.ServiceConfig = serviceConfig
	}
}

// WithTopicType sets the topic type of the client.Request
func WithTopicType(topicType string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.TopicType = topicType
	}
}

// WithTopicTypeBusiness sets the topic type as constant.TopicTypeBusiness
func WithTopicTypeBusiness() client.RequestOption {
	return WithTopicType(constant.TopicTypeBusiness)
}

// WithTopicTypeOps sets the topic type as constant.TopicTypeOPS
func WithTopicTypeOps() client.RequestOption {
	return WithTopicType(constant.TopicTypeOPS)
}

// WithTopicTypeMetrics sets the topic type as constant.TopicTypeMetrics
func WithTopicTypeMetrics() client.RequestOption {
	return WithTopicType(constant.TopicTypeMetrics)
}

// WithTopicTypeAlert sets the topic type as constant.TopicTypeAlert
func WithTopicTypeAlert() client.RequestOption {
	return WithTopicType(constant.TopicTypeAlert)
}

// WithTopicTypeLog sets the topic type as constant.TopicTypeLog
func WithTopicTypeLog() client.RequestOption {
	return WithTopicType(constant.TopicTypeLog)
}

// WithTopicTypeError sets the topic type as constant.TopicTypeError
func WithTopicTypeError() client.RequestOption {
	return WithTopicType(constant.TopicTypeError)
}

// WithTopicTypeDxc sets the topic type as constant.TopicTypeDXC
func WithTopicTypeDxc() client.RequestOption {
	return WithTopicType(constant.TopicTypeDXC)
}

// WithTopicTypeDts sets the topic type as constant.TopicTypeDTS
func WithTopicTypeDts() client.RequestOption {
	return WithTopicType(constant.TopicTypeDTS)
}

// WithTopicTypeHeartbeat sets the topic type as constant.TopicTypeHeartbeat
func WithTopicTypeHeartbeat() client.RequestOption {
	return WithTopicType(constant.TopicTypeHeartbeat)
}

// WithSemiSyncCall marks the request as a semi-synchronous call
func WithSemiSyncCall() client.RequestOption {
	return func(options *client.RequestOptions) {
		options.IsSemiSyncCall = true
	}
}

// WithORG sets the organization of the request
func WithORG(org string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.Org = org
	}
}

// WithOrgIfEmpty sets the organization of the request if organization parameter is empty
func WithOrgIfEmpty(org string) client.RequestOption {
	return func(options *client.RequestOptions) {
		if "" == options.Org {
			options.Org = org
		}
	}
}

// WithWorkspace sets the workspace of the request
func WithWorkspace(wks string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.Wks = wks
	}
}

// WithWorkspaceIfEmpty sets the workspace if the request if workspace is empty
func WithWorkspaceIfEmpty(wks string) client.RequestOption {
	return func(options *client.RequestOptions) {
		if "" == options.Wks {
			options.Wks = wks
		}
	}
}

// WithEnvironment sets the environment of the request
func WithEnvironment(env string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.Env = env
	}
}

// WithEnvironmentIfEmpty sets the environment of the request if environment is empty
func WithEnvironmentIfEmpty(env string) client.RequestOption {
	return func(options *client.RequestOptions) {
		if "" == options.Env {
			options.Env = env
		}
	}
}

// WithSU sets the SU of the request
func WithSU(su string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.Su = su
	}
}

// WithNodeID sets the node id of the request
func WithNodeID(nodeID string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.NodeID = nodeID
	}
}

// WithInstanceID sets the instance id of the request
func WithInstanceID(instanceID string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.InstanceID = instanceID
	}
}

// WithEventID sets the event id of the request
func WithEventID(eventID string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.EventID = eventID
	}
}

// WithContext sets the context.Context of the request
func WithContext(ctx context.Context) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.Context = ctx
	}
}

// WithHeader sets the header of the request
func WithHeader(header map[string]string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.HeaderLock.Lock()
		defer options.HeaderLock.Unlock()

		if nil == options.Header {
			options.Header = header
		} else {
			for k, v := range header {
				options.Header[k] = v
			}
		}
	}
}

// AddKeyToHeader adds a pair of key-value into request header
func AddKeyToHeader(key, value string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.HeaderLock.Lock()
		defer options.HeaderLock.Unlock()

		if nil == options.Header {
			options.Header = make(map[string]string)
		}
		options.Header[key] = value
	}
}

// WithElementType sets the element type of the request
func WithElementType(elementType string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.HeaderLock.Lock()
		defer options.HeaderLock.Unlock()

		if nil == options.Header {
			options.Header = make(map[string]string)
		}
		options.Header[constant.GlsElementType] = elementType
	}
}

// WithElementID sets the element ID of the request
func WithElementID(elementID string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.HeaderLock.Lock()
		defer options.HeaderLock.Unlock()

		if nil == options.Header {
			options.Header = make(map[string]string)
		}
		options.Header[constant.GlsElementID] = elementID
	}
}

// WithElementClass sets the element class of the request
func WithElementClass(elementClass string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.HeaderLock.Lock()
		defer options.HeaderLock.Unlock()

		if nil == options.Header {
			options.Header = make(map[string]string)
		}
		options.Header[constant.GlsElementClass] = elementClass
	}
}

// WithCodec sets the codec of the request
func WithCodec(codec codec.Codec) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.Codec = codec
	}
}

// WithCodecIfEmpty sets the codec of the request if codec is empty
func WithCodecIfEmpty(codec codec.Codec) client.RequestOption {
	return func(options *client.RequestOptions) {
		if nil == options.Codec {
			options.Codec = codec
		}
	}
}

// WithTimeout sets the timeout of the request
func WithTimeout(timeout time.Duration) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.Timeout = timeout
	}
}

// WithRetryWaitingMilliseconds sets the retry waiting time of the request
func WithRetryWaitingMilliseconds(retryWaitingTime time.Duration) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.RetryWaitingTime = retryWaitingTime
	}
}

// WithMaxRetryTimes sets the max retry times of the request
func WithMaxRetryTimes(maxRetryTimes int) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.MaxRetryTimes = maxRetryTimes
	}
}

// WithDeleteTransactionPropagationInformation masks delete all transaction propagation information of the request
func WithDeleteTransactionPropagationInformation(isDeleteTransactionPropagationInformation bool) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.DeleteTransactionPropagationInformation = isDeleteTransactionPropagationInformation
	}
}

// WithMaskerConfig sets the masker rules of the request
func WithMaskerConfig(masker config.Masker) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.Masker = masker
	}
}

// WithEnableLogging sets the logging flag
func WithEnableLogging(enableLogging bool) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.EnableLogging = enableLogging
	}
}

// SkipResponseAutoParse marks skip auto parse on response
func SkipResponseAutoParse() client.RequestOption {
	return func(options *client.RequestOptions) {
		options.SkipResponseAutoParse = true
	}
}

// DisableMacroModel marks disable macro model
func DisableMacroModel() client.RequestOption {
	return func(options *client.RequestOptions) {
		options.DisableMacroModel = true
	}
}

// WithMaxWaitingTime sets the max waiting time of the request
func WithMaxWaitingTime(maxWaitingTime time.Duration) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.MaxWaitingTime = maxWaitingTime
	}
}

// WithTopicTypeIfEmpty sets the topic type of the request if topic type is empty
func WithTopicTypeIfEmpty(topicType string) client.RequestOption {
	return func(options *client.RequestOptions) {
		if "" == options.TopicType {
			options.TopicType = topicType
		}
	}
}

// MarkLocalCall marks the request as local service call
func MarkLocalCall() client.RequestOption {
	return func(options *client.RequestOptions) {
		options.IsLocalCall = true
	}
}

// MarkDMQEligible marks the reuqest message enable DMQ eligible
func MarkDMQEligible() client.RequestOption {
	return func(options *client.RequestOptions) {
		options.IsDMQEligible = true
	}
}

// MarkPersistentDeliveryMode marks the request message as persistent delivery mode
func MarkPersistentDeliveryMode() client.RequestOption {
	return func(options *client.RequestOptions) {
		options.IsPersistentDeliveryMode = true
	}
}

// WithVersion sets the version of the request
func WithVersion(version string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.Version = version
	}
}

// WithServiceKey sets the service key of request
func WithServiceKey(serviceKey string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.ServiceKey = serviceKey
	}
}

// WithFallbackFunc sets the fallback fucntion of sync. call request
func WithFallbackFunc(fallbackFunc func(error) error) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.FallbackFunc = fallbackFunc
	}
}

// MarkIsEnableCircuitBreaker sets is enable circuit breaker or not
func MarkIsEnableCircuitBreaker(isEnableCircuitBreaker bool) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.IsEnabledCircuitBreaker = isEnableCircuitBreaker
	}
}

// WithHTTPRequestInfo sets the http request information of the request, only enabled in direct request model
func WithHTTPRequestInfo(address, method, contextType string) client.RequestOption {
	return func(options *client.RequestOptions) {
		options.HTTPCall = true

		options.Address = address
		// default http method is 'POST'
		if "" == method {
			options.HTTPMethod = constant.DefaultHTTPMethodPost
		} else {
			options.HTTPMethod = method
		}
		// default context type is 'application/json'
		if "" == contextType {
			options.ContentType = constant.DefaultContentTypeJSON
		} else {
			options.ContentType = contextType
		}
	}
}

// WithResponseSessionName sets the session name of response
func WithResponseSessionName(sessionName string) client.ResponseOption {
	return func(options *client.ResponseOptions) {
		options.SessionName = sessionName
	}
}

// WithReplyToAddress sets the reply to address of the response
func WithReplyToAddress(replyToAddress string) client.ResponseOption {
	return func(options *client.ResponseOptions) {
		options.ReplyToAddress = replyToAddress
	}
}

// WithResponseCodec sets the codec of the response
func WithResponseCodec(codec codec.Codec) client.ResponseOption {
	return func(options *client.ResponseOptions) {
		options.Codec = codec
	}
}

// AddKeyToResponseHeader adds a pair of key-value into response header
func AddKeyToResponseHeader(key, value string) client.ResponseOption {
	return func(options *client.ResponseOptions) {
		if nil == options.Header {
			options.Header = make(map[string]string)
		}
		options.Header[key] = value
	}
}

// WithResponseHeader sets the header of the response
func WithResponseHeader(header map[string]string) client.ResponseOption {
	return func(options *client.ResponseOptions) {
		if nil == options.Header {
			options.Header = header
		} else {
			for k, v := range header {
				options.Header[k] = v
			}
		}
	}
}

// WithResponseContext sets the context.Context of the response
func WithResponseContext(ctx context.Context) client.ResponseOption {
	return func(options *client.ResponseOptions) {
		options.Context = ctx
	}
}
