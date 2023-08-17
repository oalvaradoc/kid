package blowfish

import (
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/crypto"
	"testing"
	"time"
)

var key = []byte("123456781234567812345678")

func TestNewAES(t *testing.T) {
	cryptoIns := NewBlowFish("PKCS7", "CBC", key)
	assert.NotNil(t, cryptoIns)
}

func TestBlowFish_Encrypt(t *testing.T) {
	paddingArray := []string{"PKCS5", "ZERO", "PKCS7", "ANSIX923", "ISO10126"}
	modeArray := []string{"CBC", "CFB", "CTR", "ECB", "OFB", "GCM"}

	for _, padding := range paddingArray {
		for _, mode := range modeArray {
			cryptoIns := NewBlowFish(padding, mode, key)
			res, err := cryptoIns.Encrypt([]byte("this is a test string"))
			assert.Nil(t, err)
			t.Logf("%s-%s  Encrypted string:%s", padding, mode, string(res))
		}
	}

}

func TestBlowFish_Encrypt2(t *testing.T) {
	cryptoIns := NewBlowFish("PKCS7", "CBC", []byte(""))

	if _, err := cryptoIns.Encrypt([]byte("test")); nil == err {
		t.Error("Expect failed, but nil response")
	}
}

func TestBlowFish_Encrypt3(t *testing.T) {
	cryptoIns := NewBlowFish("PKCS7", "", []byte(""))

	if _, err := cryptoIns.Encrypt([]byte("test")); nil == err {
		t.Error("Expect failed, but nil response")
	}
}

func TestBlowFish_Encrypt4(t *testing.T) {
	cryptoIns := NewBlowFish("", "", key)

	res, err := cryptoIns.Encrypt([]byte("test"))
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestBlowFish_Decrypt(t *testing.T) {
	paddingArray := []string{"PKCS5", "ZERO", "PKCS7", "ANSIX923", "ISO10126"}
	modeArray := []string{"CBC", "CFB", "CTR", "ECB", "OFB"}

	plainText := "test"

	for _, padding := range paddingArray {
		for _, mode := range modeArray {
			cryptoIns := NewBlowFish(padding, mode, key)

			res, err := cryptoIns.Encrypt([]byte(plainText))
			assert.Nil(t, err)

			t.Logf("%s-%s Encrypted string:%s", padding, mode, string(res))
			res, err = cryptoIns.Decrypt(res)
			assert.Nil(t, err)
			assert.Equal(t, plainText, string(res))
		}
	}
}

func TestBlowFish_Decrypt2(t *testing.T) {
	cryptoIns := NewBlowFish("", "", key)
	_, err := cryptoIns.Decrypt([]byte(""))
	assert.NotNil(t, err)
}

func TestBlowFish_Decrypt3(t *testing.T) {
	cryptoIns := NewBlowFish("PKCS7", "CBC", []byte(""))
	_, err := cryptoIns.Decrypt([]byte(""))
	assert.NotNil(t, err)
}

func TestExampleNewBlowFish(t *testing.T) {
	var cryptoIns crypto.Interface
	startTime := time.Now()
	totalTimes := 10000
	var cryptoBytes []byte
	var err error
	cryptoIns = NewBlowFish("PKCS7", "CBC", []byte("1234567890ABCDEF"))
	for i := 0; i < totalTimes; i++ {
		cryptoBytes, err = cryptoIns.Encrypt([]byte("{\"frontJnlId\":\"1\",\"fromAcctId\":\"10001\",\"toAcctId\":\"20001\",\"amount\":102}"))
		assert.Nil(t, err)

		cryptoBytes, err = cryptoIns.Decrypt(cryptoBytes)
		assert.Nil(t, err)
	}
	t.Logf("blowfish.Crypto string:[%s], times [%d], total time cost:[%s]", string(cryptoBytes), totalTimes, time.Now().Sub(startTime))
}
