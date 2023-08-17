package rsa

import (
	"encoding/base64"
	"testing"
)

var pvk = `-----BEGIN RSA PRIVATE KEY-----
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

var pbk = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoldCT+F5znJQCLRrQgSZ
M9obvNcEJndLYxakKqmf/o27WtNunRBm3JGGX8fmVEx+p4AA8+U9JqGPI0uF+FqU
rf1tW1VRfkSvsxRxlEqy6CwQZTLqBskG3zOGzVewUQGPZFaex5xtn5i4woTmpfA+
05dknXPYW6i55ZMxrfIcIAVyaKJ14CwCnTV+4jPxXw7BzEpLhBDYKDDF4lvqp5Y2
1WVqwfMcjxV8quj25zrk1sepir8XjVXcJfVp8KsHa7IC3MsIoSB0ga+C/KRT29ph
j7VX0Jc3z7ri3ZHR8zGwcT13kJNVPx1vYeYSYX9OykMjSEeKVXjJFTPkWLiVumsy
wwIDAQAB
-----END PUBLIC KEY-----`

func TestNewRSAOAEP(t *testing.T) {
	rsa, err := NewRSAOAEP([]byte(pvk), []byte(pbk))
	if err != nil {
		t.Error(err)
	}
	str := "this is a test string"

	cihper, err := rsa.Encrypt([]byte(str))
	if err != nil {
		t.Error(err)
	}
	t.Logf(base64.StdEncoding.EncodeToString(cihper))

	plain, err := rsa.Decrypt(cihper)
	if err != nil {
		t.Error(err)
	}
	t.Logf(string(plain))
}