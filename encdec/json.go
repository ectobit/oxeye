package encdec

import (
	"encoding/json"
	"fmt"
)

// JSON implements EncDecoder interface using encoding/json.
type JSON struct{}

// Encode encodes provided pointer to a value to slice of bytes using JSON encoding.
func (ed *JSON) Encode(v interface{}) ([]byte, error) {
	out, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	return out, nil
}

// Decode decodes provided bytes slice into provided pointer to a value using JSON decoding.
func (ed *JSON) Decode(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}
