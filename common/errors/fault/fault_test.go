package fault

import (
	"testing"
)

var fault = Fault{
	Code:    100,
	Message: "error message",
}

func TestFault_Error(t *testing.T) {
	if "error message" != fault.Error() {
		t.Errorf("TestFault_Error failed")
	}
}

func TestFault_Result(t *testing.T) {
	result := fault.Result()

	if nil == result {
		t.Errorf("Get result failed!")
	}
}

func TestNewFault(t *testing.T) {
	f := NewFault(100, "error message")

	if f.Code != fault.Code || f.Message != fault.Message {
		t.Errorf("testing NewFault failed!")
	}
}

func TestNewFaultf(t *testing.T) {
	f := NewFaultf(100, "test [%s]", "boom")

	if f.Code != f.Code || f.Message != "test [boom]" {
		t.Errorf("testing NewFaultf failed!")
	}
}
