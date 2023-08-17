package text

import (
	"encoding/base64"
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

func TestTextCodec(t *testing.T) {
	testBs := []byte("this is a test string!")
	encoder := Encoder{}
	bs, err := encoder.Encode(testBs)
	assert.Nil(t, err)

	decoder := Decoder{}
	var newBytes []byte
	decoder.Decode(bs, &newBytes)
	assert.Equal(t, base64.StdEncoding.EncodeToString(testBs), base64.StdEncoding.EncodeToString(newBytes))

}
