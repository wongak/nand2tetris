package language

import (
	"fmt"
	"io"
)

// Command represents a valid command
type Command interface {
	fmt.Stringer

	Translate(TranslateConfig, io.Writer)
}

// TranslateConfig sets the configuration for translating a VM file
type TranslateConfig struct {
	FileName string
}

type end struct {
}

func (e end) String() string {
	return "end"
}

func (e end) Translate(_ TranslateConfig, wr io.Writer) {
	_, err := wr.Write([]byte(endOp))
	if err != nil {
		panic(err)
	}
}
