package language_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wongak/nand2tetris/pkg/hack/vm/language"
)

func TestWhitespaceScan(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	tstStr := "\n    \t\r\n    \n"
	buf.WriteString(tstStr)

	sc := language.NewScanner(buf)

	tok, lit, err := sc.Scan()
	if err != nil {
		t.Errorf("received unexpected error %v", err)
		return
	}
	if tok != language.WS {
		t.Error("expect WS token")
		return
	}
	if lit != tstStr {
		t.Errorf("unexpected literal str: %s", lit)
		return
	}

	tok, lit, err = sc.Scan()
	if tok != language.EOF {
		t.Error("expect EOF")
		return
	}
}

func TestScanPush(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	testStr := "   \n\t \r  push \t"
	buf.WriteString(testStr)

	sc := language.NewScanner(buf)

	tok, _, err := sc.Scan()
	if err != nil {
		t.Errorf("received unexpected error %v", err)
		return
	}
	if tok != language.WS {
		t.Error("expect start with WS")
		return
	}
	tok, _, err = sc.Scan()
	if err != nil {
		t.Errorf("received unexpected error on second scan %v", err)
		return
	}
	if tok != language.PUSH {
		t.Error("expect PUSH token")
		return
	}
	tok, _, err = sc.Scan()
	if err != nil {
		t.Errorf("received unexpected error on final ws: %v", err)
		return
	}
	if tok != language.WS {
		t.Error("expect final WS")
		return
	}
	tok, _, err = sc.Scan()
	if err != nil {
		t.Errorf("error finalizing: %v", err)
		return
	}
	if tok != language.EOF {
		t.Error("expect EOF token")
		return
	}
}

type testCase struct {
	input    string
	expected []expect
}

type expect struct {
	tok language.Token
	lit string
}

func execTestCase(t *testing.T, c testCase) {
	rd := strings.NewReader(c.input)
	sc := language.NewScanner(rd)

	i := 0
	for _, expect := range c.expected {
		tok := language.WS
		var lit string
		var err error
		for tok == language.WS {
			tok, lit, err = sc.Scan()
			if err != nil {
				t.Errorf("received error on token %d: %v (%s)", i, err, c.input)
				return
			}
		}
		if tok != expect.tok {
			t.Errorf("unexpected token %s index %d. expected %s (%s)", tok, i, expect.tok, c.input)
			return
		}
		if lit != expect.lit {
			t.Errorf("unexpected literal %s index %d, expected %s (%s)", lit, i, expect.lit, c.input)
			return
		}
		i++
	}
}

func ProvideMemoryAccessCommands() []testCase {
	return []testCase{
		{
			input: "push constant 15",
			expected: []expect{
				{
					tok: language.PUSH,
					lit: "push",
				},
				{
					tok: language.CONSTANT,
					lit: "constant",
				},
				{
					tok: language.VALUE,
					lit: "15",
				},
			},
		},
		{
			input: "pop temp 5",
			expected: []expect{
				{
					tok: language.POP,
					lit: "pop",
				},
				{
					tok: language.TEMP,
					lit: "temp",
				},
				{
					tok: language.VALUE,
					lit: "5",
				},
			},
		},
		{
			input: `
			// will pop from stack to temp 5
			pop temp 5`,
			expected: []expect{
				{
					tok: language.COMMENT,
					lit: "will pop from stack to temp 5",
				},
				{
					tok: language.POP,
					lit: "pop",
				},
				{
					tok: language.TEMP,
					lit: "temp",
				},
				{
					tok: language.VALUE,
					lit: "5",
				},
			},
		},
	}
}

func TestMemoryAccessCommandsLexer(t *testing.T) {
	for _, c := range ProvideMemoryAccessCommands() {
		execTestCase(t, c)
	}
}
