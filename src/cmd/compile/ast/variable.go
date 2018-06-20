package ast

type Variable struct {
	LocalValOffset           uint16
	IsGlobal                 bool
	IsFunctionParameter      bool
	IsFunctionReturnVariable bool
	BeenCaptured             bool
	Used                     bool   // use as right value
	AccessFlags              uint16 // public private or protected
	Pos                      *Position
	Expression               *Expression
	Name                     string
	Type                     *Type
	JvmDescriptor            string // jvm
}
