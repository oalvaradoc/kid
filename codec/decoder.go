package codec

// Decoder is an interface that defines the decoder's method signature
type Decoder interface {
	Decode(data []byte, v interface{}) error
}
