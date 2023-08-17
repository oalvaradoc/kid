package transaction

// GlobalTxnResultReport is a request model for DXC server report the final result of transaction.
type GlobalTxnResultReport struct {
	Head    TxnEventHeader            `json:"head"`
	Request GlobalTxnResultReportBody `json:"request"`
}

// GlobalTxnResultReportBody is a model that contains status of root transaction, and the list of all the branch transaction.
type GlobalTxnResultReportBody struct {
	//Org       string `json:"org"`
	//Az        string `json:"az"`
	//Su       string `json:"su"`
	//ServiceID string `json:"serviceId"`
	//nodeId    string `json:"nodeId"`
	RootXid string `json:"rootXid"`
	// 1-All confirm successfully
	// 2-All cancel successfully
	GTxnStat    string    `json:"gStat"`
	GTxnIsDone  bool      `json:"gIsDone"`
	GTxnStartTm string    `json:"gStartTm"`
	GTxnEndTm   string    `json:"gEndTm"`
	TxnList     []TxnInfo `json:"txnList"`
}

// TxnInfo is a model that contains the information of transaction.
type TxnInfo struct {
	//Org        string `json:"org"`
	//Az         string `json:"az"`
	//Su        string `json:"su"`
	//ServiceID  string `json:"serviceId"`
	//nodeId     string `json:"nodeId"`
	//InstanceID string `json:"instanceId"`
	ParentXid    string `json:"pXid"`
	BranchXid    string `json:"bXid"`
	TxnIsDone    bool   `json:"isDone"`
	ErrorMessage string `json:"errorMessage"`
	TxnStartTm   string `json:"startTm"`
	TxnEndTm     string `json:"endTm"`
}
