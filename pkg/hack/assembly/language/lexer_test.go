package language_test

import (
	"strings"
	"testing"

	. "github.com/wongak/nand2tetris/pkg/hack/assembly/language"
)

type testcase struct {
	input  string
	expect []expectation
}

type expectation struct {
	tok Token
	lit string
}

func runExpectations(t *testing.T, cases []testcase) {
	for _, c := range cases {
		sc := NewScanner(strings.NewReader(c.input))
		i := 0
		for tok, lit, err := sc.Scan(); tok != EOF; tok, lit, err = sc.Scan() {
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if tok == WS || tok == COMMENT {
				continue
			}
			if c.expect[i].tok != tok {
				t.Errorf("expect token %s, got %s", c.expect[i].tok, tok)
			}
			if c.expect[i].lit != lit {
				t.Errorf("expect literal %s, got %s", c.expect[i].lit, lit)
			}
			i++
		}
	}
}

func TestScanAInstruction(t *testing.T) {
	runExpectations(t, []testcase{
		{
			input: `
@14
@15
@R12
@i
@str
			`,
			expect: []expectation{
				{
					tok: AT,
					lit: "@",
				},
				{
					tok: VALUE,
					lit: "14",
				},
				{
					tok: AT,
					lit: "@",
				},
				{
					tok: VALUE,
					lit: "15",
				},
				{
					tok: AT,
					lit: "@",
				},
				{
					tok: VALUE,
					lit: "R12",
				},
				{
					tok: AT,
					lit: "@",
				},
				{
					tok: VALUE,
					lit: "i",
				},
				{
					tok: AT,
					lit: "@",
				},
				{
					tok: VALUE,
					lit: "str",
				},
			},
		},
	})
}

func TestProg(t *testing.T) {
	runExpectations(t, []testcase{
		{
			input: `
@A
D=A
@0
A=D
M=D

(END)
@END
0;JMP
			`,
			expect: []expectation{
				{
					tok: AT,
					lit: "@",
				},
				{
					tok: VALUE,
					lit: "A",
				},
				{
					tok: VALUE,
					lit: "D",
				},
				{
					tok: EQUALS,
					lit: "=",
				},
				{
					tok: VALUE,
					lit: "A",
				},
				{
					tok: AT,
					lit: "@",
				},
				{
					tok: VALUE,
					lit: "0",
				},
				{
					tok: VALUE,
					lit: "A",
				},
				{
					tok: EQUALS,
					lit: "=",
				},
				{
					tok: VALUE,
					lit: "D",
				},
				{
					tok: VALUE,
					lit: "M",
				},
				{
					tok: EQUALS,
					lit: "=",
				},
				{
					tok: VALUE,
					lit: "D",
				},
				{
					tok: LABEL_START,
					lit: "(",
				},
				{
					tok: VALUE,
					lit: "END",
				},
				{
					tok: LABEL_END,
					lit: ")",
				},
				{
					tok: AT,
					lit: "@",
				},
				{
					tok: VALUE,
					lit: "END",
				},
				{
					tok: VALUE,
					lit: "0",
				},
				{
					tok: SEMICOLON,
					lit: ";",
				},
				{
					tok: JMP,
					lit: "JMP",
				},
			},
		},
	})
}
