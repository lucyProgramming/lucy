package ast

type EnumName struct {
	Enum  *Enum
	Name  string
	Pos   *Pos
	Value int64
}

type Enum struct {
	Access   uint16 // public private or protected
	Name     string
	Pos      *Pos
	Names    []*EnumName
	NamesMap map[string]*EnumName //easy to access by name
	Init     *Expression          //should be a int expression
	Used     bool
}
