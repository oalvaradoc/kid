package util

import (
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Encode is a tool function that used for encode the object and parameter data to byte array
func Encode(jsonObj interface{}, paramData []byte) ([]byte, error) {
	result := make([]byte, 0)
	jsonBytes, err := json.Marshal(jsonObj)
	if err != nil {
		return nil, err
	}
	by, err := IntToBytes(len(jsonBytes))
	if err != nil {
		return nil, err
	}
	result = append(result, by...)
	result = append(result, jsonBytes...)
	result = append(result, paramData...)
	return result, nil
}

// Decode is a tool function that used for decode the value to the target object pointer and returns the parameter values.
func Decode(body []byte, jsonObjPtr interface{}) ([]byte, error) {
	length, err := BytesToInt(body[:4])
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body[4:length+4], jsonObjPtr); err != nil {
		return nil, err
	}
	return body[length+4:], nil
}

// SerialParams is used to serialize the parameters to the byte array
func SerialParams(params ...interface{}) ([]byte, error) {
	result := make([]byte, 0)
	for _, param := range params {
		by, err := json.Marshal(param)
		if err != nil {
			return nil, err
		}
		lBy, err := IntToBytes(len(by))
		if err != nil {
			return nil, err
		}
		result = append(result, lBy...)
		result = append(result, by...)
	}
	return result, nil
}

// DeSerialParams is used to deserialize the byte array to the parameters
func DeSerialParams(paramData []byte) (result []string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf(constant.SystemInternalError, "deserialization Panic err:%v", e)
		}
	}()
	result = make([]string, 0)
	end := 0
	for {
		length, err := BytesToInt(paramData[end : end+4])
		if err != nil {
			return nil, err
		}
		by := paramData[end+4 : end+4+length]
		result = append(result, string(by))
		end = end + 4 + length
		if end == len(paramData) {
			break
		}
	}
	return result, nil
}
