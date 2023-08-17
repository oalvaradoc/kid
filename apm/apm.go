package apm

import (
	"context"
	"git.multiverse.io/eventkit/kit/common/apm"
	"git.multiverse.io/eventkit/kit/log"
)

// ApmLogger is a global APM logger that will be initialized at startup
var ApmLogger *apm.Log

// defines all the type of APM log
const (
	LogServer = "s"
	LogClient = "c"
)

// ApmRecord wrappers the record information of APM
type ApmRecord struct {
	LogType              string
	StartTimestamp       int64
	DurationMicroseconds int64
	TraceID              string
	SpanID               string
	ParentSpanID         string
	Org                  string
	Az                   string
	Su                   string
	Wks                  string
	Env                  string
	NodeID               string
	ServiceID            string
	InstanceID           string
	TopicID              string
	SrcServiceID         string
	ErrorCode            string
	ErrorMsg             string
	Attach               string
}

// DoAPMLogging prints apm log
func DoAPMLogging(ctx context.Context, r *ApmRecord) {
	if ApmLogger == nil {
		log.Infos("The APM logger is not initialized, cannot do APM logging.")
		return
	}

	err := ApmLogger.APMlogf("%s|%d|%d|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s",
		r.LogType,
		r.StartTimestamp,
		r.DurationMicroseconds,
		r.TraceID,
		r.SpanID,
		r.ParentSpanID,
		r.Org,
		r.Az,
		r.Su,
		r.Wks,
		r.Env,
		r.NodeID,
		r.ServiceID,
		r.InstanceID,
		r.TopicID,
		r.SrcServiceID,
		r.ErrorCode,
		r.ErrorMsg,
		r.Attach)
	if err != nil {
		log.Errorf(ctx, "Failed to do APM logging %++v", err)
	}
}
