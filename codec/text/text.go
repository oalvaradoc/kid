package text

import (
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"github.com/modern-go/reflect2"
)

// Encoder is an implement of codec.Encoder used to encode non-structured data,
// such as converting strings into byte arrays or directly returning references to the original byte arrays
type Encoder struct{}

// Decoder is an implement of codec.Decoder used to decode non-structured data,
//// such as converting byte arrays into strings or directly returning references to the original byte arrays
type Decoder struct{}

// Decode decodes byte array into a string
// or directly references the original byte array back if the input parameter is a byte array
func (j *Decoder) Decode(data []byte, v interface{}) error {
	if nil == data || len(data) == 0 || util.IsNil(v) {
		return nil
	}

	switch v.(type) {
	case *string:
		*(v.(*string)) = string(data)
	case string:
		ptr := reflect2.PtrOf(v)
		*((*string)(ptr)) = string(data)
	case *[]byte:
		*(v.(*[]byte)) = data
	case []byte:
		ptr := reflect2.PtrOf(v)
		*((*[]byte)(ptr)) = data
	default:
		return errors.Errorf(constant.SystemInternalError, "failed to decode[%++v],unsupport type", v)
	}

	return nil
}

// Encode converts the string into a byte array,
// or directly references the original byte array back if the input parameter is a byte array
func (j *Encoder) Encode(v interface{}) ([]byte, error) {
	if util.IsNil(v) {
		return make([]byte, 0), nil
	}
	switch v.(type) {
	case *string:
		return []byte(*v.(*string)), nil
	case string:
		return []byte(v.(string)), nil
	case *[]byte:
		return *(v.(*[]byte)), nil
	case []byte:
		return v.([]byte), nil
	default:
		return nil, errors.Errorf(constant.SystemInternalError, "failed to encode[%++v], unsupport type", v)
	}

}
