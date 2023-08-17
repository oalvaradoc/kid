package auto

import (
	"fmt"
	"git.multiverse.io/eventkit/kit/codec"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/util"
	"git.multiverse.io/eventkit/kit/constant"
	"reflect"
	"strconv"
)

// Encoder is an implement of codec.Encoder,Automatically identify the target parameter type,
// if it is a basic type, it will be serialized according to the basic type, if it is not a basic type,
// it will be serialized according to the group coder set by ComplexObjectEncoder.
// sets ComplexObjectEncoder as the default codec.Encoder
type Encoder struct {
	ComplexObjectEncoder codec.Encoder
}

// Decoder is an implement of codec.Decoder, Automatically identify the target parameter type, if it is a basic type,
// it will be deserialized according to the basic type, if it is not a basic type, it will be deserialized according to
// the decoder set by ComplexObjectDecoder.
// sets ComplexObjectDecoder as the default codec.Decoder
type Decoder struct {
	ComplexObjectDecoder codec.Decoder
}

// Decode decodes the byte array into target object
func (j *Decoder) Decode(data []byte, v interface{}) error {
	if nil == data || len(data) == 0 || util.IsNil(v) {
		return nil
	}

	switch v.(type) {
	case *string:
		res, err := strconv.Unquote(string(data))
		if nil != err {
			return err
		}
		*(v.(*string)) = res
	case *bool:
		x, err := parseBool(string(data))
		if nil != err {
			return err
		}
		*(v.(*bool)) = x
	case *int8:
		x, err := parseInt8(string(data))
		if nil != err {
			return err
		}
		*(v.(*int8)) = x
	case *byte:
		x, err := parseUint8(string(data))
		if nil != err {
			return err
		}
		*(v.(*byte)) = x
	case *int16:
		x, err := parseInt16(string(data))
		if nil != err {
			return err
		}
		*(v.(*int16)) = x
	case *int32:
		x, err := parseInt32(string(data))
		if nil != err {
			return err
		}
		*(v.(*int32)) = x
	case *uint32:
		x, err := parseUint32(string(data))
		if nil != err {
			return err
		}
		*(v.(*uint32)) = x
	case *int64:
		x, err := parseInt64(string(data))
		if nil != err {
			return err
		}
		*(v.(*int64)) = x
	case *uint64:
		x, err := parseUint64(string(data))
		if nil != err {
			return err
		}
		*(v.(*uint64)) = x
	case *float32:
		x, err := parseFloat32(string(data))
		if nil != err {
			return err
		}
		*(v.(*float32)) = x
	case *float64:
		x, err := parseFloat64(string(data))
		if nil != err {
			return err
		}
		*(v.(*float64)) = x
	default:
		vi := reflect.ValueOf(v)
		if vi.Kind() != reflect.Ptr {
			return errors.New(constant.SystemInternalError, "The parameter of receiver(`v`) must be a pointer!")
		}

		if nil == j.ComplexObjectDecoder {
			return errors.Errorf(constant.SystemInternalError, "failed to decode[%++v], invalid type", v)
		}
		return j.ComplexObjectDecoder.Decode(data, v)
	}

	return nil
}

var doubleQuotationMark = "\""

// Encode encodes the object into byte array
func (j *Encoder) Encode(v interface{}) ([]byte, error) {
	if util.IsNil(v) {
		return make([]byte, 0), nil
	}
	switch v.(type) {
	case *string:
		return []byte(strconv.Quote(*v.(*string))), nil
	case string:
		return []byte(strconv.Quote(v.(string))), nil
	case *bool:
		return []byte(formatBool(*v.(*bool))), nil
	case bool:
		return []byte(formatBool(v.(bool))), nil
	case *int8:
		return []byte(formatInt8(*v.(*int8))), nil
	case int8:
		return []byte(formatInt8(v.(int8))), nil
	case *byte:
		return []byte(formatUnit8(*v.(*uint8))), nil
	case byte:
		return []byte(formatUnit8(v.(uint8))), nil
	case *int16:
		return []byte(formatInt16(*v.(*int16))), nil
	case int16:
		return []byte(formatInt16(v.(int16))), nil
	case *uint16:
		return []byte(formatUint16(*v.(*uint16))), nil
	case uint16:
		return []byte(formatUint16(v.(uint16))), nil
	case *int32:
		return []byte(formatInt32(*v.(*int32))), nil
	case int32:
		return []byte(formatInt32(v.(int32))), nil
	case *uint32:
		return []byte(formatUint32(*v.(*uint32))), nil
	case uint32:
		return []byte(formatUint32(v.(uint32))), nil
	case *int64:
		return []byte(formatInt64(*v.(*int64))), nil
	case int64:
		return []byte(formatInt64(v.(int64))), nil
	case *uint64:
		return []byte(formatUint64(*v.(*uint64))), nil
	case uint64:
		return []byte(formatUint64(v.(uint64))), nil
	case *float32:
		return []byte(formatFloat32(*v.(*float32))), nil
	case float32:
		return []byte(formatFloat32(v.(float32))), nil
	case *float64:
		return []byte(formatFloat64(*v.(*float64))), nil
	case float64:
		return []byte(formatFloat64(v.(float64))), nil
	default:
		if nil == j.ComplexObjectEncoder {
			return nil, errors.Errorf(constant.SystemInternalError, "failed to encode[%++v], unsupport type", v)
		}

		return j.ComplexObjectEncoder.Encode(v)
	}
}

func formatFloat32(i float32) string {
	s := fmt.Sprintf("%f", i)
	return s
}

func formatFloat64(i float64) string {
	s := fmt.Sprintf("%f", i)
	return s
}

func formatInt64(i int64) string {
	return strconv.FormatInt(i, 10)
}

func formatUint64(i uint64) string {
	return strconv.FormatInt(int64(i), 10)
}

func formatInt32(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func formatUint32(i uint32) string {
	return strconv.FormatInt(int64(i), 10)
}

func formatInt16(i int16) string {
	return strconv.FormatInt(int64(i), 10)
}

func formatUint16(i uint16) string {
	return strconv.FormatInt(int64(i), 10)
}

func formatInt8(i int8) string {
	return strconv.FormatInt(int64(i), 10)
}

func formatUnit8(i uint8) string {
	return strconv.FormatInt(int64(i), 10)
}

func parseFloat32(str string) (float32, error) {
	s, err := strconv.ParseFloat(str, 32)
	return float32(s), err
}

func parseFloat64(str string) (float64, error) {
	s, err := strconv.ParseFloat(str, 64)
	return s, err
}

func parseInt64(str string) (int64, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	return i, err
}

func parseUint64(str string) (uint64, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	return uint64(i), err
}

func parseInt32(str string) (int32, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	return int32(i), err
}

func parseUint32(str string) (uint32, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	return uint32(i), err
}

func parseInt16(str string) (int16, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	return int16(i), err
}

func parseUint16(str string) (uint16, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	return uint16(i), err
}

func parseInt8(str string) (int8, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	return int8(i), err
}

func parseUint8(str string) (uint8, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	return uint8(i), err
}

func parseBool(str string) (bool, error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False":
		return false, nil
	}
	return false, errors.Errorf(constant.SystemInternalError, "failed to parseBool:%s", str)
}

func formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// BuildAutoCodecWithJSONCodec creates a new codec.Codec and sets JSON codec(create by codec.BuildJSONCode()) as default codec
func BuildAutoCodecWithJSONCodec() codec.Codec {
	jsonCodec := codec.BuildJSONCodec()
	return &impl{
		encoder: func() codec.Encoder {
			return &Encoder{
				ComplexObjectEncoder: jsonCodec.Encoder(),
			}
		},
		decoder: func() codec.Decoder {
			return &Decoder{
				ComplexObjectDecoder: jsonCodec.Decoder(),
			}
		},
	}
}

// BuildAutoCodec create a new codec.Codec and sets customize codec.Codec as default codec
func BuildAutoCodec(complexObjectCodec codec.Codec) codec.Codec {
	return &impl{
		encoder: func() codec.Encoder {
			return &Encoder{
				ComplexObjectEncoder: complexObjectCodec.Encoder(),
			}
		},
		decoder: func() codec.Decoder {
			return &Decoder{
				ComplexObjectDecoder: complexObjectCodec.Decoder(),
			}
		},
	}
}

type impl struct {
	encoder func() codec.Encoder
	decoder func() codec.Decoder
}

func (c *impl) Encoder() codec.Encoder { return c.encoder() }
func (c *impl) Decoder() codec.Decoder { return c.decoder() }
