package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wongak/nand2tetris/pkg/hack/vm/language"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("expecting exactly one argument. vm file to translate")
		os.Exit(1)
	}

	in, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("error opening vm file: %v\n", err)
		os.Exit(1)
	}
	defer in.Close()

	fileDir := filepath.Dir(os.Args[1])
	fileBase := filepath.Base(os.Args[1])
	parts := strings.Split(fileBase, ".")
	cfgName := parts[:len(parts)-1]
	parts[len(parts)-1] = "asm"
	outFileName := filepath.Join(fileDir, strings.Join(parts, "."))

	out, err := os.OpenFile(outFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("error opening output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	cfg := language.TranslateConfig{
		FileName: strings.Join(cfgName, "."),
	}

	parser := language.NewParser(in)
	err = parser.Run()
	if err != nil {
		fmt.Printf("error while parsing vm file: %v\n", err)
		os.Exit(1)
	}
	for _, cmd := range parser.Tree() {
		cmd.Translate(cfg, out)
	}
}
