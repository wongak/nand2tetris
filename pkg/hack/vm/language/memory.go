package language

import (
	"fmt"
	"io"
	"strconv"
	"text/template"
)

// push constant i
//
// where .indexLit is the constant value to push on the stack
const memoryPushConstant = `// {{ .cmdLit }} {{ .segLit }} {{ .indexLit }}
	@{{ .indexLit }}
	D=A
	@SP
	A=M // *SP
	M=D // =i
	@SP
	M=M+1 // SP++
`

var memoryPushConstantTmpl *template.Template

// push segment i
//
// where .segSymbol is one of LCL, ARG, THIS, THAT
// and .indexLit being the index of the segment
const memoryPushSegment = `// {{ .cmdLit }} {{ .segLit }} {{ .indexLit }}
	@{{ .indexLit }}
	D=A
	@{{ .segSymbol }}
	D=D+M
	A=D
	D=M // D = *({{ .segSymbol }} + i)
	
	@SP
	A=M
	M=D // *SP=*addr
	@SP
	M=M+1 // SP++
`

var memoryPushSegmentTmpl *template.Template

const memoryPushTemp = `// {{ .cmdLit }} {{ .segLit }} {{ .indexLit }}
	@{{ .indexLit }}
	D=A
	@{{ .segSymbol }}
	D=D+A
	A=D
	D=M // D = *({{ .segSymbol }} + i)
	
	@SP
	A=M
	M=D // *SP=*addr
	@SP
	M=M+1 // SP++
`

var memoryPushTempTmpl *template.Template

const memoryPushStatic = `// {{ .cmdLit }} {{ .segLit }} {{ .indexLit }}
	@{{ .staticVar }}
	D=M
	
	@SP
	A=M
	M=D // *SP=@{{ .staticVar }}
	@SP
	M=M+1 // SP++
`

var memoryPushStaticTmpl *template.Template

const memoryPushPointer = `// {{ .cmdLit }} {{ .segLit }} {{ .indexLit }}
	@{{ .segSymbol }}
	D=M
	
	@SP
	A=M
	M=D
	@SP
	M=M+1
`

var memoryPushPointerTmpl *template.Template

const memoryPopSegment = `// {{ .cmdLit }} {{ .segLit }} {{ .indexLit }}
	@{{ .indexLit }}
	D=A
	@{{ .segSymbol}}
	D=D+M
	@R13
	M=D // addr = {{ .segSymbol }} + {{ .indexLit}}
	
	@SP
	M=M-1 // SP--
	A=M
	D=M // D=*SP
	@R13
	A=M
	M=D // *addr=*SP
`

var memoryPopSegmentTmpl *template.Template

const memoryPopTemp = `// {{ .cmdLit }} {{ .segLit }} {{ .indexLit }}
	@{{ .indexLit }}
	D=A
	@{{ .segSymbol}}
	D=D+A
	@R13
	M=D // addr = {{ .segSymbol }} + {{ .indexLit}}
	
	@SP
	M=M-1 // SP--
	A=M
	D=M // D=*SP
	@R13
	A=M
	M=D // *addr=*SP
`

var memoryPopTempTmpl *template.Template

const memoryPopStatic = `// {{ .cmdLit }} {{ .segLit }} {{ .indexLit }}
	@SP
	M=M-1 // SP--
	A=M
	D=M // D=*SP
	
	@{{ .staticVar }}
	M=D // @{{ .staticVar}} = D (*SP)
`

var memoryPopStaticTmpl *template.Template

const memoryPopPointer = `// {{ .cmdLit }} {{ .segLit }} {{ .indexLit }}
	@SP
	M=M-1
	A=M
	D=M
	
	@{{ .segSymbol }}
	M=D
`

var memoryPopPointerTmpl *template.Template

func init() {
	memoryPushConstantTmpl = template.Must(template.New("pushConstant").Parse(memoryPushConstant))
	memoryPushSegmentTmpl = template.Must(template.New("pushSegment").Parse(memoryPushSegment))
	memoryPushTempTmpl = template.Must(template.New("pushTemp").Parse(memoryPushTemp))
	memoryPushStaticTmpl = template.Must(template.New("pushStatic").Parse(memoryPushStatic))
	memoryPushPointerTmpl = template.Must(template.New("pushPointer").Parse(memoryPushPointer))
	memoryPopSegmentTmpl = template.Must(template.New("popSegment").Parse(memoryPopSegment))
	memoryPopTempTmpl = template.Must(template.New("popTemp").Parse(memoryPopTemp))
	memoryPopStaticTmpl = template.Must(template.New("popStatic").Parse(memoryPopStatic))
	memoryPopPointerTmpl = template.Must(template.New("popPointer").Parse(memoryPopPointer))
}

type (
	// Segment describes a memory segment to be accessed
	Segment struct {
		seg      Token
		segLit   string
		index    int
		indexLit string
	}

	// MemoryAccess represents a memory access command
	MemoryAccess struct {
		accessComamnd Token // push or pop
		lit           string
		seg           Segment
	}
)

// String implementing Stringer
func (m *MemoryAccess) String() string {
	return fmt.Sprintf("%s %s %d", m.accessComamnd, m.seg.seg, m.seg.index)
}

// Translate translates the VM command to assembly
func (m *MemoryAccess) Translate(t *SymbolTable, wr io.Writer) error {
	data := map[string]string{
		"cmdLit":    m.lit,
		"segLit":    m.seg.segLit,
		"indexLit":  m.seg.indexLit,
		"segSymbol": m.seg.seg.String(),
	}
	var tmpl *template.Template
	if m.accessComamnd == PUSH {
		tmpl = memoryPushSegmentTmpl
	} else {
		tmpl = memoryPopSegmentTmpl
	}
	// TEMP segment on R5 - R12
	if m.seg.seg == TEMP {
		data["segSymbol"] = "R5"
		if m.accessComamnd == PUSH {
			tmpl = memoryPushTempTmpl
		} else {
			tmpl = memoryPopTempTmpl
		}
	}
	// STATIC as assembly variable Filename.i
	if m.seg.seg == STATIC {
		data["staticVar"] = t.Static(m.seg.index)
		if m.accessComamnd == PUSH {
			tmpl = memoryPushStaticTmpl
		} else {
			tmpl = memoryPopStaticTmpl
		}
	}
	// CONSTANT
	if m.seg.seg == CONSTANT {
		tmpl = memoryPushConstantTmpl
	}
	// POINTER
	if m.seg.seg == POINTER {
		if m.seg.index == 0 {
			data["segSymbol"] = "THIS"
		} else {
			data["segSymbol"] = "THAT"
		}
		if m.accessComamnd == PUSH {
			tmpl = memoryPushPointerTmpl
		} else {
			tmpl = memoryPopPointerTmpl
		}
	}
	err := tmpl.Execute(wr, data)
	if err != nil {
		return err
	}
	return nil
}

func parseMemoryAccess(f *Function) stateFunc {
	return func(p *Parser) stateFunc {
		tok, lit, err := p.scanIgnore()
		if err != nil {
			return parseError(err)
		}
		if tok != PUSH && tok != POP {
			return parseError(fmt.Errorf("invalid token %s (%s)", tok, lit))
		}
		cmd := &MemoryAccess{
			accessComamnd: tok,
			lit:           lit,
		}
		return parseSegment(cmd, f)
	}
}

func parseSegment(cmd *MemoryAccess, f *Function) stateFunc {
	return func(p *Parser) stateFunc {
		tok, lit, err := p.scanIgnore()
		if err != nil {
			return parseError(err)
		}
		if isSegment(tok) {
			cmd.seg.seg = tok
			cmd.seg.segLit = lit
			return parseSegmentIndex(cmd, f)
		}
		if tok == CONSTANT {
			if cmd.accessComamnd == POP {
				return parseError(fmt.Errorf("invalid POP on constant"))
			}
			cmd.seg.seg = tok
			cmd.seg.segLit = lit
			return parseSegmentIndex(cmd, f)
		}
		if tok == POINTER {
			cmd.seg.seg = tok
			cmd.seg.segLit = lit
			return parsePointerIndex(cmd, f)
		}
		return parseError(fmt.Errorf("invalid token %s (%s). expect segment/pointer", tok, lit))
	}
}

func parseSegmentIndex(cmd *MemoryAccess, f *Function) stateFunc {
	return func(p *Parser) stateFunc {
		tok, lit, err := p.scanIgnore()
		if err != nil {
			return parseError(err)
		}
		if tok != VALUE {
			return parseError(fmt.Errorf("invalid token %s (%s). expect value", tok, lit))
		}
		i, err := strconv.ParseInt(lit, 10, 64)
		if err != nil {
			return parseError(fmt.Errorf("invalid value %s: %s", lit, err))
		}
		cmd.seg.index = int(i)
		cmd.seg.indexLit = lit

		return command(cmd, f)
	}
}

func parsePointerIndex(cmd *MemoryAccess, f *Function) stateFunc {
	return func(p *Parser) stateFunc {
		tok, lit, err := p.scanIgnore()
		if err != nil {
			return parseError(err)
		}
		if lit != "1" && lit != "0" {
			return parseError(fmt.Errorf("invalid token %s (%s). expect 0 or 1", tok, lit))
		}
		p.unscan()
		return parseSegmentIndex(cmd, f)
	}
}
