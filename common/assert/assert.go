package assert

import (
	"reflect"
)

type TestCommon interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
}

func NotEqual(t TestCommon, a, b interface{}) {
	t.Helper()
	if reflect.DeepEqual(a, b) {
		t.Errorf("Equal. %++v %++v", a, b)
	}
}

func Equal(t TestCommon, a, b interface{}) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Not Equal. %++v %++v", a, b)
	}
}

func True(t TestCommon, a bool) {
	t.Helper()
	if !a {
		t.Errorf("Not True %++v", a)
	}
}

func False(t TestCommon, a bool) {
	t.Helper()
	if a {
		t.Errorf("Not True %++v", a)
	}
}

func Nil(t TestCommon, a interface{}) {
	t.Helper()
	if IsNil(a) {
		t.Error("Not Nil")
	}
}

func NotNil(t TestCommon, a interface{}) {
	t.Helper()
	if IsNil(a) {
		t.Error("Is Nil")
	}
}

// IsNil returns true if input parameter is nil
func IsNil(i interface{}) bool {
	if nil == i {
		return true
	}
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
