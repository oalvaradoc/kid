package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"git.multiverse.io/eventkit/kit/common/crypto/mode"
	"git.multiverse.io/eventkit/kit/common/crypto/padding"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
)

// AES is an implement of crypto.Interface for aes encryption algorithm
type AES struct {
	padding              string
	mode                 string
	key                  []byte
	ivOrNonce            []byte
	isFixIvOrNonce       bool
	isIvOrNonceInTheBody bool
}

// NewAES creates a new AES instance
func NewAES(padding string, md string, key []byte) *AES {
	return &AES{padding: padding, mode: md, key: key, isIvOrNonceInTheBody: true}
}

// NewAESWithIVOrNonce creates a new AES with ivOrNonce
func NewAESWithIVOrNonce(padding string, md string, key []byte, ivOrNonce []byte, isIvOrNonceInTheBody ...bool) *AES {
	flag := false
	if len(isIvOrNonceInTheBody) > 0 {
		flag = isIvOrNonceInTheBody[0]
	}
	return &AES{padding: padding, mode: md, key: key, ivOrNonce: ivOrNonce, isFixIvOrNonce: true, isIvOrNonceInTheBody: flag}
}

// Encrypt encrypts plain text into cipher text
func (a *AES) Encrypt(plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if nil != err {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	plainText, err = padding.Padding(plainText, block.BlockSize(), a.padding)
	if nil != err {
		return nil, err
	}
	var ivOrNonce []byte

	if "ECB" != a.mode && !a.isFixIvOrNonce {
		ivOrNonce, err = padding.GenerateIVOrNonce("AES", a.mode)
		if nil != err {
			return nil, err
		}
	} else {
		ivOrNonce = a.ivOrNonce
	}

	switch a.mode {
	case "CBC":
		{
			blockMode := cipher.NewCBCEncrypter(block, ivOrNonce)
			blockMode.CryptBlocks(plainText, plainText)
			if a.isIvOrNonceInTheBody {
				return append(ivOrNonce, plainText...), nil
			} else {
				return plainText, nil
			}
		}
	case "CFB":
		{
			blockMode := cipher.NewCFBEncrypter(block, ivOrNonce)
			blockMode.XORKeyStream(plainText, plainText)
			if a.isIvOrNonceInTheBody {
				return append(ivOrNonce, plainText...), nil
			} else {
				return plainText, nil
			}
		}
	case "CTR":
		{
			blockMode := cipher.NewCTR(block, ivOrNonce)
			blockMode.XORKeyStream(plainText, plainText)

			if a.isIvOrNonceInTheBody {
				return append(ivOrNonce, plainText...), nil
			} else {
				return plainText, nil
			}
		}
	case "ECB":
		{
			blockMode := mode.NewECBEncrypter(block)
			encrypted := make([]byte, len(plainText))
			blockMode.CryptBlocks(encrypted, plainText)

			return encrypted, nil
		}
	case "OFB":
		{
			blockMode := cipher.NewOFB(block, ivOrNonce)
			blockMode.XORKeyStream(plainText, plainText)

			if a.isIvOrNonceInTheBody {
				return append(ivOrNonce, plainText...), nil
			} else {
				return plainText, nil
			}
		}
	case "GCM":
		{
			aesgcm, err := cipher.NewGCM(block)

			if err != nil {
				return nil, err
			}

			// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
			//nonce := make([]byte, 12)
			//if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			//	return nil, err
			//}

			cipherText := aesgcm.Seal(ivOrNonce, ivOrNonce, plainText, nil)

			return cipherText, nil

		}
	}
	return nil, nil
}

// Decrypt decrypts cipher text into plain text
func (a *AES) Decrypt(preCipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if nil != err {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	var ivOrNonce []byte
	var cipherText []byte
	if "ECB" != a.mode {
		if "GCM" == a.mode && len(preCipherText) < padding.LengthOfAes256GcmNonce {

			return nil, errors.Errorf(constant.SystemInternalError, "Invalid length of cipherText, cipherText=[%s]", string(preCipherText))
		} else if len(preCipherText) < padding.LengthOfAesCommonIv {

			return nil, errors.Errorf(constant.SystemInternalError, "Invalid length of cipherText, cipherText=[%s]", string(preCipherText))
		} else if "GCM" == a.mode {
			if a.isIvOrNonceInTheBody {
				ivOrNonce = preCipherText[0:padding.LengthOfAes256GcmNonce]
				cipherText = preCipherText[padding.LengthOfAes256GcmNonce:]
			} else {
				ivOrNonce = a.ivOrNonce
				cipherText = preCipherText
			}
		} else {
			if a.isIvOrNonceInTheBody {
				ivOrNonce = preCipherText[0:padding.LengthOfAesCommonIv]
				cipherText = preCipherText[padding.LengthOfAesCommonIv:]
			} else {
				ivOrNonce = a.ivOrNonce
				cipherText = preCipherText
			}
		}
	} else {
		cipherText = preCipherText
	}

	switch a.mode {
	case "CBC":
		{
			blockMode := cipher.NewCBCDecrypter(block, ivOrNonce)
			blockMode.CryptBlocks(cipherText, cipherText)
		}
	case "CFB":
		{
			blockMode := cipher.NewCFBDecrypter(block, ivOrNonce)
			blockMode.XORKeyStream(cipherText, cipherText)
		}
	case "CTR":
		{
			blockMode := cipher.NewCTR(block, ivOrNonce)
			blockMode.XORKeyStream(cipherText, cipherText)
		}
	case "ECB":
		{
			blockMode := mode.NewECBDecrypter(block)
			blockMode.CryptBlocks(cipherText, cipherText)
		}
	case "OFB":
		{
			blockMode := cipher.NewOFB(block, ivOrNonce)
			blockMode.XORKeyStream(cipherText, cipherText)
		}
	case "GCM":
		{
			aesgcm, err := cipher.NewGCM(block)
			if err != nil {
				return nil, err
			}

			plaintext, err := aesgcm.Open(nil, ivOrNonce, cipherText, nil)
			if err != nil {
				return nil, err
			}

			plaintext = padding.Unpadding(plaintext, a.padding)
			return plaintext, nil
		}
	}

	cipherText = padding.Unpadding(cipherText, a.padding)

	return cipherText, nil
}
