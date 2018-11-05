package cg

type Exit struct {
	CurrentCodeLength int
	BranchBytes       []byte //[2]byte
}

func (this *Exit) Init(op byte, code *AttributeCode) *Exit {
	this.CurrentCodeLength = code.CodeLength
	code.Codes[code.CodeLength] = op
	this.BranchBytes = code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	return this
}
