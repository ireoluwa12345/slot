package resp

import (
	"bytes"
	"testing"
)

func TestWriter_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewWriter(buf)

	v := Value{Typ: "string", Str: "PONG"}
	err := w.Write(v)
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}

	got := buf.String()
	want := "+PONG\r\n"
	if got != want {
		t.Errorf("Write() = %q, want %q", got, want)
	}
}

func TestWriter_WriteBulk(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewWriter(buf)

	v := Value{Typ: "bulk", Bulk: "test data"}
	err := w.Write(v)
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}

	got := buf.String()
	want := "$9\r\ntest data\r\n"
	if got != want {
		t.Errorf("Write() = %q, want %q", got, want)
	}
}

func TestWriter_WriteArray(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewWriter(buf)

	v := Value{
		Typ: "array",
		Array: []Value{
			{Typ: "string", Str: "GET"},
			{Typ: "bulk", Bulk: "key"},
		},
	}
	err := w.Write(v)
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}

	got := buf.String()
	want := "*2\r\n+GET\r\n$3\r\nkey\r\n"
	if got != want {
		t.Errorf("Write() = %q, want %q", got, want)
	}
}
