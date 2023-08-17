package msg

import (
	"encoding/json"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
	"regexp"
	"strings"
	"sync"
)

// Message is data struct of service and sed server
type Message struct {

	// Each interaction has a session unique identifier
	ID uint64

	// Used to save topic-related key-values pairs
	TopicAttribute map[string]string

	RequestURL string

	// True means this is a synchronous call message
	// False means this is an asynchronous call message
	NeedReply bool

	// If it's a synchronous call message, this field value is true
	// If it's a asynchronous call message, this field value is false
	NeedAck bool

	// sessionName is used to distinguish which session the message was sent to solace
	SessionName string

	// app properties
	// service can set this attributes pass to target service
	appProps map[string]string

	appPropsRWLock sync.RWMutex

	// message payload
	Body []byte
}

// RangeAppProps is used to range app properties with a closure function
func (msg Message) RangeAppProps(f func(string, string)) {
	msg.appPropsRWLock.Lock()
	defer func() { msg.appPropsRWLock.Unlock() }()

	for k, v := range msg.appProps {
		f(k, v)
	}
}

// GetMsgTopicType is used to get topic type of message
func (msg Message) GetMsgTopicType() string {
	if nil == msg.TopicAttribute {
		return ""
	}

	return msg.TopicAttribute[constant.TopicType]
}

// GetMsgTopicId is used to get topic id of message
func (msg Message) GetMsgTopicId() string {
	if nil == msg.TopicAttribute {
		return ""
	}

	return msg.TopicAttribute[constant.TopicID]
}

// GetSourceORG gets source org. from topic attributes
func (msg Message) GetSourceORG() string {
	if nil == msg.TopicAttribute {
		return ""
	}

	return msg.TopicAttribute[constant.TopicSourceORG]
}

// GetSourceWKS gets source workspace from topic attributes
func (msg Message) GetSourceWKS() string {
	if nil == msg.TopicAttribute {
		return ""
	}

	return msg.TopicAttribute[constant.TopicSourceWorkspace]
}

// GetSourceENV gets source environment from topic attributes
func (msg Message) GetSourceENV() string {
	if nil == msg.TopicAttribute {
		return ""
	}

	return msg.TopicAttribute[constant.TopicSourceEnvironment]
}

// GetSourceAZ gets source available zone from topic attributes
func (msg Message) GetSourceAZ() string {
	if nil == msg.TopicAttribute {
		return ""
	}

	return msg.TopicAttribute[constant.TopicSourceAZ]
}

// GetSourceServiceID gets source service ID from topic attributes
func (msg Message) GetSourceServiceID() string {
	if nil == msg.TopicAttribute {
		return ""
	}

	return msg.TopicAttribute[constant.TopicSourceServiceID]
}

// GetSourceSU gets source SU from topic attributes
func (msg Message) GetSourceSU() string {
	if nil == msg.TopicAttribute {
		return ""
	}

	return msg.TopicAttribute[constant.TopicSourceSU]
}

// GetSourceDCN gets source DCN from topic attributes
func (msg Message) GetSourceDCN() string {
	if nil == msg.TopicAttribute {
		return ""
	}

	return msg.TopicAttribute[constant.TopicSourceDCN]
}

// GetSourceNodeID gets source node id from topic attributes
func (msg Message) GetSourceNodeID() string {
	if nil == msg.TopicAttribute {
		return ""
	}

	return msg.TopicAttribute[constant.TopicSourceNodeID]
}

// GetSourceInstanceID gets source instance id from topic attributes
func (msg Message) GetSourceInstanceID() string {
	if nil == msg.TopicAttribute {
		return ""
	}

	return msg.TopicAttribute[constant.TopicSourceInstanceID]
}

// SetAppProperty sets the key-value pair into app properties
func (msg *Message) SetAppProperty(key, value string) {
	msg.appPropsRWLock.Lock()
	defer func() { msg.appPropsRWLock.Unlock() }()
	if nil == msg.appProps {
		msg.appProps = make(map[string]string, 0)
	}
	msg.appProps[key] = value
}

// SetAppProps replaces the value of app properties
func (msg *Message) SetAppProps(m map[string]string) {
	msg.appProps = m
}

// GetAppProps returns the reference of the app properties
func (msg Message) GetAppProps() map[string]string {
	return msg.appProps
}

// CloneAppProps returns the cloned app properties of Message
func (msg Message) CloneAppProps() map[string]string {
	msg.appPropsRWLock.RLock()
	defer func() { msg.appPropsRWLock.RUnlock() }()

	if nil == msg.appProps {
		return nil
	}
	cloneMaps := make(map[string]string)
	for k, v := range msg.appProps {
		cloneMaps[k] = v
	}
	return cloneMaps

}

// GetAppProperty returns the value from the app property
func (msg Message) GetAppProperty(key string) (res string, ok bool) {
	msg.appPropsRWLock.RLock()
	defer func() { msg.appPropsRWLock.RUnlock() }()

	r, b := msg.appProps[key]

	return r, b
}

// GetAppPropertyIgnoreCase returns the value from the app property and ignore the case of the key
func (msg Message) GetAppPropertyIgnoreCase(key string) (res string, ok bool) {
	msg.appPropsRWLock.RLock()
	defer func() { msg.appPropsRWLock.RUnlock() }()

	for k, v := range msg.appProps {
		if strings.EqualFold(k, key) {
			return v, true
		}
	}

	return "", false
}

// GetAppPropertyEitherSilence gets the value from the app property in order according to the two keys
func (msg Message) GetAppPropertyEitherSilence(key1, key2 string) string {
	msg.appPropsRWLock.RLock()
	defer func() { msg.appPropsRWLock.RUnlock() }()

	if nil == msg.appProps {
		return ""
	}

	if v, ok := msg.appProps[key1]; ok {
		return v
	}

	return msg.appProps[key2]
}

// GetAppPropertySilence returns the value from the app property
func (msg Message) GetAppPropertySilence(key string) string {
	msg.appPropsRWLock.RLock()
	defer func() { msg.appPropsRWLock.RUnlock() }()

	if nil == msg.appProps {
		return ""
	}
	return msg.appProps[key]
}

// GetAppPropertyIgnoreCaseSilence returns the value from the app property and ignore the case of the key
func (msg Message) GetAppPropertyIgnoreCaseSilence(key string) string {
	r, _ := msg.GetAppPropertyIgnoreCase(key)

	return r
}

// DeleteProperty deletes the app property by key
func (msg *Message) DeleteProperty(key string) {
	msg.appPropsRWLock.Lock()
	defer func() { msg.appPropsRWLock.Unlock() }()

	if nil == msg.appProps {
		return
	}

	delete(msg.appProps, key)
}

// AppPropsToString converts app properties into string
func (msg Message) AppPropsToString() string {
	msg.appPropsRWLock.RLock()
	defer func() { msg.appPropsRWLock.RUnlock() }()

	return util.MapToString(msg.appProps)
}

// TopicAttributesToString converts topic attributes into string
func (msg Message) TopicAttributesToString() string {
	return util.MapToString(msg.TopicAttribute)
}

// IsValidTopicType is used to check whether topic type is valid
func (msg Message) IsValidTopicType() bool {
	switch msg.GetMsgTopicType() {
	case constant.TopicTypeHeartbeat,
		constant.TopicTypeError,
		constant.TopicTypeAlert,
		constant.TopicTypeBusiness,
		constant.TopicTypeLog,
		constant.TopicTypeMetrics,
		constant.TopicTypeDXC,
		constant.TopicTypeDTS,
		constant.TopicTypeOPS:
		return true
	default:
		return false
	}
}

// JudgeUserLang judges the UserLang from Message header, the default user lang is constant.LangEnUS
func (msg Message) JudgeUserLang() string {
	msg.appPropsRWLock.RLock()
	defer func() { msg.appPropsRWLock.RUnlock() }()

	lang := ""
	if nil != msg.appProps {
		lang = msg.appProps[constant.UserLang]
	}

	if "" == lang {
		lang = constant.LangEnUS
	}

	return lang
}

// MessageForPrint is data struct for print Message
type MessageForPrint struct {

	// Each interaction has a session unique identifier
	ID uint64

	// Used to save topic-related key-values pairs
	TopicAttribute map[string]string

	RequestURL string

	// True means this is a synchronous call message
	// False means this is an asynchronous call message
	NeedReply bool

	// If it's a synchronous call message, this field value is true
	// If it's a asynchronous call message, this field value is false
	NeedAck bool

	// sessionName is used to distinguish which session the message was sent to solace
	SessionName string

	// app properties
	// service can set this attributes pass to target service
	AppProps map[string]string

	// message payload
	Body string
}

// String override the build-in String function, convert Message to string
func (msg Message) String() string {
	mfp := MessageForPrint{
		ID:             msg.ID,
		TopicAttribute: msg.TopicAttribute,
		NeedReply:      msg.NeedReply,
		SessionName:    msg.SessionName,
		AppProps:       msg.appProps,
		Body:           string(msg.Body),
	}
	bytes, err := json.Marshal(mfp)
	if nil != err {
		log.Errorsf("Failed to marshal message, error:%++v", err)
	}

	return string(bytes)
}

func removeLBR(text string) string {
	re := regexp.MustCompile(`\r\n|[\r\n\v\f\x{0085}\x{2028}\x{2029}\x{0000}\x{001d}]`)
	return re.ReplaceAllString(text, ` `)
}

// CustomErrorWrapperFn defines the custom error wrapper function
type CustomErrorWrapperFn func(errorCode, errorMessage string, response *Message)

// WrapperErrorResponse sets the return message when the service execution is abnormal according to the
// response template and error code
func WrapperErrorResponse(err interface{},
	lang string,
	responseTemplate string,
	customErrorWrapperFn CustomErrorWrapperFn,
	responseDataWhenError ...interface{}) (response *Message) {
	var toWrapError *errors.Error
	switch err.(type) {
	case *errors.Error:
		{
			toWrapError = err.(*errors.Error)
		}
	case errors.Error:
		{
			v := err.(errors.Error)
			toWrapError = &v
		}
	case error:
		{
			toWrapError = errors.Wrap(constant.SystemInternalError, err.(error), 0)
		}
	default:
		toWrapError = errors.Errorf(constant.SystemInternalError, "%++v", err)
	}

	// set error code && error message to response header
	errorCode, errorMsg, body := toWrapError.WrapError(lang, responseTemplate, responseDataWhenError...)
	finalErrorCode := removeLBR(errorCode)
	finalErrorMsg := removeLBR(errorMsg)
	response = &Message{
		appProps: map[string]string{
			constant.ReturnStatus:       "F",
			constant.ReturnErrorCode:    finalErrorCode,
			constant.ReturnErrorCodeOld: finalErrorCode,
			constant.ReturnErrorMsg:     finalErrorMsg,
			constant.ReturnErrorMsgOld:  finalErrorMsg,
		},
		Body: body,
	}

	if len(toWrapError.GetExtHeader()) > 0 {
		for k, v := range toWrapError.GetExtHeader() {
			response.SetAppProperty(k, v)
		}
	}
	if len(toWrapError.GetDuplicateErrorCodeTo()) > 0 {
		response.SetAppProperty(toWrapError.GetDuplicateErrorCodeTo(), finalErrorCode)
	}
	if nil != customErrorWrapperFn {
		customErrorWrapperFn(errorCode, errorMsg, response)
	}
	return response
}
