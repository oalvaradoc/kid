package factory

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"testing"
	"time"

	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/crypto/rsa"
)

var priKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAoldCT+F5znJQCLRrQgSZM9obvNcEJndLYxakKqmf/o27WtNu
nRBm3JGGX8fmVEx+p4AA8+U9JqGPI0uF+FqUrf1tW1VRfkSvsxRxlEqy6CwQZTLq
BskG3zOGzVewUQGPZFaex5xtn5i4woTmpfA+05dknXPYW6i55ZMxrfIcIAVyaKJ1
4CwCnTV+4jPxXw7BzEpLhBDYKDDF4lvqp5Y21WVqwfMcjxV8quj25zrk1sepir8X
jVXcJfVp8KsHa7IC3MsIoSB0ga+C/KRT29phj7VX0Jc3z7ri3ZHR8zGwcT13kJNV
Px1vYeYSYX9OykMjSEeKVXjJFTPkWLiVumsywwIDAQABAoIBAAgsAk+JFxuYT4UQ
p/GLz7Z3fTv1SuUwzh+vzRXEsiQbOForGH9ZiwQBY1VA98w4iYue+u1MFdby/QSW
0aidzqwvfKjDU7XaeUm3drwzQmxDg5PEi1lKF0l3C4scpeh9/pzba2S68B2/j1Vj
YUTrIg5+qXbvlO2QQcNXtIGAFYBbEUqCmfp5KNOyTSykJvxZUgqB6UFZuHfcj/LK
9i4MhuUmumlbac3zB6OFwuUGq69cNPm6k6wHRiU3IxwdKnQU9zdLsTDLG7nsgU70
aoJBmwQvRh2aGgCgJdmGFNpEqfiF9CIkhBTCOrEJinkg9hbwSruEHQANqQ47Xd9G
yS9zsJECgYEA1XmYFiMlPGEJgdbG+2p29tGD6fPjnQjkcPfWjvED1bnlCbzOVbTg
8fQlIlKAfaOa+7Pl3YJ5kz4o448lHaIQnnNnn9GcYA0xpgBPdOQJH/EXVntfjV16
f1wmSig7k3s9tmHiVh8uQA2FQZrjE6gvxx1W0C6K7ktdSsTQDzdwC3kCgYEAwq4F
SsfhMWkrYK/Bau+uMFZsYmV/u6i5X5jkutlb6LLFjcnLNvSGNU/WWCDUHGRNofi3
cKAWjGD8sFe5ZSa3kACPD2Dp7WuMvof0OVGrIkkoa0P5OqWcTpRIMY2KTycz7Y/L
mIRmuvzUAsZ/vDA/JstZme9405R+q/8ZdGmdpRsCgYEAjijCQgO3mUTZqvBXZDga
7vTJTvQOYJX6Ysx4won3zs1TnC9yjJq+rgGy9O9SB9j6raG4ctGfmpFrc1bxFZHG
VW5u1HwnEcPXiz9rqmDtPqszqnDQSfi1SbkY+oteWTFaAGmg608qYpdeZTj6/S0k
XAnKtSo5dMUVZGQ6VdfKMqkCgYEAkoxZI3/vfziCFNh5KzydzXlhQXjSfLt4QARi
Ol2hGDxrBl8fgJD17m/ZFKIxyeWfowwNWtTH5Iil70E6KHDKwbYJ+zOjJLxPSKYj
LHrT7o1Pxd93X7SHQ4fQCK1Zrlf+eRhD1N3mT2A/YI94XHudLmDpZD2moO8po+P2
j3Fp4H0CgYBOpcfblOj86w4iQWG30Vc2oyKfegyJHVlAdTM/jvRt/g5HTwi2i7Eq
Iml7vKLUw80/1sPehHC/U4ZXT7GjQ7junYw91hzOmKpA+ziH0VObicxwdMDxg0Jf
JblIgqBCH5vYxty3TYelI7cFGYdRs9lv0BeiO7es52k1flYKJZMh1Q==
-----END RSA PRIVATE KEY-----`

var pubKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoldCT+F5znJQCLRrQgSZ
M9obvNcEJndLYxakKqmf/o27WtNunRBm3JGGX8fmVEx+p4AA8+U9JqGPI0uF+FqU
rf1tW1VRfkSvsxRxlEqy6CwQZTLqBskG3zOGzVewUQGPZFaex5xtn5i4woTmpfA+
05dknXPYW6i55ZMxrfIcIAVyaKJ14CwCnTV+4jPxXw7BzEpLhBDYKDDF4lvqp5Y2
1WVqwfMcjxV8quj25zrk1sepir8XjVXcJfVp8KsHa7IC3MsIoSB0ga+C/KRT29ph
j7VX0Jc3z7ri3ZHR8zGwcT13kJNVPx1vYeYSYX9OykMjSEeKVXjJFTPkWLiVumsy
wwIDAQAB
-----END PUBLIC KEY-----`

const (
	rsa_priKey1 = `-----BEGIN RSA PRIVATE KEY-----
MIICWgIBAAKBgQDO5COnGN8388mXrkWT8JFP7eqoxZdcw7lwacccsxRd4ff4p9vB
Ohk7Zzc07r0gOg0HiHt2qAS6lj286eNgz/TUG2CYKEYFxpt1LWDI6ytwJlboNLR7
WGM0m5DJ5g7pyo59/5vfIA+7SZ65pZbNTe7DnaLhddcLWiqMkx+2VkBzdwIDAQAB
AoGAFMQNWA5FCWasy06wqSKyUyV8MihzAtqaWFAlrhnDZ5DwxMKEaiactbusbOGx
lfR9rk3ipoxCvT+rPrTzH5p/5kP+kL6Hy+FR2EbvEufnAJDg0WBWn5ImogeyIXPu
tUmecGO34lDCEgoQUkDv6mjR0SyIQB3RjV9d6/3SfSH/SsECQQDyR4lHVPq84Ii5
JTBUzZVsLwZivYMZsW/y+Mucsn38MVfEfJKzRRoKKLD+Rkf3YKuVCZs6l+ODO07/
T0H9A1kHAkEA2puPG8OSRayKMhkfDHI6n/MFqD68v6jPYxSj/snjRTqXFfv1xPa7
PVhxrp/cQFz66N1G7eQyn3cqcJ0zpgumEQJAP+z4H8YgUm28JX3Wfsmvv1e5C5yN
Vt4md6mFr9a4vy4VxlZILtzwvfV2neDVZEQxgaWDO7aP5TRk56B1/NhBSQI/RuAo
hdfilLRcGeILLv3aBAHG08WDbKBOnNEUWocaKFfWpEoMZM+Z5UnHkdZCkpuSve0A
EiDqSMlZ+Sj+ldcxAkBgjGU1ANg2LHDbcxJ/dAjT55RrFeXeD1vA5/ZESsByZKHb
ddvC/MlXbQe3bM/g6v2DjETRzzSVLnymRgf+45AR
-----END RSA PRIVATE KEY-----`
	rsa_pubKey1 = `-----BEGIN RSA PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDO5COnGN8388mXrkWT8JFP7eqo
xZdcw7lwacccsxRd4ff4p9vBOhk7Zzc07r0gOg0HiHt2qAS6lj286eNgz/TUG2CY
KEYFxpt1LWDI6ytwJlboNLR7WGM0m5DJ5g7pyo59/5vfIA+7SZ65pZbNTe7DnaLh
ddcLWiqMkx+2VkBzdwIDAQAB
-----END RSA PUBLIC KEY-----`

	rsa_priKey2 = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDiJMAqGxnHcSgFMFlUa/mbPfu0sEa4CoQnHRHlrP/RSgBC6qLh
G8wS3FD+RAiDqko7gVZgR1DytvlnlSNrcOYHVqpsKfwhiF2g2X4/GxPVyFa9gvZB
a6SG7GpcXyvHAq3MamtIdeGPDvnkWAPz3PRDpFccanRcXOqeMsPLiF76gQIDAQAB
AoGBAI1bt5saUaTvwLptnIk+7UnzFtG9lpcYS78/Vp6g40/p1/v8O1BHVes8OIyX
7lKPMdO8Z0fLjHgLlB8BhKB2c/JZiXxF/ADu3YYdvCivCphBQcYr94R/3CQgJYpA
a8eXnzYSB0rFTWgpW7CH9StqLlckvM5ea4tWuvE/k3Ce3PQhAkEA9PxZn/AP2OpG
CCf19dvqGz2E7Vrm2nxKvU8xuTjasDLuShxVc459pkqtfgOlmchaoVoBU6zUvwm/
ozyVswk31QJBAOxPiZZuXIgNs3WZSSmygc4bIwbUZ652HHyiZPCgEpuaCj/eySLx
j4XuWw+bkws7Mx58J8Ekeuq3I1web9kTGf0CQCIwXWmenPeOqjtVKFQpXqByk2x0
dScklWGZ/bx1nL9ePDcHgT1hM1PTtCaT57ZwaYV/BBRjWEVY3O+w8stLjAkCQQCl
Fguwo/jQs4GjriqGjsZQDnUx2EF2h9zu1SRfVfSp77spU6KAXvE9R38mMDFRr1HP
Aj1jmPCl+LsjJ8BLjiShAkEA46jlksq0SNqik5aIYfYLD2eCaeRYQBihBK1XBWvv
ZHPwQnB1mgd3/aldmpxchVXrcDdaIGij5dwxMiDXgaZh0A==
-----END RSA PRIVATE KEY-----`
	rsa_pubKey2 = `-----BEGIN RSA PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDiJMAqGxnHcSgFMFlUa/mbPfu0
sEa4CoQnHRHlrP/RSgBC6qLhG8wS3FD+RAiDqko7gVZgR1DytvlnlSNrcOYHVqps
KfwhiF2g2X4/GxPVyFa9gvZBa6SG7GpcXyvHAq3MamtIdeGPDvnkWAPz3PRDpFcc
anRcXOqeMsPLiF76gQIDAQAB
-----END RSA PUBLIC KEY-----`
)

func TestRSASign(t *testing.T) {
	rasForSign, err := CreateRSAInstanceForSign([]byte(priKey), []byte(pubKey), crypto.SHA256)
	assert.Nil(t, err)

	totalTimes := 10000
	data := []byte("{\"customer_info_th\":{\"thai_title\":\"นาย\",\"thai_first_name\":\"ทินกร\",\"thai_middle_name\":\"\",\"thai_last_name\":\"จรุงกลิ่น\",\"thai_full_name\":\"\"},\"customer_info_en\":{\"en_title\":\"MR.\",\"en_first_name\":\"Tinnagorn\",\"en_middle_name\":\"\",\"en_last_name\":\"Jarungklin\",\"en_full_name\":\"\"},\"birth_date\":\"1994-05-07\",\"identifier\":{\"card_number\":\"1103100250661\",\"card_type\":\"324001\",\"card_issuing_country\":\"TH\",\"card_issue_date\":\"2016-08-29\",\"card_expiry_date\":\"2025-05-06\"},\"customer_address_id_card\":{\"id_card_street_address1\":\"108/2 หมู่ที่ 4\",\"id_card_street_address2\":\"-\",\"id_card_address_subdistrict\":\"สำโรง\",\"id_card_address_district\":\"พระประแดง\",\"id_card_address_province\":\"สมุทรปราการ\",\"id_card_address_zipcode\":\"10200\",\"id_card_address_country\":\"TH\",\"id_card_address_full\":\"\"},\"customer_contact\":{\"gender\":\"M\",\"marital_status\":\"S\",\"nationality\":\"TH\",\"non_iso_nationality_description\":\"\"}}")
	var signedData []byte
	startTime := time.Now()
	for i := 0; i < totalTimes; i++ {
		signedData, err = rasForSign.Sign(data)
		assert.Nil(t, err)
	}

	t.Log("Total loop sign times:", totalTimes, "Total time cost:", time.Now().Sub(startTime), "Signed data:", base64.StdEncoding.EncodeToString(signedData))
	var ret bool
	startTime = time.Now()
	for i := 0; i < totalTimes; i++ {
		ret, err = rasForSign.VerifySign(data, signedData)
		assert.Nil(t, err)
	}
	t.Log("Total loop verify sign times:", totalTimes, "Total time cost:", time.Now().Sub(startTime), "ms")

	assert.Nil(t, err)
	t.Log("Result of verify:", ret)
}

func TestCreateCryptoInstance(t *testing.T) {
	algoArray := []string{"AES", "AES128", "AES192", "AES256", "DES", "3DES", "BLOWFISH"}

	for _, algo := range algoArray {
		_, err := CreateCryptoInstance(algo, "ZERO", "CBC", []byte("123456"))
		assert.Nil(t, err)
	}
}

func TestCreateCryptoInstance2(t *testing.T) {
	_, err := CreateCryptoInstance("", "ZERO", "CBC", []byte("123456"))
	assert.NotNil(t, err)
}

func TestCreateRSAInstance(t *testing.T) {
	priKeyBuffer := bytes.NewBuffer(make([]byte, 0))
	pubKeyBuffer := bytes.NewBuffer(make([]byte, 0))

	rsa.GenRsaKey(priKeyBuffer, pubKeyBuffer, 1024)

	_, err := CreateRSAInstance(priKeyBuffer.Bytes(), pubKeyBuffer.Bytes(), "OAEP")

	assert.Nil(t, err)
}

func TestCreateRSAInstance2(t *testing.T) {
	_, err := CreateRSAInstance(nil, nil)

	assert.NotNil(t, err)
}

func TestCreateRSAInstanceForSign(t *testing.T) {
	priKeyBuffer := bytes.NewBuffer(make([]byte, 0))
	pubKeyBuffer := bytes.NewBuffer(make([]byte, 0))

	_, err := CreateRSAInstanceForSign(priKeyBuffer.Bytes(), pubKeyBuffer.Bytes(), crypto.SHA256)

	assert.Nil(t, err)
}

func TestCreateRSAInstanceForSign2(t *testing.T) {

	_, err := CreateRSAInstanceForSign(nil, nil, crypto.SHA256)

	assert.NotNil(t, err)
}

func TestCreateRSAOAEPInstance(t *testing.T) {
	rsa, err := CreateRSAInstance([]byte(rsa_priKey1), []byte(rsa_pubKey2), "OAEP")
	if err != nil {
		t.Error(err)
	}
	res, err := rsa.Encrypt([]byte("another message"))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(base64.StdEncoding.EncodeToString(res)))

	ciphertext, err := base64.StdEncoding.DecodeString("BMYxt9TDW9FHWT4YUi92brkzDSQXT7vjlV2+JrM/34ubJ6xXjZUxUwJjw0tqT3lYNjXHR8P0UGCkA5jInudgYjn5zAWvyhKZk09inTBqaySqQwNCdVkbY5pkO33dgwTRsMCRvGZqbHVRjhwXzt8boBdNhUwVB8l23eCOBG8HR3Q=")
	if err != nil {
		t.Error(err)
	}
	res, err = rsa.Decrypt(ciphertext)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(res))
}
