package ast

import "github.com/756445638/lucy/src/cmd/compile/jvm/class_json"

type VariableDefinition struct {
	LocalValOffset      byte
	IsGlobal            bool
	BeenCaptured        uint8
	isFunctionParameter bool
	Used                bool   // use as right value
	AccessFlags         uint16 // public private or protected
	Pos                 *Pos
	Expression          *Expression
	NameWithType

	Signature *class_json.FieldSignature
}

type Const struct {
	VariableDefinition
	Value interface{}
}
