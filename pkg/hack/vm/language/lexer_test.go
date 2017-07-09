package language_test

import (
	"bytes"
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
