package encdec

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

// JSONIter implements EncDecoder interface using github.com/json-iterator/go.
type JSONIter struct{}

// Encode encodes provided pointer to a value to slice of bytes using JSON encoding.
func (ed *JSONIter) Encode(v interface{}) ([]byte, error) {
	out, err := jsoniter.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	return out, nil
}

// Decode decodes provided bytes slice into provided pointer to a value using JSON decoding.
func (ed *JSONIter) Decode(data []byte, v interface{}) error {
	if err := jsoniter.Unmarshal(data, v); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}
