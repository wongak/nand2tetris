package language

import "text/template"

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

const memoryPopSegment = `// {{ .cmdLit }} {{ .segLit }} {{ .indexLit }}
@{{ .indexLit }}
D=A
@{{ .segSymbol}}
D=D+M
@R13
D=M // addr = {{ .segSymbol }} + {{ .indexLit}}

@SP
M=M-1 // SP--
A=M
D=M // D=*SP
@R13
M=D // *addr=*SP
`

var memoryPopSegmentTmpl *template.Template

func init() {
	memoryPushConstantTmpl = template.Must(template.New("pushConstant").Parse(memoryPushConstant))

	memoryPushSegmentTmpl = template.Must(template.New("pushSegment").Parse(memoryPushSegment))

	memoryPopSegmentTmpl = template.Must(template.New("popSegment").Parse(memoryPopSegment))
}
