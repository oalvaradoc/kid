package des

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/crypto"
	"testing"
	"time"
)

func TestNewDES(t *testing.T) {
	cryptoIns := NewDES("PKCS7", "CBC", []byte("12345678"))

	assert.NotNil(t, cryptoIns)
}

func TestDES_Encrypt(t *testing.T) {
	paddingArray := []string{"PKCS5", "ZERO", "PKCS7", "ANSIX923", "ISO10126"}
	modeArray := []string{"CBC", "CFB", "CTR", "ECB", "OFB"}

	for _, padding := range paddingArray {
		for _, mode := range modeArray {
			cryptoIns := NewDES(padding, mode, []byte("12345678"))

			res, err := cryptoIns.Encrypt([]byte("test"))
			assert.Nil(t, err)
			t.Logf("Encrypted string:%s", string(res))
		}
	}

}

func TestDES_Encrypt2(t *testing.T) {
	cryptoIns := NewDES("PKCS7", "CBC", []byte(""))

	_, err := cryptoIns.Encrypt([]byte("test"))
	assert.NotNil(t, err)
}

func TestDES_Encrypt3(t *testing.T) {
	cryptoIns := NewDES("PKCS7", "", []byte(""))

	_, err := cryptoIns.Encrypt([]byte("test"))
	assert.NotNil(t, err)
}

func TestDES_Encrypt4(t *testing.T) {
	cryptoIns := NewDES("", "", []byte("12345678"))

	res, err := cryptoIns.Encrypt([]byte("test"))
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestDES_Decrypt(t *testing.T) {
	paddingArray := []string{"PKCS5", "ZERO", "PKCS7", "ANSIX923", "ISO10126"}
	modeArray := []string{"CBC", "CFB", "CTR", "ECB", "OFB"}

	plainText := "test"

	for _, padding := range paddingArray {
		for _, mode := range modeArray {
			cryptoIns := NewDES(padding, mode, []byte("12345678"))

			res, err := cryptoIns.Encrypt([]byte(plainText))
			assert.Nil(t, err)

			t.Logf("Encrypted string:%s", string(res))
			res, err = cryptoIns.Decrypt(res)
			assert.Nil(t, err)
			assert.Equal(t, plainText, string(res))
		}
	}
}

func TestDES_Decrypt2(t *testing.T) {
	cryptoIns := NewDES("", "", []byte("12345678"))
	_, err := cryptoIns.Decrypt([]byte(""))
	assert.NotNil(t, err)
}

func TestDES_Decrypt3(t *testing.T) {
	cryptoIns := NewDES("PKCS7", "CBC", []byte(""))
	_, err := cryptoIns.Decrypt([]byte(""))
	assert.NotNil(t, err)
}

func TestExampleNewDES(t *testing.T) {
	var cryptoIns crypto.Interface
	startTime := time.Now()
	totalTimes := 10000
	var cryptoBytes []byte
	var err error
	cryptoIns = NewDES("PKCS7", "CBC", []byte("12345678"))
	for i := 0; i < totalTimes; i++ {
		cryptoBytes, err = cryptoIns.Encrypt([]byte("{\"frontJnlId\":\"1\",\"fromAcctId\":\"10001\",\"toAcctId\":\"20001\",\"amount\":102}"))
		assert.Nil(t, err)

		cryptoBytes, err = cryptoIns.Decrypt(cryptoBytes)
		assert.Nil(t, err)
	}
	t.Logf("des.Crypto string:[%s], times [%d], total time cost:[%s]", string(cryptoBytes), totalTimes, time.Now().Sub(startTime))
}
