package util

import (
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"sync"
	"time"
)

// Future is a controller that used for non-blocking executor.
type Future struct {
	isFinished bool
	result     interface{}
	resultChan chan interface{}
	l          sync.Mutex
}

// GetResult is a blocking function for get the result of future execution with timeout
func (f *Future) GetResult(timeoutMilliseconds int) interface{} {
	f.l.Lock()
	defer f.l.Unlock()
	if f.isFinished {
		return f.result
	}
	timer := AcquireTimer(time.Millisecond * time.Duration(timeoutMilliseconds))
	defer ReleaseTimer(timer)

	select {
	case <-timer.C:
		f.isFinished = true
		f.result = nil
		return errors.Errorf(constant.SystemInternalError, "get result timeout: %d ms", timeoutMilliseconds)
	case f.result = <-f.resultChan:
		f.isFinished = true
		return f.result
	}

}

// SetResult is used to set the result of future execution
func (f *Future) SetResult(result interface{}) {
	if f.isFinished {
		return
	}
	f.resultChan <- result
	close(f.resultChan)
}

// NewFuture returns a new future.
func NewFuture() *Future {
	return &Future{
		isFinished: false,
		result:     nil,
		resultChan: make(chan interface{}, 1),
	}
}
