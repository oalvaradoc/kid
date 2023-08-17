package manager

import (
	"context"
	"encoding/base64"
	"git.multiverse.io/eventkit/kit/client"
	"git.multiverse.io/eventkit/kit/client/mesh"
	"git.multiverse.io/eventkit/kit/codec"
	codec_json "git.multiverse.io/eventkit/kit/codec/json"
	codec_text "git.multiverse.io/eventkit/kit/codec/text"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/model/transaction"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/contexts"
	"git.multiverse.io/eventkit/kit/handler/config"
	"git.multiverse.io/eventkit/kit/log"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// IsFlag checks whether the flag is in the flag set
func IsFlag(flagSet int, flag int) bool {
	return (flagSet & flag) != 0
}

// NewTxnManager creates a new transaction manager
func NewTxnManager(ctx context.Context, transactionConfig *config.Transaction, client client.Client) TxnManager {
	return &DefaultTxnManager{
		client:            client,
		Ctx:               ctx,
		transactionConfig: transactionConfig,
	}
}

// DefaultTxnManager is a transaction that contains remote call instance and transaction configuration
type DefaultTxnManager struct {
	client            client.Client
	Ctx               context.Context
	transactionConfig *config.Transaction
}

// TxnManager defines interfaces related to transaction management
type TxnManager interface {
	TxnBegin(handlerContexts *contexts.HandlerContexts, serverAddress string, compensableFlagSet int, paramData []byte, serviceName string, headers map[string]string) (string, error)

	TxnJoin(handlerContexts *contexts.HandlerContexts, serverAddress string, rootXid string, parentXid string, compensableFlagSet int, paramData []byte, serviceName string, headers map[string]string) (string, error)

	TxnEnd(handlerContexts *contexts.HandlerContexts, serverAddress string, rootXid string, parentXid string, branchXid string, ok bool, tryReturnError *errors.Error, rootKVToSecondStageHeaders map[string]string) (int, error)
}

// TxnBegin does the transaction begin logic for branch transaction
func (d *DefaultTxnManager) TxnBegin(handlerContexts *contexts.HandlerContexts, serverAddress string, compensableFlagSet int, paramData []byte, serviceName string, headers map[string]string) (rootXid string, err error) {
	serverAddressElem := strings.Split(serverAddress, constant.ParticipantAddressSplitChar)

	isDirectRequest := strings.EqualFold(constant.CommDirect, d.transactionConfig.CommType)

	if (isDirectRequest && len(serverAddressElem) != 4) || (!isDirectRequest && len(serverAddressElem) != 9) {
		return "", errors.Errorf(constant.SystemInternalError, "cannot parse participant address:[%s] in transaction begin, context:[%++v]", serverAddress, d.Ctx)
	}

	isExistConfirmMethod := IsFlag(compensableFlagSet, constant.ConfirmFlag)
	isExistCancelMethod := IsFlag(compensableFlagSet, constant.CancelFlag)

	participantAddress := generateParticipantAddress(d.transactionConfig, isExistConfirmMethod, isExistCancelMethod, isDirectRequest, serverAddressElem)

	headerStr := ""
	finalRequestHeaders := make(map[string]string)
	if nil != headers {
		currentSu := headers[constant.CurrentSU]
		finalRequestHeaders[constant.CurrentSU] = currentSu
	}

	if config.GetConfigs().Transaction.SaveHeaders {
		for k, v := range headers {
			finalRequestHeaders[k] = v
		}
	}

	headerbytes, e := json.Marshal(finalRequestHeaders)
	if nil != e {
		log.Errorsf("Failed to marshal headers:%++v, context:[%++v]", e, d.Ctx)
	}
	headerStr = base64.StdEncoding.EncodeToString(headerbytes)

	// use SpanId as RootXid
	rootXid = handlerContexts.SpanContexts.SpanID
	// use TranceId as ParentXid
	parentXid := handlerContexts.SpanContexts.TraceID
	//make up RootTxnBeginRequest
	transactionRequest := transaction.RootTxnBeginRequest{
		Head: transaction.TxnEventHeader{
			Service: "rootTxnBeginRequest",
		},
		Request: transaction.RootTxnBeginRequestBody{
			ParticipantAddress: participantAddress,
			RequestTime:        util.CurrentTime(),
			ParentXid:          parentXid,
			RootXid:            rootXid,
			BranchXid:          rootXid,
			ServiceName:        serviceName,
			Headers:            headerStr,
		},
	}
	reqBody, err := util.Encode(transactionRequest, paramData)
	if err != nil {
		err = errors.Errorf(constant.SystemInternalError, "encodeReq fail, err:[%++v], context:[%++v]", err, d.Ctx)
		return
	}
	timeoutMilliseconds := handlerContexts.SpanContexts.TimeoutMilliseconds - d.transactionConfig.MaxServiceConsumeMilliseconds
	if timeoutMilliseconds < 0 {
		timeoutMilliseconds = handlerContexts.SpanContexts.TimeoutMilliseconds
	}
	log.Debugf(d.Ctx, "TxnBegin root transaction, request body base64:[%s], timeout:[%d]", base64.StdEncoding.EncodeToString(reqBody), timeoutMilliseconds)
	start := time.Now()
	response := &transaction.RootTxnBeginResponse{}

	request := mesh.NewMeshRequest(reqBody)
	requestOptions := []client.RequestOption{
		mesh.WithTimeout(time.Duration(timeoutMilliseconds) * time.Millisecond), // timeout
		mesh.WithMaxRetryTimes(0), // retry times
		mesh.WithHeader(map[string]string{
			constant.ClientSDKVersion: constant.ClientSDKVersionV2,
		}),
		mesh.WithCodec(codec.BuildCustomCodec(
			&codec_text.Encoder{}, // request with text
			&codec_json.Decoder{}, // response with json
		)),
	}

	if isDirectRequest {
		// DXC server address + "|" + Transaction Begin URL + "|" + Transaction Join URL + "|" + Transaction End URL
		requestOptions = append(requestOptions,
			[]client.RequestOption{
				mesh.WithHTTPRequestInfo(
					serverAddressElem[0]+serverAddressElem[1],
					constant.DefaultHTTPMethodPost,
					"",
				),
			}...)
	} else {
		// ORG + "|" + WORKSPACE + "|" + ENVIRONMENT + "|" + DCN + "|" + NODE ID + "|" + Instance ID + "|" + Transaction Begin TopicID + "|" + Transaction Join TopicID + "|" + Transaction End TopicID
		requestOptions = append(requestOptions,
			[]client.RequestOption{
				mesh.WithTopicTypeDxc(),                    // mark topic type to DXC
				mesh.WithORG(serverAddressElem[0]),         // org id
				mesh.WithWorkspace(serverAddressElem[1]),   // workspace
				mesh.WithEnvironment(serverAddressElem[2]), // environment
				mesh.WithSU(serverAddressElem[3]),          // su
				mesh.WithNodeID(serverAddressElem[4]),      // node id
				mesh.WithInstanceID(serverAddressElem[5]),  // instance id
				mesh.WithEventID(serverAddressElem[6]),     // dst event id
			}...)
	}

	if isMacroService(d.transactionConfig, isDirectRequest, serverAddressElem) {
		requestOptions = requestOptions[:len(requestOptions)-8]
		requestOptions = append(requestOptions,
			[]client.RequestOption{
				mesh.WithHTTPRequestInfo(
					d.transactionConfig.TransactionServer.MacroServiceAddressURL+d.transactionConfig.TransactionServer.TxnBeginURLPath,
					constant.DefaultHTTPMethodPost,
					"",
				),
			}...)
	}

	request.WithOptions(requestOptions...)

	// call server to do end
	_, err = d.client.SyncCall(d.Ctx, request, response)
	useTime := int64(time.Now().Sub(start)) / 1e6
	handlerContexts.SpanContexts.With(contexts.TimeoutMilliseconds(handlerContexts.SpanContexts.TimeoutMilliseconds - int(useTime)))

	if err != nil {
		//err = errors.New(constant.SystemInternalError, err)
		return
	}

	log.Debugf(d.Ctx, "TxnBegin root transaction, response:[%++v]", response)

	if response.ErrorCode != 0 {
		err = errors.Errorf(constant.SystemInternalError, "Transaction begin failed, Msg:[%s]!", response.ErrorMsg)
		return
	}
	rootXid = response.Data.RootXid
	return
}

// TxnJoin does the transaction join logic for branch transaction
func (d *DefaultTxnManager) TxnJoin(handlerContexts *contexts.HandlerContexts, serverAddress string, rootXid string, parentXid string, compensableFlagSet int, paramData []byte, serviceName string, headers map[string]string) (branchXid string, err error) {
	var participantAddress string
	serverAddressElem := strings.Split(serverAddress, constant.ParticipantAddressSplitChar)
	isDirectRequest := strings.EqualFold(constant.CommDirect, d.transactionConfig.CommType)

	if (isDirectRequest && len(serverAddressElem) != 4) || (!isDirectRequest && len(serverAddressElem) != 9) {
		return "", errors.Errorf(constant.SystemInternalError, "cannot parse participant address:[%s] in transaction begin, context:[%++v]", serverAddress, d.Ctx)
	}

	isExistConfirmMethod := IsFlag(compensableFlagSet, constant.ConfirmFlag)
	isExistCancelMethod := IsFlag(compensableFlagSet, constant.CancelFlag)

	participantAddress = generateParticipantAddress(d.transactionConfig, isExistConfirmMethod, isExistCancelMethod, isDirectRequest, serverAddressElem)

	headerStr := ""
	finalRequestHeaders := make(map[string]string)
	if nil != headers {
		currentSu := headers[constant.CurrentSU]
		finalRequestHeaders[constant.CurrentSU] = currentSu
	}

	if config.GetConfigs().Transaction.SaveHeaders {
		for k, v := range headers {
			finalRequestHeaders[k] = v
		}
	}

	headerbytes, e := json.Marshal(finalRequestHeaders)
	if nil != e {
		log.Errorsf("Failed to marshal headers:%++v, context:[%++v]", e, d.Ctx)
	}
	headerStr = base64.StdEncoding.EncodeToString(headerbytes)

	// use SpanId as BranchXid
	branchXid = handlerContexts.SpanContexts.SpanID
	transactionRequest := transaction.BranchTxnJoinRequest{
		Head: transaction.TxnEventHeader{
			Service: "branchTxnJoinRequest",
		},
		Request: transaction.BranchTxnJoinRequestBody{
			ParticipantAddress: participantAddress,
			RootXid:            rootXid,
			ParentXid:          parentXid,
			BranchXid:          branchXid,
			RequestTime:        util.CurrentTime(),
			ServiceName:        serviceName,
			Headers:            headerStr,
		},
	}
	reqBody, err := util.Encode(transactionRequest, paramData)
	if err != nil {
		err = errors.Errorf(constant.SystemInternalError, "encodeReq fail, err:[%++v], context:[%++v]", err, d.Ctx)
		return
	}
	timeoutMilliseconds := handlerContexts.SpanContexts.TimeoutMilliseconds - d.transactionConfig.MaxServiceConsumeMilliseconds
	if timeoutMilliseconds < 0 {
		timeoutMilliseconds = handlerContexts.SpanContexts.TimeoutMilliseconds
	}
	log.Debugf(d.Ctx, "TxnJoin branch transaction, request body base64:[%s], timeout:[%d]", base64.StdEncoding.EncodeToString(reqBody), timeoutMilliseconds)
	start := time.Now()
	response := &transaction.BranchTxnJoinResponse{}
	request := mesh.NewMeshRequest(reqBody)
	requestOptions := []client.RequestOption{
		mesh.WithTimeout(time.Duration(timeoutMilliseconds) * time.Millisecond), // timeout
		mesh.WithMaxRetryTimes(0), // retry times
		mesh.WithHeader(map[string]string{
			constant.ClientSDKVersion: constant.ClientSDKVersionV2,
		}),
		mesh.WithCodec(codec.BuildCustomCodec(
			&codec_text.Encoder{}, // request with text
			&codec_json.Decoder{}, // response with json
		)),
	}

	if isDirectRequest {
		// DXC server address + "|" + Transaction Begin URL + "|" + Transaction Join URL + "|" + Transaction End URL
		requestOptions = append(requestOptions, []client.RequestOption{
			mesh.WithHTTPRequestInfo(
				serverAddressElem[0]+serverAddressElem[2],
				constant.DefaultHTTPMethodPost,
				"",
			),
		}...)
	} else {
		// ORG + "|" + WORKSPACE + "|" + ENVIRONMENT + "|" + DCN + "|" + NODE ID + "|" + Instance ID + "|" + Transaction Begin TopicID + "|" + Transaction Join TopicID + "|" + Transaction End TopicID
		requestOptions = append(requestOptions, []client.RequestOption{
			mesh.WithTopicTypeDxc(),                    // mark topic type to DXC
			mesh.WithORG(serverAddressElem[0]),         // org id
			mesh.WithWorkspace(serverAddressElem[1]),   // workspace
			mesh.WithEnvironment(serverAddressElem[2]), // environment
			mesh.WithSU(serverAddressElem[3]),          // su
			mesh.WithNodeID(serverAddressElem[4]),      // node id
			mesh.WithInstanceID(serverAddressElem[5]),  // instance id
			mesh.WithEventID(serverAddressElem[7]),     // dst event id
		}...)
	}

	if isMacroService(d.transactionConfig, isDirectRequest, serverAddressElem) {
		requestOptions = requestOptions[:len(requestOptions)-8]
		requestOptions = append(requestOptions,
			[]client.RequestOption{
				mesh.WithHTTPRequestInfo(
					d.transactionConfig.TransactionServer.MacroServiceAddressURL+d.transactionConfig.TransactionServer.TxnJoinURLPath,
					constant.DefaultHTTPMethodPost,
					"",
				),
			}...)
	}

	request.WithOptions(requestOptions...)

	// call server to do end
	_, err = d.client.SyncCall(d.Ctx, request, response)

	useTime := int64(time.Now().Sub(start)) / 1e6
	handlerContexts.SpanContexts.With(contexts.TimeoutMilliseconds(handlerContexts.SpanContexts.TimeoutMilliseconds - int(useTime)))

	if err != nil {
		err = errors.New(constant.SystemInternalError, err)
		return
	}

	log.Debugf(d.Ctx, "TxnJoin branch transaction, response:[%++v]", response)

	if response.ErrorCode != 0 {
		err = errors.Errorf(constant.SystemInternalError, "Transaction join failed:[%s], context:[%++v]", response.ErrorMsg, d.Ctx)
		return
	}
	branchXid = response.Data.BranchXid
	return
}

// TxnEnd does the transaction end logic for root transaction
func (d *DefaultTxnManager) TxnEnd(handlerContexts *contexts.HandlerContexts, serverAddress string, rootXid string, parentXid string, branchXid string, ok bool, tryReturnError *errors.Error, rootKVToSecondStageHeaders map[string]string) (errCode int, err error) {
	serverAddressElem := strings.Split(serverAddress, constant.ParticipantAddressSplitChar)
	isDirectRequest := strings.EqualFold(constant.CommDirect, d.transactionConfig.CommType)

	if (isDirectRequest && len(serverAddressElem) != 4) || (!isDirectRequest && len(serverAddressElem) != 9) {
		return -1, errors.Errorf(constant.SystemInternalError, "cannot parse participant address:[%s] in transaction begin, context:[%++v]", serverAddress, d.Ctx)
	}

	transactionRequest := transaction.TxnEndRequest{
		Head: transaction.TxnEventHeader{
			Service: "txnEndRequest",
		},
		Request: transaction.TxnEndRequestBody{
			RootXid:                       rootXid,
			BranchXid:                     branchXid,
			ParentXid:                     parentXid,
			Ok:                            ok,
			RequestTime:                   util.CurrentTime(),
			TryFailedIgnoreCallbackCancel: d.transactionConfig.TryFailedIgnoreCallbackCancel,
		},
	}

	timeoutMilliseconds := handlerContexts.SpanContexts.TimeoutMilliseconds - d.transactionConfig.MaxServiceConsumeMilliseconds
	if timeoutMilliseconds < 0 {
		timeoutMilliseconds = handlerContexts.SpanContexts.TimeoutMilliseconds
	}

	log.Debugf(d.Ctx, "TxnEnd transactionRequest:[%++v], timeout:[%d]", transactionRequest, timeoutMilliseconds)

	response := &transaction.TxnEndResponse{}
	request := mesh.NewMeshRequest(transactionRequest)

	rootErrorCode := ""
	rootErrorMsg := ""
	if tryReturnError != nil {
		rootErrorCode = tryReturnError.ErrorCode
		rootErrorMsg = tryReturnError.Error()
	}
	headers := make(map[string]string)
	headers[constant.ClientSDKVersion] = constant.ClientSDKVersionV2
	headers[constant.RootErrorCode] = rootErrorCode
	headers[constant.RootErrorMsg] = rootErrorMsg
	for k, v := range rootKVToSecondStageHeaders {
		headers[k] = v
	}

	requestOptions := []client.RequestOption{
		mesh.WithTimeout(time.Duration(timeoutMilliseconds) * time.Millisecond),
		mesh.WithMaxRetryTimes(0),
		mesh.WithHeader(headers),
		mesh.WithCodec(codec.BuildJSONCodec()),
	}

	if isDirectRequest {
		// DXC server address + "|" + Transaction Begin URL + "|" + Transaction Join URL + "|" + Transaction End URL
		requestOptions = append(requestOptions, []client.RequestOption{
			mesh.WithHTTPRequestInfo(
				serverAddressElem[0]+serverAddressElem[3],
				constant.DefaultHTTPMethodPost,
				"",
			),
		}...)
	} else {
		// ORG + "|" + WORKSPACE + "|" + ENVIRONMENT + "|" + DCN + "|" + NODE ID + "|" + Instance ID + "|" + Transaction Begin TopicID + "|" + Transaction Join TopicID + "|" + Transaction End TopicID
		requestOptions = append(requestOptions, []client.RequestOption{
			mesh.WithTopicTypeDxc(),
			mesh.WithORG(serverAddressElem[0]),         // org id
			mesh.WithWorkspace(serverAddressElem[1]),   // workspace
			mesh.WithEnvironment(serverAddressElem[2]), // environment
			mesh.WithSU(serverAddressElem[3]),          // su
			mesh.WithNodeID(serverAddressElem[4]),      // node id
			mesh.WithInstanceID(serverAddressElem[5]),  // instance id
			mesh.WithEventID(serverAddressElem[8]),     //	dst event id
		}...)
	}

	if d.transactionConfig.IsMacroService {
		if !isDirectRequest {
			if strings.ReplaceAll(strings.ReplaceAll(serverAddressElem[5], "-", ""), "_", "") == config.GetConfigs().Service.InstanceID {
				requestOptions = requestOptions[:len(requestOptions)-8]
				requestOptions = append(requestOptions,
					[]client.RequestOption{
						mesh.WithHTTPRequestInfo(
							d.transactionConfig.TransactionServer.MacroServiceAddressURL+d.transactionConfig.TransactionServer.TxnEndURLPath,
							constant.DefaultHTTPMethodPost,
							"",
						),
					}...)
			}
		}
	}

	request.WithOptions(requestOptions...)

	// call server to do end
	_, err = d.client.SyncCall(d.Ctx, request, response)

	if err != nil {
		if constant.SystemRemoteCallTimeout == errors.GetErrorCode(err) {
			errCode = constant.TxnEndFailedTimeOut
			return
		}
		errCode = -1
		return
	}

	log.Debugf(d.Ctx, "TxnEnd response:[%++v]", response)

	if response.ErrorCode != 0 {
		err = errors.Errorf(constant.SystemInternalError, "txn do end failed:[%s],context:[%++v]", response.ErrorMsg, d.Ctx)
		errCode = response.ErrorCode
		return
	}
	return
}

func generateParticipantAddress(transactionConfig *config.Transaction, isExistConfirmMethod bool, isExistCancelMethod bool, isDirectRequest bool, serverAddressElem []string) (participantAddress string) {
	branchConfirmAddress := ""
	branchCancelAddress := ""

	if isDirectRequest {
		//get branchConfirmAddress
		if isExistConfirmMethod {
			branchConfirmAddress = transactionConfig.TransactionClient.ConfirmAddressURL
		}
		//get branchCancelAddress
		if isExistCancelMethod {
			branchCancelAddress = transactionConfig.TransactionClient.CancelAddressURL
		}

		//make up participantAddress
		participantAddress =
			transactionConfig.CommType + "|" +
				transactionConfig.TransactionClient.ParticipantAddress + "|" +
				branchConfirmAddress + "|" +
				branchCancelAddress
	} else {
		//get branchConfirmAddress
		if isExistConfirmMethod {
			branchConfirmAddress = transactionConfig.TransactionClient.ConfirmEventID
		}

		//get branchCancelAddress
		if isExistCancelMethod {
			branchCancelAddress = transactionConfig.TransactionClient.CancelEventID
		}

		participantAddress =
			transactionConfig.CommType + "|" +
				branchConfirmAddress + "|" +
				branchCancelAddress
	}

	if isMacroService(transactionConfig, isDirectRequest, serverAddressElem) {
		//get branchConfirmAddress
		if isExistConfirmMethod {
			branchConfirmAddress = transactionConfig.TransactionClient.ConfirmAddressURL
		}
		//get branchCancelAddress
		if isExistCancelMethod {
			branchCancelAddress = transactionConfig.TransactionClient.CancelAddressURL
		}

		//make up participantAddress
		participantAddress =
			constant.CommDirect + "|" +
				transactionConfig.TransactionClient.ParticipantAddress + "|" +
				branchConfirmAddress + "|" +
				branchCancelAddress
	}
	return
}

func isMacroService(transactionConfig *config.Transaction, isDirectRequest bool, serverAddressElem []string) bool {
	if transactionConfig.IsMacroService {
		if !isDirectRequest {
			dxcServerInstance := serverAddressElem[5]
			if strings.ReplaceAll(strings.ReplaceAll(dxcServerInstance, "-", ""), "_", "") == config.GetConfigs().Service.InstanceID {
				return true
			}
		}
	}
	return false
}
