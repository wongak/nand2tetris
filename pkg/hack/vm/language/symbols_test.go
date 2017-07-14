package language_test

import (
	"testing"

	"github.com/wongak/nand2tetris/pkg/hack/vm/language"
)

func TestSymbolStatic(t *testing.T) {
	tbl, err := language.NewSymbolTable().RegisterFile("testFile")
	if err != nil {
		t.Fatal(err)
	}
	expect := "testFile.0"
	got := tbl.Static(12)
	if got != expect {
		t.Errorf("expect symbol, %s got %s", expect, got)
	}
	expect = "testFile.1"
	got = tbl.Static(123)
	if got != expect {
		t.Errorf("expect symbol, %s got %s", expect, got)
	}
	expect = "testFile.0"
	got = tbl.Static(12)
	if got != expect {
		t.Errorf("expect symbol, %s got %s", expect, got)
	}

	expect = "testFile.2"
	got = tbl.Static(1234)
	if got != expect {
		t.Errorf("expect symbol, %s got %s", expect, got)
	}

	expect = "testFile.2"
	got = tbl.Static(1234)
	if got != expect {
		t.Errorf("expect symbol, %s got %s", expect, got)
	}
}
