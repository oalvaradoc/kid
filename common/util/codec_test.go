package util

import (
	"fmt"
	event "git.multiverse.io/eventkit/kit/common/model/transaction"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
	"testing"
	"time"
)

type Req struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestEncodeAndDecode(t *testing.T) {
	// encode
	request := event.RootTxnBeginRequest{
		Head: event.TxnEventHeader{
			Service: "RegisterRootTransaction",
		},
		Request: event.RootTxnBeginRequestBody{
			ParticipantAddress: uuid.NewV4().String(),
			RequestTime:        time.Now().String(),
			ParentXid:          uuid.NewV4().String(),
			RootXid:            uuid.NewV4().String(),
			BranchXid:          uuid.NewV4().String(),
			ServiceName:        "ServiceName",
		},
	}
	paramData := []byte("hello world")
	requestBody, err := Encode(request, paramData)
	if err != nil {
		panic(err)
	}
	requestA := event.RootTxnBeginRequest{}
	// decode
	paramBytes, err := Decode(requestBody, &requestA)
	if err != nil {
		panic(err)
	}
	fmt.Println(request.Request.RootXid == requestA.Request.RootXid)
	fmt.Println(string(paramData) == string(paramBytes))
}

func TestParamsSerialAndDeSerial(t *testing.T) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	req := &Req{
		Name: "KanYuXia",
		Age:  24,
	}
	param2 := "hello world"
	bytes, err := SerialParams(req, param2)
	if err != nil {
		panic(err)
	}
	params, err := DeSerialParams(bytes)
	if err != nil {
		panic(err)
	}
	r := new(Req)
	if err = json.Unmarshal([]byte(params[0]), r); err != nil {
		panic(err)
	}
	b := "1"
	if err = json.Unmarshal([]byte(params[1]), &b); err != nil {
		panic(err)
	}
	fmt.Println(req.Name == r.Name)
	fmt.Println(b == param2)
}
