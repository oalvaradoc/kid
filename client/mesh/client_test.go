package mesh

import (
	"context"
	"fmt"
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/sed/callback"
	"testing"
	"time"
)

func TestNewMeshRequestIsNil(t *testing.T) {
	t.Skip()
	callback.SetSedServerAddr("http://127.0.0.1:19801")
	client := NewMeshClient()

	request := NewMeshRequest(nil)
	request.WithOptions(
		WithTopicTypeOps(),             // mark topic type to OPS
		WithORG("ORG001"),              // org id
		WithWorkspace("workspace"),     // workspace
		WithEnvironment("environment"), // environment
		WithSU("SU001"),                // su
		WithEventID("Event001"),        // dst event id
		WithMaxRetryTimes(2),           // retry times
		WithMaxWaitingTime(4*time.Second),
		WithTimeout(2*time.Second),
		WithHeader(map[string]string{
			"test": constant.ClientSDKVersionV2,
		}),
	)

	// sync call
	_, err := client.SyncCall(context.Background(), request, nil)
	if nil != err {
		t.Error(err)
	}
}

func TestClusterMembers(t *testing.T) {
	t.Skip()
	callback.SetSedServerAddr("http://127.0.0.1:19801")
	c := NewMeshClient()
	request := NewMeshRequest(nil, func(options *client.RequestOptions) {
		options.EventID = "Rdbmysql00001_ClusterMembers"
		options.TopicType = "OPD"
		options.Header = map[string]string{}
	})
	var response *msg.Message
	_, err := c.SyncCall(context.Background(), request, response)
	if err != nil {
		t.Errorf("Failed to execute, error=%++v", err)
	}
	fmt.Println(response)
}
