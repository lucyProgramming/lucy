package ast

type VariableDefinition struct {
	LocalValOffset      uint16
	IsGlobal            bool
	BeenCaptured        bool
	CaptureLevel        uint8
	IsFunctionParameter bool
	Used                bool   // use as right value
	AccessFlags         uint16 // public private or protected
	Pos                 *Pos
	Expression          *Expression
	NameWithType
	Descriptor string
}
