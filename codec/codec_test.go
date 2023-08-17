package codec

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
)

type StructA struct {
	A string
	B string
	C int
}

var JSONBytes = []byte(`{"A":"testA1",
"B":"testB1",
"C": 1}`)

var XMLBytes = []byte(`<xml>
<A>testA2</A>
<B>testB2</B>
<C>2</C>
</xml>`)

func TestJSONCodec(t *testing.T) {
	codec := BuildJSONCodec()
	structA := &StructA{}
	codec.Decoder().Decode(nil, nil)
	jsonBytes, _ := codec.Encoder().Encode(nil)
	fmt.Println(structA, "|", string(jsonBytes))
}

func TestXmlCodec(t *testing.T) {
	structA := &StructA{}
	codec := BuildXMLCodec()
	xmlBytes, _ := codec.Encoder().Encode(structA)
	codec.Decoder().Decode(XMLBytes, &structA)
	fmt.Println(structA, "|", string(xmlBytes))
}

func TestTextCodec(t *testing.T) {
	codec := BuildTextCodec()
	body := "test"
	textBytes, _ := codec.Encoder().Encode(body)
	var res string
	err := codec.Decoder().Decode(textBytes, res)
	if nil != err {
		t.Errorf("test xml codec error:%++v", err)
		return
	}
	assert.Equal(t, body, res)
}

func findTypeWrapper(v interface{}) {
	findType(v)
}

func findType(v interface{}) {
	switch v.(type) {
	case *string:
		fmt.Println("*string")
	case string:
		fmt.Println("string")
	case *[]byte:
		fmt.Println("*[]byte")
	case []byte:
		fmt.Println("[]byte")
	default:
		fmt.Println("other")
	}
}

func TestFindType(t *testing.T) {
	str := "test"
	bs := []byte(str)
	findTypeWrapper(str)
	findTypeWrapper(&str)
	findTypeWrapper(bs)
	findTypeWrapper(&bs)
}
