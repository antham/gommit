package reference

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// symbolicRefPathStmt parse a following statement master~2
// to split branch name and refrence path
type symbolicRefPathStmt struct {
	branchName string
	refPath    []int
}

// parser represents a parser.
type parser struct {
	s   *scanner
	buf struct {
		tok token
		lit string
		n   int
	}
}

// newParser returns a new instance of parser.
func newParser(r io.Reader) *parser {
	return &parser{s: newScanner(r)}
}

// ParseSymbolicReferencePath creates a symbolicRefPathStmt
func (p *parser) parseSymbolicReferencePath() (*symbolicRefPathStmt, error) {
	branchName, err := p.parseBranchName()

	if err != nil {
		return nil, err
	}

	refPath, err := p.parseRefPath()

	if err != nil {
		return nil, err
	}

	return &symbolicRefPathStmt{
		branchName: branchName,
		refPath:    refPath,
	}, nil
}

// parseFuncArg parse one function argument
func (p *parser) parseBranchName() (string, error) {
	var tok token
	var previousTok token
	var lit string
	var buf string
	var subBuf string

	for {
		tok, lit = p.scan()

		err := p.checkBranchNameToken(tok)

		if err != nil {
			return "", err
		}

		if len(buf) > 4 {
			subBuf = strings.ToLower(buf[len(buf)-5:])
		}

		err = p.checkBranchNameConditions(tok, previousTok, buf, subBuf)

		if err != nil {
			return "", err
		}

		switch tok {
		case tilde, caret:
			p.unscan()
			return buf, nil
		case eof:
			return buf, nil
		}

		buf += lit
		previousTok = tok
	}
}

func (p *parser) checkBranchNameToken(tok token) error {
	switch tok {
	case colon:
		return fmt.Errorf("branch name must not contains any colon")
	case slash:
		return fmt.Errorf("branch name must not end with a slash")
	case space:
		return fmt.Errorf("branch name contains a space character")
	case control:
		return fmt.Errorf("branch name contains a control character")
	}

	return nil
}

func (p *parser) checkBranchNameConditions(currentToken token, previousToken token, currentBuf string, subBuf string) error {
	if currentToken == dot && currentBuf == "" {
		return fmt.Errorf("branch name must not start with a dot")
	}

	if currentToken == dot && previousToken == dot {
		return fmt.Errorf("branch name must not contains any double dot")
	}

	if (currentToken == tilde || currentToken == caret || currentToken == eof) && subBuf == ".lock" {
		return fmt.Errorf("branch name cannot end with .lock")
	}

	return nil
}

func (p *parser) parseRefPath() ([]int, error) {
	var tok token
	var buf = []int{}

	for {
		tok, _ = p.scan()

		if tok == caret {
			err := p.parseCaretLevel(&buf)

			if err != nil {
				return nil, err
			}
		}

		if tok == tilde {
			err := p.parseTildeLevel(&buf)

			if err != nil {
				return nil, err
			}
		}

		if tok == eof {
			return buf, nil
		}

		if tok != caret && tok != tilde {
			return buf, fmt.Errorf("must be a caret or a tilde and may be followed with a number")
		}
	}
}

func (p *parser) parseCaretLevel(buf *[]int) error {
	nTok, nLit := p.scan()

	if nTok != number {
		p.unscan()
		*buf = append(*buf, 1)
		return nil
	}

	level, err := strconv.Atoi(nLit)

	if err != nil {
		return err
	}

	if level != 2 && level != 1 {
		return fmt.Errorf("level associated to a caret must be 1 or 2")
	}

	*buf = append(*buf, level)

	return nil
}

func (p *parser) parseTildeLevel(buf *[]int) error {
	nTok, nLit := p.scan()

	if nTok != number {
		p.unscan()
		*buf = append(*buf, 1)
		return nil
	}

	level, err := strconv.Atoi(nLit)

	if err != nil {
		return err
	}

	for i := 0; i < level; i++ {
		*buf = append(*buf, 1)
	}

	return nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *parser) scan() (tok token, lit string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	tok, lit = p.s.scan()

	p.buf.tok, p.buf.lit = tok, lit
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *parser) unscan() { p.buf.n = 1 }
