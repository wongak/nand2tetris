package language

import (
	"bufio"
	"bytes"
	"io"
)

// Token is a lexical token
type Token int

const (
	ILLEGAL Token = iota
	EOF
	WS

	VALUE

	PUSH
	POP

	ADD
	SUB
	NEG
	EQ
	GT
	LT
	AND
	OR
	NOT

	CONSTANT
	STATIC
	LOCAL
	ARGUMENT
	THIS
	THAT
	TEMP
	POINTER
)

var eof = rune(0)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// Scanner can scan tokens
type Scanner struct {
	r *bufio.Reader
}

// NewScanner creates a new scanner, which reads from the given Reader r
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r: bufio.NewReader(r),
	}
}

func (s *Scanner) read() (rune, error) {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		if err == io.EOF {
			return eof, nil
		}
		return eof, err
	}
	return ch, nil
}

func (s *Scanner) unread() error {
	return s.r.UnreadRune()
}

func (s *Scanner) Scan() (tok Token, lit string, err error) {
	ch, err := s.read()
	if err != nil {
		return ILLEGAL, "", err
	}

	if isWhitespace(ch) {
		err := s.unread()
		if err != nil {
			return ILLEGAL, "", err
		}
		return s.scanWhitespace()
	} else if isLetter(ch) {
		err := s.unread()
		if err != nil {
			return ILLEGAL, "", err
		}
		return s.scanIdent()
	}

	switch ch {
	case eof:
		return EOF, "", nil
	}

	return ILLEGAL, string(ch), nil
}

func (s *Scanner) scanWhitespace() (tok Token, lit string, err error) {
	var buf bytes.Buffer
	ch, err := s.read()
	if err != nil {
		return ILLEGAL, "", err
	}
	buf.WriteRune(ch)

	for {
		if ch, err := s.read(); err != nil {
			return ILLEGAL, "", err
		} else if ch == eof {
			break
		} else if !isWhitespace(ch) {
			err = s.unread()
			if err != nil {
				return ILLEGAL, "", err
			}
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String(), nil
}

func (s *Scanner) scanIdent() (tok Token, lit string, err error) {
	var buf bytes.Buffer
	ch, err := s.read()
	if err != nil {
		return ILLEGAL, "", err
	}
	buf.WriteRune(ch)

	for {
		if ch, err := s.read(); err != nil {
			return ILLEGAL, "", err
		} else if ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) {
			err = s.unread()
			if err != nil {
				return ILLEGAL, "", err
			}
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	switch buf.String() {
	case "push":
		return PUSH, buf.String(), nil
	}

	return VALUE, buf.String(), nil
}
