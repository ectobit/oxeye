package encdec_test

import (
	"bytes"
	"reflect"
	"testing"

	"go.ectobit.com/oxeye/encdec"
)

var _ encdec.EncDecoder = (*encdec.JSONIter)(nil)

func TestJSONIterEncode(t *testing.T) {
	t.Parallel()

	ed := &encdec.JSONIter{}

	got, err := ed.Encode(newMsg())
	if err != nil {
		t.Error(err)
	}

	if want := jsonData(); !bytes.Equal(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestJSONIterDecode(t *testing.T) {
	t.Parallel()

	ed := &encdec.JSONIter{}

	got := &msg{} //nolint:exhaustivestruct

	if err := ed.Decode(jsonData(), got); err != nil {
		t.Error(err)
	}

	if want := newMsg(); !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func BenchmarkJSONIterEncode(b *testing.B) {
	ed := &encdec.JSONIter{}
	msg := newMsg()

	for n := 0; n < b.N; n++ {
		if _, err := ed.Encode(msg); err != nil {
			panic(err)
		}
	}
}

func BenchmarkJSONIterDecode(b *testing.B) {
	ed := &encdec.JSONIter{}
	data := jsonData()
	msg := newMsg()

	for n := 0; n < b.N; n++ {
		if err := ed.Decode(data, msg); err != nil {
			panic(err)
		}
	}
}
