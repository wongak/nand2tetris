package language

import (
	"fmt"
	"sync/atomic"
)

// SymbolTable holds all known symbols for a VM file
type SymbolTable struct {
	fileName string

	static *staticVars
	// conditions
	condIndex int64
}

type staticVars struct {
	mapping map[int]int64
	j       int64
}

func newStaticVars() *staticVars {
	return &staticVars{
		mapping: make(map[int]int64),
		j:       -1,
	}
}

func (s *staticVars) index(index int) int64 {
	if _, ok := s.mapping[index]; !ok {
		s.mapping[index] = atomic.AddInt64(&s.j, 1)
	}
	return s.mapping[index]
}

// NewSymbolTable creates a new symbol table for a VM file
func NewSymbolTable(fileName string) *SymbolTable {
	return &SymbolTable{
		fileName: fileName,

		static:    newStaticVars(),
		condIndex: -1,
	}
}

// Static returns the symbol for the static variable index i
func (t *SymbolTable) Static(index int) string {
	return fmt.Sprintf("%s.%d", t.fileName, t.static.index(index))
}

// Condition returns a condition label per symbol table
func (t *SymbolTable) Condition() string {
	return fmt.Sprintf("%s.if.%d", t.fileName, atomic.AddInt64(&t.condIndex, 1))
}
