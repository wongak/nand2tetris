package language

import (
	"fmt"
	"io"
)

// Parser is a hack VM language parser
type Parser struct {
	s   *Scanner
	buf struct {
		tok         Token
		lit         string
		isUnscanned bool
	}
	i int

	err  error
	tree []Command
}

// ParserContext gives all states a context
type ParserContext struct {
	file     *File
	function *Function
}

// parser state machine
type stateFunc func(p *Parser, ctx ParserContext) (ParserContext, stateFunc)

// NewParser creates a new parser on the given Reader r
func NewParser(r io.Reader) *Parser {
	return &Parser{
		s: NewScanner(r),
		i: -1,

		tree: make([]Command, 0),
	}
}

func (p *Parser) scan() (tok Token, lit string, err error) {
	p.i++
	if p.buf.isUnscanned {
		p.buf.isUnscanned = false
		return p.buf.tok, p.buf.lit, nil
	}

	tok, lit, err = p.s.Scan()
	if err != nil {
		return ILLEGAL, "", err
	}

	p.buf.tok = tok
	p.buf.lit = lit

	return
}

func (p *Parser) unscan() {
	p.i--
	p.buf.isUnscanned = true
}

// scanIgnore ignores all non-command tokens
func (p *Parser) scanIgnore() (tok Token, lit string, err error) {
	for {
		tok, lit, err = p.scan()
		if err != nil {
			return ILLEGAL, "", err
		}
		if tok == WS || tok == COMMENT {
			p.i--
			continue
		}
		return
	}
}

// Run starts the parser
func (p *Parser) Run(table *SymbolTable, fileName string) error {
	ctx := ParserContext{
		file: &File{name: fileName},
	}
	_, err := table.RegisterFile(fileName)
	if err != nil {
		return err
	}

	for state := top; state != nil; {
		ctx, state = state(p, ctx)
	}
	if p.err != nil {
		return fmt.Errorf("parse error on token %d (index %d): %v", p.i, p.s.i, p.err)
	}
	return nil
}

// Tree returns the normalized parse tree
func (p *Parser) Tree() []Command {
	return p.tree
}

func parseError(err error) stateFunc {
	return func(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
		p.err = err
		return ctx, nil
	}
}

// top is the top level parser state machine
func top(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
	tok, lit, err := p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	switch true {
	case tok == EOF:
		return ctx, nil
	case isMemoryAccessCommand(tok):
		p.unscan()
		return ctx, parseMemoryAccess
	case isArithmeticCommand(tok):
		p.unscan()
		return ctx, parseArithmetic
	case tok == LABEL:
		p.unscan()
		return ctx, parseLabel
	case tok == IFGOTO:
		p.unscan()
		return ctx, parseIfGoto
	case tok == GOTO:
		p.unscan()
		return ctx, parseGoto
	case tok == FUNCTION:
		p.unscan()
		return ctx, parseFunction
	case tok == CALL:
		p.unscan()
		return ctx, parseCall
	case tok == RETURN:
		p.unscan()
		return ctx, parseReturn
	}

	return ctx, parseError(fmt.Errorf("invalid token %s (%s)", tok, lit))
}

func command(cmd Command) stateFunc {
	return func(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
		p.tree = append(p.tree, cmd)
		return ctx, top
	}
}
