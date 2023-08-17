package auto

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestAutoCodec(t *testing.T) {
	codec := BuildAutoCodecWithJSONCodec()
	// ---------string---------
	bodyStr := "test"
	textBytes, err := codec.Encoder().Encode(bodyStr)
	if nil != err {
		t.Errorf("test auto codec error:%++v", err)
		return
	}

	var resString string
	err = codec.Decoder().Decode(textBytes, &resString)
	if nil != err {
		t.Errorf("test auto codec error:%++v", err)
		return
	}
	assert.Equal(t, bodyStr, resString)
	// ---------string---------
	// ---------*string---------
	bodyStrForPointer := "test"
	textBytes, err = codec.Encoder().Encode(&bodyStrForPointer)
	if nil != err {
		t.Errorf("test auto codec error:%++v", err)
		return
	}
	var resStringForPointer string
	err = codec.Decoder().Decode(textBytes, &resStringForPointer)
	if nil != err {
		t.Errorf("test auto codec error:%++v", err)
		return
	}
	assert.Equal(t, bodyStrForPointer, resStringForPointer)
	// ---------*string---------
	// ---------bool---------
	var bv = true
	textBytes, err = codec.Encoder().Encode(bv)
	if nil != err {
		t.Errorf("test auto codec error:%++v", err)
		return
	}
	var resBool bool
	err = codec.Decoder().Decode(textBytes, &resBool)
	if nil != err {
		t.Errorf("test auto codec error:%++v", err)
		return
	}
	assert.Equal(t, bv, resBool)
	// ---------bool---------
	// ---------*bool---------
	var bvForPointer = true
	textBytes, err = codec.Encoder().Encode(&bvForPointer)
	if nil != err {
		t.Errorf("test auto codec error:%++v", err)
		return
	}
	var resBoolForPointer bool
	err = codec.Decoder().Decode(textBytes, &resBoolForPointer)
	if nil != err {
		t.Errorf("test auto codec error:%++v", err)
		return
	}
	assert.Equal(t, bvForPointer, resBoolForPointer)
	// ---------*bool---------

	// ---------[]bool---------
	var bvArray = []bool{true, true, false, false, true}
	textBytes, err = codec.Encoder().Encode(bvArray)
	if nil != err {
		t.Errorf("test auto codec error:%++v", err)
		return
	}
	resBoolArray := make([]bool, 0)
	err = codec.Decoder().Decode(textBytes, &resBoolArray)
	if nil != err {
		t.Errorf("test auto codec error:%++v", err)
		return
	}
	assert.Equal(t, bvArray, resBoolArray)
	// ---------[]bool---------

	// ---------*[]bool---------
	var bvArrayForPointer = []bool{true, true, false, false, true}
	textBytes, err = codec.Encoder().Encode(&bvArrayForPointer)
	if nil != err {
		t.Errorf("test auto codec error:%++v", err)
		return
	}
	resBoolArrayForPointer := make([]bool, 0)
	err = codec.Decoder().Decode(textBytes, &resBoolArrayForPointer)
	if nil != err {
		t.Errorf("test auto codec error:%++v", err)
		return
	}
	assert.Equal(t, bvArrayForPointer, resBoolArrayForPointer)
	// ---------*[]bool---------

}
