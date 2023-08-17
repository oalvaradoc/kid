package codec

import (
	"git.multiverse.io/eventkit/kit/codec/json"
	"git.multiverse.io/eventkit/kit/codec/text"
	"git.multiverse.io/eventkit/kit/codec/xml"
)

// Codec is an interface that defines the create function of Encoder/Decoder
type Codec interface {
	Encoder() Encoder
	Decoder() Decoder
}

type impl struct {
	encoder func() Encoder
	decoder func() Decoder
}

func (c *impl) Encoder() Encoder { return c.encoder() }
func (c *impl) Decoder() Decoder { return c.decoder() }

// BuildCustomCodec creates a new codec with customize Encoder and Decoder
func BuildCustomCodec(encoder Encoder, decoder Decoder) Codec {
	return &impl{
		encoder: func() Encoder { return encoder },
		decoder: func() Decoder { return decoder },
	}
}

// BuildJSONCodec creates a new Codec that used to JSON marshal/unmarshal
func BuildJSONCodec() Codec {
	return &impl{
		encoder: func() Encoder { return &json.Encoder{} },
		decoder: func() Decoder { return &json.Decoder{} },
	}
}

// BuildXMLCodec creates a new Codec that used to XML marshal/unmarshal
func BuildXMLCodec() Codec {
	return &impl{
		encoder: func() Encoder { return &xml.Encoder{} },
		decoder: func() Decoder { return &xml.Decoder{} },
	}
}

// BuildTextCodec creates a new Codec that used to byte array marshal/unmarshal
func BuildTextCodec() Codec {
	return &impl{
		encoder: func() Encoder { return &text.Encoder{} },
		decoder: func() Decoder { return &text.Decoder{} },
	}
}
