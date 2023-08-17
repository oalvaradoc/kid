package errors

import (
	"bytes"
	"fmt"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/common/json"
	"git.multiverse.io/eventkit/kit/log"
	"github.com/beego/i18n"
	"os"
	"reflect"
	"runtime"
	"strings"
	"text/template"
)

type uncaughtPanic struct{ message string }

func (p uncaughtPanic) Error() string {
	return p.message
}

// Error is a custom type of build-in error interface that contains additional field like error code、
// error arguments、error stack and so on.
type Error struct {
	ErrorCode string
	ErrorArgs []interface{}
	Err       error
	stack     []uintptr
	frames    []StackFrame
	prefix    string
	extKV     map[string]string
	extHeader    map[string]string
	duplicateErrorCodeTo string
}

// MaxStackDepth The maximum number of stackframes on any error.
var MaxStackDepth = 50

// NewWithExtKV is a wrapper to wraps the New function for pass external key-value
func NewWithExtKV(extKV map[string]string, errorCode string, e interface{}) *Error {
	re := New(errorCode, e)
	re.extKV = extKV

	return re
}

// NewWithExtKVDpErrorCode is a wrapper to wraps the New function for pass external key-value
func NewWithExtKVDpErrorCode(duplicateErrorCodeTo string, extKV map[string]string, errorCode string, e interface{}) *Error {
	re := New(errorCode, e)
	re.extKV = extKV
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}

// NewWithExtKVAndHeader is a wrapper to wraps the New function for pass external key-value
func NewWithExtKVAndHeader(extHeader, extKV map[string]string, errorCode string, e interface{}) *Error {
	re := New(errorCode, e)
	re.extKV = extKV
	re.extHeader = extHeader
	return re
}


// NewWithExtKVAndHeaderDpErrorCode is a wrapper to wraps the New function for pass external key-value
func NewWithExtKVAndHeaderDpErrorCode(duplicateErrorCodeTo string, extHeader, extKV map[string]string, errorCode string, e interface{}) *Error {
	re := New(errorCode, e)
	re.extKV = extKV
	re.extHeader = extHeader
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}

// NewWithExtHeader is a wrapper to wraps the New function for pass external key-value
func NewWithExtHeader(extHeader map[string]string, errorCode string, e interface{}) *Error {
	re := New(errorCode, e)
	re.extHeader = extHeader
	return re
}

// NewWithExtHeaderDpErrorCode is a wrapper to wraps the New function for pass external key-value
func NewWithExtHeaderDpErrorCode(duplicateErrorCodeTo string, extHeader map[string]string, errorCode string, e interface{}) *Error {
	re := New(errorCode, e)
	re.extHeader = extHeader
	re.duplicateErrorCodeTo = duplicateErrorCodeTo
	return re
}


// NewDpErrorCode is a wrapper to wraps the New function for pass external key-value
func NewDpErrorCode(duplicateErrorCodeTo string, errorCode string, e interface{}) *Error {
	re := New(errorCode, e)
	re.duplicateErrorCodeTo = duplicateErrorCodeTo
	return re
}

// New makes an Error from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The stacktrace will point to the line of code that
// called New.
func New(errorCode string, e interface{}) *Error {
	var err error

	switch e := e.(type) {
	case error:
		err = e
	default:
		err = fmt.Errorf("%v", e)
	}

	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2, stack[:])
	return &Error{
		ErrorCode: errorCode,
		Err:       err,
		stack:     stack[:length],
	}
}

// WrapWithExtKV is a wrapper to wraps the Wrap function for pass external key-value
func WrapWithExtKV(extKV map[string]string, errorCode string, e interface{}, skip int, args ...interface{}) *Error {
	re := Wrap(errorCode, e, skip, args ...)
	re.extKV = extKV

	return re
}

// WrapWithExtKVDpErrorCode is a wrapper to wraps the Wrap function for pass external key-value
func WrapWithExtKVDpErrorCode(duplicateErrorCodeTo string, extKV map[string]string, errorCode string, e interface{}, skip int, args ...interface{}) *Error {
	re := Wrap(errorCode, e, skip, args ...)
	re.extKV = extKV
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}


// WrapWithExtKVAndHeader is a wrapper to wraps the Wrap function for pass external key-value
func WrapWithExtKVAndHeader(extHeader, extKV map[string]string, errorCode string, e interface{}, skip int, args ...interface{}) *Error {
	re := Wrap(errorCode, e, skip, args ...)
	re.extKV = extKV
	re.extHeader = extHeader

	return re
}

// WrapWithExtKVAndHeaderDpErrorCode is a wrapper to wraps the Wrap function for pass external key-value
func WrapWithExtKVAndHeaderDpErrorCode(duplicateErrorCodeTo string, extHeader, extKV map[string]string, errorCode string, e interface{}, skip int, args ...interface{}) *Error {
	re := Wrap(errorCode, e, skip, args ...)
	re.extKV = extKV
	re.extHeader = extHeader
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}

// WrapWithExtHeader is a wrapper to wraps the Wrap function for pass external key-value
func WrapWithExtHeader(extHeader map[string]string, errorCode string, e interface{}, skip int, args ...interface{}) *Error {
	re := Wrap(errorCode, e, skip, args ...)
	re.extHeader = extHeader

	return re
}

// WrapWithExtHeaderDpErrorCode is a wrapper to wraps the Wrap function for pass external key-value
func WrapWithExtHeaderDpErrorCode(duplicateErrorCodeTo string, extHeader map[string]string, errorCode string, e interface{}, skip int, args ...interface{}) *Error {
	re := Wrap(errorCode, e, skip, args ...)
	re.extHeader = extHeader
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}

// WrapDpErrorCode is a wrapper to wraps the Wrap function for pass external key-value
func WrapDpErrorCode(duplicateErrorCodeTo string, errorCode string, e interface{}, skip int, args ...interface{}) *Error {
	re := Wrap(errorCode, e, skip, args ...)
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}

// Wrap makes an Error from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The skip parameter indicates how far up the stack
// to start the stacktrace. 0 is from the current call, 1 from its caller, etc.
func Wrap(errorCode string, e interface{}, skip int, args ...interface{}) *Error {
	//if e == nil {
	//	return nil
	//}

	var err error

	switch e := e.(type) {
	case *Error:
		return e
	case error:
		err = e
	default:
		err = fmt.Errorf("%v", e)
	}

	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2+skip, stack[:])
	return &Error{
		ErrorCode: errorCode,
		ErrorArgs: args,
		Err:       err,
		stack:     stack[:length],
	}
}

// WrapPrefixWithExtKV is a wrapper to wraps the WrapPrefix function for pass external key-value
func WrapPrefixWithExtKV(extKV map[string]string, errorCode string, e interface{}, prefix string, skip int, args ...interface{}) *Error {
	re := WrapPrefix(errorCode, e, prefix, skip, args ...)
	re.extKV = extKV

	return re
}


// WrapPrefixWithExtKVDpErrorCode is a wrapper to wraps the WrapPrefix function for pass external key-value
func WrapPrefixWithExtKVDpErrorCode(duplicateErrorCodeTo string, extKV map[string]string, errorCode string, e interface{}, prefix string, skip int, args ...interface{}) *Error {
	re := WrapPrefix(errorCode, e, prefix, skip, args ...)
	re.extKV = extKV
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}

// WrapPrefixWithExtKVAndHeader is a wrapper to wraps the WrapPrefix function for pass external key-value
func WrapPrefixWithExtKVAndHeader(extHeader, extKV map[string]string, errorCode string, e interface{}, prefix string, skip int, args ...interface{}) *Error {
	re := WrapPrefix(errorCode, e, prefix, skip, args ...)
	re.extKV = extKV
	re.extHeader = extHeader

	return re
}

// WrapPrefixWithExtHeader is a wrapper to wraps the WrapPrefix function for pass external key-value
func WrapPrefixWithExtHeader(extHeader map[string]string, errorCode string, e interface{}, prefix string, skip int, args ...interface{}) *Error {
	re := WrapPrefix(errorCode, e, prefix, skip, args ...)
	re.extHeader = extHeader

	return re
}


// WrapPrefixWithExtHeaderDpErrorCode is a wrapper to wraps the WrapPrefix function for pass external key-value
func WrapPrefixWithExtHeaderDpErrorCode(duplicateErrorCodeTo string, extHeader map[string]string, errorCode string, e interface{}, prefix string, skip int, args ...interface{}) *Error {
	re := WrapPrefix(errorCode, e, prefix, skip, args ...)
	re.extHeader = extHeader
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}

// WrapPrefixDpErrorCode is a wrapper to wraps the WrapPrefix function for pass external key-value
func WrapPrefixDpErrorCode(duplicateErrorCodeTo string, errorCode string, e interface{}, prefix string, skip int, args ...interface{}) *Error {
	re := WrapPrefix(errorCode, e, prefix, skip, args ...)
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}

// WrapPrefix makes an Error from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The prefix parameter is used to add a prefix to the
// error message when calling Error(). The skip parameter indicates how far
// up the stack to start the stacktrace. 0 is from the current call,
// 1 from its caller, etc.
func WrapPrefix(errorCode string, e interface{}, prefix string, skip int, args ...interface{}) *Error {
	//if e == nil {
	//	return nil
	//}

	err := Wrap(errorCode, e, 1+skip)

	if err.prefix != "" {
		prefix = fmt.Sprintf("%s: %s", prefix, err.prefix)
	}

	return &Error{
		Err:    err.Err,
		ErrorArgs: args,
		stack:  err.stack,
		prefix: prefix,
	}

}

// Is detects whether the error is equal to a given error. Errors
// are considered equal by this function if they are the same object,
// or if they both contain the same error inside an errors.Error.
func Is(e error, original error) bool {

	if e == original {
		return true
	}

	if e, ok := e.(*Error); ok {
		return Is(e.Err, original)
	}

	if original, ok := original.(*Error); ok {
		return Is(e, original.Err)
	}

	return false
}

// ErrorfWithExtKV is a wrapper to wraps the Errorf function for pass external key-value
func ErrorfWithExtKV(extKV map[string]string, errorCode string, format string, a ...interface{}) *Error {
	re := Errorf(errorCode, format, a ...)
	re.extKV = extKV

	return re
}

// ErrorfWithExtKVDpErrorCode is a wrapper to wraps the Errorf function for pass external key-value
func ErrorfWithExtKVDpErrorCode(duplicateErrorCodeTo string, extKV map[string]string, errorCode string, format string, a ...interface{}) *Error {
	re := Errorf(errorCode, format, a ...)
	re.extKV = extKV
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}

// ErrorfWithExtKVAndHeader is a wrapper to wraps the Errorf function for pass external key-value
func ErrorfWithExtKVAndHeader(extHeader, extKV map[string]string, errorCode string, format string, a ...interface{}) *Error {
	re := Errorf(errorCode, format, a ...)
	re.extKV = extKV
	re.extHeader = extHeader

	return re
}


// ErrorfWithExtKVAndHeaderDpErrorCode is a wrapper to wraps the Errorf function for pass external key-value
func ErrorfWithExtKVAndHeaderDpErrorCode(duplicateErrorCodeTo string, extHeader, extKV map[string]string, errorCode string, format string, a ...interface{}) *Error {
	re := Errorf(errorCode, format, a ...)
	re.extKV = extKV
	re.extHeader = extHeader
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}

// ErrorfWithExtHeader is a wrapper to wraps the Errorf function for pass external key-value
func ErrorfWithExtHeader(extHeader map[string]string, errorCode string, format string, a ...interface{}) *Error {
	re := Errorf(errorCode, format, a ...)
	re.extHeader = extHeader

	return re
}


// ErrorfWithExtHeaderDpErrorCode is a wrapper to wraps the Errorf function for pass external key-value
func ErrorfWithExtHeaderDpErrorCode(duplicateErrorCodeTo string, extHeader map[string]string, errorCode string, format string, a ...interface{}) *Error {
	re := Errorf(errorCode, format, a ...)
	re.extHeader = extHeader
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}


// ErrorfDpErrorCode is a wrapper to wraps the Errorf function for pass external key-value
func ErrorfDpErrorCode(duplicateErrorCodeTo string, errorCode string, format string, a ...interface{}) *Error {
	re := Errorf(errorCode, format, a ...)
	re.duplicateErrorCodeTo = duplicateErrorCodeTo

	return re
}

// Errorf creates a new error with the given message. You can use it
// as a drop-in replacement for fmt.Errorf() to provide descriptive
// errors in return values.
func Errorf(errorCode string, format string, a ...interface{}) *Error {
	return Wrap(errorCode, fmt.Errorf(format, a...), 1, a...)
}

// Error returns the underlying error's message.
func (err *Error) Error() string {
	if nil == err || nil == err.Err {
		// "Error is nil, please check your error type and go to `https://golang.org/doc/faq#nil_error` know more golang features"
		return "<nil>"
	}
	msg := err.Err.Error()
	if err.prefix != "" {
		msg = fmt.Sprintf("%s: %s", err.prefix, msg)
	}

	return msg
}

// Stack returns the callstack formatted the same way that go does
// in runtime/debug.Stack()
func (err *Error) Stack() []byte {
	buf := bytes.Buffer{}

	for _, frame := range err.StackFrames() {
		buf.WriteString(frame.String())
	}

	return buf.Bytes()
}

// Callers satisfies the bugsnag ErrorWithCallerS() interface
// so that the stack can be read out.
func (err *Error) Callers() []uintptr {
	return err.stack
}

// ErrorStack returns a string that contains both the
// error message and the callstack.
func (err *Error) ErrorStack() string {
	return err.TypeName() + " " + err.Error() + "\n" + string(err.Stack())
}

// StackFrames returns an array of frames containing information about the
// stack.
func (err *Error) StackFrames() []StackFrame {
	if err.frames == nil {
		err.frames = make([]StackFrame, len(err.stack))

		for i, pc := range err.stack {
			err.frames[i] = NewStackFrame(pc)
		}
	}

	return err.frames
}

// TypeName returns the type this error. e.g. *errors.stringError.
func (err *Error) TypeName() string {
	if _, ok := err.Err.(uncaughtPanic); ok {
		return "panic"
	}
	return reflect.TypeOf(err.Err).String()
}

// GetI18nMessage gets the i18n error message with error code and user lang.
func (err *Error) GetI18nMessage(lang string) string {
	errorMsg := i18n.Tr(lang, "errors."+err.ErrorCode, err.ErrorArgs)
	if errorMsg == "" || strings.Contains(errorMsg, "!(EXTRA") || errorMsg == err.ErrorCode {
		errorMsg = err.Error()
	}
	return errorMsg
}

// GetExtHeader gets the ext header of Error
func (err *Error) GetExtHeader() map[string]string{
	return err.extHeader
}

func (err *Error) GetDuplicateErrorCodeTo() string {
	return err.duplicateErrorCodeTo
}

// WrapError wraps the error using response template and error code
func (err *Error) WrapError(lang string, tpl string, responseDataWhenError ...interface{}) (code string, msg string, body []byte) {
	buf := new(bytes.Buffer)
	errorCode := err.ErrorCode
	if len(errorCode) == 0 {
		log.Errorsf("Error code is empty, setting `constant.SystemInternalError` to return!")
		errorCode = constant.SystemInternalError
	}
	errorMsg := err.GetI18nMessage(lang)
	errorResponse, e := template.New("error").Parse(tpl)
	if nil != e {
		return GetFinalErrorCode(constant.SystemInternalError), e.Error(), []byte(e.Error())
	}
	instanceID := os.Getenv("INSTANCE_ID")
	errorStack := ""
	if "" != instanceID {
		errorStack = fmt.Sprintf("[%s]:%s", instanceID, err.ToStackString())
	} else {
		errorStack = err.ToStackString()
	}
	errInfoMap := make(map[string]string)
	errInfoMap["errorCode"] = GetFinalErrorCode(errorCode)

	errorMsgBytes, e1 := json.Marshal(errorMsg)
	if nil == e1 {
		errInfoMap["errorMsg"] = string(errorMsgBytes[1 : len(errorMsgBytes)-1])
	}

	errorStackBytes, e2 := json.Marshal(errorStack)
	if nil == e2 {
		errInfoMap["errorStack"] = string(errorStackBytes[1 : len(errorStackBytes)-1])
	}
	if len(responseDataWhenError) > 0 && responseDataWhenError[0] != nil {
		errInfoMap["data"] = fmt.Sprintf("%v", responseDataWhenError[0])
	} else {
		errInfoMap["data"] = "null"
	}

	if len(err.extKV) > 0 {
		for k,v := range err.extKV {
			errInfoMap[k] = v
		}
	}

	errorResponse.Execute(buf, errInfoMap)

	return GetFinalErrorCode(errorCode), errorMsg, buf.Bytes()
}

// ToStackString prints the stack information into string
func (err *Error) ToStackString() string {
	return ErrorToString(err)
}

// IsTimeoutError is used to check whether the error is timeout.
func (err *Error) IsTimeoutError() bool {
	return err.ErrorCode == constant.SystemRemoteCallTimeout ||
		err.ErrorCode == constant.SystemMeshRequestReplyTimeout ||
		err.ErrorCode == constant.SystemCallbackAppTimeout
}

// GetErrorCode returns the error code if the input type is Error
func GetErrorCode(err error) string {
	if err == nil {
		return ""
	}

	switch err := err.(type) {
	case *Error:
		return err.ErrorCode
	default:
		return ""
	}

}

// IsTimeoutError  is used to check whether the error is timeout.
func IsTimeoutError(err error) bool {
	var errorCode = GetErrorCode(err)

	return errorCode == constant.SystemRemoteCallTimeout ||
		errorCode == constant.SystemMeshRequestReplyTimeout ||
		errorCode == constant.SystemCallbackAppTimeout
}

// ErrorToString formats the input error into string
func ErrorToString(err error) string {
	var str string
	if err != nil {
		switch err.(type) {
		case *Error:
			if nil != err.(*Error) && nil != err.(*Error).Err {
				str = err.(*Error).ErrorStack()
			} else {
				var buf [2 << 10]byte
				stackFrame := string(buf[:runtime.Stack(buf[:], true)])
				str = fmt.Sprintf("Input Error instance invalid, stack frame=%s, please check,err=%++v", stackFrame, err)
			}
		default:
			str = err.Error()
		}
	}

	return str
}
