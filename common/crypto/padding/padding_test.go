package padding

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"testing"
)

type paddingFunc func(cipherText []byte, blockSize int) ([]byte, error)
type unPaddingFunc func(originData []byte) []byte

func TestPaddingAndUnPaddingUsingFunc(t *testing.T) {
	paddingFunctions := []paddingFunc{PKCS5Padding, PKCS7Padding, ZeroPadding, ANSIX923Padding, ISO10126Padding}
	unPaddingFunctions := []unPaddingFunc{PKCS5UnPadding, PKCS7UnPadding, ZeroUnPadding, ANSIX923UnPadding, ISO10126UnPadding}
	cipherText := []byte("this is a test string")
	for i := range paddingFunctions {
		plainText, err := paddingFunctions[i](cipherText, 100)
		assert.Nil(t, err)
		assert.Equal(t, len(plainText), 100)
		assert.NotEqual(t, plainText, cipherText)

		cipherText2 := unPaddingFunctions[i](plainText)
		assert.Equal(t, cipherText, cipherText2)
	}
}

func TestPaddingAndUnPadding(t *testing.T) {
	paddingList := []string{"PKCS5", "ZERO", "PKCS7", "ANSIX923", "ISO10126"}
	cipherText := []byte("this is a test string")
	for _, paddingName := range paddingList {
		plainText, err := Padding(cipherText, 100, paddingName)
		assert.Nil(t, err)
		if "PKCS5" != paddingName {
			t.Logf("padding name:%s", paddingName)
			assert.Equal(t, len(plainText), 100)
		}
		assert.NotEqual(t, plainText, cipherText)

		cipherText2 := Unpadding(plainText, paddingName)
		assert.Equal(t, cipherText, cipherText2)
	}
}

func TestGenerateIVOrNonce(t *testing.T) {
	algoList := []string{"DES", "3DES", "BLOWFISH"}
	for _, s := range algoList {
		v, err := GenerateIVOrNonce(s, "")
		assert.Nil(t, err)
		assert.NotNil(t, v)
	}

	v, err := GenerateIVOrNonce("AES", "")
	assert.Nil(t, err)
	assert.NotNil(t, v)

	v, err = GenerateIVOrNonce("AES", "GCM")
	assert.Nil(t, err)
	assert.NotNil(t, v)

}
