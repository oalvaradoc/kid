package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
)

type RSAOAEP struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewRSAOAEP(privateKey, publicKey []byte) (res *RSAOAEP, err error) {
	res = &RSAOAEP{}
	if privateKey != nil {
		block, _ := pem.Decode(privateKey)
		if block == nil {
			return nil, errors.Errorf(constant.SystemInternalError, "private key is bad")
		}
		res.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.Wrap(constant.SystemInternalError, err, 0)
		}
	}

	if publicKey != nil {
		block, _ := pem.Decode(publicKey)
		if block == nil {
			return nil, errors.Errorf(constant.SystemInternalError, "public key is bad")
		}
		pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, errors.Wrap(constant.SystemInternalError, err, 0)
		}
		res.publicKey = pubInterface.(*rsa.PublicKey)
	}

	return res, err
}

func (r *RSAOAEP) Encrypt(orgidata []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, r.publicKey, orgidata, nil)
}

func (r *RSAOAEP) Decrypt(cipertext []byte) ([]byte, error) {
	return r.privateKey.Decrypt(nil, cipertext, &rsa.OAEPOptions{Hash: crypto.SHA256})
}
