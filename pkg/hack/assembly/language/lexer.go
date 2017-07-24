package language

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// Token is a lexical token
type Token int

const (
	ILLEGAL Token = iota
	EOF
	WS
	COMMENT

	VALUE

	NULL

	AT
	EQUALS

	LABEL_START
	LABEL_END

	SEMICOLON
	JGT
	JEQ
	JGE
	JLT
	JNE
	JLE
	JMP
)

func (t Token) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case WS:
		return "WS"
	case COMMENT:
		return "COMMENT"
	case VALUE:
		return "VALUE"
	case NULL:
		return "NULL"
	case AT:
		return "AT"
	case EQUALS:
		return "EQUALS"
	case LABEL_START:
		return "LABEL_START"
	case LABEL_END:
		return "LABEL_END"
	case SEMICOLON:
		return "SEMICOLON"
	case JGT:
		return "JGT"
	case JEQ:
		return "JEQ"
	case JGE:
		return "JGE"
	case JLT:
		return "JLT"
	case JNE:
		return "JNE"
	case JLE:
		return "JLE"
	case JMP:
		return "JMP"
	default:
		return ""
	}
}

func mapIdent(str string) Token {
	switch str {
	case "JGT":
		return JGT
	case "JEQ":
		return JEQ
	case "JGE":
		return JGE
	case "JLT":
		return JLT
	case "JNE":
		return JNE
	case "JLE":
		return JLE
	case "JMP":
		return JMP
	}

	return VALUE
}

var eof = rune(0)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// returns true if the character belongs to the class of allowed identifier characters
func isIdent(ch rune) bool {
	if isLetter(ch) {
		return true
	}
	if isDigit(ch) {
		return true
	}
	if ch == '.' || ch == '-' || ch == '_' {
		return true
	}
	return false
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

type Scanner struct {
	r *bufio.Reader

	i int
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r: bufio.NewReader(r),

		i: -1,
	}
}

func (s *Scanner) read() (rune, error) {
	s.i++
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
	s.i--
	return s.r.UnreadRune()
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
		} else if !isIdent(ch) {
			err = s.unread()
			if err != nil {
				return ILLEGAL, "", err
			}
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return mapIdent(buf.String()), buf.String(), nil
}

func (s *Scanner) scanDigit() (tok Token, lit string, err error) {
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
		} else if !isDigit(ch) {
			err = s.unread()
			if err != nil {
				return ILLEGAL, "", err
			}
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return VALUE, buf.String(), nil
}

func (s *Scanner) scanComment() (tok Token, lit string, err error) {
	var buf bytes.Buffer
	ch, err := s.read()
	if err != nil {
		return ILLEGAL, "", err
	}
	buf.WriteRune(ch)

	for {
		if ch, err := s.read(); err != nil {
			return ILLEGAL, "", err
		} else if ch == eof || ch == '\n' || ch == '\r' {
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return COMMENT, buf.String(), nil
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
	} else if isIdent(ch) {
		err := s.unread()
		if err != nil {
			return ILLEGAL, "", err
		}
		return s.scanIdent()
	}

	switch ch {
	case eof:
		return EOF, "", nil
	case '/':
		next, err := s.read()
		if err != nil {
			return ILLEGAL, "", err
		}
		if next != '/' {
			return ILLEGAL, "", fmt.Errorf("invalid single / on index %d", s.i)
		}
		_, _, err = s.scanWhitespace()
		if err != nil {
			return ILLEGAL, "", err
		}
		return s.scanComment()
	case '@':
		return AT, "@", nil
	case '=':
		return EQUALS, "=", nil
	case ';':
		return SEMICOLON, ";", nil
	case '(':
		return LABEL_START, "(", nil
	case ')':
		return LABEL_END, ")", nil
	}
	return ILLEGAL, string(ch), nil
}
