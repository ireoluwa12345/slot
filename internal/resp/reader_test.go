package resp

import (
	"bytes"
	"testing"
)

func TestReader_ReadBulk(t *testing.T) {
	data := "$5\r\nhello\r\n"
	r := NewReader(bytes.NewBufferString(data))

	v, err := r.Read()
	if err != nil {
		t.Errorf("Read() error = %v", err)
	}

	if v.Typ != "bulk" {
		t.Errorf("Read().Typ = %q, want %q", v.Typ, "bulk")
	}
	if v.Bulk != "hello" {
		t.Errorf("Read().Bulk = %q, want %q", v.Bulk, "hello")
	}
}

func TestReader_ReadBulkEmpty(t *testing.T) {
	data := "$0\r\n\r\n"
	r := NewReader(bytes.NewBufferString(data))

	v, err := r.Read()
	if err != nil {
		t.Errorf("Read() error = %v", err)
	}

	if v.Typ != "bulk" {
		t.Errorf("Read().Typ = %q, want %q", v.Typ, "bulk")
	}
	if v.Bulk != "" {
		t.Errorf("Read().Bulk = %q, want %q", v.Bulk, "")
	}
}

func TestReader_ReadArray(t *testing.T) {
	data := "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	r := NewReader(bytes.NewBufferString(data))

	v, err := r.Read()
	if err != nil {
		t.Errorf("Read() error = %v", err)
	}

	if v.Typ != "array" {
		t.Errorf("Read().Typ = %q, want %q", v.Typ, "array")
	}
	if len(v.Array) != 2 {
		t.Errorf("Read().Array len = %d, want %d", len(v.Array), 2)
	}
	if v.Array[0].Bulk != "foo" {
		t.Errorf("Read().Array[0].Bulk = %q, want %q", v.Array[0].Bulk, "foo")
	}
	if v.Array[1].Bulk != "bar" {
		t.Errorf("Read().Array[1].Bulk = %q, want %q", v.Array[1].Bulk, "bar")
	}
}

func TestReader_ReadArrayEmpty(t *testing.T) {
	data := "*0\r\n"
	r := NewReader(bytes.NewBufferString(data))

	v, err := r.Read()
	if err != nil {
		t.Errorf("Read() error = %v", err)
	}

	if v.Typ != "array" {
		t.Errorf("Read().Typ = %q, want %q", v.Typ, "array")
	}
	if len(v.Array) != 0 {
		t.Errorf("Read().Array len = %d, want %d", len(v.Array), 0)
	}
}

func TestReader_ReadNestedArray(t *testing.T) {
	data := "*2\r\n*2\r\n$3\r\none\r\n$3\r\ntwo\r\n$3\r\nthr\r\n"
	r := NewReader(bytes.NewBufferString(data))

	v, err := r.Read()
	if err != nil {
		t.Errorf("Read() error = %v", err)
	}

	if v.Typ != "array" {
		t.Errorf("Read().Typ = %q, want %q", v.Typ, "array")
	}
	if len(v.Array) != 2 {
		t.Errorf("Read().Array len = %d, want %d", len(v.Array), 2)
	}
	if v.Array[0].Typ != "array" {
		t.Errorf("Read().Array[0].Typ = %q, want %q", v.Array[0].Typ, "array")
	}
	if v.Array[0].Array[0].Bulk != "one" {
		t.Errorf("Read().Array[0].Array[0].Bulk = %q, want %q", v.Array[0].Array[0].Bulk, "one")
	}
}

func TestReader_ReadInvalidType(t *testing.T) {
	data := "X\r\n"
	r := NewReader(bytes.NewBufferString(data))

	_, err := r.Read()
	if err != nil {
		t.Errorf("Read() error = %v", err)
	}
}
