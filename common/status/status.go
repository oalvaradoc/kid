package status

import (
	"git.multiverse.io/eventkit/kit/log"
	"sync"
	"time"
)

var serverStatus int
var serverStatusLock sync.RWMutex

var clientStatus int
var clientStatusLock sync.RWMutex

// ClientVersion is a global variable that keeps the current SDK version
var ClientVersion string

// ClientProtocolLevel is a global variable keeps the protocol level that current SDK support.
var ClientProtocolLevel int

// ServerVersion is global variable that keeps the version number of SED server, Only used when the service enable FSM.
var ServerVersion string

// ServerProtocolLevel is global variable that keeps the protocol level of SED server, Only used when the service enable FSM.
var ServerProtocolLevel int

// The following elements is represent all of the client side FSM status.
const (
	ClientStop = iota
	ClientPreStop
	ClientStartingInit
	ClientStartingRunningHook
	ClientStarted
)

var clientStatusString = map[int]string{
	ClientStop:                "STOP",
	ClientPreStop:             "PRE_STOP",
	ClientStartingInit:        "STARTING-INIT",
	ClientStartingRunningHook: "STARTING-RUNNING-HOOK",
	ClientStarted:             "STARTED",
}

// The following elements is represent all of the server side FSM status.
const (
	ServerStop = iota
	ServerPreStop
	ServerStartingInit
	ServerStartingPreCanSend
	ServerStartingCanSend
	ServerStarted
)

var serverStatusString = map[int]string{
	ServerStop:               "STOP",
	ServerPreStop:            "PRE_STOP",
	ServerStartingInit:       "STARTING-INIT",
	ServerStartingPreCanSend: "STARTING-PRE-CAN-SEND",
	ServerStartingCanSend:    "STARTING-CAN-SEND",
	ServerStarted:            "STARTED",
}

// SetServerStatus is used to setting server status
func SetServerStatus(status int) {
	serverStatusLock.Lock()
	defer serverStatusLock.Unlock()

	if serverStatus != status {
		log.Infosf(">>>SERVER:Change server status from `%s` to `%s`", serverStatusString[serverStatus], serverStatusString[status])
	}

	serverStatus = status
}

// GetServerStatus returns the status value of server side
func GetServerStatus() int {
	serverStatusLock.RLock()
	defer serverStatusLock.RUnlock()
	return serverStatus
}

// SetClientStatus is used to setting client status
func SetClientStatus(status int) {
	clientStatusLock.Lock()
	defer clientStatusLock.Unlock()

	if clientStatus != status {
		log.Infosf(">>>CLIENT:Change client status from `%s` to `%s`", clientStatusString[clientStatus], clientStatusString[status])
	}

	clientStatus = status
}

// GetClientStatus returns the status value of client side
func GetClientStatus() int {
	clientStatusLock.RLock()
	defer clientStatusLock.RUnlock()

	return clientStatus
}

// GetClientStatusString returns the string value of client status
func GetClientStatusString() string {
	return clientStatusString[GetClientStatus()]
}

// GetServerStatusString returns the string value of server status
func GetServerStatusString() string {
	return serverStatusString[GetServerStatus()]
}

// ConvertServerStatusListToString is a tool function for convert server status list as a string
func ConvertServerStatusListToString(statusList ...int) string {
	res := ""
	for _, status := range statusList {
		res += "`" + serverStatusString[status] + "`,"
	}

	if len(res) > 0 {
		res = res[:len(res)-1]
	}
	return res
}

// ConvertClientStatusListToString is a tool function for convert client status list as a string
func ConvertClientStatusListToString(statusList ...int) string {
	res := ""
	for _, status := range statusList {
		res += "`" + clientStatusString[status] + "`,"
	}

	if len(res) > 0 {
		res = res[:len(res)-1]
	}
	return res
}

// IsInStatus is a tool function for check the status whether is in the target status list.
func IsInStatus(originStatus int, statusList ...int) bool {
	for _, status := range statusList {
		if originStatus == status {
			return true
		}
	}

	return false
}

// Response is used to client and server to obtain each other's status information.
type Response struct {
	Status        int
	Version       string
	ProtocolLevel int
}

// WaitingForServerStatus is used to waiting the server status change to the specified value
func WaitingForServerStatus(serverSideStatusFSMShutdown chan struct{}, checkInterval time.Duration, maxWaitTime time.Duration, statusList ...int) {
	log.Infosf("Start waiting for Server's status to [%s]", ConvertServerStatusListToString(statusList...))
	t := time.NewTicker(checkInterval)
	startTime := time.Now()
	for {
		select {
		case <-serverSideStatusFSMShutdown:
			log.Infosf("Server side FSM shutdown, stop waiting and exit!")
			return
		case <-t.C:
			if maxWaitTime != -1 && time.Now().Sub(startTime) > maxWaitTime {
				log.Errorsf("Server status Waiting timeout %s, exit!", maxWaitTime.String())
				return
			}
			if IsInStatus(GetServerStatus(), statusList...) {
				log.Infosf("Server's status is `%s`, pass!", GetServerStatusString())
				return
			}
			log.Infosf("Server's status is `%s`, waiting...", GetServerStatusString())
		}
	}
}

// WaitingForClientStatus is used to waiting the client status change to the specified value
func WaitingForClientStatus(clientSideStatusFSMShutdown chan struct{}, checkInterval time.Duration, maxWaitTime time.Duration, statusList ...int) {
	log.Infosf("Start waiting for Client's status to [%s]", ConvertClientStatusListToString(statusList...))
	t := time.NewTicker(checkInterval)
	startTime := time.Now()
	for {
		select {
		case <-clientSideStatusFSMShutdown:
			log.Infosf("Client side FSM shutdown, stop waiting and exit!")
			return
		case <-t.C:
			if maxWaitTime != -1 && time.Now().Sub(startTime) > maxWaitTime {
				log.Errorsf("Server status Waiting timeout %s, exit!", maxWaitTime.String())
				return
			}
			if IsInStatus(GetClientStatus(), statusList...) {
				log.Infosf("Client's status is `%s`, pass!", GetClientStatusString())
				return
			}
			log.Infosf("Client's status is `%s`, waiting...", GetClientStatusString())
		}
	}
}
