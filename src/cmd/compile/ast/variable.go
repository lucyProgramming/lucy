package ast

type VariableDefinition struct {
	LocalValOffset      uint16
	IsGlobal            bool
	BeenCaptured        bool
	IsFunctionParameter bool
	Used                bool   // use as right value
	AccessFlags         uint16 // public private or protected
	Pos                 *Pos
	Expression          *Expression
	Name                string
	Typ                 *VariableType
	Descriptor          string
}
