package factory

import (
	oriCrypto "crypto"

	"git.multiverse.io/eventkit/kit/common/crypto"
	tripledes "git.multiverse.io/eventkit/kit/common/crypto/3des"
	"git.multiverse.io/eventkit/kit/common/crypto/aes"
	"git.multiverse.io/eventkit/kit/common/crypto/blowfish"
	"git.multiverse.io/eventkit/kit/common/crypto/des"
	"git.multiverse.io/eventkit/kit/common/crypto/rsa"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
)

// CreateCryptoInstance is a factory that creates crypto.Interface instance with specified parameters for symmetric encryption.
func CreateCryptoInstance(algo string, cryptoPadding string, cryptoMode string, key []byte) (crypto.Interface, error) {
	var cryptoIns crypto.Interface

	if "AES" == algo || "AES128" == algo || "AES192" == algo || "AES256" == algo {
		cryptoIns = aes.NewAES(cryptoPadding, cryptoMode, key)
	} else if "DES" == algo {
		cryptoIns = des.NewDES(cryptoPadding, cryptoMode, key)
	} else if "3DES" == algo {
		cryptoIns = tripledes.NewTripleDES(cryptoPadding, cryptoMode, key)
	} else if "BLOWFISH" == algo {
		cryptoIns = blowfish.NewBlowFish(cryptoPadding, cryptoMode, key)
	} else {
		return nil, errors.Errorf(constant.SystemInternalError, "Cannot support algo type:[%s]", algo)
	}

	return cryptoIns, nil

}

// CreateRSAInstance is a factory that creates crypto.Interface instance with specified parameters for asymmetric encryption.
// args[0] = algo
func CreateRSAInstance(privateKey []byte, publicKey []byte, args ...string) (crypto.Interface, error) {
	if nil == privateKey && nil == publicKey {
		return nil, errors.Errorf(constant.SystemInternalError, "Please input private key or public key ")
	}
	if len(args) > 0 {
		if args[0] == "OAEP" {
			return rsa.NewRSAOAEP(privateKey, publicKey)
		} else if args[0] == "PKCS1" {
			return rsa.NewRSA(privateKey, publicKey), nil
		} else {
			return nil, errors.Errorf(constant.SystemInternalError, "Cannot support algo type:[%s]", args[0])
		}
	}
	return rsa.NewRSA(privateKey, publicKey), nil
}

// CreateRSAInstanceForSign is a factory that creates crypto.Interface instance with specified parameters for signature.
func CreateRSAInstanceForSign(privateKey []byte, publicKey []byte, hashType oriCrypto.Hash) (crypto.SignInterface, error) {
	if (nil == privateKey && nil == publicKey) || 0 == hashType {
		return nil, errors.Errorf(constant.SystemInternalError, "input parameter invalid, privateKey=[%s]-publicKey=[%s]-hashType[%d]", privateKey, publicKey, hashType)
	}

	return rsa.NewRSAForSign(privateKey, publicKey, hashType), nil
}
