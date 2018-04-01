package cg

type JumpBackPatch struct {
	CurrentCodeLength int
	Bs                []byte
}

func (j *JumpBackPatch) FromCode(op byte, code *AttributeCode) *JumpBackPatch {
	j.CurrentCodeLength = code.CodeLength
	code.Codes[code.CodeLength] = op
	j.Bs = code.Codes[code.CodeLength+1 : code.CodeLength+3]
	code.CodeLength += 3
	return j
}
