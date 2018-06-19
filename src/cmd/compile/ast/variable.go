package ast

type VariableDefinition struct {
	LocalValOffset           uint16
	IsGlobal                 bool
	IsFunctionParameter      bool
	IsFunctionReturnVariable bool
	BeenCaptured             bool
	Used                     bool   // use as right value
	AccessFlags              uint16 // public private or protected
	Pos                      *Pos
	Expression               *Expression
	Name                     string
	Type                     *VariableType
	JvmDescriptor            string // jvm
}
