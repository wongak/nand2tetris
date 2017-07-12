package language

import (
	"fmt"
	"io"
	"sync/atomic"
)

var compIndex int64 = -1

func nextIndex() int64 {
	n := atomic.AddInt64(&compIndex, 1)
	return n
}

func createCompIndex(cfg TranslateConfig) string {
	return fmt.Sprintf("comp.%s.%d", cfg.FileName, nextIndex())
}

// Arithmetic represents an arithmetic/logical command
type Arithmetic struct {
	cmd Token
	lit string
}

// String implements the Stringer
func (a *Arithmetic) String() string {
	return fmt.Sprintf("%s", a.cmd)
}

// Translate implementing the Command
func (m *Arithmetic) Translate(cfg TranslateConfig, wr io.Writer) {
	data := map[string]string{
		"cmdLit": m.lit,
	}
	tmpl := arithmeticOpTmpl
	switch m.cmd {
	case ADD:
		data["operation"] = "D=D+M"
	case SUB:
		data["operation"] = "D=M-D"
	case NEG:
		tmpl = arithmeticSingleOpTmpl
		data["operation"] = "D=-D"
	case AND:
		data["operation"] = "D=D&M"
	case OR:
		data["operation"] = "D=D|M"
	case NOT:
		data["operation"] = "D=!D"
		tmpl = arithmeticSingleOpTmpl
	case EQ:
		data["labelSet"] = createCompIndex(cfg)
		tmpl = logicalCompTmpl
		data["comp"] = "JEQ" // true if pop1 - pop2 = 0
	case GT:
		data["labelSet"] = createCompIndex(cfg)
		tmpl = logicalCompTmpl
		data["comp"] = "JLT" // true if pop1 - pop2 < 0
	case LT:
		data["labelSet"] = createCompIndex(cfg)
		tmpl = logicalCompTmpl
		data["comp"] = "JGT"
	}
	err := tmpl.Execute(wr, data)
	if err != nil {
		panic(err)
	}
}

func parseArithmetic(p *Parser) stateFunc {
	tok, lit, err := p.scanIgnore()
	if err != nil {
		return parseError(err)
	}
	if !isArithmeticCommand(tok) {
		return parseError(fmt.Errorf("invalid token %s (%s). epxect arithmetic/logical cmd", tok, lit))
	}
	cmd := &Arithmetic{
		cmd: tok,
		lit: lit,
	}
	return command(cmd)
}
