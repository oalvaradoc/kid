package callback

import (
	"encoding/json"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
	"reflect"
	"time"
)

var configVersion int64
var kitVersion string

// GetConfigVersion return the version number of config
func GetConfigVersion() int64 {
	return configVersion
}

// SetConfigVersion sets the version number of config
func SetConfigVersion(version int64) {
	configVersion = version
}

func SetKitVersion(v string) {
	kitVersion = v
}

// HeartBeatBody is a model for heartbeat
type HeartBeatBody struct {
	ConfigVersion int64 `json:"configVersion"`
}

var heartBeatConfig *HeartBeat

// SetHeartBeatConfig sets the config of heartbeat
func SetHeartBeatConfig(beat *HeartBeat) {
	heartBeatConfig = beat
}

// GetCurrentHeartBeatConfig returns the config of heartbeat
func GetCurrentHeartBeatConfig() *HeartBeat {
	return heartBeatConfig
}

// startRunHeartBeat is an internal program for looping up heartbeat messages.
// default session set to "default" if not specified
// default interval set to 15 if not specified
// default heartbeat collect topic set to "rapm000" if not specified
func startRunHeartBeat(beat *HeartBeat) {
	heartBeatConfig = beat
	time.Sleep(time.Duration(heartBeatConfig.Interval) * time.Second)
	// forever loop
	for {
		session := heartBeatConfig.Session
		if session == "" {
			session = "default"
		}

		heartBeatTopicAttribute := map[string]string{
			constant.TopicType:          constant.TopicTypeHeartbeat,
			constant.TopicID:            heartBeatConfig.EventID,
			constant.TopicDestinationSU: "*",
		}
		message := &msg.Message{
			SessionName:    session,
			TopicAttribute: heartBeatTopicAttribute,
		}
		body := HeartBeatBody{
			ConfigVersion: GetConfigVersion(),
		}
		message.SetAppProps(map[string]string{
			"service.lang.type":        "golang",
			"version.eventkit.kit":     kitVersion,
		})
		bodyBs, err := json.Marshal(body)
		if nil != err {
			log.Errorsf("Failed to marshal heart body, error:%s", errors.Wrap(constant.SystemInternalError, err, 0).ErrorStack())
		} else {
			message.Body = bodyBs
			if err = Publish(message); err != nil {
				log.Errorsf("SendHeartbeat meet err=%s", errors.Wrap(constant.SystemInternalError, err, 0).ErrorStack())
			}
		}
		time.Sleep(time.Duration(heartBeatConfig.Interval) * time.Second)
	}
}

// StartHeartBeat is use for service start a heart beat goroutine
func StartHeartBeat(heartBeatEventID string, heartBeatInterval int) {
	log.Infosf("Start heartbeat with interval:%d seconds, target event ID:%s", heartBeatInterval, heartBeatEventID)
	go startRunHeartBeat(&HeartBeat{
		EventID:  heartBeatEventID,
		Interval: heartBeatInterval,
	})
}

// HeartBeat is a internal struct for store session name and topic id and interval
type HeartBeat struct {
	Session  string
	EventID  string
	Interval int
}

// Equals returns whether the self and other are equals
func (h HeartBeat) Equals(o *HeartBeat) bool {
	return reflect.DeepEqual(&h, o)
}
