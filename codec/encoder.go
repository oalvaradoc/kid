package codec

// Encoder is an interface that defines the encoder's method signature
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
}
