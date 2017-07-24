package language

import (
	"fmt"
	"io"
	"strconv"
)

type AInstruction struct {
	Address string
}

func (a *AInstruction) String() string {
	return fmt.Sprintf("@%s", a.Address)
}

func (a *AInstruction) Translate(t *SymbolTable, wr io.Writer) error {
	var address int
	if ai, err := strconv.ParseInt(a.Address, 10, 64); err == nil {
		address = int(ai)
	} else {
		address = t.Label(a.Address)
	}
	_, err := wr.Write([]byte("0"))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(wr, "%015b\n", address)
	if err != nil {
		return err
	}

	t.RegisterInstruction()
	return nil
}

func parseAInstruction(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
	tok, lit, err := p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != AT {
		panic("internal error")
	}

	tok, lit, err = p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != VALUE {
		return ctx, parseError(fmt.Errorf("invalid token %s (%s) for A-Instruction. Expect VALUE.", tok, lit))
	}
	a := &AInstruction{
		Address: lit,
	}
	return ctx, command(a)
}
