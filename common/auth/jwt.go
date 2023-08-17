package auth

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
)

// LoadPublicKey loads a public key from PEM/DER/JWK-encoded data.
func LoadPublicKey(data []byte) (interface{}, string, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return cert.PublicKey, cert.SerialNumber.Text(16), nil
	}

	return nil, "", fmt.Errorf("square/go-jose: parse error, got '%s'", err1)
}

//PublicKeyFileName is the file name of public key file
const PublicKeyFileName = "server.crt"

//TokenBearer All valid token should starts with this string
const TokenBearer = "Bearer "

//TokenBearerLength is the length of TokenBearer
const TokenBearerLength = len(TokenBearer)

//JwtVerifier is a JWT verifier that used to verify the
type JwtVerifier struct {
	publicKey    interface{}
	SerialNumber string
}

//ReadFileData gets all the data from disk file
func ReadFileData(workDir string, fileName string) ([]byte, error) {
	_filepath := filepath.Join(workDir, fileName)
	if _, err := os.Stat(_filepath); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(_filepath)
	if err != nil {
		return nil, fmt.Errorf("ReadFile [%s]fail,as:%s", _filepath, err.Error())
	}
	return data, nil
}

//NewJwtVerifier create a new JwtVerifier
func NewJwtVerifier(sslDir string) (*JwtVerifier, error) {
	topDir, err := filepath.Abs(sslDir)
	if err != nil {
		return nil, err
	}

	keyBytes, err := ReadFileData(topDir, PublicKeyFileName)
	if err != nil {
		return nil, fmt.Errorf("NewMessageProcessor ReadFile [%s] fail, as[%s]", PublicKeyFileName, err.Error())
	}

	pub, id, err := LoadPublicKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("NewMessageProcessor LoadPublicKey [%s] fail, as[%s]", PublicKeyFileName, err.Error())
	}

	return &JwtVerifier{
		publicKey:    pub,
		SerialNumber: id,
	}, nil
}

//verify verifies and parses the token
func (p *JwtVerifier) verify(tokenStr string) ([]byte, error) {
	obj, err := jose.ParseSigned(tokenStr)
	if err != nil {
		return []byte{}, fmt.Errorf("Verify parse token fail,as: %s", err.Error())
	}

	plaintext, err := obj.Verify(p.publicKey)
	if err != nil {
		return []byte{}, fmt.Errorf("Verify invalid signature,as: %s", err.Error())
	}
	return plaintext, nil
}

//VerifyToken This function will verify whether the token request over is a valid token
//mainly verify these points
//one is whether the length of the token meets the requirements
//two is whether the token encryption is correct
//three is whether the token in the valid time
func (p *JwtVerifier) VerifyToken(tokenStr string) (*Passport, error) {
	tokenStr = strings.TrimSpace(tokenStr)
	if len(tokenStr) <= TokenBearerLength {
		return nil, fmt.Errorf("invalid token")
	}

	if tokenStr[0:TokenBearerLength] == TokenBearer {
		tokenStr = tokenStr[TokenBearerLength:]
	}

	payload, err := p.verify(tokenStr)
	if err != nil {
		return nil, fmt.Errorf("Verify token fail,as: %s", err.Error())
	}

	claim := &ClaimSet{}
	err = json.Unmarshal([]byte(payload), claim)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal payload fail,as: %s", err.Error())
	}

	if int64(*claim.Expiry) < time.Now().Unix() {
		return nil, fmt.Errorf("token expiry")
	}

	return claim.Passport, nil
}

// EnvDefaultString is used to get string value from OS environment,
// returns default string if the key is not find in the OS environment.
func EnvDefaultString(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
