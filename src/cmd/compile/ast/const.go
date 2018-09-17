package ast

type Constant struct {
	Variable
	Value interface{} // value base on type
}

func (c *Constant) mkDefaultValue() {
	switch c.Type.Type {
	case VariableTypeBool:
		c.Value = false
	case VariableTypeByte:
		c.Value = byte(0)
	case VariableTypeShort:
		c.Value = int32(0)
	case VariableTypeChar:
		c.Value = int32(0)
	case VariableTypeInt:
		c.Value = int32(0)
	case VariableTypeLong:
		c.Value = int64(0)
	case VariableTypeFloat:
		c.Value = float32(0)
	case VariableTypeDouble:
		c.Value = float64(0)
	case VariableTypeString:
		c.Value = ""
	}
}
