package transaction

// RootTxnBeginRequestBody is a request model for root transaction begin the global transaction.
type RootTxnBeginRequestBody struct {
	ParticipantAddress string `json:"participantAddress"`
	RequestTime        string `json:"requestTime"`
	ParentXid          string `json:"parentXid"`
	RootXid            string `json:"rootXid"`
	BranchXid          string `json:"branchXid"`
	ServiceName        string `json:"serviceName"` // v2
	Headers            string `json:"headers"`     // v2
}

// RootTxnBeginRequest is a model that contains TxnEventHeader and RootTxnBeginRequestBody
type RootTxnBeginRequest struct {
	Head    TxnEventHeader          `json:"head"`
	Request RootTxnBeginRequestBody `json:"request"`
}

// RootTxnBeginResponseBody is a response model for root transaction begin the global transaction.
type RootTxnBeginResponseBody struct {
	RootXid      string `json:"rootXid"`
	ResponseTime string `json:"responseTime"`
}

// RootTxnBeginResponse is a common model that contains `error code` and `error message` and RootTxnBeginResponseBody
type RootTxnBeginResponse struct {
	ErrorCode int                      `json:"errorCode"`
	ErrorMsg  string                   `json:"errorMsg"`
	Data      RootTxnBeginResponseBody `json:"data"`
}

// BranchTxnJoinRequestBody is a request model for branch transaction join the global transaction.
type BranchTxnJoinRequestBody struct {
	ParticipantAddress string `json:"participantAddress"`
	RequestTime        string `json:"requestTime"`
	ParentXid          string `json:"parentXid"`
	RootXid            string `json:"rootXid"`
	BranchXid          string `json:"branchXid"`
	ServiceName        string `json:"serviceName"` // v2
	Headers            string `json:"headers"`     // v2
}

// BranchTxnJoinRequest is a model that contains TxnEventHeader and BranchTxnJoinRequestBody
type BranchTxnJoinRequest struct {
	Head    TxnEventHeader           `json:"head"`
	Request BranchTxnJoinRequestBody `json:"request"`
}

// BranchTxnJoinResponseBody is a response model for branch transaction join the global transaction.
type BranchTxnJoinResponseBody struct {
	BranchXid    string `json:"branchXid"`
	ResponseTime string `json:"responseTime"`
}

// BranchTxnJoinResponse is a common model that contains `error code` and `error message` and BranchTxnJoinResponseBody
type BranchTxnJoinResponse struct {
	ErrorCode int                       `json:"errorCode"`
	ErrorMsg  string                    `json:"errorMsg"`
	Data      BranchTxnJoinResponseBody `json:"data"`
}

// TxnEndRequestBody is a request model for root transaction end the global transaction.
type TxnEndRequestBody struct {
	RootXid                       string `json:"rootXid"`
	BranchXid                     string `json:"branchXid"`
	ParentXid                     string `json:"parentXid"`
	Ok                            bool   `json:"ok"`
	RequestTime                   string `json:"requestTime"`
	TryFailedIgnoreCallbackCancel bool   `json:"tryFailedIgnoreCallbackCancel"`
}

// TxnEndRequest is a model that contains TxnEventHeader and TxnEndRequestBody
type TxnEndRequest struct {
	Head    TxnEventHeader    `json:"head"`
	Request TxnEndRequestBody `json:"request"`
}

// TxnEndResponseBody is a response model for root transaction end the global transaction.
type TxnEndResponseBody struct {
	ResponseTime string `json:"responseTime"`
}

// TxnEndResponse is a common model that contains `error code` and `error message` and TxnEndResponseBody
type TxnEndResponse struct {
	ErrorCode int                `json:"errorCode"`
	ErrorMsg  string             `json:"errorMsg"`
	Data      TxnEndResponseBody `json:"data"`
}

// AbnormalTxnProcessingRequestBody is a request model for report the abnormal transaction.
type AbnormalTxnProcessingRequestBody struct {
	RootXid     string `json:"rootXid"`
	BranchXid   string `json:"branchXid"`
	Operation   string `json:"operation"`
	RequestTime string `json:"requestTime"`
}

// AbnormalTxnProcessingRequest is a model that contains TxnEventHeader and AbnormalTxnProcessingRequestBody
type AbnormalTxnProcessingRequest struct {
	Head    TxnEventHeader                   `json:"head"`
	Request AbnormalTxnProcessingRequestBody `json:"request"`
}

// AbnormalTxnProcessingResponseBody is a response model for report the abnormal transaction.
type AbnormalTxnProcessingResponseBody struct {
	ResponseTime string `json:"responseTime"`
}
