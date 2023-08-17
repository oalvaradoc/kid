package callback

import (
	"encoding/base64"
	"fmt"
	"git.multiverse.io/eventkit/kit/common/model"
	event "git.multiverse.io/eventkit/kit/common/model/transaction"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/handler/base"
	"git.multiverse.io/eventkit/kit/log"
	jsoniter "github.com/json-iterator/go"
	"runtime/debug"
	"strconv"
)

// defines the error message template of transaction SDK
const (
	ConfirmRecoverErr = "Client confirm failed, error=%++v, Stack: %s"
	ConfirmDecodeErr  = "Local transaction confirm: request decode failed, err=%s"
	ConfirmExecuteErr = "Atomic transaction confirm fail, err=%s"
	CancelRecoverErr  = "Client cancel failed, error=%++v, Stack: %s"
	CancelDecodeErr   = "Local transaction cancel: request decode failed, err=%s"
	CancelExecuteErr  = "Atomic transaction cancel fail, err=%s"
	ParseHeaderErr    = "Parse header fail, err=%s"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// TransactionCallbackHandler is a handler that proxy the confirm/cancel method execution of transaction service
type TransactionCallbackHandler struct {
	base.Handler
}

// CallbackConfirm locally executes confirm method
func (c *TransactionCallbackHandler) CallbackConfirm(requestBytes []byte) (responseBody []byte, e error) {
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf(ConfirmRecoverErr, err, string(debug.Stack()))
			response := model.BuildErrorResponse(msg)
			responseBody = c.Resp(response)
		}
	}()

	request := &event.AtomicTxnCallbackRequest{}
	paramData, err := util.Decode(requestBytes, request)
	if err != nil {
		msg := fmt.Sprintf(ConfirmDecodeErr, err)
		response := model.BuildErrorResponse(msg)
		responseBody = c.Resp(response)
		return
	}

	// remove repeat header key
	originHeaders := c.GetRequestHeader()
	headers := request.Request.Headers
	if headers != "" {
		mp := make(map[string]string)
		bytes, err := base64.StdEncoding.DecodeString(headers)
		if err != nil {
			msg := fmt.Sprintf(ParseHeaderErr, err)
			response := model.BuildErrorResponse(msg)
			responseBody = c.Resp(response)
			return
		}
		err = json.Unmarshal(bytes, &mp)
		if err != nil {
			msg := fmt.Sprintf(ParseHeaderErr, err)
			response := model.BuildErrorResponse(msg)
			responseBody = c.Resp(response)
			return
		}
		for key, value := range mp {
			if _, ok := originHeaders[key]; !ok {
				originHeaders[key] = value
			}
		}
	}

	txnCallback := NewTxnCallback()
	errorCode, err := txnCallback.Confirm(c.Ctx, c.RemoteCallInc, request.Request.ServiceName, paramData, originHeaders, c.GetTopicAttributes())
	if err != nil {
		msg := fmt.Sprintf(ConfirmExecuteErr, err)
		var response *model.CommonResponse
		if errorCode != 0 {
			response = model.BuildErrorResponseWithErrorCode(errorCode, msg)
		} else {
			response = model.BuildErrorResponse(msg)
		}
		responseBody = c.Resp(response)
		return
	}
	respBody := event.AtomicTxnConfirmResponseBody{
		RootXid:      request.Request.RootXid,
		BranchXid:    request.Request.BranchXid,
		ResponseTime: util.CurrentTime(),
	}
	response := model.BuildResponse(respBody)
	responseBody = c.Resp(response)
	return
}

// CallbackCancel locally executes cancel method
func (c *TransactionCallbackHandler) CallbackCancel(requestBytes []byte) (responseBody []byte, e error) {
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf(CancelRecoverErr, err, string(debug.Stack()))
			response := model.BuildErrorResponse(msg)
			responseBody = c.Resp(response)
		}
	}()

	request := &event.AtomicTxnCallbackRequest{}
	paramData, err := util.Decode(requestBytes, request)
	if err != nil {
		msg := fmt.Sprintf(CancelDecodeErr, err)
		response := model.BuildErrorResponse(msg)
		responseBody = c.Resp(response)
		return
	}

	// remove repeat header key
	originHeaders := c.GetRequestHeader()
	headers := request.Request.Headers
	if headers != "" {
		mp := make(map[string]string)
		bytes, err := base64.StdEncoding.DecodeString(headers)
		if err != nil {
			msg := fmt.Sprintf(ParseHeaderErr, err)
			response := model.BuildErrorResponse(msg)
			responseBody = c.Resp(response)
			return
		}
		err = json.Unmarshal(bytes, &mp)
		if err != nil {
			msg := fmt.Sprintf(ParseHeaderErr, err)
			response := model.BuildErrorResponse(msg)
			responseBody = c.Resp(response)
			return
		}
		for key, value := range mp {
			if _, ok := originHeaders[key]; !ok {
				originHeaders[key] = value
			}
		}
	}

	txnCallback := NewTxnCallback()
	errorCode, err := txnCallback.Cancel(c.Ctx, c.RemoteCallInc, request.Request.ServiceName, paramData, originHeaders, c.GetTopicAttributes())
	if err != nil {
		msg := fmt.Sprintf(CancelExecuteErr, err)
		var response *model.CommonResponse
		if errorCode != 0 {
			response = model.BuildErrorResponseWithErrorCode(errorCode, msg)
		} else {
			response = model.BuildErrorResponse(msg)
		}
		responseBody = c.Resp(response)
		return
	}
	respBody := event.AtomicTxnConfirmResponseBody{
		RootXid:      request.Request.RootXid,
		BranchXid:    request.Request.BranchXid,
		ResponseTime: util.CurrentTime(),
	}
	response := model.BuildResponse(respBody)
	return c.Resp(response), nil
}

// Resp builds the response for transaction confirm/cancel callback
func (c *TransactionCallbackHandler) Resp(resp *model.CommonResponse) []byte {
	appProps := make(map[string]string)

	if resp.ErrorCode == 0 {
		appProps["RetMsgCode"] = "0"
		appProps["RetMessage"] = ""
	} else {
		log.Errorf(c.Ctx, "transaction callback failed, errMsg:%s", resp.ErrorMsg)
		appProps["RetMsgCode"] = strconv.Itoa(resp.ErrorCode)
		appProps["RetMessage"] = resp.ErrorMsg
	}
	bytes, _ := json.Marshal(resp)
	log.Debugf(c.Ctx, "transaction callback, response=%++v", base64.StdEncoding.EncodeToString(bytes))
	return bytes
}
