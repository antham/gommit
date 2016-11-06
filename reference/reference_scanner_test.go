package reference

import (
	"bytes"
	"testing"
)

func TestReadColon(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString(":"))
	tok, data := scanner.scan()

	if tok != colon || data != ":" {
		t.Error("Token must be a colon got", data)
	}
}

func TestReadTilde(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString("~"))
	tok, data := scanner.scan()

	if tok != tilde || data != "~" {
		t.Error("Token must be a tilde got", data)
	}
}

func TestReadCaret(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString("^"))
	tok, data := scanner.scan()

	if tok != caret || data != "^" {
		t.Error("Token must be a caret got", data)
	}
}

func TestReadDot(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString("."))
	tok, data := scanner.scan()

	if tok != dot || data != "." {
		t.Error("Token must be a dot got", data)
	}
}

func TestReadSlash(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString("/"))
	tok, data := scanner.scan()

	if tok != slash || data != "/" {
		t.Error("Token must be a slash got", data)
	}
}

func TestReadEOF(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString(string(rune(0))))
	tok, data := scanner.scan()

	if tok != eof || data != "" {
		t.Error("Token must mark end of file got", data)
	}
}

func TestReadNumber(t *testing.T) {
	scanner := newScanner(bytes.NewBufferString("1234"))
	tok, data := scanner.scan()

	if tok != number || data != "1234" {
		t.Error("Token must mark end of file got", data)
	}
}
