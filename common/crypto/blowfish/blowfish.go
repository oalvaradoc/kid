package blowfish

import (
	"crypto/cipher"
	"git.multiverse.io/eventkit/kit/common/crypto/mode"
	"git.multiverse.io/eventkit/kit/common/crypto/padding"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"golang.org/x/crypto/blowfish"
)

// BlowFish is an implement of crypto.Interface for blow-fish encryption algorithm
type BlowFish struct {
	padding string
	mode    string
	key     []byte
}

// NewBlowFish creates a new BlowFish instance
func NewBlowFish(padding string, mode string, key []byte) *BlowFish {
	return &BlowFish{padding: padding, mode: mode, key: key}
}

// Encrypt encrypts plain text into cipher text
func (a *BlowFish) Encrypt(plainText []byte) ([]byte, error) {
	block, err := blowfish.NewCipher(a.key)
	if nil != err {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	plainText, err = padding.Padding(plainText, block.BlockSize(), a.padding)
	if nil != err {
		return nil, err
	}

	var iv []byte

	if "ECB" != a.mode {
		iv, err = padding.GenerateIVOrNonce("BLOWFISH", a.mode)
		if nil != err {
			return nil, err
		}
	}

	switch a.mode {
	case "CBC":
		{
			blockMode := cipher.NewCBCEncrypter(block, iv)
			blockMode.CryptBlocks(plainText, plainText)

			return append(iv, plainText...), nil
		}
	case "CFB":
		{
			blockMode := cipher.NewCFBEncrypter(block, iv)
			blockMode.XORKeyStream(plainText, plainText)

			return append(iv, plainText...), nil
		}
	case "CTR":
		{
			blockMode := cipher.NewCTR(block, iv)
			blockMode.XORKeyStream(plainText, plainText)

			return append(iv, plainText...), nil
		}
	case "ECB":
		{
			encrypted := make([]byte, len(plainText))
			blockMode := mode.NewECBEncrypter(block)
			blockMode.CryptBlocks(encrypted, plainText)

			return encrypted, nil
		}
	case "OFB":
		{
			blockMode := cipher.NewOFB(block, iv)
			blockMode.XORKeyStream(plainText, plainText)

			return append(iv, plainText...), nil
		}
	}
	return nil, nil
}

// Decrypt decrypts cipher text into plain text
func (a *BlowFish) Decrypt(preCipherText []byte) ([]byte, error) {
	block, err := blowfish.NewCipher(a.key)
	if nil != err {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	var iv []byte
	var cipherText []byte

	if "ECB" != a.mode {
		if len(preCipherText) < padding.LengthOfBlowfishCommonIv {
			return nil, errors.Errorf(constant.SystemInternalError, "Invalid length of cipherText, cipherText=[%s]", string(preCipherText))
		}

		iv = preCipherText[0:padding.LengthOfBlowfishCommonIv]
		cipherText = preCipherText[padding.LengthOfBlowfishCommonIv:]
	} else {
		cipherText = preCipherText
	}

	switch a.mode {
	case "CBC":
		{
			blockMode := cipher.NewCBCDecrypter(block, iv)
			blockMode.CryptBlocks(cipherText, cipherText)
		}
	case "CFB":
		{
			blockMode := cipher.NewCFBDecrypter(block, iv)
			blockMode.XORKeyStream(cipherText, cipherText)
		}
	case "CTR":
		{
			blockMode := cipher.NewCTR(block, iv)
			blockMode.XORKeyStream(cipherText, cipherText)
		}
	case "ECB":
		{
			blockMode := mode.NewECBDecrypter(block)
			blockMode.CryptBlocks(cipherText, cipherText)
		}
	case "OFB":
		{
			blockMode := cipher.NewOFB(block, iv)
			blockMode.XORKeyStream(cipherText, cipherText)
		}
	}

	cipherText = padding.Unpadding(cipherText, a.padding)
	return cipherText, nil
}
