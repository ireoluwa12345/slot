package resp

import (
	"bytes"
	"testing"
)

func TestValue_MarshalString(t *testing.T) {
	v := Value{
		Typ: "string",
		Str: "OK",
	}
	got := v.Marshal()
	want := []byte("+OK\r\n")
	if !bytes.Equal(got, want) {
		t.Errorf("MarshalString() = %q, want %q", got, want)
	}
}

func TestValue_MarshalBulk(t *testing.T) {
	v := Value{
		Typ:  "bulk",
		Bulk: "hello",
	}
	got := v.Marshal()
	want := []byte("$5\r\nhello\r\n")
	if !bytes.Equal(got, want) {
		t.Errorf("MarshalBulk() = %q, want %q", got, want)
	}
}

func TestValue_MarshalBulkEmpty(t *testing.T) {
	v := Value{
		Typ:  "bulk",
		Bulk: "",
	}
	got := v.Marshal()
	want := []byte("$0\r\n\r\n")
	if !bytes.Equal(got, want) {
		t.Errorf("MarshalBulk() empty = %q, want %q", got, want)
	}
}

func TestValue_MarshalArray(t *testing.T) {
	v := Value{
		Typ: "array",
		Array: []Value{
			{Typ: "bulk", Bulk: "foo"},
			{Typ: "bulk", Bulk: "bar"},
		},
	}
	got := v.Marshal()
	want := []byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")
	if !bytes.Equal(got, want) {
		t.Errorf("MarshalArray() = %q, want %q", got, want)
	}
}

func TestValue_MarshalArrayEmpty(t *testing.T) {
	v := Value{
		Typ:   "array",
		Array: []Value{},
	}
	got := v.Marshal()
	want := []byte("*0\r\n")
	if !bytes.Equal(got, want) {
		t.Errorf("MarshalArray() empty = %q, want %q", got, want)
	}
}

func TestValue_MarshalError(t *testing.T) {
	v := Value{
		Typ: "error",
		Str: "ERR something went wrong",
	}
	got := v.Marshal()
	want := []byte("-ERR something went wrong\r\n")
	if !bytes.Equal(got, want) {
		t.Errorf("MarshalError() = %q, want %q", got, want)
	}
}

func TestValue_MarshalNull(t *testing.T) {
	v := Value{
		Typ: "null",
	}
	got := v.Marshal()
	want := []byte("$-1\r\n")
	if !bytes.Equal(got, want) {
		t.Errorf("MarshalNull() = %q, want %q", got, want)
	}
}

func TestValue_MarshalUnknown(t *testing.T) {
	v := Value{
		Typ: "unknown",
	}
	got := v.Marshal()
	want := []byte{}
	if !bytes.Equal(got, want) {
		t.Errorf("MarshalUnknown() = %q, want %q", got, want)
	}
}
