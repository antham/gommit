package reference

// token represents a entity extracted from string parsing
type token int

const (
	illegal token = iota
	eof

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
