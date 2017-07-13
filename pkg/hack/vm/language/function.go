package language

import (
	"fmt"
	"io"
	"strconv"
)

// Function implements the VM function command
type Function struct {
	name        string
	lit         string
	numLocal    int
	numLocalLit string
}

// String implements the Stringer
func (f *Function) String() string {
	return fmt.Sprintf("%s %s %s", FUNCTION, f.name, f.numLocalLit)
}

func (f *Function) Translate(t *SymbolTable, wr io.Writer) {
}

func parseFunction(p *Parser) stateFunc {
	tok, lit, err := p.scanIgnore()
	if err != nil {
		return parseError(err)
	}
	if tok != FUNCTION {
		panic("internal error")
	}
	f := &Function{
		lit: lit,
	}

	tok, lit, err = p.scanIgnore()
	if err != nil {
		return parseError(err)
	}
	if tok != VALUE {
		return parseError(fmt.Errorf("invalid token %s (%s), expect function name", tok, lit))
	}
	f.name = lit

	tok, lit, err = p.scanIgnore()
	if err != nil {
		return parseError(err)
	}
	if tok != VALUE {
		return parseError(fmt.Errorf("invalid token %s (%s). expect num local vars", tok, lit))
	}
	num, err := strconv.ParseInt(lit, 10, 64)
	if err != nil {
		return parseError(fmt.Errorf("invalid number of local vars %s: %v", lit, err))
	}
	f.numLocalLit = lit
	f.numLocal = int(num)

	return parseFunctionBody(f)
}

func parseFunctionBody(f *Function) stateFunc {
	return func(p *Parser) stateFunc {
		tok, lit, err := p.scanIgnore()
		if err != nil {
			return parseError(err)
		}
		switch true {
		case tok == EOF:
			return parseError(fmt.Errorf("unexpected EOF before function return"))
		case isMemoryAccessCommand(tok):
			p.unscan()
			return parseMemoryAccess(f)
		case isArithmeticCommand(tok):
			p.unscan()
			return parseArithmetic(f)
		case tok == LABEL:
			p.unscan()
			return parseLabel(f)
		case tok == RETURN:
			p.unscan()
			return parseReturn(f)
		}

		return parseError(fmt.Errorf("invalid token %s (%s) in function body", tok, lit))
	}
}

func parseReturn(f *Function) stateFunc {
	return func(p *Parser) stateFunc {
		_, _, err := p.scanIgnore()
		if err != nil {
			return parseError(err)
		}

		return top
	}
}
