package language_test

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	. "github.com/wongak/nand2tetris/pkg/hack/assembly/language"
)

func TestParseAInstruction(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	input := strings.NewReader(`
	@R2
	@111
	@SCREEN
	@NEW
	`)
	p := NewParser(input)
	err := p.Run()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	cmds := p.Tree()
	tbl := NewSymbolTable()
	for _, cmd := range cmds {
		t.Logf("cmd: %s", cmd)
		err = cmd.Translate(tbl, buf)
		if err != nil {
			t.Errorf("unexpected translate error: %v", err)
			return
		}
	}
	expect := []string{
		"0000000000000010",
		"0" + fmt.Sprintf("%015b", 111),
		"0" + fmt.Sprintf("%015b", 0x4000),
		"0" + fmt.Sprintf("%015b", 16),
	}
	sc := bufio.NewScanner(strings.NewReader(buf.String()))
	for i, e := range expect {

		if !sc.Scan() {
			t.Errorf("expect scan %d to succeed", i)
			return
		}
		if sc.Text() != e {
			t.Errorf("unexpected translation. expect\n%s, got\n%s", e, sc.Text())
		}
	}
	if err := sc.Err(); err != nil {
		t.Errorf("unexpected error on scan: %v", err)
		return
	}
	if tbl.CurrentInstruction() != 4 {
		t.Errorf("expect current instruction 4, got %d", tbl.CurrentInstruction())
		return
	}
}
