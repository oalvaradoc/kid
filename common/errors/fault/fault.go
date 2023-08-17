package fault

import (
	"fmt"
	"git.multiverse.io/eventkit/kit/common/util"
)

// Fault is a common response model without data, only contains error code and error message
type Fault struct {
	Message string
	Code    int
}

// Error gets the error of Fault
func (f Fault) Error() string {
	return f.Message
}

// Result returns a Result that contains the error code and error messages
func (f Fault) Result() *util.Result {
	return util.NewFailResult(f.Code, f.Message)
}

// NewFault creates new Fault with error code and error message
func NewFault(code int, msg string) *Fault {
	return &Fault{
		Message: msg,
		Code:    code,
	}
}

// NewFaultf creates new Fault with error code and error message
func NewFaultf(code int, format string, args ...interface{}) *Fault {
	return &Fault{
		Message: fmt.Sprintf(format, args...),
		Code:    code,
	}
}
