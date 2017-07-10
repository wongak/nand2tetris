package language

import "fmt"

// Command represents a valid command
type Command interface {
	fmt.Stringer

	Translate(TranslateConfig) string
}

type TranslateConfig struct {
	fileName string
}
