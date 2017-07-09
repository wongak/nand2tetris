package language

import "io"

// Parser is a hack VM language parser
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token
		lit string
		n   int
	}
}

// NewParser creates a new parser on the given Reader r
func NewParser(r io.Reader) *Parser {
	return &Parser{
		s: NewScanner(r),
	}
}
