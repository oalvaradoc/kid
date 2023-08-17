package logging

import (
	"context"
	"fmt"
	"time"

	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
)

// Wrapper is an wrapper for logging
type Wrapper struct{}

// Before will logging the request
func (t *Wrapper) Before(ctx context.Context, request interface{}, opts interface{}) (context.Context, error) {
	if log.IsEnable(log.DebugLevel) {
		startTime := time.Now()
		defer func() {
			now := time.Now()
			if now.Sub(startTime).Milliseconds() > 10 {
				log.Infof(ctx, "######time cost in `Logging` wrapper(Before) greater than 10 ms:%++v", now.Sub(startTime))
			}
		}()

		if !util.IsNil(opts) {
			requestOptions := opts.(*client.RequestOptions)
			if requestOptions.EnableLogging {
				log.Debugf(ctx, "logging.Before, request:[%++v]", request)
			}
		}
	}
	newCtx := context.WithValue(ctx, constant.KeyDownstreamServiceLoggingStartTime, time.Now())
	return newCtx, nil
}

// After will logging the response
func (t *Wrapper) After(ctx context.Context, request interface{}, responseMeta interface{}, opts interface{}) (context.Context, error) {
	if log.IsEnable(log.DebugLevel) {
		startTime := time.Now()
		defer func() {
			now := time.Now()
			if now.Sub(startTime).Milliseconds() > 10 {
				log.Infof(ctx, "######time cost in `Logging` wrapper(After) greater than 10 ms:%++v", now.Sub(startTime))
			}
		}()

		if !util.IsNil(opts) {
			requestOptions := opts.(*client.RequestOptions)
			if requestOptions.EnableLogging {
				downstreamServiceLoggingStartTime := ctx.Value(constant.KeyDownstreamServiceLoggingStartTime)
				if value, ok := downstreamServiceLoggingStartTime.(time.Time); ok {
					timeCost := time.Now().Sub(value)
					log.Debugsw(fmt.Sprintf("logging.After, request:[%++v]-responseMeta:[%++v]", request, responseMeta), buildTraceKeyValues(ctx, timeCost)...)
				} else {
					log.Debugf(ctx, "logging.After, request:[%++v]-responseMeta:[%++v]", request, responseMeta)
				}
			}
		}
	}
	return ctx, nil
}

func buildTraceKeyValues(ctx context.Context, timeCost ...time.Duration) []interface{} {
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
	if len(timeCost) > 0 {
		res = append(res, "downstreamTimeCost", fmt.Sprintf("%++v", timeCost[0]))
	}
	return res
}

func checkNil(arg interface{}) interface{} {
	if arg == nil {
		return ""
	}
	return arg
}

func (t Wrapper) String() string {
	return constant.WrapperLogging
}
