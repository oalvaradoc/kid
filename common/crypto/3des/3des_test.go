package tripledes

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/crypto"
	"testing"
	"time"
)

var key = []byte("123456789123456789123456")

func TestNewTripleDES(t *testing.T) {
	cryptoIns := NewTripleDES("PKCS7", "CBC", key)
	assert.NotNil(t, cryptoIns)
}

func TestTripleDES_Encrypt(t *testing.T) {
	paddingArray := []string{"PKCS5", "ZERO", "PKCS7", "ANSIX923", "ISO10126"}
	modeArray := []string{"CBC", "CFB", "CTR", "ECB", "OFB"}

	for _, padding := range paddingArray {
		for _, mode := range modeArray {
			cryptoIns := NewTripleDES(padding, mode, key)
			res, err := cryptoIns.Encrypt([]byte("test"))

			assert.Nil(t, err)
			t.Logf("Encrypted string:%s", string(res))
		}
	}

}

func TestTripleDES_Encrypt2(t *testing.T) {
	cryptoIns := NewTripleDES("PKCS7", "GCM", []byte(""))
	_, err := cryptoIns.Encrypt([]byte("this is a test string"))

	assert.NotNil(t, err)
}

func TestTripleDES_Encrypt3(t *testing.T) {
	cryptoIns := NewTripleDES("PKCS7", "", []byte(""))

	_, err := cryptoIns.Encrypt([]byte("this is a test string"))
	assert.NotNil(t, err)
}

func TestTripleDES_Encrypt4(t *testing.T) {
	cryptoIns := NewTripleDES("", "", key)

	res, err := cryptoIns.Encrypt([]byte("this is a test string"))
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestTripleDES_Decrypt(t *testing.T) {
	paddingArray := []string{"PKCS5", "ZERO", "PKCS7", "ANSIX923", "ISO10126"}
	modeArray := []string{"CBC", "CFB", "CTR", "ECB", "OFB"}

	plainText := "test"

	for _, padding := range paddingArray {
		for _, mode := range modeArray {
			cryptoIns := NewTripleDES(padding, mode, key)

			res, err := cryptoIns.Encrypt([]byte(plainText))
			assert.Nil(t, err)
			t.Logf("Encrypted string:%s", string(res))

			res, err = cryptoIns.Decrypt(res)
			assert.Nil(t, err)
			assert.Equal(t, plainText, string(res))
		}
	}
}

func TestTripleDES_Decrypt2(t *testing.T) {
	cryptoIns := NewTripleDES("", "", key)
	_, err := cryptoIns.Decrypt([]byte(""))
	assert.NotNil(t, err)
}

func TestTripleDES_Decrypt3(t *testing.T) {
	cryptoIns := NewTripleDES("PKCS7", "CBC", []byte(""))
	_, err := cryptoIns.Decrypt([]byte(""))
	assert.NotNil(t, err)
}

func TestExampleNewTripleDES(t *testing.T) {
	var cryptoIns crypto.Interface
	startTime := time.Now()
	totalTimes := 1
	var cryptoBytes []byte
	var err error
	cryptoIns = NewTripleDES("PKCS7", "CBC", key)
	for i := 0; i < totalTimes; i++ {
		cryptoBytes, err = cryptoIns.Encrypt([]byte("{\"jnlId\":\"123456\",\"fromAcctId\":\"10001\",\"toAcctId\":\"20001\",\"amount\":102}"))
		assert.Nil(t, err)

		t.Logf("Crypto string:[%s]", string(cryptoBytes))

		cryptoBytes, err = cryptoIns.Decrypt(cryptoBytes)
		assert.Nil(t, err)
	}
	t.Logf("3des.Crypto string:[%s], times [%d], total time cost:[%s]", string(cryptoBytes), totalTimes, time.Now().Sub(startTime))
}
