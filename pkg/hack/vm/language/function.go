package language

import (
	"fmt"
	"io"
	"strconv"
	"text/template"
)

const functionAsm = `// {{ .cmdLit }} {{ .numLocalLit }}
({{ .functionLabel }})
	// init local
{{range .localVars}}
	// push local 0 (arg {{ . }})
	@SP
	A=M
	M=0
	@SP
	M=M+1
{{- end }}
`

var functionAsmTmpl *template.Template

func init() {
	functionAsmTmpl = template.Must(template.New("functionAsm").Parse(functionAsm))
}

// Function implements the VM function command
type Function struct {
	lit string

	name string

	numLocal    int
	numLocalLit string
}

// String implements the Stringer
func (f *Function) String() string {
	return fmt.Sprintf("%s %s %s", FUNCTION, f.name, f.numLocalLit)
}

// Translate creates assembly for the function definition
func (f *Function) Translate(t *SymbolTable, wr io.Writer) error {
	ft, err := t.RegisterFunction(f.name)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"cmdLit":        f.lit,
		"numLocalLit":   f.numLocalLit,
		"functionLabel": ft.FunctionLabel(),
	}
	lcls := make([]int, f.numLocal)
	for i := 0; i < f.numLocal; i++ {
		lcls[i] = i
	}
	data["localVars"] = lcls
	err = functionAsmTmpl.Execute(wr, data)
	return err
}

func parseFunction(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
	tok, lit, err := p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != FUNCTION {
		panic("internal error")
	}
	f := &Function{
		lit: lit,
	}

	tok, lit, err = p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != VALUE {
		return ctx, parseError(fmt.Errorf("invalid token %s (%s), expect function name", tok, lit))
	}
	f.name = lit

	tok, lit, err = p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != VALUE {
		return ctx, parseError(fmt.Errorf("invalid token %s (%s). expect num local vars", tok, lit))
	}
	num, err := strconv.ParseInt(lit, 10, 64)
	if err != nil {
		return ctx, parseError(fmt.Errorf("invalid number of local vars %s: %v", lit, err))
	}
	f.numLocalLit = lit
	f.numLocal = int(num)

	ctx.function = f
	return ctx, command(f)
}

const callAsm = `// {{ .cmdLit }} {{ .nameLit }} {{ .numArgsLit }}
	// push return address
	@{{ .returnLabel }}
	D=A
	@SP
	A=M
	M=D // *SP=return address
	@SP
	M=M+1 // SP++
	// push LCL
	@LCL
	D=M
	@SP
	A=M
	M=D // *SP=LCL
	@SP
	M=M+1 // SP++
	// push ARG
	@ARG
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1 // SP++
	// push THIS
	@THIS
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1 // SP++
	// push THAT
	@THAT
	D=M
	@SP
	A=M
	M=D
	@SP
	M=M+1
	// set arg
	@{{ .argDelta }}
	D=A
	@SP
	D=M-D
	@ARG
	M=D
	// set LCL
	@SP
	D=M
	@LCL
	M=D // LCL =SP

	// goto {{ .nameLit }}
	@{{ .functionLabel }}
	0;JMP
({{ .returnLabel }})
`

var callAsmTmpl *template.Template

func init() {
	callAsmTmpl = template.Must(template.New("callAsm").Parse(callAsm))
}

// Call implements the call command (call a function)
type Call struct {
	lit string

	name string

	numArgs    int
	numArgsLit string
}

// NewCall creates a separate call command
func NewCall(name string, numArgs int) *Call {
	return &Call{
		lit: "call",

		name: name,

		numArgs:    numArgs,
		numArgsLit: strconv.FormatInt(int64(numArgs), 10),
	}
}

// String implements Stringer
func (c *Call) String() string {
	return fmt.Sprintf("%s %s", CALL, c.name)
}

// Translate creates the assembly to call a function
func (c *Call) Translate(t *SymbolTable, wr io.Writer) error {
	ft := t.FunctionTable(c.name)

	data := map[string]interface{}{
		"cmdLit":        c.lit,
		"nameLit":       c.name,
		"numArgsLit":    c.numArgsLit,
		"returnLabel":   ft.ReturnLabel(),
		"functionLabel": ft.FunctionLabel(),
		"argDelta":      5 + c.numArgs,
	}

	err := callAsmTmpl.Execute(wr, data)
	return err
}

func parseCall(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
	tok, lit, err := p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != CALL {
		panic("internal error")
	}

	c := &Call{
		lit: lit,
	}

	tok, lit, err = p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != VALUE {
		return ctx, parseError(fmt.Errorf("invalid token %s (%s), expect function identifier", tok, lit))
	}
	c.name = lit

	tok, lit, err = p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != VALUE {
		return ctx, parseError(fmt.Errorf("invalid token %s (%s), expect number of args", tok, lit))
	}
	n, err := strconv.ParseInt(lit, 10, 64)
	if err != nil {
		return ctx, parseError(fmt.Errorf("invalid token %s (%s), number of args: %v", tok, lit, err))
	}
	c.numArgsLit = lit
	c.numArgs = int(n)

	return ctx, command(c)
}

const returnAsm = `// {{ .cmdLit }}
	@LCL
	D=M
	@R13
	M=D // endFrame = LCL

	@5
	D=D-A
	A=D
	D=M
	@R14
	M=D // retAddr = *(endFrame -5)

	@SP
	M=M-1
	A=M
	D=M
	@ARG
	A=M
	M=D // *ARG = pop()

	@ARG
	D=M+1
	@SP
	M=D // SP = ARG +1

	@R13
	D=M
	D=D-1
	A=D
	D=M // *(endFrame - 1)
	@THAT
	M=D

	@R13
	D=M
	@2
	D=D-A
	A=D
	D=M // *(endFrame - 2)
	@THIS
	M=D

	@R13
	D=M
	@3
	D=D-A
	A=D
	D=M // *(endFrame - 3)
	@ARG
	M=D

	@R13
	D=M
	@4
	D=D-A
	A=D
	D=M // *(endFrame - 4)
	@LCL
	M=D

	@R14
	A=M
	0;JMP
`

var returnAsmTmpl *template.Template

func init() {
	returnAsmTmpl = template.Must(template.New("returnAsm").Parse(returnAsm))
}

// Return implements the return command
type Return struct {
	lit string

	function *Function
}

// String implements Stringer
func (r *Return) String() string {
	return RETURN.String()
}

// Translate creates the assembly for the return command
func (r *Return) Translate(t *SymbolTable, wr io.Writer) error {
	data := map[string]string{
		"cmdLit": r.lit,
	}

	err := returnAsmTmpl.Execute(wr, data)
	return err
}

func parseReturn(p *Parser, ctx ParserContext) (ParserContext, stateFunc) {
	tok, lit, err := p.scanIgnore()
	if err != nil {
		return ctx, parseError(err)
	}
	if tok != RETURN {
		panic("internal error")
	}
	if ctx.function == nil {
		return ctx, parseError(fmt.Errorf("invalid return without function context"))
	}

	ret := &Return{
		lit:      lit,
		function: ctx.function,
	}

	return ctx, command(ret)
}
