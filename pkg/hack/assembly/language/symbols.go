package language

type SymbolTable struct {
	addrs  map[string]int
	labels map[string]int

	instruction int
	addr        int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		addrs: map[string]int{
			"R0":  0,
			"R1":  1,
			"R2":  2,
			"R3":  3,
			"R4":  4,
			"R5":  5,
			"R6":  6,
			"R7":  7,
			"R8":  8,
			"R9":  9,
			"R10": 10,
			"R11": 11,
			"R12": 12,
			"R13": 13,
			"R14": 14,
			"R15": 15,

			"SP":   0,
			"LCL":  1,
			"ARG":  2,
			"THIS": 3,
			"THAT": 4,

			"SCREEN": 0x4000,
			"KBD":    0x6000,
		},
		labels: make(map[string]int),

		addr: 0x0010,
	}
}

func (t *SymbolTable) Label(str string) int {
	if label, ok := t.labels[str]; ok {
		return label
	}
	if s, ok := t.addrs[str]; ok {
		return s
	}
	t.addrs[str] = t.addr
	t.addr++
	return t.addrs[str]
}

func (t *SymbolTable) RegisterLabel(str string) int {
	current := t.instruction
	t.labels[str] = current
	return current
}

func (t *SymbolTable) RegisterInstruction() {
	t.instruction++
}

func (t *SymbolTable) CurrentInstruction() int {
	return t.instruction
}
