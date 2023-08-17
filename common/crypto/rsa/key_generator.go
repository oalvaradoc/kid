package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"io"
	"os"
)

// GenRsaKey generates a pair of RSA key
func GenRsaKey(priWriter io.Writer, pubWriter io.Writer, bits int) error {
	//generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return errors.Wrap(constant.SystemInternalError, err, 0)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	err = pem.Encode(priWriter, block)
	if err != nil {
		return errors.Wrap(constant.SystemInternalError, err, 0)
	}

	// generate public key
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return errors.Wrap(constant.SystemInternalError, err, 0)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	err = pem.Encode(pubWriter, block)
	if err != nil {
		return errors.Wrap(constant.SystemInternalError, err, 0)
	}

	return nil
}

// GenRsaKeyToFile generates a pair of RSA key and writes to the disk file
func GenRsaKeyToFile(priKeyPath string, pubKeyPath string, bits int) error {
	priFile, err := os.Create(priKeyPath)
	if err != nil {
		return errors.Wrap(constant.SystemInternalError, err, 0)
	}

	pubFile, err := os.Create(pubKeyPath)
	if err != nil {
		return errors.Wrap(constant.SystemInternalError, err, 0)
	}

	return GenRsaKey(priFile, pubFile, bits)
}
