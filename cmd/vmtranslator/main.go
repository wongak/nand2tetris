package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wongak/nand2tetris/pkg/hack/vm/language"
)

var (
	headless bool
	verbose  bool
)

func main() {
	flag.BoolVar(&headless, "hl", false, "headless mode")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("expecting exactly one argument. vm file to translate")
		os.Exit(1)
	}

	inputFileName := flag.Args()[0]

	info, err := os.Stat(inputFileName)
	if err != nil {
		fmt.Printf("error on stat input: %v\n", err)
		os.Exit(1)
	}

	table := language.NewSymbolTable()

	var outFileName string
	if info.IsDir() {
		abs, err := filepath.Abs(info.Name())
		if err != nil {
			fmt.Printf("error resolving out path: %v", err)
			os.Exit(1)
		}

		outFileName = filepath.Join(abs, filepath.Base(abs)+".asm")
	} else {
		fileDir := filepath.Dir(inputFileName)
		fileBase := filepath.Base(inputFileName)
		parts := strings.Split(fileBase, ".")
		parts[len(parts)-1] = "asm"
		outFileName = filepath.Join(fileDir, strings.Join(parts, "."))
	}
	if verbose {
		fmt.Printf("writing %s\n", outFileName)
	}

	out, err := os.OpenFile(outFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("error opening output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	if !headless {
		_, err = out.WriteString(`// BOOT
@256
D=A
@SP
M=D // SP = 256
`)
		if err != nil {
			fmt.Printf("error writing bootstrap: %v\n", err)
			os.Exit(1)
		}
		sysinit := language.NewCall("Sys.init", 0)
		err = sysinit.Translate(table, out)
		if err != nil {
			fmt.Printf("error writing bootstrap: %v\n", err)
			os.Exit(1)
		}

	}

	if info.IsDir() {
		dir, err := os.Open(inputFileName)
		if err != nil {
			fmt.Printf("error opening dir: %v\n", err)
			os.Exit(1)
		}
		defer dir.Close()
		content, err := dir.Readdir(-1)
		if err != nil {
			fmt.Printf("error reading dir: %v\n", err)
			os.Exit(1)
		}
		for _, f := range content {
			if f.Name() == "." || f.Name() == ".." {
				continue
			}
			parts := strings.Split(f.Name(), ".")
			if parts[len(parts)-1] != "vm" {
				continue
			}
			if verbose {
				fmt.Print(f.Name(), " ")
			}
			err = parseFile(out, table, filepath.Join(dir.Name(), f.Name()))
			if err != nil {
				fmt.Printf("parse error on file %s: %v\n", f.Name(), err)
				os.Exit(1)
			}
			if verbose {
				fmt.Println(".")
			}
		}
	} else {
		err = parseFile(out, table, inputFileName)
		if err != nil {
			fmt.Printf("parse error on file %s: %v\n", inputFileName, err)
			os.Exit(1)
		}

	}

	if headless {
		_, err = out.WriteString(`// END
(END)
@END
0;JMP
`)
		if err != nil {
			fmt.Printf("error writing halt: %v\n", err)
			os.Exit(1)
		}
	}
}

func parseFile(wr io.Writer, table *language.SymbolTable, fileName string) error {
	fileBase := filepath.Base(fileName)
	parts := strings.Split(fileBase, ".")
	cfgName := parts[:len(parts)-1]
	symbolTableFileName := strings.Join(cfgName, ".")

	in, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening vm file: %v", err)
	}
	defer in.Close()

	if verbose {
		fmt.Printf("parsing %s...\n", symbolTableFileName)
	}
	p := language.NewParser(in)

	err = p.Run(table, symbolTableFileName)
	if err != nil {
		return err
	}

	for _, cmd := range p.Tree() {
		if verbose {
			fmt.Printf("  %+v\n", cmd)
		}
		err = cmd.Translate(table, wr)
		if err != nil {
			return err
		}

	}

	return nil
}
