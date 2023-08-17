package logging

import (
	"context"
	"fmt"
	"time"

	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/log"
	masker "git.multiverse.io/eventkit/kit/masker/comm"
)

// Interceptor is an interceptor that used to write request/response messages into the log file.
type Interceptor struct{}

// PreHandle writes request message into the log file
func (l *Interceptor) PreHandle(ctx context.Context, request *msg.Message) error {
	var gerr error
	startTime := time.Now()
	defer func() {
		now := time.Now()
		if now.Sub(startTime).Milliseconds() > 10 {
			log.Infof(ctx, "######time cost in `Logging` interceptor(PreHandle) greater than 10 ms:%++v", now.Sub(startTime))
		}
	}()
	handlerContexts := contexts.HandlerContextsFromContext(ctx)
	if nil != handlerContexts {
		handlerContexts.StartHandleTime = time.Now()
	}
	configs := config.GetConfigs()

	targetRequest := request
	if nil != configs && nil != targetRequest {
		targetRequest, gerr = masker.MaskMessageIfNecessary(ctx, configs.Log.HeaderMaskRules, configs.Log.RequestHeaderMaskRules, configs.Log.MaskRules, configs.Log.RequestBodyMaskRules, request)
		if nil != gerr {
			log.Errorf(ctx, "Failed to do the header mask for request, error:%++v", gerr)
		}
	}

	log.Infosw(fmt.Sprintf("%++v", targetRequest), buildTraceKeyValues(ctx, "Request")...)
	return nil
}

// PostHandle writes response message into the log file
func (l *Interceptor) PostHandle(ctx context.Context, request *msg.Message, response *msg.Message) error {
	var gerr error
	startTime := time.Now()
	defer func() {
		now := time.Now()
		if now.Sub(startTime).Milliseconds() > 10 {
			log.Infof(ctx, "######time cost in `Logging` interceptor(PostHandle) greater than 10 ms:%++v", now.Sub(startTime))
		}
	}()
	handlerContexts := contexts.HandlerContextsFromContext(ctx)
	configs := config.GetConfigs()

	targetResponse := response
	if nil != configs && nil != targetResponse {
		targetResponse, gerr = masker.MaskMessageIfNecessary(ctx, configs.Log.HeaderMaskRules, configs.Log.ResponseHeaderMaskRules, configs.Log.MaskRules, configs.Log.ResponseBodyMaskRules, response)
		if nil != gerr {
			log.Errorf(ctx, "Failed to do the header mask for request, error:%++v", gerr)
		}
	}

	if nil != handlerContexts && 0 != handlerContexts.StartHandleTime.Second() {
		timeCost := time.Now().Sub(handlerContexts.StartHandleTime)
		log.Infosw(fmt.Sprintf("%++v", targetResponse), buildTraceKeyValues(ctx, "Response", timeCost)...)
	} else {
		log.Infosw(fmt.Sprintf("%++v", targetResponse), buildTraceKeyValues(ctx, "Response")...)
	}

	return nil
}

func buildTraceKeyValues(ctx context.Context, messageType string, timeCost ...time.Duration) []interface{} {
	res := make([]interface{}, 0)
	if traceID := checkNil(ctx.Value(log.KeyTraceID)); traceID != "" {
		res = append(res, log.KeyTraceID, traceID)
	}
	if parentID := checkNil(ctx.Value(log.KeyParentID)); parentID != "" {
		res = append(res, log.KeyParentID, parentID)
	}
	if spanID := checkNil(ctx.Value(log.KeySpanID)); spanID != "" {
		res = append(res, log.KeySpanID, spanID)
	}
	if serviceID := checkNil(ctx.Value(log.KeyServiceID)); serviceID != "" {
		res = append(res, log.KeyServiceID, serviceID)
	}

	if topicID := checkNil(ctx.Value(log.KeyTopicID)); topicID != "" {
		res = append(res, log.KeyTopicID, topicID)
	}
	res = append(res, "MessageType", messageType)

	if len(timeCost) > 0 {
		res = append(res, "timeCost", fmt.Sprintf("%++v", timeCost[0]))
	}
	return res
}

func checkNil(arg interface{}) interface{} {
	if arg == nil {
		return ""
	}
	return arg
}

func (l Interceptor) String() string {
	return constant.InterceptorLogging
}
