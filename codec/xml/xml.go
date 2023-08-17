package xml

import (
	"encoding/xml"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
)

// Encoder is an implement of codec.Encoder used to decode XML byte array into a object
type Encoder struct{}

// Decoder is an implement of codec.Decoder used to encode a object into XML array
type Decoder struct{}

// Decode decodes byte array into a object
func (j *Decoder) Decode(data []byte, v interface{}) error {
	if nil == data || len(data) == 0 || util.IsNil(v) {
		return nil
	}

	err := xml.Unmarshal(data, v)
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
	bs, err := xml.Marshal(v)
	if nil != err {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	return bs, nil
}
