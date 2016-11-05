package reference

import (
	"bytes"
	"testing"
)

func TestReadColon(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString(":"))
	token, data := scanner.scan()

	if token != colon || data != ":" {
		t.Error("Token must be a colon got", data)
	}
}

func TestReadTilde(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString("~"))
	token, data := scanner.scan()

	if token != tilde || data != "~" {
		t.Error("Token must be a tilde got", data)
	}
}

func TestReadCaret(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString("^"))
	token, data := scanner.scan()

	if token != caret || data != "^" {
		t.Error("Token must be a caret got", data)
	}
}

func TestReadDot(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString("."))
	token, data := scanner.scan()

	if token != dot || data != "." {
		t.Error("Token must be a dot got", data)
	}
}

func TestReadSlash(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString("/"))
	token, data := scanner.scan()

	if token != slash || data != "/" {
		t.Error("Token must be a slash got", data)
	}
}

func TestReadEOF(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString(string(rune(0))))
	token, data := scanner.scan()

	if token != eof || data != "" {
		t.Error("Token must mark end of file got", data)
	}
}

func TestReadNumber(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString("1234"))
	token, data := scanner.scan()

	if token != number || data != "1234" {
		t.Error("Token must mark end of file got", data)
	}
}
