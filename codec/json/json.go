package json

import (
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	jsoniter "github.com/json-iterator/go"
)

// Decoder is an implement of codec.Decoder used to decode JSON byte array into a object
type Decoder struct{}

// Encoder is an implement of codec.Encoder used to encode a object into JSON array
type Encoder struct{}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Decode decodes byte array into a object
func (j *Decoder) Decode(data []byte, v interface{}) error {
	if nil == data || len(data) == 0 || util.IsNil(v) {
		return nil
	}
	err := json.Unmarshal(data, v)
	if nil != err {
		return errors.Wrap(constant.SystemInternalError, err, 0)
	}
	return nil
}

// Encode encodes a object into byte array
func (j *Encoder) Encode(v interface{}) ([]byte, error) {
	if util.IsNil(v) {
		return make([]byte, 0), nil
	}
	bs, err := json.Marshal(v)
	if nil != err {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}
	return bs, nil
}
