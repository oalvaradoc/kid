package compensable

import (
	"sync"
	"sync/atomic"
)

// Compensable is a controller that contains all the transaction information.
type Compensable struct {
	sync.Mutex
	TryMethod          string
	ConfirmMethod      string
	CancelMethod       string
	ServiceName        string
	done               uint32
	CompensableFlagSet int
	IsPropagator       bool
}

func (c *Compensable) doSlow(f func() error) error {
	c.Lock()
	defer c.Unlock()
	if c.done == 0 {
		err := f()
		if nil == err {
			atomic.StoreUint32(&c.done, 1)
		}
		return err
	}

	return nil
}

// DoUntilHasSucceeded processes the closure function until has successed.
func (c *Compensable) DoUntilHasSucceeded(f func() error) error {
	if atomic.LoadUint32(&c.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		return c.doSlow(f)
	}

	return nil
}
