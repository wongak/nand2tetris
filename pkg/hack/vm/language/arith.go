package language

import (
	"fmt"
	"io"
	"text/template"
)

const arithmeticOp = `// {{ .cmdLit }}
	@SP
	M=M-1 // SP--
	A=M
	D=M // D=*SP
	
	@SP
	M=M-1 // SP--
	A=M
	{{ .operation }}
	
	M=D
	@SP
	M=M+1 // SP++
`

var arithmeticOpTmpl *template.Template

const arithmeticSingleOp = `// {{ .cmdLit }}
	@SP
	M=M-1 // SP--
	A=M
	D=M // D=*SP

	{{ .operation }}

	M=D
	@SP
	M=M+1 // SP++
`

var arithmeticSingleOpTmpl *template.Template

const logicalComp = `// {{ .cmdLit }}
	@SP
	M=M-1 // SP--
	A=M
	D=M // D=*SP
	
	@SP
	M=M-1 // SP--
	A=M
	D=D-M // D=D-*SP
	
	@R13
	M=-1 // true
	@{{ .labelSet }}
	D;{{ .comp }} // jump if true
	@R13
	M=0 // false
	
({{ .labelSet }})
	@R13
	D=M // D=result
	@SP
	A=M
	M=D
	@SP
	M=M+1
`

var logicalCompTmpl *template.Template

func init() {
	arithmeticOpTmpl = template.Must(template.New("arithmeticOp").Parse(arithmeticOp))
	arithmeticSingleOpTmpl = template.Must(template.New("arithmeticSingleOp").Parse(arithmeticSingleOp))
	logicalCompTmpl = template.Must(template.New("logicalComp").Parse(logicalComp))
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
func (m *Arithmetic) Translate(t *SymbolTable, wr io.Writer) error {
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
		data["labelSet"] = t.Condition()
		tmpl = logicalCompTmpl
		data["comp"] = "JEQ" // true if pop1 - pop2 = 0
	case GT:
		data["labelSet"] = t.Condition()
		tmpl = logicalCompTmpl
		data["comp"] = "JLT" // true if pop1 - pop2 < 0
	case LT:
		data["labelSet"] = t.Condition()
		tmpl = logicalCompTmpl
		data["comp"] = "JGT"
	}
	err := tmpl.Execute(wr, data)
	if err != nil {
		return err
	}
	return nil
}

func parseArithmetic(f *Function) stateFunc {
	return func(p *Parser) stateFunc {
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
		return command(cmd, f)
	}
}
