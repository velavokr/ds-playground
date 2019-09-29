package utils

import (
	"errors"
	"testing"
)

func TestWriteAll(t *testing.T) {
	str := "abcdefghij"
	for i := 0; i <= len(str); i++ {
		w := &writer{res: []byte{}}
		err := WriteAll(w, []byte(str)[:i])
		if err != nil {
			t.Fatal(err.Error())
		}
		if string(w.res) != str[:i] {
			t.Fatal()
		}
	}
}

type writer struct {
	res []byte
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.res = append(w.res, p[0])
	return 1, nil
}

func TestQuote(t *testing.T) {
	tests := []struct {
		name string
		args []byte
		want string
	}{
		{"nil", nil, "{nil}"},
		{"empty", []byte{}, `{slice=[],hex=,ascii=""}`},
		{"a", []byte("a"), `{slice=[97],hex=61,ascii="a"}`},
		{"CRLF0", []byte("\r\n\x00"), `{slice=[13 10 0],hex=0d0a00,ascii="\r\n\x00"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Quote(tt.args); got != tt.want {
				t.Errorf("Quote(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestSprint(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want string
	}{
		{"empty", []interface{}{}, ""},
		{"str", []interface{}{"a"}, "a"},
		{"str2", []interface{}{"a", "b"}, "a b"},
		{"err", []interface{}{errors.New("hello")}, "hello"},
		{"nil", []interface{}{nil}, "<nil>"},
		{"bytes", []interface{}{[]byte("hello")},
			`{slice=[104 101 108 108 111],hex=68656c6c6f,ascii="hello"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sprint(tt.args...); got != tt.want {
				t.Errorf("Sprint(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
