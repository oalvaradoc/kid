package register

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/compensable"
	"git.multiverse.io/eventkit/kit/handler/base"
	"testing"
)

type Request struct{}

type Response struct{}

var myCompensable = compensable.Compensable{
	TryMethod:     "TryMethod",
	ConfirmMethod: "ConfirmMethod",
	CancelMethod:  "CancelMethod",
	ServiceName:   "myService",
}

var myCompensable1 = compensable.Compensable{
	TryMethod:     "TryMethod",
	ConfirmMethod: "ConfirmMethod",
	CancelMethod:  "CancelMethod",
}

var _myCompensable = compensable.Compensable{
	TryMethod:     "MyTryMethod",
	ConfirmMethod: "MyConfirmMethod",
	CancelMethod:  "MyCancelMethod",
}

var myCompensable2 = compensable.Compensable{
	TryMethod:    "TryMethod",
	CancelMethod: "CancelMethod",
}

var myCompensable3 = compensable.Compensable{
	TryMethod:     "TryMethod",
	ConfirmMethod: "ConfirmMethod",
}

var myCompensable4 = compensable.Compensable{
	TryMethod: "TryMethod",
}

var errorCompensable1 = compensable.Compensable{
	TryMethod:     "Does not exist",
	ConfirmMethod: "ConfirmMethod",
	CancelMethod:  "CancelMethod",
}

var errorCompensable2 = compensable.Compensable{
	TryMethod:     "TryMethod",
	ConfirmMethod: "Does not exist",
	CancelMethod:  "CancelMethod",
}

var errorCompensable3 = compensable.Compensable{
	TryMethod:     "TryMethod",
	ConfirmMethod: "ConfirmMethod",
	CancelMethod:  "Does not exist",
}

var errorCompensable4 = compensable.Compensable{
	TryMethod:     "TryMethod",
	ConfirmMethod: "ConfirmMethod",
	CancelMethod:  "CancelMethod2",
}

var errorCompensable5 = compensable.Compensable{
	TryMethod:     "TryMethod",
	ConfirmMethod: "ConfirmMethod2",
	CancelMethod:  "CancelMethod",
}

var errorCompensable6 = compensable.Compensable{
	TryMethod:     "TryMethod",
	ConfirmMethod: "ConfirmMethod3",
	CancelMethod:  "CancelMethod",
}

var errorCompensable7 = compensable.Compensable{
	TryMethod:     "TryMethod",
	ConfirmMethod: "ConfirmMethod",
	CancelMethod:  "CancelMethod3",
}

type MyHandler struct {
	base.Handler
}

func (t MyHandler) MyTryMethod(request *Request) (response *Response, err *errors.Error) {
	return nil, nil
}

func (t MyHandler) MyConfirmMethod(request *Request) (response *Response, err *errors.Error) {
	return nil, nil
}

func (t MyHandler) MyCancelMethod(request *Request) (response *Response, err *errors.Error) {
	return nil, nil
}

func (t *MyHandler) TryMethod(request *Request) (response *Response, err *errors.Error) {
	return nil, nil
}

func (t *MyHandler) ConfirmMethod(request *Request) (response *Response, err *errors.Error) {
	return nil, nil
}

func (t *MyHandler) ConfirmMethod2() (response *Response, err *errors.Error) {
	return nil, nil
}

func (t *MyHandler) ConfirmMethod3(request string) (response *Response, err *errors.Error) {
	return nil, nil
}

func (t *MyHandler) CancelMethod(request *Request) (response *Response, err *errors.Error) {
	return nil, nil
}

func (t *MyHandler) CancelMethod2(request *Request) (err *errors.Error) {
	return nil
}

func (t *MyHandler) CancelMethod3(request string) (err *errors.Error) {
	return nil
}

func TestCompensableService(t *testing.T) {
	err := CompensableService(&MyHandler{}, myCompensable)
	assert.Nil(t, err)

	err = CompensableService(&MyHandler{}, myCompensable1)
	assert.Nil(t, err)

	err = CompensableService(MyHandler{}, _myCompensable)
	assert.Nil(t, err)

	err = CompensableService(&MyHandler{}, myCompensable2)
	assert.Nil(t, err)

	err = CompensableService(&MyHandler{}, myCompensable3)
	assert.Nil(t, err)

	err = CompensableService(&MyHandler{}, myCompensable4)
	assert.Nil(t, err)

	err = CompensableService(&MyHandler{}, errorCompensable1)
	assert.NotNil(t, err)
	t.Logf("the result of error:%++v", err)

	err = CompensableService(&MyHandler{}, errorCompensable2)
	assert.NotNil(t, err)
	t.Logf("the result of error:%++v", err)

	err = CompensableService(&MyHandler{}, errorCompensable3)
	assert.NotNil(t, err)
	t.Logf("the result of error:%++v", err)

	err = CompensableService(&MyHandler{}, errorCompensable4)
	assert.NotNil(t, err)
	t.Logf("the result of error:%++v", err)

	err = CompensableService(&MyHandler{}, errorCompensable5)
	assert.NotNil(t, err)
	t.Logf("the result of error:%++v", err)

	err = CompensableService(&MyHandler{}, errorCompensable6)
	assert.NotNil(t, err)
	t.Logf("the result of error:%++v", err)

	err = CompensableService(&MyHandler{}, errorCompensable7)
	assert.NotNil(t, err)
	t.Logf("the result of error:%++v", err)

	v := GetCompensableService(myCompensable.ServiceName)
	assert.NotNil(t, v)
	t.Logf("the result of service:%++v", v)

}
