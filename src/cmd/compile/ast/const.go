package ast

type Constant struct {
	Variable
	Value interface{} // value base on type
}

func (c *Constant) mkDefaultValue() {
	switch c.Type.Type {
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
