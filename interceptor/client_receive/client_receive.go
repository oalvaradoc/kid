package client_receive

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/log"
	"github.com/clbanning/mxj/v2"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"strings"
)

// Interceptor is an interceptor for unwrap the downstream server response body
type Interceptor struct{}

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	defaultResponseAutoParseKeyMapping = map[string]string{
		constant.ResponseAutoParseTypeMappingKey: constant.ResponseAutoDefaultParseType,
		constant.ErrorCodeMappingKey:             constant.DefaultErrorCodeKey,
		constant.ErrorMsgMappingKey:              constant.DefaultErrorMsgKey,
		constant.ResponseDataBodyMappingKey:      constant.DefaultDataBodyKey,
	}
)

// PreHandle does nothing
func (l *Interceptor) PreHandle(ctx context.Context, request *msg.Message) error {
	// DO NOTHING
	return nil
}

func toIntOrDefault(ori string, defaultValue int) int {
	result, err := strconv.Atoi(ori)
	if nil != err {
		return defaultValue
	}

	return result
}

func getMapValueOrDefault(m map[string]string, key string, defaultValue string) string {
	result := util.GetMapValueIgnoreCase(m, key)
	if "" != result {
		return result
	}

	return defaultValue
}

func interface2String(o interface{}) string {
	switch o.(type) {
	case string:
		return o.(string)
	default:
		return fmt.Sprintf("%v", o)
	}

	return ""
}

func interface2Map(o interface{}) map[string]interface{} {
	switch o.(type) {
	case map[string]interface{}:
		return o.(map[string]interface{})
	}

	return make(map[string]interface{})
}

func handleForJSON(responseAutoParseKeyMapping map[string]string, request *msg.Message, response *msg.Message) error {
	if len(response.Body) > 0 {
		errorCodeKey := getMapValueOrDefault(responseAutoParseKeyMapping, constant.ErrorCodeMappingKey, constant.DefaultErrorCodeKey)
		errorCodeAny := json.Get(response.Body, errorCodeKey)
		if nil != errorCodeAny.LastError() {
			return errors.Errorf(constant.SystemInternalError, "Try to parse error code[key=%s] failed, error:%++v, response body:[%s]",
				errorCodeKey, errorCodeAny.LastError(), string(response.Body))
		}
		errorCode := toIntOrDefault(errorCodeAny.ToString(), -1)

		errorMessageKey := getMapValueOrDefault(responseAutoParseKeyMapping, constant.ErrorMsgMappingKey, constant.DefaultErrorMsgKey)
		errorMessageAng := json.Get(response.Body, errorMessageKey)
		if nil != errorMessageAng.LastError() {
			return errors.Errorf(constant.SystemInternalError, "Try to parse error message[key=%s] failed, error:%++v, response body:[%s]",
				errorMessageKey, errorMessageAng.LastError(), string(response.Body))
		}

		if 0 != errorCode {
			return errors.New(errorCodeAny.ToString(), errorMessageAng.ToString())
		}

		dataKey := getMapValueOrDefault(responseAutoParseKeyMapping, constant.ResponseDataBodyMappingKey, constant.DefaultDataBodyKey)
		data := json.Get(response.Body, dataKey)
		lastError := data.LastError()
		if nil != lastError {
			return errors.Errorf(constant.SystemInternalError, "parse JSON failed, error:%++v", lastError)
		}

		stream := jsoniter.NewStream(jsoniter.ConfigDefault, nil, 32)
		data.WriteTo(stream)
		dataBytes := stream.Buffer()
		response.Body = dataBytes
	}

	return nil
}

func handleForXML(responseAutoParseKeyMapping map[string]string, request *msg.Message, response *msg.Message) error {
	if len(response.Body) > 0 {
		errorCodeKey := getMapValueOrDefault(responseAutoParseKeyMapping, constant.ErrorCodeMappingKey, constant.DefaultErrorCodeKey)
		mxjm, err := mxj.NewMapXml(response.Body)
		if nil != err {
			return errors.Errorf(constant.SystemInternalError, "parse XML failed ,error:%++v", err)
		}

		errorCodeAny, err := mxjm.ValueForKey(errorCodeKey)

		if nil != err {
			return errors.Errorf(constant.SystemInternalError, "parse XML with key[%s] failed, error=%++v", errorCodeKey, err)
		}
		errorCode := toIntOrDefault(interface2String(errorCodeAny), -1)

		errorMessageKey := getMapValueOrDefault(responseAutoParseKeyMapping, constant.ErrorMsgMappingKey, constant.DefaultErrorMsgKey)
		errorMessageAng, err := mxjm.ValueForKey(errorMessageKey)
		if nil != err {
			return errors.Errorf(constant.SystemInternalError, "parse XML with key[%s] failed, error=%++v", errorMessageKey, err)
		}
		if 0 != errorCode {
			return errors.New(interface2String(errorCodeAny), interface2String(errorMessageAng))
		}

		dataKey := getMapValueOrDefault(responseAutoParseKeyMapping, constant.ResponseDataBodyMappingKey, constant.DefaultDataBodyKey)
		responseBodyValue, err := mxjm.ValueForKey(dataKey)
		if nil != err {
			return errors.Errorf(constant.SystemInternalError, "parse XML with key[%s] failed, error=%++v\n", dataKey, err)
		}
		nmv := mxj.Map(interface2Map(responseBodyValue))
		bodyBytes, err := nmv.Xml()
		if nil != err {
			return errors.Errorf(constant.SystemInternalError, "new sub xml failed, error:%++v", err)
		}
		response.Body = bodyBytes
	}

	return nil
}

// PostHandle performs the unwrap logic according the response auto parse key mapping
func (l *Interceptor) PostHandle(ctx context.Context, request *msg.Message, response *msg.Message) error {
	if nil == response {
		log.Errorf(ctx, "response is empty, skip client-receive parse logic")
		return nil
	}
	if response.GetAppPropertySilence(constant.ReturnStatus) == "F" {
		return errors.New(response.GetAppPropertySilence(constant.ReturnErrorCodeOld), response.GetAppPropertySilence(constant.ReturnErrorMsgOld))
	}
	errorCodeString := response.GetAppPropertyIgnoreCaseSilence(constant.ReturnErrorCode)
	if "" != errorCodeString {
		if 0 != toIntOrDefault(errorCodeString, -1) {
			return errors.New(errorCodeString, response.GetAppPropertyIgnoreCaseSilence(constant.ReturnErrorMsg))
		}

	}

	var responseAutoParseKeyMapping map[string]string
	if flag := ctx.Value(constant.SkipResponseAutoParseKeyMappingFlagKey); nil != flag {
		// skip response auto parse
		if true == flag.(bool) {
			return nil
		}
	}
	if m := ctx.Value(constant.ResponseAutoParseKeyMappingKey); nil != m {
		responseAutoParseKeyMapping = m.(map[string]string)
	} else {
		handlerContexts := contexts.HandlerContextsFromContext(ctx)
		if nil == handlerContexts {
			responseAutoParseKeyMapping = defaultResponseAutoParseKeyMapping
		} else {
			responseAutoParseKeyMapping = handlerContexts.ResponseAutoParseKeyMapping
		}
	}

	log.Debugf(ctx, "start execute post handle of `client_receive`, response auto parse key mapping=[%++v]", responseAutoParseKeyMapping)

	tp := getMapValueOrDefault(responseAutoParseKeyMapping, constant.ResponseAutoParseTypeMappingKey, constant.ResponseAutoDefaultParseType)
	switch strings.ToUpper(tp) {
	case constant.ResponseAutoParseTypeXML:
		{
			return handleForXML(responseAutoParseKeyMapping, request, response)
		}
	default:
		{
			return handleForJSON(responseAutoParseKeyMapping, request, response)
		}
	}
}

func (l Interceptor) String() string {
	return constant.InterceptorClientReceive
}
