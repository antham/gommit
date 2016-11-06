package reference

// token represents a entity extracted from string parsing
type token int

const (
	eof token = iota

	colon
	tilde
	caret
	dot
	slash
	control
	space
	number
	char
)
