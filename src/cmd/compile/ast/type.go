package ast

const (
	VARIABLE_TYPE_BOOL = iota
	VARIABLE_TYPE_BYTE
	VARIABLE_TYPE_INT
	VARIABLE_TYPE_FLOAT
	VARIABLE_TYPE_STRING
	VARIALBE_TYPE_FUNCTION
	VARIALBE_TYPE_ENUM        //enum
	VARIABLE_TYPE_CLASS       //new Person()
	VARIABLE_TYPE_COMBINATION // []int
)

type AccessProperty struct {
	Access int // public private or protected
}
type VariableType struct {
	Typ             int
	Name            string // class name or function name
	CombinationType *CombinationType
}

const (
	COMBINATION_TYPE_ARRAY = iota
)

type CombinationType struct {
	Typ         int
	Combination VariableType
}

//把树型转化为可读字符串
func (c *CombinationType) TypeString(ret *string) {
	if c.Typ == COMBINATION_TYPE_ARRAY {
		*ret += "[]"
	}
	c.Combination.TypeString(ret)
}

//可读的类型信息
func (v *VariableType) TypeString(ret *string) {
	switch v.Typ {
	case VARIABLE_TYPE_BOOL:
		*ret = "bool"
	case VARIABLE_TYPE_BYTE:
		*ret = "byte"
	case VARIABLE_TYPE_INT:
		*ret = "int"
	case VARIABLE_TYPE_FLOAT:
		*ret = "float"
	case VARIALBE_TYPE_FUNCTION:
		*ret = "function"
	case VARIABLE_TYPE_CLASS:
		*ret = v.Name
	case VARIALBE_TYPE_ENUM:
		*ret = v.Name
	case VARIABLE_TYPE_COMBINATION:
		v.CombinationType.TypeString(ret)
	}
}

func (v *VariableType) Equal(e *VariableType) bool {
	if v.Typ != e.Typ {
		return false
	}
	if v.CombinationType.Typ != e.CombinationType.Typ {
		return false
	}
	return v.CombinationType.Combination.Equal(&e.CombinationType.Combination)
}
