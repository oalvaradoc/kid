package aes

import (
	"encoding/base64"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/crypto"
	"git.multiverse.io/eventkit/kit/common/crypto/padding"
	"git.multiverse.io/eventkit/kit/common/crypto/rsa"
	"testing"
	"time"
)

var key = []byte("1234567890ABCDEF1234567890ABCDEF")

func TestAES_Decrypt(t *testing.T) {
	var cryptoIns crypto.Interface
	startTime := time.Now()
	totalTimes := 100000
	var cryptoBytes []byte
	var err error
	cryptoIns = NewAES("PKCS7", "GCM", key)

	str := "this is a test string"
	for i := 0; i < totalTimes; i++ {

		cryptoBytes, err = cryptoIns.Encrypt([]byte(str))
		assert.Nil(t, err)

		cryptoBytes, err = cryptoIns.Decrypt(cryptoBytes)
		assert.Nil(t, err)
	}
	t.Logf("aes.Crypto times [%d], total time cost:[%s],string:[%s]", totalTimes, time.Now().Sub(startTime), string(cryptoBytes))
}
func TestAes256_Encrypt(t *testing.T) {
	var cryptoIns crypto.Interface
	cryptoIns = NewAES("PKCS5", "GCM", key)
	str := "this is a test string"
	cryptoBytes, err := cryptoIns.Encrypt([]byte(str))
	assert.Nil(t, err)

	t.Logf("base64ed string:[%s]", string(base64.StdEncoding.EncodeToString(cryptoBytes)))
}

func TestAes256_Decrypt_With_Input(t *testing.T) {
	base64String := "ZaSwc7P+3BxSmc4pxHe9qOnpKF4i3u5r/b5reaJNvt3EbipoVfAOg5GTEQDVB4jP8yU4VCSznf8kC/wmZhNwnjAOWPgwHTNFipBkkmqjtJKDC7/vvWwa22smSTsovaGgM+d4d1N5Mlj/U/RmLB3x5M2Sc5eVaGphgtcoPnocSlKg/Pn3Jad2ohYn9w6Xqs6xvPhJ3w423kJW4YOi"
	var cryptoIns crypto.Interface
	var keyBytes, _ = base64.StdEncoding.DecodeString("YjU1YzBiOWNkNGQ1OTgwNmIxZjVlMTgzNTBhMTgzNzA=")
	cryptoIns = NewAES("PKCS7", "GCM", keyBytes)
	bs, err := base64.StdEncoding.DecodeString(base64String)
	assert.Nil(t, err)

	res, err := cryptoIns.Decrypt(bs)
	assert.Nil(t, err)
	t.Logf("The result is:[%s]", string(res))
}

func TestExampleNewAES(t *testing.T) {
	var cryptoIns crypto.Interface
	startTime := time.Now()
	totalTimes := 200
	var cryptoBytes []byte
	var err error
	cryptoIns = NewAES("PKCS7", "GCM", key)

	str := "{\"frontJnlId\":\"1\",\"fromAcctId\":\"10001\",\"toAcctId\":\"20001\",\"amount\":102}"
	for i := 0; i < totalTimes; i++ {
		cryptoBytes, err = cryptoIns.Encrypt([]byte(str))
		assert.Nil(t, err)

		cryptoBytes, err = cryptoIns.Decrypt(cryptoBytes)
		assert.Nil(t, err)
	}

	t.Logf("aes.Crypto times [%d], total time cost:[%s],string:[%s]", totalTimes, time.Now().Sub(startTime), string(cryptoBytes))
}

var otpRequestPubKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzyLGZS8Isnmesk+OE0i2
QyaYliYtTwOxTQRBZDug5EGJDt1A46IKsAOsJE9Bkhx4PgbZrl6MGYUiJi92oqWs
U9ov3ZBpRbc6T4V8Chi6Xka9b8c0OeDVFuB5Rct3AdW2bcITpcm9Cjd4KrGE9ISN
6ZxDCnsi/zdlGPdIPPAuuoNlSYEV4JnJdfY71d3Zh9c22HLcPaskCC6rxm3f0Sev
2lt3cYPUNeY7sS4QcdKnjlHKQSZsEoHRW0FfGhus94Q8WT+2WnfhMXMxjaeizNtL
U01UndzwlN4VgIc0pdqoLZRG5icMSyJyoIV0ZgJdWvuln0xDLfKfAOjmoJeZAcQ/
zQIDAQAB
-----END PUBLIC KEY-----`

func TestExampleNewAES2(t *testing.T) {
	s := ""
	for i := 32; i <= 126; i++ {
		s += string(byte(i))
	}

	t.Logf(s)

	var cryptoIns crypto.Interface
	var err error
	key, err := base64.StdEncoding.DecodeString("ODMx7p7ylpaSWvO43ctm8A==")
	assert.Nil(t, err)
	cryptoIns = NewAES("PKCS7", "ECB", key)

	cryptoBytes, err := cryptoIns.Encrypt([]byte("{\"headerReq\":{\"cardNo\":\"\",\"ip\":\"52047,dtac-T,th,10.72.54.30,02:00:00:00:00:00\",\"reqBy\":\"0658209937\",\"reqChannel\":\"PTANDROID\",\"reqDtm\":\"2019-10-15 15:10:28.156\",\"reqID\":\"ugKWCcYOZ0METCjH1hAsn\",\"service\":\"RequestOTPService\",\"sofType\":\"\"},\"otpReq\":{\"lang\":\"th\",\"mobileNumber\":\"0658209937\",\"otpType\":\"AUT\",\"appId\":\"dFlB97hv74IPPGEKqwFDUhxIowPhVZz6ZT4CFxWNg2Vzu6cQMMLOdDjjP\",\"userId\":\"\",\"uuId\":\"eb091648-b046-4c17-a899-f2a0c29505f9\"}}"))
	if nil != err {
		panic(err)
	}

	t.Log(base64.StdEncoding.EncodeToString(cryptoBytes))

	rsaIns := rsa.NewRSA(nil, []byte(otpRequestPubKey))
	_, err = rsaIns.Encrypt(key)
	assert.Nil(t, err)
	text := base64.StdEncoding.EncodeToString([]byte("sq394WZ0OdpjHTDW"))

	t.Log("text", text)
}

var key2 = []byte("123456781234567812345678")

func TestNewAES(t *testing.T) {
	cryptoIns := NewAES("PKCS7", "CBC", key2)

	assert.NotNil(t, cryptoIns)
}

func TestNewAESWithIVOrNonce(t *testing.T) {
	ivOrNonce, err := padding.GenerateIVOrNonce("AES", "CBC")
	if nil != err {
		t.Errorf("generate IV OR Nonce failed, error=%++v", err)
	}
	cryptoIns := NewAESWithIVOrNonce("PKCS7", "CBC", key2, ivOrNonce)

	assert.NotNil(t, cryptoIns)
}

func TestAES_Encrypt(t *testing.T) {
	paddingArray := []string{"PKCS5", "ZERO", "PKCS7", "ANSIX923", "ISO10126"}
	modeArray := []string{"CBC", "CFB", "CTR", "ECB", "OFB", "GCM"}

	for _, padding := range paddingArray {
		for _, mode := range modeArray {
			if "PKCS5" == padding && ("CBC" == mode || "ECB" == mode) {
				continue
			}
			cryptoIns := NewAES(padding, mode, key2)
			t.Logf("padding:[%s], mode:[%s]", padding, mode)
			res, err := cryptoIns.Encrypt([]byte("this is a test string"))
			assert.Nil(t, err)
			t.Logf("Encrypted string:%s", string(res))
		}
	}
}

func TestAES_EncryptWithIVOrNonce(t *testing.T) {
	paddingArray := []string{ "ZERO", "PKCS7", "ANSIX923", "ISO10126"}
	modeArray := []string{"CBC", "CFB", "CTR", "ECB", "OFB", "GCM"}
	//modeArray := []string{ "CTR", "ECB","OFB", "GCM"}
	ivOrNonceMap := make(map[string][]byte, 0)

	for _, mode := range modeArray {
		ivOrNonce, err := padding.GenerateIVOrNonce("AES", mode)
		assert.Nil(t, err)

		ivOrNonceMap[mode] = ivOrNonce
	}

	for _, paddingStr := range paddingArray {
		for _, mode := range modeArray {
			cryptoIns := NewAESWithIVOrNonce(paddingStr, mode, key2, ivOrNonceMap[mode])
			res, err := cryptoIns.Encrypt([]byte("this is a test string"))
			assert.Nil(t, err)

			t.Logf("Encrypted string,[%s]-[%s]:%s", paddingStr, mode, string(res))
		}
	}
}

func TestAES_Encrypt2(t *testing.T) {
	cryptoIns := NewAES("PKCS7", "CBC", []byte(""))

	if _, err := cryptoIns.Encrypt([]byte("this is a test string")); nil == err {
		t.Error("Expect failed, but nil response")
	}
}

func TestAES_Encrypt3(t *testing.T) {
	cryptoIns := NewAES("PKCS7", "", []byte(""))

	_, err := cryptoIns.Encrypt([]byte("this is a test string"))
	assert.NotNil(t, err)
}

func TestAES_Encrypt4(t *testing.T) {
	cryptoIns := NewAES("", "", key2)

	_, err := cryptoIns.Encrypt([]byte("this is a test string"))
	assert.Nil(t, err)
}

func TestAES_Decrypt2(t *testing.T) {
	cryptoIns := NewAES("", "", key2)
	_, err := cryptoIns.Decrypt([]byte(""))
	assert.NotNil(t, err)
}

func TestAES_Decrypt3(t *testing.T) {
	keyBytes, err := base64.StdEncoding.DecodeString("WVRZMFpqZGpOekZtTVRNd05EazJPR1JoWXpRd01qY3k=")
	assert.Nil(t, err)

	t.Log(string(keyBytes))
	cryptoIns := NewAES("ZERO", "GCM", []byte("NzYzMTZmOWM2M2ZhY2ZhMTRmM2U2YWQy"))
	res, err := cryptoIns.Encrypt([]byte("12345678"))
	assert.Nil(t, err)
	t.Log(base64.StdEncoding.EncodeToString(res))

	_, err = cryptoIns.Decrypt([]byte("YkhiNDYxdDJMd095Xjc9kGB0zieBJs2Qymo6su6xCH0oQdIjco24bgvLthI="))
	assert.NotNil(t, err)
}
