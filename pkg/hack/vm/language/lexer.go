package language

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// Token is a lexical token
type Token int

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

	case PUSH:
		return "PUSH"
	case POP:
		return "POP"

	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case NEG:
		return "NEG"
	case EQ:
		return "EQ"
	case GT:
		return "GT"
	case LT:
		return "LT"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"

	case CONSTANT:
		return "CONSTANT"
	case STATIC:
		return "STATIC"
	case LCL:
		return "LCL"
	case ARG:
		return "ARG"
	case THIS:
		return "THIS"
	case THAT:
		return "THAT"
	case TEMP:
		return "TEMP"
	case POINTER:
		return "POINTER"

	case LABEL:
		return "LABEL"
	case IFGOTO:
		return "IFGOTO"
	case GOTO:
		return "GOTO"

	case FUNCTION:
		return "FUNCTION"
	case CALL:
		return "CALL"
	case RETURN:
		return "RETURN"

	default:
		return "unknown token"
	}
}

const (
	ILLEGAL Token = iota
	EOF
	WS
	COMMENT

	VALUE

	// commands

	// memory access
	PUSH
	POP

	// arithmetic/logical commands
	ADD
	SUB
	NEG
	EQ
	GT
	LT
	AND
	OR
	NOT

	// memory segments
	CONSTANT
	STATIC
	LCL
	ARG
	THIS
	THAT
	TEMP
	POINTER

	// branching
	LABEL
	IFGOTO
	GOTO

	// functions
	FUNCTION
	CALL
	RETURN
)

var eof = rune(0)

func isMemoryAccessCommand(tok Token) bool {
	return tok == PUSH || tok == POP
}

func isArithmeticCommand(tok Token) bool {
	return tok == ADD ||
		tok == SUB ||
		tok == NEG ||
		tok == EQ ||
		tok == GT ||
		tok == LT ||
		tok == AND ||
		tok == OR ||
		tok == NOT
}

func isSegment(tok Token) bool {
	return tok == STATIC ||
		tok == LCL ||
		tok == ARG ||
		tok == THIS ||
		tok == THAT ||
		tok == TEMP
}

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

// Scanner can scan tokens
type Scanner struct {
	r *bufio.Reader

	i int
}

// NewScanner creates a new scanner, which reads from the given Reader r
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
	} else if isDigit(ch) {
		err := s.unread()
		if err != nil {
			return ILLEGAL, "", err
		}
		return s.scanDigit()
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
			return ILLEGAL, "", fmt.Errorf("invalid comment starting character / on index %d", s.i)
		}
		_, _, err = s.scanWhitespace()
		if err != nil {
			return ILLEGAL, "", err
		}
		return s.scanComment()
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

func mapIdent(str string) Token {
	switch str {
	case "push":
		return PUSH
	case "pop":
		return POP

	case "add":
		return ADD
	case "sub":
		return SUB
	case "neg":
		return NEG
	case "eq":
		return EQ
	case "gt":
		return GT
	case "lt":
		return LT
	case "and":
		return AND
	case "or":
		return OR
	case "not":
		return NOT

	case "constant":
		return CONSTANT
	case "static":
		return STATIC
	case "local":
		return LCL
	case "argument":
		return ARG
	case "this":
		return THIS
	case "that":
		return THAT
	case "temp":
		return TEMP
	case "pointer":
		return POINTER

	case "label":
		return LABEL
	case "goto":
		return GOTO
	case "if-goto":
		return IFGOTO

	case "function":
		return FUNCTION
	case "call":
		return CALL
	case "return":
		return RETURN
	}

	return VALUE
}
