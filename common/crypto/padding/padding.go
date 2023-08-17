package padding

import (
	"bytes"
	"crypto/rand"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"io"
)

// PKCS5Padding pkcs5 is a subset algorithm of pkcs7. There is no difference in concept, but the blockSize is fixed at 8 bytes,
// that is, the data will always be cut into 8-byte data blocks, and then the length to be filled is calculated.
// The pkcs7 padding length blockSize is 1~255 bytes
func PKCS5Padding(cipherText []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...), nil
}

// PKCS5UnPadding removes the padding data
func PKCS5UnPadding(originData []byte) []byte {
	return lengthUnPadding(originData)
}

// ZeroPadding removes the padding data
func ZeroPadding(cipherText []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{0}, padding)
	return append(cipherText, padText...), nil
}

// ZeroUnPadding uses "0" as the padding method of filling data, that is to say, when grouping, the length of the last
// group of plaintext does not reach the packet length, then use "0" to make up
func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

// PKCS7Padding takes the last character paddingChar, the decimal ord (paddingChar) of the
// ASCII code of this character is the padding data length paddingSize
func PKCS7Padding(cipherText []byte, blocksize int) ([]byte, error) {
	padding := blocksize - len(cipherText)%blocksize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...), nil
}

func lengthUnPadding(originData []byte) []byte {
	length := len(originData)
	unpadding := int(originData[length-1])
	return originData[:(length - unpadding)]
}

// PKCS7UnPadding removes the padding data
func PKCS7UnPadding(originData []byte) []byte {
	return lengthUnPadding(originData)
}

// ANSIX923Padding The last byte of the padding sequence is filled with paddingSize, and the others are filled with 0.
func ANSIX923Padding(cipherText []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(cipherText)%blockSize

	padText := bytes.Repeat([]byte{0}, padding-1)
	padText = append(padText, byte(padding))

	return append(cipherText, padText...), nil
}

// ANSIX923UnPadding removes the padding data
func ANSIX923UnPadding(originData []byte) []byte {
	return lengthUnPadding(originData)
}

// ISO10126Padding The last byte of the padding sequence is filled with paddingSize, and the others are filled with random numbers.
func ISO10126Padding(cipherText []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(cipherText)%blockSize
	padText := make([]byte, padding-1)

	if _, err := io.ReadFull(rand.Reader, padText); err != nil {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}
	padText = append(padText, byte(padding))

	return append(cipherText, padText...), nil
}

// ISO10126UnPadding removes the padding data
func ISO10126UnPadding(originData []byte) []byte {
	return lengthUnPadding(originData)
}

// Define all supported IVOrNonce length
const (
	LengthOfAesCommonIv      = 16
	LengthOfAes256GcmNonce   = 12
	LengthOfDesCommonIv      = 8
	LengthOf3desCommonIv     = 8
	LengthOfBlowfishCommonIv = 8
	FixBlockSizeOfPKCS5      = 8
)

// Padding fills the cipher text according to the block size and padding type
func Padding(cipherText []byte, blockSize int, padding string) ([]byte, error) {
	var err error
	switch padding {
	case "PKCS5":
		cipherText, err = PKCS5Padding(cipherText, 8)
	case "ZERO":
		cipherText, err = ZeroPadding(cipherText, blockSize)
	case "PKCS7":
		cipherText, err = PKCS7Padding(cipherText, blockSize)
	case "ANSIX923":
		cipherText, err = ANSIX923Padding(cipherText, blockSize)
	case "ISO10126":
		cipherText, err = ISO10126Padding(cipherText, blockSize)
	}

	return cipherText, err
}

// Unpadding cancel filling according to filling type
func Unpadding(cipherText []byte, padding string) []byte {
	switch padding {
	case "PKCS5":
		cipherText = PKCS5UnPadding(cipherText)
	case "ZERO":
		cipherText = ZeroUnPadding(cipherText)
	case "PKCS7":
		cipherText = PKCS7UnPadding(cipherText)
	case "ANSIX923":
		cipherText = ANSIX923UnPadding(cipherText)
	case "ISO10126":
		cipherText = ISO10126UnPadding(cipherText)
	}

	return cipherText
}

func generateAESIVOrNonce(cryptoMode string) ([]byte, error) {
	if "GCM" == cryptoMode {
		nonce := make([]byte, LengthOfAes256GcmNonce)
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, errors.Wrap(constant.SystemInternalError, err, 0)
		}

		return nonce, nil
	}

	iv := make([]byte, LengthOfAesCommonIv)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	return iv, nil
}

func generateDESIV() ([]byte, error) {
	iv := make([]byte, LengthOfDesCommonIv)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	return iv, nil
}

func generate3DESIV() ([]byte, error) {
	iv := make([]byte, LengthOf3desCommonIv)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	return iv, nil
}

func generateBLOWFISHIV() ([]byte, error) {
	iv := make([]byte, LengthOfBlowfishCommonIv)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.Wrap(constant.SystemInternalError, err, 0)
	}

	return iv, nil
}

// GenerateIVOrNonce creates IVOrNonce according to algorithm and filling type
func GenerateIVOrNonce(algo string, cryptoMode string) ([]byte, error) {
	switch algo {
	case "AES":
		{
			return generateAESIVOrNonce(cryptoMode)
		}
	case "DES":
		{
			return generateDESIV()
		}
	case "3DES":
		{
			return generate3DESIV()
		}
	case "BLOWFISH":
		{
			return generateBLOWFISHIV()
		}
	}
	return nil, nil
}
