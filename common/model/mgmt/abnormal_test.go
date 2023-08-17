package mgmt

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestDate_MarshalJSON(t *testing.T) {
	transaction := Transaction{
		RootXID:      "rootXid",
		BranchXID:    "branchXid",
		ParentXID:    "parentXid",
		Stat:         "0",
		ORG:          "org",
		AZ:           "az",
		SU:           "su",
		NodeID:       "nodeId",
		ServiceID:    "serviceId",
		InstanceID:   "instanceId",
		Params:       "params",
		Headers:      "headers",
		Reason:       "reason",
		Environment:  "env",
		Workspace:    "wks",
		TxnStartTime: Date(time.Now()),
		TxnEndTime:   Date(time.Now()),
	}
	transactions := make([]Transaction, 0)
	transactions = append(transactions, transaction)
	request := ReportAbnormalTxnRequest{
		RootXID:            "rootXid",
		Stat:               "0",
		Reason:             "reason",
		GlobalTxnStartTime: Date(time.Now()),
		GlobalTxnEndTime:   Date(time.Now()),
		Transactions:       transactions,
	}
	bytes, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
	var newReq ReportAbnormalTxnRequest
	err = json.Unmarshal(bytes, &newReq)
	if err != nil {
		panic(err)
	}
	var result time.Time
	result = time.Time(newReq.GlobalTxnEndTime)
	fmt.Println(result)
}

func TestDate_String_1(t *testing.T) {
	request := PageQueryAbnormalTxnRequest{
		Page:        1,
		PageNum:     2,
		Stats:       []string{"1", "2"},
		ORG:         "",
		AZ:          "",
		SU:          "",
		NodeID:      "",
		Environment: "",
		Workspace:   "",
		ServiceID:   "",
		InstanceID:  "",
		Begin:       Date(time.Now()),
		End:         Date(time.Now()),
	}
	bytes, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
func TestDate_String(t *testing.T) {
	jsonStr := `
{"rootXid":"rootXid","stat":"0","reason":"reason","globalTxnStartTime":"2021-12-07 17:25:05.424217","globalTxnEndTime":"2021-12-07 17:25:05.424218","transactions":[{"parentXid":"parentXid","rootXid":"rootXid","branchXid":"branchXid","stat":"0","org":"org","az":"az","environment":"env","workspace":"wks","su":"su","nodeId":"nodeId","serviceId":"serviceId","instanceId":"instanceId","reason":"reason","params":"params","headers":"headers","txnStartTime":"2021-12-07 17:25:05.423781","txnEndTime":"2021-12-07 17:25:05.424216"}]}
`
	var request ReportAbnormalTxnRequest
	err := json.Unmarshal([]byte(jsonStr), &request)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%++v", request)
}
