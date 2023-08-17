package util

import (
	"fmt"
)

// SUCCESS is a code of result that represents the API executed successfully.
const SUCCESS = 0

// Result is a struct that used to define a result to the API.
type Result struct {
	Code    int         `json:"errorCode"`
	Message string      `json:"errorMsg"`
	Data    interface{} `json:"response"`
}

// Success check whether the code of result is SUCCESS
func (r *Result) Success() bool {
	return r.Code == SUCCESS
}

// NewSuccessResult returns a result that means success, the code of result is SUCCESS
func NewSuccessResult(data interface{}, msg ...interface{}) *Result {
	message := "Done"
	if len(msg) > 0 {
		message = fmt.Sprint(msg...)
	}
	return &Result{
		Code:    SUCCESS,
		Data:    data,
		Message: message,
	}
}

// NewFailResult returns a Result that contains the error code and error messages
func NewFailResult(errorCode int, msg string) *Result {
	return &Result{
		Code:    errorCode,
		Data:    nil,
		Message: msg,
	}
}

// NewFailResultf returns a Result that contains the error code and error messages
func NewFailResultf(errorCode int, format string, args ...interface{}) *Result {
	return &Result{
		Code:    errorCode,
		Data:    nil,
		Message: fmt.Sprintf(format, args...),
	}
}

// TaskResult is a model for define a task result
type TaskResult struct {
	TaskID    int64
	Status    string
	Action    string
	ErrorCode int
	Progress  float64
}

// OperateResult is a model for define an operate result
type OperateResult struct {
	Success int
	Fail    int
}
