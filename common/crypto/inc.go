package crypto

// Interface define the interface of crypto
type Interface interface {
	// Encrypt encrypts plain text into cipher text
	Encrypt(plainText []byte) ([]byte, error)

	// Decrypt decrypts cipher text into plain text
	Decrypt(cipherText []byte) ([]byte, error)
}

// SignInterface define the interface of signature
type SignInterface interface {
	Sign(data []byte) ([]byte, error)
	VerifySign(data, signData []byte) (bool, error)
}
