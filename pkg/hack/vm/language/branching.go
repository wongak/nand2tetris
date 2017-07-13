package language

import (
	"fmt"
	"io"
	"text/template"
)

const labelAsm = `// {{ .cmdLit }} {{ .name }}
({{ .symbol }})
`

var labelAsmTmpl *template.Template

func init() {
	labelAsmTmpl = template.Must(template.New("labelAsm").Parse(labelAsm))
}

// Label represents the label command
type Label struct {
	name string
	lit  string
}

// String implementing the Stringer
func (l *Label) String() string {
	return fmt.Sprintf("%s %s", LABEL, l.name)
}

// Translate generates assembly code for the label command
func (l *Label) Translate(t *SymbolTable, wr io.Writer) error {
	return nil
}

func parseLabel(f *Function) stateFunc {
	return func(p *Parser) stateFunc {
		tok, lit, err := p.scanIgnore()
		if err != nil {
			return parseError(err)
		}
		if tok != LABEL {
			panic("internal error")
		}
		if f == nil {
			return parseError(fmt.Errorf("invalid label without function context"))
		}
		l := &Label{
			lit: lit,
		}
		return parseLabelName(l, f)
	}
}

func parseLabelName(l *Label, f *Function) stateFunc {
	return func(p *Parser) stateFunc {
		tok, lit, err := p.scanIgnore()
		if err != nil {
			return parseError(err)
		}
		if tok != VALUE {
			return parseError(fmt.Errorf("expect label value, got token %s (%s)", tok, lit))
		}
		l.name = lit
		return command(l, f)
	}
}

// IfGoto implements the if-goto command
type IfGoto struct {
	lit   string
	label string
}

// String implements Stringer
func (g *IfGoto) String() string {
	return fmt.Sprintf("%s %s", IFGOTO, g.label)
}

// Translate generates assembly code for if-goto
func (g *IfGoto) Translate(t *SymbolTable, wr io.Writer) error {
	return nil
}

func parseIfGoto(f *Function) stateFunc {
	return func(p *Parser) stateFunc {
		if f == nil {
			return parseError(fmt.Errorf("invalid if-goto in non-function context"))
		}
		tok, lit, err := p.scanIgnore()
		if err != nil {
			return parseError(err)
		}
		if tok != IFGOTO {
			panic("internal error")
		}
		g := &IfGoto{
			lit: lit,
		}

		tok, lit, err = p.scanIgnore()
		if err != nil {
			return parseError(err)
		}
		if tok != VALUE {
			return parseError(fmt.Errorf("invalid token %s (%s) after if-goto. expect label"))
		}
		g.label = lit

		return command(g, f)
	}
}
