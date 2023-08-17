package apm

import (
	"context"
	"encoding/base64"
	"fmt"
	"git.multiverse.io/eventkit/kit/apm"
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/log"
	"time"
)

// Wrapper is an wrapper for APM
type Wrapper struct{}

// Before prints APM log for macro model
func (t *Wrapper) Before(ctx context.Context, request interface{}, opts interface{}) (context.Context, error) {
	if util.IsNil(opts) {
		log.Info(ctx, "APM-Before:options is empty, skip APM logging!")
		return ctx, nil
	}

	requestOptions := opts.(*client.RequestOptions)
	if !requestOptions.IsLocalCall {
		log.Debugf(ctx, "The request is not local call, skip APM logging!")
		return ctx, nil
	}

	configs := config.GetConfigs()
	if nil == configs || !configs.Apm.Enable {
		return ctx, nil
	}

	csStartTimestamp := time.Now()
	nCtx := context.WithValue(ctx, constant.CsStartTimestamp, csStartTimestamp)

	return nCtx, nil
}

// After prints APM log for macro model
func (t *Wrapper) After(ctx context.Context, request interface{}, responseMeta interface{}, opts interface{}) (context.Context, error) {
	if util.IsNil(opts) {
		log.Info(ctx, "APM-After:options is empty, skip APM logging!")
		return ctx, nil
	}
	requestOptions := opts.(*client.RequestOptions)
	if !requestOptions.IsLocalCall {
		log.Debugf(ctx, "The request is not local call, skip APM logging!")
		return ctx, nil
	}

	startTime := time.Now()
	defer func() {
		now := time.Now()
		if now.Sub(startTime).Milliseconds() > 10 {
			log.Infof(ctx, "######time cost in `APM` wrapper(After) greater than 10 ms:%++v", now.Sub(startTime))
		}
	}()

	configs := config.GetConfigs()
	if nil == configs || !configs.Apm.Enable {
		return ctx, nil
	}

	errCode := ""
	errMsg := ""

	handlerContexts := contexts.HandlerContextsFromContext(ctx)
	if nil == handlerContexts {
		return ctx, errors.Errorf(constant.SystemInternalError, "Cannot found handler contexts in context")
	}

	if nil == handlerContexts.SpanContexts {
		return ctx, errors.Errorf(constant.SystemInternalError, "Cannot found span contexts in handler contexts")
	}

	csStartTime := ctx.Value(constant.CsStartTimestamp)

	if util.IsNil(csStartTime) {
		log.Info(ctx, "client start timestamp is empty, skip APM logging")
		return ctx, nil
	}
	csStartTimestamp := csStartTime.(time.Time)
	if !util.IsNil(responseMeta) {
		responseMetaObj := responseMeta.(client.ResponseMeta)
		errCode = util.GetEither(responseMetaObj.Header(), constant.ReturnErrorCode, constant.ReturnErrorCodeOld)
		errMsg = util.GetEither(responseMetaObj.Header(), constant.ReturnErrorMsg, constant.ReturnErrorMsgOld)
	}

	if len(errMsg) > 0 {
		errMsg = base64.StdEncoding.EncodeToString([]byte(errMsg))
	}

	su := requestOptions.Su
	if len(su) == 0 {
		su = handlerContexts.Su
	}
	// do the apm logging(client)
	apm.DoAPMLogging(ctx, &apm.ApmRecord{
		LogType:              apm.LogClient,
		StartTimestamp:       csStartTimestamp.UnixNano() / 1000,
		DurationMicroseconds: time.Now().Sub(csStartTimestamp).Microseconds(),
		TraceID:              handlerContexts.SpanContexts.TraceID,
		SpanID:               handlerContexts.SpanContexts.SpanID,
		ParentSpanID:         handlerContexts.SpanContexts.ParentSpanID,
		Org:                  handlerContexts.Org,
		Az:                   handlerContexts.Az,
		Su:                   su,
		Wks:                  handlerContexts.Wks,
		Env:                  handlerContexts.Env,
		NodeID:               handlerContexts.NodeID,
		ServiceID:            handlerContexts.ServiceID,
		InstanceID:           handlerContexts.InstanceID,
		TopicID:              requestOptions.EventID,
		SrcServiceID:         handlerContexts.ServiceID,
		ErrorCode:            errCode,
		ErrorMsg:             errMsg,
		Attach:               base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("packetSize:%d", 0))),
	})

	return ctx, nil
}

func (t Wrapper) String() string {
	return constant.WrapperApm
}
