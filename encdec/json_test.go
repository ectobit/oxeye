package encdec_test

import (
	"bytes"
	"reflect"
	"testing"

	"go.ectobit.com/oxeye/encdec"
)

var _ encdec.EncDecoder = (*encdec.JSON)(nil)

func TestJSONEncode(t *testing.T) {
	t.Parallel()

	ed := &encdec.JSON{}

	got, err := ed.Encode(newMsg())
	if err != nil {
		t.Error(err)
	}

	if want := jsonData(); !bytes.Equal(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestJSONDecode(t *testing.T) {
	t.Parallel()

	ed := &encdec.JSON{}

	got := &msg{} //nolint:exhaustruct

	if err := ed.Decode(jsonData(), got); err != nil {
		t.Error(err)
	}

	if want := newMsg(); !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func BenchmarkJSONEncode(b *testing.B) {
	ed := &encdec.JSON{}
	msg := newMsg()

	for n := 0; n < b.N; n++ {
		if _, err := ed.Encode(msg); err != nil {
			panic(err)
		}
	}
}

func BenchmarkJSONDecode(b *testing.B) {
	ed := &encdec.JSON{}
	data := jsonData()
	msg := newMsg()

	for n := 0; n < b.N; n++ {
		if err := ed.Decode(data, msg); err != nil {
			panic(err)
		}
	}
}

func jsonData() []byte {
	return []byte(`{"name":"John Doe"}`)
}
