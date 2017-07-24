package language

import (
	"fmt"
	"io"
)

type Command interface {
	fmt.Stringer
	Translate(*SymbolTable, io.Writer) error
}

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

type ParserContext struct {
}

type stateFunc func(p *Parser, ctx ParserContext) (ParserContext, stateFunc)

func NewParser(r io.Reader) *Parser {
	return &Parser{
		s:    NewScanner(r),
		i:    -1,
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

	return tok, lit, nil
}

func (p *Parser) unscan() {
	p.i--
	p.buf.isUnscanned = true
}

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

func (p *Parser) Tree() []Command {
	return p.tree
}

func (p *Parser) Run() error {
	ctx := ParserContext{}

	for state := top; state != nil; {
		ctx, state = state(p, ctx)
	}
	if p.err != nil {
		return fmt.Errorf("parse error on token %d (index %d): %v", p.i, p.s.i, p.err)
	}
	return nil
}

func parseError(err error) stateFunc {
	return func(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
		p.err = err
		return ctx, nil
	}
}

func top(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
	tok, lit, err := p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	switch true {
	case tok == EOF:
		return ctx, nil
	case tok == AT:
		p.unscan()
		return ctx, parseAInstruction
	}

	return ctx, parseError(fmt.Errorf("invalid token %s (%s)", tok, lit))
}

func command(cmd Command) stateFunc {
	return func(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
		p.tree = append(p.tree, cmd)
		return ctx, top
	}
}
