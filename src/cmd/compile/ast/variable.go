package ast

import "github.com/756445638/lucy/src/cmd/compile/jvm/class_json"

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
	Signature  *class_json.FieldSignature
	Descriptor string
}

func (v *VariableDefinition) mkTypRight() {
	if v.Typ.isPrimitive() {
		return
	}
	switch v.Typ.Typ {
	case VARIABLE_TYPE_CLASS:
		v.Typ.Typ = VARIABLE_TYPE_OBJECT
	case VARIABLE_TYPE_ARRAY:
		v.Typ.Typ = VARIABLE_TYPE_ARRAY_INSTANCE
	default:
		panic("......")
	}

}
