package language

import (
	"fmt"
	"io"
	"text/template"
)

const labelAsm = `// {{ .cmdLit }} {{ .name }}
({{ .label }})
`

var labelAsmTmpl *template.Template

func init() {
	labelAsmTmpl = template.Must(template.New("labelAsm").Parse(labelAsm))
}

func init() {
	labelAsmTmpl = template.Must(template.New("labelAsm").Parse(labelAsm))
}

// Label represents the label command
type Label struct {
	name string
	lit  string

	function *Function
}

// String implementing the Stringer
func (l *Label) String() string {
	return fmt.Sprintf("%s %s", LABEL, l.name)
}

// Translate generates assembly code for the label command
func (l *Label) Translate(t *SymbolTable, wr io.Writer) error {
	var ft *functionTable
	if l.function == nil {
		ft = t.FunctionTable("")
	} else {
		ft = t.FunctionTable(l.function.name)
	}

	data := map[string]string{
		"cmdLit": l.lit,
		"name":   l.name,
		"label":  ft.Label(l.name),
	}

	err := labelAsmTmpl.Execute(wr, data)

	return err
}

func parseLabel(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
	tok, lit, err := p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != LABEL {
		panic("internal error")
	}
	l := &Label{
		lit:      lit,
		function: ctx.function,
	}
	return ctx, parseLabelName(l)
}

func parseLabelName(l *Label) stateFunc {
	return func(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
		tok, lit, err := p.scanIgnore()
		if err != nil {
			return ctx, parseError(err)
		}
		if tok != VALUE {
			return ctx, parseError(fmt.Errorf("expect label value, got token %s (%s)", tok, lit))
		}
		l.name = lit
		return ctx, command(l)
	}
}

const ifGotoAsm = `// {{ .cmdLit }} {{ .label }}
	// pop
	@SP
	M=M-1 // SP--
	A=M
	D=M // D=*SP

	@{{ .jumpLabel }}
	D;JNE
`

var ifGotoAsmTmpl *template.Template

func init() {
	ifGotoAsmTmpl = template.Must(template.New("ifGotoAsm").Parse(ifGotoAsm))
}

// IfGoto implements the if-goto command
type IfGoto struct {
	lit   string
	label string

	function *Function
}

// String implements Stringer
func (g *IfGoto) String() string {
	return fmt.Sprintf("%s %s", IFGOTO, g.label)
}

// Translate generates assembly code for if-goto
func (g *IfGoto) Translate(t *SymbolTable, wr io.Writer) error {
	var funcT *functionTable
	if g.function == nil {
		funcT = t.FunctionTable("")
	} else {
		funcT = t.FunctionTable(g.function.name)
	}

	data := map[string]string{
		"cmdLit":    g.lit,
		"label":     g.label,
		"jumpLabel": funcT.Label(g.label),
	}

	err := ifGotoAsmTmpl.Execute(wr, data)

	return err
}

func parseIfGoto(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
	tok, lit, err := p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != IFGOTO {
		panic("internal error")
	}
	g := &IfGoto{
		lit: lit,

		function: ctx.function,
	}

	tok, lit, err = p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != VALUE {
		return ctx, parseError(fmt.Errorf("invalid token %s (%s) after if-goto. expect label"))
	}
	g.label = lit

	return ctx, command(g)
}

const gotoAsm = `// {{ .cmdLit }} {{ .label }}
	@{{ .jumpLabel }}
	0;JMP
`

var gotoAsmTmpl *template.Template

func init() {
	gotoAsmTmpl = template.Must(template.New("gotoAsm").Parse(gotoAsm))
}

// Goto implements the goto command
type Goto struct {
	lit   string
	label string

	function *Function
}

// String implements the Stringer
func (g *Goto) String() string {
	return fmt.Sprintf("%s %s", GOTO, g.label)
}

// Translate generates assembly code for goto
func (g *Goto) Translate(t *SymbolTable, wr io.Writer) error {
	var ft *functionTable
	if g.function == nil {
		ft = t.FunctionTable("")
	} else {
		ft = t.FunctionTable(g.function.name)
	}

	data := map[string]string{
		"cmdLit":    g.lit,
		"label":     g.label,
		"jumpLabel": ft.Label(g.label),
	}

	err := gotoAsmTmpl.Execute(wr, data)
	return err
}

func parseGoto(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
	tok, lit, err := p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != GOTO {
		panic("internal error")
	}
	g := &Goto{
		lit: lit,

		function: ctx.function,
	}

	tok, lit, err = p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != VALUE {
		return ctx, parseError(fmt.Errorf("invalid token %s (%s) after goto. expect label", tok, lit))
	}
	g.label = lit

	return ctx, command(g)
}
