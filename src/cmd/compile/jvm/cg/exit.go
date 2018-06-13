package cg

type Exit struct {
	CurrentCodeLength int
	BranchBytes       []byte
}

func (j *Exit) FromCode(op byte, code *AttributeCode) *Exit {
	j.CurrentCodeLength = code.CodeLength
	code.Codes[code.CodeLength] = op
	j.BranchBytes = code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	return j
}
