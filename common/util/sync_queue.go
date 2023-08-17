package util

import (
	"sync"

	"gopkg.in/eapache/queue.v1"
)

// SyncQueue FIFO queue
type SyncQueue struct {
	lock    sync.Mutex
	popable *sync.Cond
	buffer  *queue.Queue
	closed  bool
}

// NewSyncQueue is used to create a new SyncQueue
func NewSyncQueue() *SyncQueue {
	ch := &SyncQueue{
		buffer: queue.New(),
	}
	ch.popable = sync.NewCond(&ch.lock)
	return ch
}

// Pop an item from SyncQueue, will block if SyncQueue is empty
func (q *SyncQueue) Pop() (v interface{}) {
	c := q.popable
	buffer := q.buffer

	q.lock.Lock()
	for buffer.Length() == 0 && !q.closed {
		c.Wait()
	}

	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
	}

	q.lock.Unlock()
	return
}

// TryPop is try to pop an item from SyncQueue, will return immediately with bool=false if SyncQueue is empty
func (q *SyncQueue) TryPop() (v interface{}, ok bool) {
	buffer := q.buffer

	q.lock.Lock()

	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
		ok = true
	} else if q.closed {
		ok = true
	}

	q.lock.Unlock()
	return
}

// Push an item to SyncQueue. Always returns immediately without blocking
func (q *SyncQueue) Push(v interface{}) {
	q.lock.Lock()
	if !q.closed {
		q.buffer.Add(v)
		q.popable.Signal()
	}
	q.lock.Unlock()
}

// Len is used to get the length of SyncQueue
func (q *SyncQueue) Len() (l int) {
	q.lock.Lock()
	l = q.buffer.Length()
	q.lock.Unlock()
	return
}

// Close is used to close SyncQueue
//
// After close, Pop will return nil without block, and TryPop will return v=nil, ok=True
func (q *SyncQueue) Close() {
	q.lock.Lock()
	if !q.closed {
		q.closed = true
		q.popable.Signal()
	}
	q.lock.Unlock()
}
