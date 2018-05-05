package ast

type Const struct {
	VariableDefinition
	Value interface{} // value base on type
}

func (c *Const) mkDefaultValue() {
	switch c.Typ.Typ {
	case VARIABLE_TYPE_BOOL:
		c.Value = false
	case VARIABLE_TYPE_BYTE:
		c.Value = byte(0)
	case VARIABLE_TYPE_SHORT:
		c.Value = int32(0)
	case VARIABLE_TYPE_INT:
		c.Value = int32(0)
	case VARIABLE_TYPE_LONG:
		c.Value = int64(0)
	case VARIABLE_TYPE_FLOAT:
		c.Value = float32(0)
	case VARIABLE_TYPE_DOUBLE:
		c.Value = float64(0)
	case VARIABLE_TYPE_STRING:
		c.Value = ""
	}
}
