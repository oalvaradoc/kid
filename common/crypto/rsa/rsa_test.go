package rsa

import (
	"crypto"
	"encoding/base64"
	"git.multiverse.io/eventkit/kit/common/assert"
	"git.multiverse.io/eventkit/kit/common/crypto/aes"
	"testing"
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

var skmPriKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAqFSQk6Gze4KNyhpvHXu+8FrzKWTaH+Od+j9eLxNyyQvrjDcW
+TiT2Ln60poNoLWc7EL78Xg8MqKat7VsEgQOe/QmCIsre7kr9CP6/Dt1LI817TBw
PcdGKVR+iOZ14ZdXzB8rmICGfoTAMarqLJlzBF2mJYULWT9MFbQoJ5yjOppzqZjc
+A0E5FpVh4TrvQhh1zd9jonN8R851i5f91dNQ76QEm3DgqWFkGFRnkVZ5A7cOyZh
EwQBKxfta1ue86YYGiMMO/r5yTOmkk/AiTwWJlwJJFi9V7zXWtm5MwSfaaDY50bn
fVYbeSeNNqqpcUGPvbR8FNQWynBFr8xa3PDNVQIDAQABAoIBAE3VJDCPIS1n1WXw
yRYJ5OTAORUXw9/g4GgYqtT7miSp9VUVF/NOnNYmUHrWrpxzvUZlRpeFb9g95Woy
YfEGnSflYTysFQQVP+SYSnIcj/Z1lYrBzfRS0vdDUWq9nR4dW4RPmVnfe9C+Uxvk
GnlazprjLnLEzNWMdgLHFZYTUEpE5efK7rsczBlnik4/GT2Gsdj/Yo7QCIw2tbih
9NZiNr7ld7dDl9B7gbxJZeDYO7uzIdXdhP2to0FolpNH5I7eWhHNNrkW5TKRrgvv
bclUSuWW3O6wfzE4qs6VAESVBueVWCMXITy5pJ3jkR7EVrX9SS03iSOaZcqICOL5
uKcgURUCgYEA3tjg7/Q/Meg4eXLPNDoumFs1lyXXnodquDzc0MWoUwWFucct32nC
mFe3/jE1o9ulfTKj5lJORnrTEIsPhBs0wNAFBlczU3mx6VXReTuaab8WY1GzlbeX
4fq5XsbD+yx+10C2ID/1QqXB21effQgeFQa4rFS7ReeiHpfwBf9zpi8CgYEAwV9q
RpjfgzNAB6ym/OOshPtrFwaM71Uqcw6WjwfbYzm7/LwR8zfwR/aQJkCh4c+d5Xu6
pKC8wAb9NyLDI3+NtbypU4UVhV4sbzjhyBHgy4kKcI/Z/4IjYw3LMKr5+nIQ4nt2
J+UHsTeOJvmIYJuxBjq5bI+myjRFWb9UL8Ee57sCgYBAmXIrXQxstTqZyjRSmYMk
W1xfonKs2+iN2+bPBl1TI8iuIBUmLIxiiRsnLrCz/Votvt5QSA+00qoYo5ct3o0e
T68FNYYFbsOqNlxw3lxWxzQAOpDql7wJoBrYZJovV6i1UWb6VlAMr+xQX0g2gIHn
6njiS/W1v/35DGZh9rlZtQKBgQCs61MQ8HGnVHQkqLrnF/1VKbL48y2ic8ky/E+c
dc00rRMzDUcL9PDUmWMMIe3hDRTIet1LjEVdfqJ+5IIVw2GIq73LZw34pl8b0oTs
sTgRKmoAgFLUDp7wXAxgZ/SEhe4daYQeZst7KQ/gQHI42eDyjh70On1PAnElsVdq
IMsvMwKBgQCtrGcg71WwIXYbCne/qLanmbE6I14tEG5YA5u9FidRxvvNcc9I5sld
i1NNncXdI41JVE8jwmTbrdME9VmAB1e9AoJypeTlRE5agPO+des5Uf3xPNZhrF1p
BsFg79t0IDc7uFg5VhJd5PGSKAwcp9rcgFl/po1GbWVjoiSAHoLXSw==
-----END RSA PRIVATE KEY-----`

var skmPubKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqFSQk6Gze4KNyhpvHXu+
8FrzKWTaH+Od+j9eLxNyyQvrjDcW+TiT2Ln60poNoLWc7EL78Xg8MqKat7VsEgQO
e/QmCIsre7kr9CP6/Dt1LI817TBwPcdGKVR+iOZ14ZdXzB8rmICGfoTAMarqLJlz
BF2mJYULWT9MFbQoJ5yjOppzqZjc+A0E5FpVh4TrvQhh1zd9jonN8R851i5f91dN
Q76QEm3DgqWFkGFRnkVZ5A7cOyZhEwQBKxfta1ue86YYGiMMO/r5yTOmkk/AiTwW
JlwJJFi9V7zXWtm5MwSfaaDY50bnfVYbeSeNNqqpcUGPvbR8FNQWynBFr8xa3PDN
VQIDAQAB
-----END PUBLIC KEY-----`

func TestNewRSA(t *testing.T) {
	rsa := NewRSA([]byte(priKey), []byte(pubKey))
	assert.NotNil(t, rsa)
}

func TestNewRSA2(t *testing.T) {
	rsa := NewRSA(nil, nil)
	assert.NotNil(t, rsa)
}

func TestRSA_Encrypt(t *testing.T) {
	rsa := NewRSA([]byte(priKey), []byte(pubKey))
	assert.NotNil(t, rsa)

	res, err := rsa.Encrypt([]byte("test"))
	assert.Nil(t, err)

	t.Logf("encrypted string:[%s]", string(res))
}

func TestRSA_Decrypt(t *testing.T) {
	rsa := NewRSA([]byte(priKey), []byte(pubKey))
	assert.NotNil(t, rsa)

	res, err := rsa.Encrypt([]byte("test"))
	assert.Nil(t, err)

	t.Logf("encrypted string:[%s]", string(res))
	plainText, err := rsa.Decrypt(res)
	assert.Nil(t, err)

	t.Logf("plainText=%s", string(plainText))
}

func TestRSA_Sign(t *testing.T) {
	hashTypeArray := []crypto.Hash{crypto.SHA256, crypto.SHA512}
	for _, hashType := range hashTypeArray {
		rasForSign := NewRSAForSign([]byte(priKey), []byte(pubKey), hashType)

		data := []byte("{\"customer_info_th\":{\"thai_title\":\"นาย\",\"thai_first_name\":\"ทินกร\",\"thai_middle_name\":\"\",\"thai_last_name\":\"จรุงกลิ่น\",\"thai_full_name\":\"\"},\"customer_info_en\":{\"en_title\":\"MR.\",\"en_first_name\":\"Tinnagorn\",\"en_middle_name\":\"\",\"en_last_name\":\"Jarungklin\",\"en_full_name\":\"\"},\"birth_date\":\"1994-05-07\",\"identifier\":{\"card_number\":\"1103100250661\",\"card_type\":\"324001\",\"card_issuing_country\":\"TH\",\"card_issue_date\":\"2016-08-29\",\"card_expiry_date\":\"2025-05-06\"},\"customer_address_id_card\":{\"id_card_street_address1\":\"108/2 หมู่ที่ 4\",\"id_card_street_address2\":\"-\",\"id_card_address_subdistrict\":\"สำโรง\",\"id_card_address_district\":\"พระประแดง\",\"id_card_address_province\":\"สมุทรปราการ\",\"id_card_address_zipcode\":\"10200\",\"id_card_address_country\":\"TH\",\"id_card_address_full\":\"\"},\"customer_contact\":{\"gender\":\"M\",\"marital_status\":\"S\",\"nationality\":\"TH\",\"non_iso_nationality_description\":\"\"}}")
		signedData, err := rasForSign.Sign(data)
		t.Logf("signed data:%s", string(signedData))

		assert.Nil(t, err)
	}
}

func TestRSA_VerifySign(t *testing.T) {
	hashTypeArray := []crypto.Hash{crypto.SHA256, crypto.SHA512}
	for _, hashType := range hashTypeArray {
		rasForSign := NewRSAForSign([]byte(priKey), []byte(pubKey), hashType)

		data := []byte("{\"customer_info_th\":{\"thai_title\":\"นาย\",\"thai_first_name\":\"ทินกร\",\"thai_middle_name\":\"\",\"thai_last_name\":\"จรุงกลิ่น\",\"thai_full_name\":\"\"},\"customer_info_en\":{\"en_title\":\"MR.\",\"en_first_name\":\"Tinnagorn\",\"en_middle_name\":\"\",\"en_last_name\":\"Jarungklin\",\"en_full_name\":\"\"},\"birth_date\":\"1994-05-07\",\"identifier\":{\"card_number\":\"1103100250661\",\"card_type\":\"324001\",\"card_issuing_country\":\"TH\",\"card_issue_date\":\"2016-08-29\",\"card_expiry_date\":\"2025-05-06\"},\"customer_address_id_card\":{\"id_card_street_address1\":\"108/2 หมู่ที่ 4\",\"id_card_street_address2\":\"-\",\"id_card_address_subdistrict\":\"สำโรง\",\"id_card_address_district\":\"พระประแดง\",\"id_card_address_province\":\"สมุทรปราการ\",\"id_card_address_zipcode\":\"10200\",\"id_card_address_country\":\"TH\",\"id_card_address_full\":\"\"},\"customer_contact\":{\"gender\":\"M\",\"marital_status\":\"S\",\"nationality\":\"TH\",\"non_iso_nationality_description\":\"\"}}")
		var signedData []byte
		signedData, err := rasForSign.Sign(data)

		assert.Nil(t, err)

		ret, err := rasForSign.VerifySign(data, signedData)
		assert.Nil(t, err)

		assert.True(t, ret)
	}
}

func TestRAS_Encrypt(t *testing.T) {
	rsa := NewRSA([]byte(priKey), []byte(pubKey))
	res, err := rsa.Encrypt([]byte("test"))

	assert.Nil(t, err)
	t.Logf("base64ed string:[%s]", string(base64.StdEncoding.EncodeToString(res)))
}

func TestRSADecrypt_With_Input(t *testing.T) {
	inputBase64edString := "lo30JXdnDXr9Lotlkyr+HjX7gVFog2/pliaV9vWzfw7AEG3aQOHe9UV6ISoHWNLVOIWtAhYJ0xaAZDGOw/bXhuHS94i3i6QnoDi6v/UKpMi0Q7JUZrM8ifW0FRSRw3f0dAlN9SEJANYJ/bZKfTst/gIaMQXGyHXVZBZsC3qauW8lssjcPYpFN6xnrIXLCTa9MzhJvb2pi7hpoek1RuNxTbTF4om/oQbwiEh33LDAUh/uGCwTOaTIviD8wSyasAJiQvRL3D989qrSibyePK0ERTIQP3/qLlOzvnU47d5xyo/tuMNsLPmn8QWM0/hC+eT+y9PlIRuHXGGEWL7c/JYRUA=="
	rsa := NewRSA([]byte(priKey), []byte(pubKey))
	assert.NotNil(t, rsa)
	input, err := base64.StdEncoding.DecodeString(inputBase64edString)

	assert.Nil(t, err)

	plainText, err := rsa.Decrypt(input)
	assert.Nil(t, err)
	t.Logf("plainText=%s", string(plainText))
}

func TestRSAEncrypt_Mock_SKM(t *testing.T) {
	var requestStr = `{"services":[{"keyType":"AES256","serviceId":"upstreamServices2","systemType":"APP","effectTime":"2021-03-21 23:43:47","expireTime":"2021-03-23 23:43:47","operator":"withdraw-sed-127.0.0.1","requestTime":"2021-03-21 23:43:47"},{"keyType":"AES256","serviceId":"withdraw","systemType":"APP","effectTime":"2021-03-21 23:43:47","expireTime":"2021-03-23 23:43:47","operator":"withdraw-sed-127.0.0.1","requestTime":"2021-03-21 23:43:47"},{"keyType":"AES256","serviceId":"multiverse","systemType":"APP","effectTime":"2021-03-21 23:43:47","expireTime":"2021-03-23 23:43:47","operator":"withdraw-sed-127.0.0.1","requestTime":"2021-03-21 23:43:47"}]}`
	rsa := NewRSA([]byte(skmPriKey), []byte(skmPubKey))
	assert.NotNil(t, rsa)

	base64edAesEncryptedKey := "oE2cySagOdumRC5NguKvwyemFTDBQ5cyY9WJRg90zaoUOhiUrlI3zeX8yvJtQbuJD+NE2aRi+bogH9nVuM8RNWOQ7QXMmzT/sw6BPMjqpxqLYrz4kXN6ntw14aHn12UliSdCAAH02Ypyz1UvcNP17fzv2edt6lzF89Uo5eeHim9yXzgmeoDPwzefkbUo2Is61mHZKT8B5+bNyZPKZhE4JI81sTvsbYfK2Ygu9DqbFbiBDXrlDVwVHwwROOzdRNua/X9h9Tfm+YylMiGhYV2+hI24KcRMAijUwOXqRt9okOC/y6Jw1OB2yRCcI17z62wZ8zE5VRxnBh2Yy2pQuUeU2A=="
	aesEncryptedKeyBs, err := base64.StdEncoding.DecodeString(base64edAesEncryptedKey)
	assert.Nil(t, err)

	aesKey, err := rsa.Decrypt(aesEncryptedKeyBs)
	assert.Nil(t, err)
	t.Logf("aes key:[%s]", string(aesKey))

	cryptoIns := aes.NewAESWithIVOrNonce("PKCS5", "GCM", aesKey, []byte("123456789012"))
	originalBs, err := cryptoIns.Encrypt([]byte(requestStr))
	assert.Nil(t, err)

	request := base64.StdEncoding.EncodeToString(originalBs)
	t.Logf("length of request_str=[%d], length of plain_text[%d],request:[%s]", len(requestStr), len(originalBs), request)
}

func TestRSADecrypt_Mock_SKM(t *testing.T) {
	base64edAesEncryptedKey := "kko7RHaw+5fkRlmin9CPkPUoYw1eUekM/ZRg1+wfH8UX3btvdNu13J7h22vYEU08zwVhI8Ov+ENK8i4gumorncRETdfezsQbVOby+kgvYAWqsevptcrdJD0UENqFdKJnx4WtdltjCZ6OUKeVPyCvL0FUbkq2h6fFzCRK8flCY7wbh8P1cY2j5whl8sdoCl8av7hAZveIWjHK/Pv2b4c9SAD1H6LGrv2taFFhUOSu/9FyjdTwtApfGyeQeKceTcSxLmbX+JPedc2PHuR73l9NoxWX2WW4wJY9AY+0hm4oH8UFvc8dHzfXWPczVaRAiHVFm+0HUtLKv7aMhVx3e8sVow=="
	aesEncryptedKeyBs, err := base64.StdEncoding.DecodeString(base64edAesEncryptedKey)
	assert.Nil(t, err)

	rsa := NewRSA([]byte(skmPriKey), []byte(skmPubKey))
	assert.NotNil(t, rsa)

	aesKey, err := rsa.Decrypt(aesEncryptedKeyBs)
	assert.Nil(t, err)

	t.Logf("aes key:[%s]", string(aesKey))

	base64edAesEncryptedBody := "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"
	body, err := base64.StdEncoding.DecodeString(base64edAesEncryptedBody)
	assert.Nil(t, err)

	cryptoIns := aes.NewAES("PKCS5", "GCM", aesKey)
	plainText, err := cryptoIns.Decrypt(body)
	assert.Nil(t, err)

	t.Logf("plainText:[%s]", string(plainText))
}
