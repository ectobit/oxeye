// Package encdec contains common encoder/decoder interface and several implementation of it.
package encdec

// EncDecoder defines common methods to encode and decode values into and from bytes slice.
type EncDecoder interface {
	// Encode encodes provided pointer to a value to slice of bytes.
	Encode(v interface{}) ([]byte, error)
	// Decode decodes provided bytes slice into provided pointer to a value.
	Decode(b []byte, v interface{}) error
}
