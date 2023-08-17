package tripledes

import (
	"crypto/cipher"
	"crypto/des"
	"git.multiverse.io/eventkit/kit/common/crypto/mode"
	"git.multiverse.io/eventkit/kit/common/crypto/padding"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
)

// TripleDES is an implement of crypto.Interface for 3des encryption algorithm
type TripleDES struct {
	padding string
	mode    string
	key     []byte
}

// NewTripleDES creates a new TripleDES instance
func NewTripleDES(padding string, mode string, key []byte) *TripleDES {
	return &TripleDES{padding: padding, mode: mode, key: key}
}

// Encrypt encrypts plain text into cipher text
func (a *TripleDES) Encrypt(plainText []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(a.key)
	if nil != err {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	plainText, err = padding.Padding(plainText, block.BlockSize(), a.padding)
	if nil != err {
		return nil, err
	}

	var iv []byte

	if "ECB" != a.mode {
		iv, err = padding.GenerateIVOrNonce("3DES", a.mode)
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
func (a *TripleDES) Decrypt(preCipherText []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(a.key)
	if nil != err {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}
	var iv []byte
	var cipherText []byte

	if "ECB" != a.mode {
		if len(preCipherText) < padding.LengthOf3desCommonIv {
			return nil, errors.Errorf(constant.SystemInternalError, "Invalid length of plainText, plainText=[%s]", string(preCipherText))
		}

		iv = preCipherText[0:padding.LengthOf3desCommonIv]
		cipherText = preCipherText[padding.LengthOf3desCommonIv:]
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
