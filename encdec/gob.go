package encdec

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Gob implements EncDecoder interface using encoding/gob.
type Gob struct {
	// mu  sync.Mutex
	b   bytes.Buffer
	enc *gob.Encoder
}

// NewGob create new gob encoder/decoder.
func NewGob() *Gob {
	ed := &Gob{} //nolint:exhaustivestruct

	ed.enc = gob.NewEncoder(&ed.b)

	return ed
}

// Encode encodes provided pointer to a value to slice of bytes using gob encoding.
func (ed *Gob) Encode(v interface{}) ([]byte, error) {
	ed.b.Reset()

	if err := ed.enc.Encode(v); err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	return ed.b.Bytes(), nil
}

// Decode decodes provided bytes slice into provided pointer to a value using gob decoding.
func (ed *Gob) Decode(data []byte, v interface{}) error {
	dec := gob.NewDecoder(bytes.NewBuffer(data))

	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}
