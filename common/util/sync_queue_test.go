package util

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestNewSyncQueue(t *testing.T) {
	queue := NewSyncQueue()
	assert.NotNil(t, queue)
}

func TestPop(t *testing.T) {
	queue := NewSyncQueue()
	queue.Push("test string")
	v := queue.Pop()
	assert.NotNil(t, v)
}

func TestTryPop(t *testing.T) {
	queue := NewSyncQueue()
	r, ok := queue.TryPop()
	assert.False(t, ok)
	assert.Nil(t, r)
}

func TestPush(t *testing.T) {
	queue := NewSyncQueue()
	assert.Equal(t, queue.Len(), 0)

	queue.Push("test")
	assert.Equal(t, queue.Len(), 1)
}

func TestLen(t *testing.T) {
	queue := NewSyncQueue()
	assert.Equal(t, queue.Len(), 0)

	queue.Push("test")
	assert.Equal(t, queue.Len(), 1)

	queue.Pop()
	assert.Equal(t, queue.Len(), 0)
}
