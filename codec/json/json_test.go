package json

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

type TestStruct struct {
	A string
	B string
	C int
	D float64
}

func TestJSONCodec(t *testing.T) {
	testStruct := TestStruct{
		A: "test field A",
		B: "test field B",
		C: 10,
		D: 12.45,
	}
	encoder := Encoder{}
	bs, err := encoder.Encode(testStruct)
	assert.Nil(t, err)

	decoder := Decoder{}
	newStruct := TestStruct{}
	decoder.Decode(bs, &newStruct)
	assert.Equal(t, testStruct.A, newStruct.A)
	assert.Equal(t, testStruct.B, newStruct.B)
	assert.Equal(t, testStruct.C, newStruct.C)
	assert.Equal(t, testStruct.D, newStruct.D)
}
