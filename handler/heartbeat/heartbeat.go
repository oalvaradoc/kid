package heartbeat

import (
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/log"
	"git.multiverse.io/eventkit/kit/sed/callback"
)

// StartHeartbeat creates and starts a coroutine to report heartbeat information regularly
func StartHeartbeat(eventID string, intervalSeconds int) error {
	if "" == eventID {
		return errors.Errorf(constant.SystemInternalError, "start heartbeat failed, eventID is empty.")
	}

	if intervalSeconds <= 0 {
		return errors.Errorf(constant.SystemInternalError, "start heartbeat failed, interval seconds must greater then zero.")
	}
	callback.StartHeartBeat(eventID, intervalSeconds)

	config.RegisterConfigOnChangeHookFunc("heartbeat", rotateHeartbeatWhenConfigChanged, false)

	return nil
}

func rotateHeartbeatWhenConfigChanged(oldConfig *config.ServiceConfigs, newConfig *config.ServiceConfigs) error {
	if nil == newConfig {
		return nil
	}
	currentHeartBeatConfig := callback.GetCurrentHeartBeatConfig()
	if currentHeartBeatConfig.EventID != newConfig.Heartbeat.TopicName ||
		currentHeartBeatConfig.Interval != newConfig.Heartbeat.IntervalSeconds {
		heartbeat := &callback.HeartBeat{
			EventID:  newConfig.Heartbeat.TopicName,
			Interval: newConfig.Heartbeat.IntervalSeconds,
		}

		log.Infosf("New heartbeat config:%++v", heartbeat)
		callback.SetHeartBeatConfig(heartbeat)
	}

	return nil
}
