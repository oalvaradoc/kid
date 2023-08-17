package mgmt

import "time"

const TimeFormat = "2006-01-02 15:04:05"

type Date time.Time

func (t *Date) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = Date(now)
	return
}

func (t Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Date) String() string {
	return time.Time(t).Format(TimeFormat)
}

// ReportAbnormalTxnRequest is a request for DXC server report abnormal transaction information when failed.
type ReportAbnormalTxnRequest struct {
	RootXID            string        `json:"rootXid"`
	Stat               string        `json:"stat"` // 0:trying, 1-confirming, 2-confirmed, 3-cancelling, 4-cancelled
	Reason             string        `json:"reason"`
	GlobalTxnStartTime Date          `json:"globalTxnStartTime"`
	GlobalTxnEndTime   Date          `json:"globalTxnEndTime"`
	Transactions       []Transaction `json:"transactions"`
}

// Transaction is a model that contains all transaction information for DXC
type Transaction struct {
	RootXID      string `json:"rootXid"`
	BranchXID    string `json:"branchXid"`
	ParentXID    string `json:"parentXid"`
	Stat         string `json:"stat"`
	ORG          string `json:"org"`
	AZ           string `json:"az"`
	SU           string `json:"su"`
	NodeID       string `json:"nodeId"`
	ServiceID    string `json:"serviceId"`
	InstanceID   string `json:"instanceId"`
	Params       string `json:"params"`
	Headers      string `json:"headers"`
	Reason       string `json:"reason"`
	Environment  string `json:"environment"`
	Workspace    string `json:"workspace"`
	TxnStartTime Date   `json:"txnStartTime"`
	TxnEndTime   Date   `json:"txnEndTime"`
}

// PageQueryAbnormalTxnRequest is a request model for DXC management paging query abnormal transaction records
type PageQueryAbnormalTxnRequest struct {
	Page        int      `json:"page"`
	PageNum     int      `json:"pageNum"`
	Stats       []string `json:"stats"`
	ORG         string   `json:"org"`
	AZ          string   `json:"az"`
	SU          string   `json:"su"`
	NodeID      string   `json:"nodeId"`
	Environment string   `json:"environment"`
	Workspace   string   `json:"workspace"`
	ServiceID   string   `json:"serviceId"`
	InstanceID  string   `json:"instanceId"`
	Begin       Date     `json:"begin"`
	End         Date     `json:"end"`
}

// PageQueryAbnormalTxnResponse  is a response model for DXC management paging query abnormal transaction records
type PageQueryAbnormalTxnResponse struct {
	TotalNum     int64         `json:"totalNum"`
	AbnormalTxns []AbnormalTxn `json:"abnormalTxns"`
}

// AbnormalTxn is a model that contains all abnormal transaction information
type AbnormalTxn struct {
	ID                 int64  `json:"id"`
	RootXID            string `json:"rootXid"`
	Stat               string `json:"stat"`
	ORG                string `json:"org"`
	AZ                 string `json:"az"`
	SU                 string `json:"su"`
	NodeID             string `json:"nodeId"`
	Environment        string `json:"environment"`
	Workspace          string `json:"workspace"`
	ServiceID          string `json:"serviceId"`
	InstanceID         string `json:"instanceId"`
	Reason             string `json:"reason"`
	GlobalTxnStartTime Date   `json:"globalTxnStartTime"`
	GlobalTxnEndTime   Date   `json:"globalTxnEndTime"`
}

// DetailQueryAbnormalTxnRequest is a request model for enquiring the detail information of specified abnormal transaction.
type DetailQueryAbnormalTxnRequest struct {
	RootXID string `json:"rootXid"`
}

// DetailQueryAbnormalTxnResponse is a response model for enquiring the detail information of specified abnormal transaction.
type DetailQueryAbnormalTxnResponse struct {
	AbnormalTxn AbnormalTxn `json:"abnormalTxn"`
	DetailTxns  []DetailTxn `json:"detailTxns"`
}

// DetailTxn is a model that contains all the detail information of abnormal transaction.
type DetailTxn struct {
	ID            int64  `json:"id"`
	AbnormalTxnID int64  `json:"abnormalTxnId"`
	RootXid       string `json:"rootXid"`
	BranchXid     string `json:"branchXid"`
	ParentXid     string `json:"parentXid"`
	ORG           string `json:"org"`
	AZ            string `json:"az"`
	SU            string `json:"su"`
	NodeID        string `json:"nodeId"`
	Environment   string `json:"environment"`
	Workspace     string `json:"workspace"`
	ServiceID     string `json:"serviceId"`
	InstanceID    string `json:"instanceId"`
	Params        string `json:"params"`
	Reason        string `json:"reason"`
	Stat          string `json:"stat"`
	Headers       string `json:"headers"`
	TxnStartTime  Date   `json:"txnStartTime"`
	TxnEndTime    Date   `json:"txnEndTime"`
}

// UpdateAbnormalTxnStatRequest is a request model for DXC server update transaction status
type UpdateAbnormalTxnStatRequest struct {
	RootXID string `json:"rootXid"`
	Stat    string `json:"stat"`
}

type QueryPanicTxnRequest struct {
	NodeIp string `json:"NodeIp"`
}

type QueryPanicTxnResponse struct {
	Keys []string `json:"keys"`
}

type DetailQueryPanicTxnRequest struct {
	Key string `json:"key"`
}

type DetailQueryPanicTxnResponse struct {
	AbnormalTxn AbnormalTxn `json:"abnormalTxn"`
	DetailTxns  []DetailTxn `json:"detailTxns"`
}

type DeletePanicTxnRequest struct {
	NodeIp string `json:"NodeIp"`
}

type DeletePanicTxnResponse struct {
	Success bool `json:"success"`
}
