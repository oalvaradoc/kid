package callback

import (
	"testing"
)

func TestRunHookHandleFunc(t *testing.T) {
	RegisterInitHookFunc("test1", func() error {
		t.Logf("run the hook function of test1")
		return nil
	}, true)

	RegisterInitHookFunc("test2", func() error {
		t.Logf("run the hook function of test2")
		return nil
	}, true)

	RunHookHandleFunc()
}
