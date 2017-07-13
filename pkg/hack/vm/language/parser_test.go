package language_test

import (
	"strings"
	"testing"

	"github.com/wongak/nand2tetris/pkg/hack/vm/language"
)

func TestMemAccessParser(t *testing.T) {
	code := `
// test code
// parsing mem acces
push constant 1
push constant 15
pop argument 0
push constant 23
pop argument 1
	`
	p := language.NewParser(strings.NewReader(code))
	err := p.Run()
	if err != nil {
		t.Errorf("unexpected error on parse: %v", err)
		return
	}

}

func TestMemPopOnConstantReturnsError(t *testing.T) {
	code := `
// test
pop constant 12
	`
	p := language.NewParser(strings.NewReader(code))
	err := p.Run()
	if err == nil {
		t.Error("expect error")
		return
	}
	if !strings.Contains(err.Error(), "invalid POP on constant") {
		t.Errorf("expect pop on constant err, got: %v", err)
		return
	}
}

func TestPointerInvalidValue(t *testing.T) {
	code := `
// test
push constant 12
push pointer 3
	`
	p := language.NewParser(strings.NewReader(code))
	err := p.Run()
	if err == nil {
		t.Error("expect error")
		return
	}
	if !strings.Contains(err.Error(), "invalid token VALUE (3). expect 0 or 1") {
		t.Errorf("expect pop on constant err, got: %v", err)
		return
	}
}

type parserTestCase struct {
	input    string
	expected []func(language.Command) bool
}

func execParserTestCase(t *testing.T, c parserTestCase) {
	rd := strings.NewReader(c.input)
	p := language.NewParser(rd)

	err := p.Run()
	if err != nil {
		t.Errorf("error on parsing: %v", err)
		return
	}

	cmds := p.Tree()
	for i, expect := range c.expected {
		if !expect(cmds[i]) {
			t.Errorf("error on expectation %d (%+v)", i, cmds[i])
		}
	}
}

func TestLabelParser(t *testing.T) {
	tst := parserTestCase{
		input: `
	function test 0
	push constant 1
	label ABC
	return
			`,
		expected: []func(language.Command) bool{
			func(cmd language.Command) bool {
				if cmd.String() != "PUSH CONSTANT 1" {
					return false
				}
				return true
			},
			func(cmd language.Command) bool {
				if cmd.String() != "LABEL ABC" {
					return false
				}
				return true
			},
		},
	}
	execParserTestCase(t, tst)
}
