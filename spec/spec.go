package spec

// Opcode ...
type Opcode struct {
	Word    string `json:"word"`
	WordAlt string `json:"word_alt"`
	Opcode  string `json:"opcode"`
	Input   string `json:"input"`
	Output  string `json:"output"`
	Short   string `json:"short"`
}

// Script ...
type Script map[string]Opcode
