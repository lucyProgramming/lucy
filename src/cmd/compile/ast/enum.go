package ast

type EnumName struct {
	Enum  *Enum
	Name  string
	Pos   *Pos
	Value int64
}

type Enum struct {
	Access       uint16 // public private or protected
	Name         string
	Pos          *Pos
	Names        []*EnumName
	NamesMap     map[string]*EnumName
	Init         *Expression //should be a int expression
	Used         bool
	VariableType *VariableType
}

func (e *Enum) mkVariableType() {
	e.VariableType = &VariableType{}
	e.VariableType.Typ = VARIABLE_TYPE_ENUM
	e.VariableType.Resource = &VariableTypeResource{}
	e.VariableType.Resource.Enum = e
}
