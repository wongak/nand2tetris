package language

import (
	"fmt"
	"sync/atomic"
)

// SymbolTable holds all known symbols for a VM file
type SymbolTable struct {
	files               map[string]*fileTable
	functionDefinitions map[string]struct{}
	functions           map[string]*functionTable
}

type fileTable struct {
	fileName string
	static   *staticVars
	// conditions
	condIndex int64
}

type functionTable struct {
	fileName     string
	functionName string
	flabels      map[string]string

	callIndex int64
}

type staticVars struct {
	mapping map[int]int64
	j       int64
}

func newFileTable(fileName string) *fileTable {
	return &fileTable{
		fileName:  fileName,
		static:    newStaticVars(),
		condIndex: -1,
	}
}

func newFunctionTable(functionName string) *functionTable {
	return &functionTable{
		functionName: functionName,
		flabels:      make(map[string]string),

		callIndex: -1,
	}
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

// NewSymbolTable creates a new symbol table for a VM program
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		files:               make(map[string]*fileTable),
		functionDefinitions: make(map[string]struct{}),
		functions:           make(map[string]*functionTable),
	}
}

// RegisterFile registers a new file table
func (t *SymbolTable) RegisterFile(fileName string) (*fileTable, error) {
	if _, ok := t.files[fileName]; ok {
		return nil, fmt.Errorf("file %s already registered", fileName)
	}
	t.files[fileName] = newFileTable(fileName)
	return t.files[fileName], nil
}

// FileTable returns the file table for the given file name
func (t *SymbolTable) FileTable(fileName string) (*fileTable, error) {
	f, ok := t.files[fileName]
	if !ok {
		return nil, fmt.Errorf("no symbol table for unknown file %s", fileName)
	}
	return f, nil
}

// RegisterFunction registers a new function table
func (t *SymbolTable) RegisterFunction(fName string) (*functionTable, error) {
	if _, ok := t.functionDefinitions[fName]; ok {
		return nil, fmt.Errorf("function %s already registered", fName)
	}
	t.functionDefinitions[fName] = struct{}{}
	return t.FunctionTable(fName), nil
}

// FunctionTable returns the function table for the given name
func (t *SymbolTable) FunctionTable(fName string) *functionTable {
	ft, ok := t.functions[fName]
	if !ok {
		ft = newFunctionTable(fName)
		t.functions[fName] = ft
	}
	return ft
}

// Static returns the symbol for the static variable index i
func (t *fileTable) Static(index int) string {
	return fmt.Sprintf("%s.%d", t.fileName, t.static.index(index))
}

// Condition returns a condition label per symbol table
func (t *fileTable) Condition() string {
	return fmt.Sprintf("%s.if.%d", t.fileName, atomic.AddInt64(&t.condIndex, 1))
}

// RegisterLabel registers a function scoped label
func (t *functionTable) RegisterLabel(label string) error {
	if _, ok := t.flabels[label]; ok {
		return fmt.Errorf("label %s already registered.")
	}
	t.flabels[label] = t.Label(label)
	return nil
}

// Label generates a function label
func (t *functionTable) Label(label string) string {
	return fmt.Sprintf("%s$%s", t.functionName, label)
}

// FunctionLabel generates a label for the function definition
func (t *functionTable) FunctionLabel() string {
	return t.functionName
}

// ReturnLabel generates a return label for a function call
func (t *functionTable) ReturnLabel() string {
	return fmt.Sprintf("%s$ret.%d", t.FunctionLabel(), atomic.AddInt64(&t.callIndex, 1))
}
