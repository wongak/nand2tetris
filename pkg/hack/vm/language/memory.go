package language

import (
	"fmt"
	"io"
	"strconv"
	"text/template"
)

// Segment describes a memory segment to be accessed
type Segment struct {
	seg      Token
	segLit   string
	index    int
	indexLit string
}

// MemoryAccess represents a memory access command
type MemoryAccess struct {
	accessComamnd Token // push or pop
	lit           string
	seg           Segment
}

// String implementing Stringer
func (m *MemoryAccess) String() string {
	return fmt.Sprintf("%s %s %d", m.accessComamnd, m.seg.seg, m.seg.index)
}

// Translate translates the VM command to assembly
func (m *MemoryAccess) Translate(cfg TranslateConfig, wr io.Writer) {
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
		data["staticVar"] = cfg.FileName + "." + m.seg.indexLit
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
		panic(err)
	}
}

func parseMemoryAccess(p *Parser) stateFunc {
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
	return parseSegment(cmd)
}

func parseSegment(cmd *MemoryAccess) stateFunc {
	return func(p *Parser) stateFunc {
		tok, lit, err := p.scanIgnore()
		if err != nil {
			return parseError(err)
		}
		if isSegment(tok) {
			cmd.seg.seg = tok
			cmd.seg.segLit = lit
			return parseSegmentIndex(cmd)
		}
		if tok == CONSTANT {
			if cmd.accessComamnd == POP {
				return parseError(fmt.Errorf("invalid POP on constant"))
			}
			cmd.seg.seg = tok
			cmd.seg.segLit = lit
			return parseSegmentIndex(cmd)
		}
		if tok == POINTER {
			cmd.seg.seg = tok
			cmd.seg.segLit = lit
			return parsePointerIndex(cmd)
		}
		return parseError(fmt.Errorf("invalid token %s (%s). expect segment/pointer", tok, lit))
	}
}

func parseSegmentIndex(cmd *MemoryAccess) stateFunc {
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

		return command(cmd)
	}
}

func parsePointerIndex(cmd *MemoryAccess) stateFunc {
	return func(p *Parser) stateFunc {
		tok, lit, err := p.scanIgnore()
		if err != nil {
			return parseError(err)
		}
		if lit != "1" && lit != "0" {
			return parseError(fmt.Errorf("invalid token %s (%s). expect 0 or 1", tok, lit))
		}
		p.unscan()
		return parseSegmentIndex(cmd)
	}
}
