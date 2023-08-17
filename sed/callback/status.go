package callback

import (
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/status"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"net/http"
	"runtime/debug"
	"time"
)

var (
	// HTTPClient global http client object
	clientForGetStatus = &fasthttp.Client{
		// MaxConnsPerHost  default is 512, increase to 16384
		MaxConnsPerHost: 16384,

		// Disable idempotent calls attempts when remote call abnormal.
		// default max attempts is 5 times.
		MaxIdemponentCallAttempts: 1,
	}
)

// postMessage provides a callback application to the message via a callback function
func currentServerSideStatus() (int, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	req := fasthttp.AcquireRequest()
	// when the program exits, immediately release the resources to the connection pool.
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(http.MethodGet)
	req.SetRequestURI(serverAddr + ServerStatusGetPath)
	resp := fasthttp.AcquireResponse()
	// when the program exits, immediately release the resources to the connection pool.
	defer fasthttp.ReleaseResponse(resp)
	clientForGetStatus.ReadTimeout = time.Second * 15
	clientForGetStatus.WriteTimeout = time.Second * 15

	if err := clientForGetStatus.DoTimeout(req, resp, 15*time.Second); err != nil {
		return 0, errors.Wrap(constant.SystemInternalError, err, 0)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return 0, errors.Errorf(constant.SystemInternalError, "postMessage httpResponse.StatusCode != 200, %s", resp.String())
	}
	body := resp.Body()

	currentStatusResponse := &status.Response{}
	if err := json.Unmarshal(body, currentStatusResponse); nil != err {
		return 0, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	status.ServerProtocolLevel = currentStatusResponse.ProtocolLevel

	return currentStatusResponse.Status, nil
}

// clientSideStatusFSM is the executor of client FSM, This function will periodically obtain the status of the server
func clientSideStatusFSM() {
	defer func() {
		if e := recover(); e != nil {
			log.Errorsf("clientSideStatusFSM panic:%v, please check[%s]", e, string(debug.Stack()))
		}
	}()

	log.Infosf("Start client side status FSM, client callback API=[%s]...", serverAddr+ServerStatusGetPath)
	t := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-t.C:
			if currentServerStatus, getCurrentServerStatusError := currentServerSideStatus(); nil != getCurrentServerStatusError {
				log.Warnsf("get server status failed:%v, set server status to `STOP` ", getCurrentServerStatusError)
				status.SetServerStatus(status.ServerStop)
			} else {
				status.SetServerStatus(currentServerStatus)
				log.Debugsf("get server status successfully, current server status is: %s", status.GetServerStatusString())
			}
		}
	}
}

// StartClientSideStatusFSM is used to start client side status FSM
func StartClientSideStatusFSM() error {
	go clientSideStatusFSM()

	WaitingForServerStatusToCanSend()

	status.SetClientStatus(status.ClientStartingRunningHook)

	RunHookHandleFunc()

	status.SetClientStatus(status.ClientStarted)

	return nil
}

// WaitingForServerStatusToCanSend is used to waiting the server side status change to `CanSend`
func WaitingForServerStatusToCanSend() {
	status.WaitingForServerStatus(make(chan struct{}, 0), time.Second*1, -1, status.ServerStartingCanSend, status.ServerStarted)
}

// StopClient is used to stop the client side executor.
func StopClient(maxClientStopTimeout time.Duration, clientSideStatusFSM bool) {
	status.SetClientStatus(status.ClientPreStop)
	defer status.SetClientStatus(status.ClientStop)

	if !clientSideStatusFSM {
		log.Infos("Skip client side status FSM, set server status to `STOP`...")
		status.SetServerStatus(status.ServerStop)
	}

	status.WaitingForServerStatus(make(chan struct{}, 0), time.Second*1, maxClientStopTimeout, status.ServerStop)
}
