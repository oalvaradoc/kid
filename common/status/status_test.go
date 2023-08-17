package status

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestConvertServerStatusListToString(t *testing.T) {
	str := ConvertServerStatusListToString([]int{ServerPreStop, ServerStartingInit, ServerStartingCanSend}...)
	t.Logf("the result of server status list to string:%s", str)
	assert.Equal(t, str, "`PRE_STOP`,`STARTING-INIT`,`STARTING-CAN-SEND`")
}

func TestConvertClientStatusListToString(t *testing.T) {
	str := ConvertClientStatusListToString([]int{ClientStop, ClientPreStop, ClientStartingInit}...)
	t.Logf("the result of client status list to string:%s", str)
	assert.Equal(t, str, "`STOP`,`PRE_STOP`,`STARTING-INIT`")
}

func TestIsInStatus(t *testing.T) {
	assert.True(t, IsInStatus(ClientPreStop, []int{ClientStop, ClientPreStop, ClientStartingInit}...))
	assert.False(t, IsInStatus(ClientStarted, []int{ClientStop, ClientStartingInit}...))
}
