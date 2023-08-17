package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"hash"
)

// RSA is an implement of crypto.Interface for rsa encryption algorithm
type RSA struct {
	privateKey []byte
	publicKey  []byte
	hashType   crypto.Hash
}

// NewRSA creates a new RSA
func NewRSA(privateKey []byte, publicKey []byte) *RSA {
	return &RSA{privateKey: privateKey, publicKey: publicKey}
}

// NewRSAForSign creates a new RSA for signature
func NewRSAForSign(privateKey []byte, publicKey []byte, hashType crypto.Hash) *RSA {
	return &RSA{privateKey: privateKey, publicKey: publicKey, hashType: hashType}
}

// Encrypt encrypts plain text into cipher text
func (r *RSA) Encrypt(orgidata []byte) ([]byte, error) {
	block, _ := pem.Decode(r.publicKey)
	if block == nil {
		return nil, errors.Errorf(constant.SystemInternalError, "public key is bad")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	pub := pubInterface.(*rsa.PublicKey)

	return rsa.EncryptPKCS1v15(rand.Reader, pub, orgidata)

}

// Decrypt decrypts cipher text into plain text
func (r *RSA) Decrypt(cipertext []byte) ([]byte, error) {
	block, _ := pem.Decode(r.privateKey)
	if block == nil {
		return nil, errors.Errorf(constant.SystemInternalError, "public key is bad")

	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {

		return nil, errors.Wrap(constant.SystemInternalError, err, 0)

	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipertext)

}

// Sign is used to signature the input data
func (r *RSA) Sign(data []byte) ([]byte, error) {
	var h hash.Hash
	switch r.hashType {
	case crypto.SHA256:
		h = sha256.New()
	case crypto.SHA512:
		h = sha512.New()
	}

	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(r.privateKey)
	if block == nil {
		return nil, errors.Errorf(constant.SystemInternalError, "private key error")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, r.hashType, hashed)
	if err != nil {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	return signature, nil
}

// VerifySign is used to verify whether the data is a valid signature data.
func (r *RSA) VerifySign(data, signData []byte) (bool, error) {
	block, _ := pem.Decode(r.publicKey)
	if block == nil {
		return false, errors.Errorf(constant.SystemInternalError, "public key error")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, errors.Wrap(constant.SystemInternalError, err, 0)
	}
	var hashed []byte
	switch r.hashType {
	case crypto.SHA256:
		h := sha256.Sum256(data)
		hashed = h[:]
	case crypto.SHA512:
		h := sha512.Sum512(data)
		hashed = h[:]
	}

	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), r.hashType, hashed, signData)
	if err != nil {
		return false, nil
	}

	return true, nil

}
