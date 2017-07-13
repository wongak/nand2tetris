package language

import (
	"fmt"
	"io"
)

// Command represents a valid command
type Command interface {
	fmt.Stringer

	Translate(*SymbolTable, io.Writer) error
}

const endOp = `// END
(END)
	@END
	0,JMP
`
