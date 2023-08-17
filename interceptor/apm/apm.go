package apm

import (
	"context"
	"encoding/base64"
	"fmt"
	"git.multiverse.io/eventkit/kit/apm"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/log"
	"strconv"
	"time"
)

// Interceptor is an interceptor for APM log print
type Interceptor struct{}

// PreHandle writes the client side APM log
func (l *Interceptor) PreHandle(ctx context.Context, request *msg.Message) error {
	configs := config.GetConfigs()
	if nil == configs || !configs.Apm.Enable || nil == request {
		return nil
	}

	if "1" == request.GetAppPropertyIgnoreCaseSilence(constant.TxnIsLocalCall) {
		srStartTimestampString := strconv.FormatInt(time.Now().UnixNano()/1000, 10)
		request.SetAppProperty(constant.SrStartTimestamp, srStartTimestampString)
	}
	return nil
}

// PostHandle writes the server side APM log
func (l *Interceptor) PostHandle(ctx context.Context, request *msg.Message, response *msg.Message) error {
	configs := config.GetConfigs()
	if nil == request || nil == configs || !configs.Apm.Enable {
		return nil
	}

	if "1" == request.GetAppPropertyIgnoreCaseSilence(constant.TxnIsLocalCall) {
		startTime := time.Now()
		defer func() {
			now := time.Now()
			if now.Sub(startTime).Milliseconds() > 10 {
				log.Infof(ctx, "######time cost in `APM` interceptor(PostHandle) greater than 10 ms:%++v", now.Sub(startTime))
			}
		}()
		handlerContexts := contexts.HandlerContextsFromContext(ctx)
		if nil == handlerContexts {
			log.Errorf(ctx, "Cannot found handler contexts in context")
			return nil
		}

		if nil == handlerContexts.SpanContexts {
			log.Errorf(ctx, "Cannot found span contexts in handler contexts")
			return nil
		}

		srStartTimestampString := request.GetAppPropertyIgnoreCaseSilence(constant.SrStartTimestamp)
		if "" == srStartTimestampString {
			log.Info(ctx, "Cannot found the SR(server receive) start timestamp")
			return nil
		}
		srStartTimestamp, err := strconv.ParseInt(srStartTimestampString, 10, 64)
		if nil != err {
			log.Errorf(ctx, "parse CS start timestamp[%s] to int64 failed,err=%v", srStartTimestampString, err)
			// skip apm logging
			return nil
		}
		//durationString := strconv.FormatInt((time.Now().UnixNano()-srStartTimestamp)/1000, 10)
		su := util.GetEither(request.TopicAttribute, constant.TopicDestinationSU, constant.TopicDestinationDCN)
		if len(su) == 0 {
			su = handlerContexts.Su
		}

		errCode := ""
		errMsg := ""
		if !util.IsNil(response) {
			errCode = response.GetAppPropertyEitherSilence(constant.ReturnErrorCode, constant.ReturnErrorCodeOld)
			errMsg = response.GetAppPropertyEitherSilence(constant.ReturnErrorMsg, constant.ReturnErrorMsgOld)
		}

		if len(errMsg) > 0 {
			errMsg = base64.StdEncoding.EncodeToString([]byte(errMsg))
		}
		// do the apm logging(client)
		apm.DoAPMLogging(ctx, &apm.ApmRecord{
			LogType:              apm.LogServer,
			StartTimestamp:       srStartTimestamp,
			DurationMicroseconds: time.Now().UnixNano()/1000 - srStartTimestamp,
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
			TopicID:              request.GetMsgTopicId(),
			SrcServiceID:         handlerContexts.ServiceID,
			ErrorCode:            errCode,
			ErrorMsg:             errMsg,
			Attach:               base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("packetSize:%d", len(request.Body)))),
		})
	}
	return nil
}

func (l Interceptor) String() string {
	return constant.InterceptorApm
}
