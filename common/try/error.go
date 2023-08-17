package try

import (
	"fmt"
)

const (
	// ErrorCodeLen defines the length of error code
	ErrorCodeLen = 6

	// SuccCode set this code when executed successfully
	SuccCode = 0

	// DefaultCode default code
	DefaultCode = 999999

	// ModuleDefault the default value for module
	ModuleDefault = 10

	// ModuleGls means the GLS module
	ModuleGls = 66
)

// Exception is a model that warps the result of API execution
type Exception struct {
	Code    uint
	Message string
	Data    *interface{}
}

// Error formats the error to string
func (e *Exception) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Successful creates a new wrapped result that represents the successfully to execute API
func Successful(data *interface{}) *Exception {
	return &Exception{
		Code:    SuccCode,
		Message: "O.K.",
		Data:    data,
	}
}

// Throwble creates a new wrapped result that using the default code
func Throwble(err *error) *Exception {
	return Error(DefaultCode, *err)
}

// Errorf creates a new wrapped result that using the input error code and error message
func Errorf(code int, f string, args ...interface{}) *Exception {
	return &Exception{
		Code:    uint(code),
		Message: fmt.Sprintf(f, args...),
		Data:    nil,
	}
}

// Append creates a new wrapped result that using the input error code and error message,
// the built-in error result will append to the error message
func Append(code int, err error, f string, args ...interface{}) *Exception {
	return &Exception{
		Code:    uint(code),
		Message: fmt.Sprintf("%s %s", fmt.Sprintf(f, args...), err.Error()),
		Data:    nil,
	}
}

// Error creates a new wrapped result using error code and build-in error
func Error(c int, e error) *Exception {
	if e == nil {
		return nil
	}
	switch e.(type) {
	case *Exception:
		return e.(*Exception)
	default:
		return &Exception{
			Code:    uint(c),
			Message: e.Error(),
			Data:    nil,
		}
	}
}
