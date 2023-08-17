package constant


// Define all topic-related constants and platform-related parameter keys
const (
	TopicType = "type"

	TopicSourceORG         = "srcOrg"
	TopicSourceWorkspace   = "srcWks"
	TopicSourceEnvironment = "srcEnv"
	TopicSourceAZ          = "srcAz"
	TopicSourceSU          = "srcSu"
	TopicSourceDCN         = "srcDcn"
	TopicSourceServiceID   = "srcSvrId"
	TopicSourceNodeID      = "srcNodeId"
	TopicSourceInstanceID  = "srcInsId"

	TopicDestinationORG         = "dstOrg"
	TopicDestinationWorkspace   = "dstWks"
	TopicDestinationEnvironment = "dstEnv"
	TopicDestinationSU          = "dstSu"
	TopicDestinationDCN         = "dstDcn"
	TopicDestinationVersion     = "dstVersion"
	TopicDestinationNodeID      = "dstNodeId"
	TopicDestinationInstanceID  = "dstInsId"

	TopicORG         = "org"
	TopicWorkspace   = "wks"
	TopicEnvironment = "env"
	TopicAZ          = "az"
	TopicSU          = "su"
	TopicServiceID   = "serviceId"
	TopicVersion     = "version"
	TopicNodeID      = "nodeId"
	TopicInstanceID  = "insId"
	TopicWildcard    = "wildcard"

	TopicID            = "id"
	TopicShareName     = "shareName"
	TopicShareFilterID = "shareFilterID"

	Enable            = "1"
	RrReplyTo         = "_REPLY_TO"
	TxnIsLocalCall    = "isLocalCall"
	TxnIsSemiSyncCall = "_IS_SEMI_SYNC_CALL"
	IsPayloadCrypto   = "_IS_PAYLOAD_CRYPTO"
	CryptoAlgo        = "_CRYPTO_ALGO"
	CryptoPadding     = "_CRYPTO_PADDING"
	CryptoMode        = "_CRYPTO_MODE"
	CryptoKeyVersion  = "_CRYPTO_KEY_VERSION"
	FixedKeyID        = "_FIXED_KEY_ID"
	DeliveryMode      = "_DELIVERY_MODE"
	DmqEligible       = "_DMQ_ELIGIBLE"
	DstTopicID        = "_DST_TOPIC_ID"
	DiscardResponse   = "_DISCARD_RESPONSE"
	CsStartTimestamp  = "_CS_START_TIMESTAMP"
	SrStartTimestamp  = "_SR_START_TIMESTAMP"
	IsNeedLookup      = "_is_need_lookup"

	TargetSU        = "_TARGET_SU"
	GlsElementType  = "_GLS_ELEMENT_TYPE"
	GlsElementClass = "_GLS_ELEMENT_CLASS"
	GlsElementID    = "_GLS_ELEMENT_ID"

	TargetSUOld        = "_TARGET_DCN"
	GlsElementTypeOld  = "_DLS_ELEMENT_TYPE"
	GlsElementClassOld = "_DLS_ELEMENT_CLASS"
	GlsElementIDOld    = "_DLS_ELEMENT_ID"

	FastFailedErrorCode    = "_FAST_FAILED_ERROR_CODE"
	FastFailedErrorMessage = "_FAST_FAILED_ERROR_MESSAGE"
	To1                    = "_TO1"
	To2                    = "_TO2"
	To3                    = "_TO3"

	//topic_type
	TopicTypeHeartbeat    = "HBT"
	TopicTypeError        = "ERR"
	TopicTypeAlert        = "ALT"
	TopicTypeBusiness     = "TRN"
	TopicTypeLog          = "LOG"
	TopicTypeMetrics      = "MTR"
	TopicTypeDXC          = "DXC"
	TopicTypeDTS          = "DTS"
	TopicTypeOPS          = "OPS"
	TopicTypeP2P          = "P2P"
	TopicTypeSessionInBox = "#P2P"
)

//Define common error message structures
const (
	ReturnErrorCode = "errorCode"
	ReturnErrorMsg  = "errorMsg"

	ReturnErrorCodeOld = "RetMsgCode"
	ReturnErrorMsgOld  = "RetMessage"
	ReturnStatus       = "RetStatus"
)

//Define common error code
const (
	Success                    = "0"
	SystemInternalError        = "SY99999999"
	SystemRemoteCallTimeout    = "SY99999998"
	SystemErrConnectionClosed  = "SY99999997"
	SystemErrConnectionAborted = "SY99999996"
	SystemErrConnectionRefused = "SY99999995"
	SystemErrConnectionReset   = "SY99999994"

	SystemMeshRequestReplyTimeout = "SY99999993"

	SystemCallbackAppTimeout           = "SY99999992"
	SystemCallbackAppConnectionClosed  = "SY99999991"
	SystemCallbackAppConnectionAborted = "SY99999990"
	SystemCallbackAppConnectionRefused = "SY99999989"
	SystemCallbackAppConnectionReset   = "SY99999988"

	UpstreamServiceMessageDecodeError   = "SY99999987"
	UpstreamServiceMessageEncodeError   = "SY99999986"
	DownstreamServiceMessageDecodeError = "SY99999985"
	DownstreamServiceMessageEncodeError = "SY99999984"
	ValidationError                     = "SY99999983"

	NewTransactionProxyError = "SY99999982"

	TransactionBeginError                = "SY99999981"
	TransactionJoinError                 = "SY99999980"
	TransactionEndCallbackConfirmError   = "SY99999979"
	TransactionEndCallbackCancelError    = "SY99999978"
	TransactionEndCallbackConfirmTimeout = "SY99999977"
	TransactionEndCallbackCancelTimeout  = "SY99999976"
	TransactionEndOtherError             = "SY99999975"

	CannotFoundHandlerWithURLError     = "SY99999974"
	InvalidEventTypeError              = "SY99999973"
	CannotFoundHandlerWithEventIDError = "SY99999972"
)

// Define trace id related keys, contains old version key
const (
	KeyTraceID            = "traceId"
	KeySpanID             = "spanId"
	KeyParentSpanID       = "parentSpanId"
	KeyParentSpanIDForLog = "parentId"
	KeyServiceID          = "serviceID"
	KeyTopicID            = "topicID"
	KeyST = "st"					// start time

	KeyTraceIDOld      = "GlobalBizSeqNo"
	KeySpanIDOld       = "SrcBizSeqNo"
	KeyParentSpanIDOld = "ParentBizSeqNo"
	KeyDownstreamServiceLoggingStartTime = "DownstreamServiceLoggingStartTime"
)

// Define date/time format related keys
const (
	TimeStamp       = "2006-01-02 15:04:05"
	DateDashFormat  = "2006-01-02"
	Dateformat      = "20060102"
	TimeColonFormat = "15:04:05"
)

// Define response template related keys
const (
	DefaultResponseTemplate         = `{"errorCode":"{{.errorCode}}","errorMsg":"{{.errorMsg}}","response":{{.data}}}`
	ErrorCodeMappingKey             = "errorCodeKey"
	ErrorMsgMappingKey              = "errorMsgKey"
	ResponseDataBodyMappingKey      = "responseDataKey"
	ResponseAutoParseTypeMappingKey = "type"
	ResponseAutoParseTypeJSON       = "JSON"
	ResponseAutoParseTypeXML        = "XML"
	ResponseAutoDefaultParseType    = ResponseAutoParseTypeJSON

	DefaultErrorCodeKey = "errorCode"
	DefaultErrorMsgKey  = "errorMsg"
	DefaultDataBodyKey  = "response"

	ResponseAutoParseKeyMappingKey         = "ResponseAutoParseKeyMapping"
	SkipResponseAutoParseKeyMappingFlagKey = "SkipResponseAutoParseKeyMappingFlagKey"
)

// Define user lang related keys
const (
	UserLang = "UserLang" // ISO International Language Code

	LangEnUS = "en-US" //Define language
	LangZhCN = "zh-CN"
	LangZhTW = "zh-TW"
	LangThTH = "th-TH"
)

// Define http content related keys
const (
	HTTPContentTypeKey     = "Content-Type"
	DefaultContentTypeJSON = "application/json"
	DefaultHTTPMethodPost  = "POST"
	RequestTypeHTTP        = "http"
)

// define trace id related keys
const (
	RootXIDKey              = "TxnRootXId"
	ParentXIDKey            = "TxnParentXId"
	BranchXIDKey            = "TxnBranchXId"
	TransactionAgentAddress = "TxnAddress"
	CurrentSU               = "_currentSu"

	RootXIDKeyOld              = "ROOT_XID"
	ParentXIDKeyOld            = "PARENT_XID"
	BranchXIDKeyOld            = "BRANCH_XID"
	TransactionAgentAddressOld = "DTS_AGENT_ADDRESS"
)

// Used to return the function name of the handler name, the name of the handler can be set after the method is implemented in each handler
const (
	FunctionForGetHandlerName = "HandlerName"
)

const (
	// DefaultTimeoutMilliseconds default timeout will be set if request downstream if the timeout is not specified
	DefaultTimeoutMilliseconds = 30 * 1000

	// Default redis pool size
	DefaultRedisPoolSize = 10
)

// Define the key type when generating the serial number
const (
	TraceIDType = "0"
	SpanIDType  = "1"
)

//pre handle name
const (
	PreHandleMethodName  = "PreHandle"
	ValidationMethodName = "Validation"
)

// define transaction related error codes
const (
	//Success = 0

	// for transaction start
	ErrorCodeOffset      = 230000
	TransactionErrorCode = ErrorCodeOffset

	// 230001:uncaught exception error
	UnCaughtExceptionError = ErrorCodeOffset + 1

	// 230002:internal error
	InternalError = ErrorCodeOffset + 2

	// 230100:common error
	CommonError = ErrorCodeOffset + 100 //

	// 230101:invalid parameter error
	InvalidParameter = CommonError + 1

	// 230102:Transaction begin failed, root transaction already exists
	TxnBeginRootXidAlreadyExists = CommonError + 2

	// 230103:Transaction join failed, root transaction does not exist
	TxnJoinFailedCannotFindRootXid = CommonError + 3

	// 230104:Transaction registration failed, branch transaction already exists
	TxnJoinFailedBranchXidAlreadyExists = CommonError + 4

	// 230105:Transaction registration failed, transaction status error
	TxnJoinFailedGlobalTxnStateError = CommonError + 5

	// 230106:Global transaction result report failed, root transaction does not exist
	DoEndFailedCannotFindRootXid = CommonError + 6

	// 230107:Failed to report global transaction results, the transaction status is wrong
	DoEndFailedGlobalTxnStateError = CommonError + 7

	// 230108:Global transaction do end failed, branch transaction confirm failed
	TxnEndFailedBranchConfirmFailed = CommonError + 8

	// 230109:Global transaction do end failed, branch transaction cancel failed
	TxnEndFailedBranchCancelFailed = CommonError + 9 // :

	// 230110:Branch transaction confirm failed, branch transaction does not exist
	BranchTxnConfirmFailedCannotFindBranchXid = CommonError + 10

	// 230111:Branch transaction confirm failed, status of branch transaction is invalid
	BranchTxnConfirmFailedBranchTxnStateError = CommonError + 11

	// 230112:Branch transaction cancel failed, branch transaction does not exist
	BranchTxnCancelFailedCannotFindBranchXid = CommonError + 12

	// 230113:Branch transaction cancel failed, status of branch transaction if invalid
	BranchTxnCancelFailedBranchTxnStateError = CommonError + 13

	// 230114:Global transaction do end failed, transactions not all callback success
	TxnEndFailedBranchesNotAllCallbackSuccess = CommonError + 14

	// 230115:Global transaction do end failed, invoke server timeout
	TxnEndFailedTimeOut = CommonError + 15

	// 230116:Do metrics record failed
	DoMetricsRecordFailed = CommonError + 16

	// 230117:Do abnormal transaction record
	DoAbnormalTransactionFailed = CommonError + 17

	// 230118:Global transaction query failed
	GlobalTransactionQueryFailed = CommonError + 18

	// 230119:Sync Abnormal transaction record stat failed
	SyncAbnormalTransactionRecordStatFailed = CommonError + 19

	// 230120:Global transaction page query
	GlobalTransactionPageQueryFailed = CommonError + 20
)

// define communicate type related keys
const (
	CommDirect                  = "direct"
	CommMesh                    = "mesh"
	ParticipantAddressSplitChar = "|"
	TopicIDSplitChat            = "/"
	AttributesPrefix            = "_attr."
)

// define client SDK version related keys
const (
	ClientSDKVersion   = "ClientSDKVersion"
	ClientSDKVersionV1 = "v1"
	ClientSDKVersionV2 = "v2"
)

// defined context related keys
const (
	ContextTransactionKey          = "ContextTransactionKey"
	HandlerContextsKey             = "HANDLER_CONTEXTS"
)

const (
	PassThroughHeaderKeyListKey = "_pthklist"
)

// Define all additional config key
const (
	ExtConfigTransaction                       = "transaction"
	ExtConfigDownstreamService                 = "downstreamService"
	ExtEventKeyMap                             = "eventKeyMap"
	ExtConfigService                           = "service"
	ExtConfigCustomClient                      = "customClient"
	ExtConfigCustomDefaultInterceptors         = "customDefaultInterceptors"
	ExtConfigCustomSedClientOptions            = "customSedClientOptions"
	ExtConfigCustomCallbackHandleWrapper       = "customCallbackHandleWrapper"
	ExtConfigCustomResponseTemplate            = "customResponseTemplate"
	ExtConfigCustomCallWrappers                = "customCallWrappers"
	ExtConfigCustomCallInterceptors            = "customCallInterceptors"
	ExtConfigCustomResponseAutoParseKeyMapping = "responseAutoParseKeyMapping"
	ExtConfigDefaultEnableValidation           = "defaultEnableValidation"
	ExtConfigDefaultUserLang                   = "defaultUserLang"
	ExtConfigCustomUserLangKey                 = "customUserLangKey"
	ExtEnableExecutorLogging                   = "enableExecutorLogging"
)

// Define default URL path of transaction begin、join、end
const (
	TxnBeginURLPath           = "/v1/txn_mgt/txn_begin"
	TxnJoinURLPath            = "/v1/txn_mgt/txn_join"
	TxnEndURLPath             = "/v1/txn_mgt/txn_end"
	TxnMacroServiceAddressURL = "http://127.0.0.1:9999"
)

// Define confirm or cancel flag
const (
	ConfirmFlag = 1 << iota
	CancelFlag
)

// Define all the http method
const (
	HTTPMethodGet = 1 << iota
	HTTPMethodHead
	HTTPMethodOptions
	HTTPMethodPost
	HTTPMethodPut
	HTTPMethodPatch
	HTTPMethodDelete
)

// Define all the interceptors name
const (
	InterceptorApm            = "APM"
	InterceptorAudit          = "AUDIT"
	InterceptorClientReceive  = "CLIENT-RECEIVE"
	InterceptorLogging        = "LOGGING"
	InterceptorServerResponse = "SERVER-RESPONSE"
	InterceptorTrace          = "TRACE"
	InterceptorTransaction    = "TRANSACTION"
)

// Define all the wrappers name
const (
	WrapperApm        = "APM"
	WrapperAddressing = "ADDRESSING"
	WrapperLogging    = "LOGGING"
	WrapperTrace      = "TRACE"
)

const (
	ArrPreffix = '['
	ArrSuffix  = ']'
	ObjPreffix = '{'
	ObjSuffix  = '}'
	Seperator  = ','

	DefaultOverlayString = "*"

	DefaultCircuitBreakerMonitorPort = 8888
)

const MarkAsErrorResponseKey = "MarkAsErrorResponse"

const (
	RootErrorCode                = "rootErrorCode"
	RootErrorMsg                 = "rootErrorMsg"
	RootKVToSecondStageKeyPrefix = "_dxc_."
)


var (
	AddressingSuccessful        = "0"
	InvalidParameters           = "GLSE661400"
	InvalidRequest              = "GLSE661401"
	CacheDoesNotHaveInitialize  = "GLSE661402"
	RecordsNotFound             = "GLSE661403"
	GLSSystemInternalError         = "GLSE661499"
)

const (
	// EventHandlerVersion is the current SDK version number
	EventHandlerVersion = "v1.0.0"

	// ProtocolLevel is protocol level that the current SDK supported
	ProtocolLevel = 2
)
