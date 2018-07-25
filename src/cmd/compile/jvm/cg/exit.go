package cg

type Exit struct {
	CurrentCodeLength int
	BranchBytes       []byte // [2]byte
}

func (exit *Exit) Init(op byte, code *AttributeCode) *Exit {
	exit.CurrentCodeLength = code.CodeLength
	code.Codes[code.CodeLength] = op
	exit.BranchBytes = code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	return exit
}
