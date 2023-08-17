package transaction

// AtomicTxnConfirmResponseBody is a response model for DXC server callback to confirm the local transaction
type AtomicTxnConfirmResponseBody struct {
	RootXid      string `json:"rootXid"`
	BranchXid    string `json:"branchXid"`
	ResponseTime string `json:"responseTime"`
}

// AtomicTxnConfirmResponse is a common response with `errorCode` and `errorMsg` and AtomicTxnConfirmResponseBody
type AtomicTxnConfirmResponse struct {
	ErrorCode int                          `json:"errorCode"`
	ErrorMsg  string                       `json:"errorMsg"`
	Data      AtomicTxnConfirmResponseBody `json:"data"`
}

// AtomicTxnCancelResponseBody is a response model for DXC server callback to cancel the local transaction
type AtomicTxnCancelResponseBody struct {
	RootXid      string `json:"rootXid"`
	BranchXid    string `json:"branchXid"`
	ResponseTime string `json:"responseTime"`
}

// AtomicTxnCancelResponse is a common response with `errorCode` and `errorMsg` and AtomicTxnCancelResponseBody
type AtomicTxnCancelResponse struct {
	ErrorCode int                         `json:"errorCode"`
	ErrorMsg  string                      `json:"errorMsg"`
	Data      AtomicTxnCancelResponseBody `json:"data"`
}

// AtomicTxnCallbackRequest is a root request model for DXC Server do the callback that contains TxnEventHeader and AtomicTxnCallbackRequestBody
type AtomicTxnCallbackRequest struct {
	Head    TxnEventHeader               `json:"head"`
	Request AtomicTxnCallbackRequestBody `json:"request"`
}

// AtomicTxnCallbackRequestBody is a request body for DXC Server do the callback.
type AtomicTxnCallbackRequestBody struct {
	RootXid     string `json:"rootXid"`
	ParentXid   string `json:"parentXid"`
	BranchXid   string `json:"branchXid"`
	RequestTime string `json:"requestTime"`
	ServiceName string `json:"serviceName"`
	Headers     string `json:"headers"`
}
