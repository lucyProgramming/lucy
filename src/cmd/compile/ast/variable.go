package ast

type Variable struct {
	IsBuildIn                bool
	IsGlobal                 bool
	IsFunctionParameter      bool
	IsFunctionReturnVariable bool
	BeenCaptured             bool
	Used                     bool   // use as right value
	AccessFlags              uint16 // public private or protected
	Pos                      *Pos
	Expression               *Expression
	Name                     string
	Type                     *Type
	LocalValOffset           uint16 // offset in stack frame
	JvmDescriptor            string // jvm
}
