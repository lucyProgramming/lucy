package ast

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"

type Constant struct {
	IsGlobal               bool
	IsBuildIn              bool
	Used                   bool
	Pos                    *Pos
	Type                   *Type
	Name                   string
	DefaultValueExpression *Expression
	AccessFlags            uint16
	Comment                string
	Value                  interface{} // value base on type
}

func (this *Constant) isPublic() bool {
	return this.AccessFlags|cg.AccFieldPublic != 0
}

func (this *Constant) mkDefaultValue() {
	switch this.Type.Type {
	case VariableTypeBool:
		this.Value = false
	case VariableTypeByte:
		this.Value = int64(0)
	case VariableTypeShort:
		this.Value = int64(0)
	case VariableTypeChar:
		this.Value = int64(0)
	case VariableTypeInt:
		this.Value = int64(0)
	case VariableTypeLong:
		this.Value = int64(0)
	case VariableTypeFloat:
		this.Value = float32(0)
	case VariableTypeDouble:
		this.Value = float64(0)
	case VariableTypeString:
		this.Value = ""
	}
}
