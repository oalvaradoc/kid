package callback

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestNewHandlerOptions(t *testing.T) {
	opts := NewHandlerOptions()
	assert.Equal(t, opts.Port, DefaultPort)
	assert.Equal(t, opts.CommType, DefaultCommType)
	assert.Equal(t, opts.ServerAddress, DefaultServerAddress)
	assert.Equal(t, opts.CallbackPort, DefaultCallbackPort)
	assert.Equal(t, opts.EnableClientSideStatusFSM, DefaultEnableClientSideStatusFSM)
	assert.Equal(t, len(opts.ExtConfigs), 0)
}

func TestWithPort(t *testing.T) {
	opts := NewHandlerOptions()
	opt := WithPort(1000)
	opt(&opts)

	assert.Equal(t, opts.Port, 1000)
}

func TestWithCommType(t *testing.T) {
	opts := NewHandlerOptions()
	opt := WithCommType("comm")
	opt(&opts)

	assert.Equal(t, opts.CommType, "comm")
}

func TestWithServerAddress(t *testing.T) {
	opts := NewHandlerOptions()
	opt := WithServerAddress("127.0.0.2")
	opt(&opts)

	assert.Equal(t, opts.ServerAddress, "127.0.0.2")
}

func TestWithCallbackPort(t *testing.T) {
	opts := NewHandlerOptions()
	opt := WithCallbackPort(2000)
	opt(&opts)

	assert.Equal(t, opts.CallbackPort, 2000)
}

func TestWithEnableClientSideStatusFSM(t *testing.T) {
	opts := NewHandlerOptions()
	opt := WithEnableClientSideStatusFSM(true)
	opt(&opts)

	assert.Equal(t, opts.EnableClientSideStatusFSM, true)
}

func TestWithExtConfigs(t *testing.T) {
	opts := NewHandlerOptions()
	opt := WithExtConfigs(map[string]interface{}{"test-key": "test-value"})
	opt(&opts)

	assert.Equal(t, len(opts.ExtConfigs), 1)
	assert.Equal(t, opts.ExtConfigs["test-key"], "test-value")
}

func TestAddExtConfig(t *testing.T) {
	opts := NewHandlerOptions()
	opt := AddExtConfig("test-key", "test-value")
	opt(&opts)

	assert.Equal(t, len(opts.ExtConfigs), 1)
	assert.Equal(t, opts.ExtConfigs["test-key"], "test-value")
}
