package audit

import (
	"context"
	"git.multiverse.io/eventkit/kit/common/auth"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/interceptor"
	"git.multiverse.io/eventkit/kit/log"
	"time"
)

// NewAuditInterceptor creates a new audit interceptor with operators mapping and excludeTopics mapping
func NewAuditInterceptor(operators map[string]string, excludeTopics []string, sslDir ...string) interceptor.Interceptor {
	// initialize JWT.
	//auth.InitJwt(sslDir...)

	if operators == nil {
		operators = map[string]string{}
	}
	if excludeTopics == nil {
		excludeTopics = []string{}
	}
	excludeTopicMap := make(map[string]bool)
	for _, topic := range excludeTopics {
		excludeTopicMap[topic] = true
	}
	return &auditInterceptor{
		operator: operators,
		excludes: excludeTopicMap,
	}
}

type auditInterceptor struct {
	operator map[string]string
	excludes map[string]bool
}

func (a auditInterceptor) exclude(topicID string) bool {
	_, ok := a.excludes[topicID]
	return ok
}

// PreHandle does nothing
func (a auditInterceptor) PreHandle(ctx context.Context, request *msg.Message) error {
	return nil
}

// PostHandle writes the audit information into the log file.
func (a auditInterceptor) PostHandle(ctx context.Context, request *msg.Message, response *msg.Message) error {
	startTime := time.Now()
	defer func() {
		now := time.Now()
		if now.Sub(startTime).Milliseconds() > 10 {
			log.Infof(ctx, "######time cost in `Audit` interceptor(PostHandle) greater than 10 ms:%++v", now.Sub(startTime))
		}
	}()
	if nil == request || nil == response || a.exclude(request.GetMsgTopicId()) {
		return nil
	}

	errCode, _ := response.GetAppPropertyIgnoreCase(constant.ReturnErrorCode)
	errMsg, _ := response.GetAppPropertyIgnoreCase(constant.ReturnErrorMsg)

	tokenStr, _ := request.GetAppProperty("Authorization")
	passport, err := auth.UMVerifier.VerifyToken(tokenStr)
	if err != nil {
		passport = &auth.Passport{}
		passport.TenantCode = "UnKnow"
		passport.AccountName = "UnKnow"
		//doing nothing...
	}
	org := passport.TenantCode
	user := passport.AccountName

	optionArg := map[string]string{
		"org":       org,
		"errorCode": errCode,
		"errorMsg":  errMsg,
		"response":  string(response.Body),
	}
	topicID := request.GetMsgTopicId()
	if operator, ok := a.operator[topicID]; ok {
		log.AuditError(ctx, user, operator, "", request.Body, []byte(""), log.MetaDataWithMap(optionArg))
	} else {
		log.AuditError(ctx, user, topicID, "", request.Body, []byte(""), log.MetaDataWithMap(optionArg))

	}
	return nil
}

func (a auditInterceptor) String() string {
	return constant.InterceptorAudit
}
