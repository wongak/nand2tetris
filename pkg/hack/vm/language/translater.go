package language

import "text/template"

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

const endOp = `// END
(END)
	@END
	0,JMP
`
