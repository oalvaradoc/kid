package compensable

import (
	"testing"
)

func TestCompensable(t *testing.T) {
	compensable := Compensable{}
	f := func() error {
		t.Logf("run registered function...")
		return nil
	}
	compensable.DoUntilHasSucceeded(f)
	compensable.DoUntilHasSucceeded(f)
}
