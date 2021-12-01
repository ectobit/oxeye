package encdec_test

import (
	"bytes"
	"reflect"
	"testing"

	"go.ectobit.com/oxeye/encdec"
)

var _ encdec.EncDecoder = (*encdec.Gob)(nil)

func TestGobEncode(t *testing.T) {
	t.Parallel()

	ed := encdec.NewGob()

	got, err := ed.Encode(newMsg())
	if err != nil {
		t.Error(err)
	}

	if want := gobData(); !bytes.Equal(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestGobDecode(t *testing.T) {
	t.Parallel()

	encdec := encdec.NewGob()

	got := &msg{} //nolint:exhaustivestruct

	if err := encdec.Decode(gobData(), got); err != nil {
		t.Error(err)
	}

	if want := newMsg(); !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}

	if err := encdec.Decode(gobData(), got); err != nil {
		t.Error(err)
	}
}

func BenchmarkGobEncode(b *testing.B) {
	ed := encdec.NewGob()
	msg := newMsg()

	for n := 0; n < b.N; n++ {
		if _, err := ed.Encode(msg); err != nil {
			panic(err)
		}
	}
}

func BenchmarkGobDecode(b *testing.B) {
	ed := encdec.NewGob()
	data := gobData()
	msg := newMsg()

	for n := 0; n < b.N; n++ {
		if err := ed.Decode(data, msg); err != nil {
			panic(err)
		}
	}
}

func gobData() []byte {
	return []byte{26, 255, 129, 3, 1, 1, 3, 109, 115, 103, 1, 255, 130, 0, 1, 1, 1, 4, 78, 97, 109, 101, 1, 12, 0, 0, 0, 13, 255, 130, 1, 8, 74, 111, 104, 110, 32, 68, 111, 101, 0} //nolint:lll
}
