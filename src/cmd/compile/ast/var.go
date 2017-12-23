package ast

import (
	"fmt"
	"regexp"
)

type PackageLoader interface {
	LoadPackage(name string) (*Package, error)
}

var (
	THIS                 = "this"
	small_float          = 0.0001
	negative_small_float = -small_float
	Nodes                *[]*Node //
	packageAliasReg      *regexp.Regexp
	PackageLoad          PackageLoader
	JAVA_ROOT_CLASS      = "java/lang/Object"
)

type NameWithPos struct {
	Name string
	Pos  *Pos
}

func init() {
	var err error
	packageAliasReg, err = regexp.Compile(`^[a-zA-Z][[a-zA-Z1-9\_]+$`)
	if err != nil {
		panic(err)
	}
}

func notFoundError(pos *Pos, typ, name string) error {
	return fmt.Errorf("%s %s named %s not found", errMsgPrefix(pos), typ, name)
}

func errMsgPrefix(pos *Pos) string {
	return fmt.Sprintf("%s:%d:%d", pos.Filename, pos.StartLine, pos.StartColumn)
}

func errsNotEmpty(errs []error) bool {
	return errs != nil && len(errs) > 0
}
func checkEnum(enums []*Enum) []error {
	ret := make([]error, 0)
	for _, v := range enums {
		if len(v.Names) == 0 {
			continue
		}
		is, typ, value, err := v.Init.getConstValue()
		if err != nil || is == false || typ != EXPRESSION_TYPE_INT {
			ret = append(ret, fmt.Errorf("enum type must inited by integer"))
			continue
		}
		for k, vv := range v.Names {
			vv.Value = int64(k) + value.(int64)
		}
	}
	return ret
}

func mkVoidVariableTypes(pos *Pos, length ...int) []*VariableType {
	l := 1
	if len(length) > 0 && length[0] > 0 {
		l = length[0]
	}
	ret := make([]*VariableType, l)
	for k := range ret {
		ret[k] = mkVoidVariableType(pos)
	}
	return ret
}
func mkVoidVariableType(pos *Pos) *VariableType {
	return &VariableType{
		Typ: VARIABLE_TYPE_VOID,
		Pos: pos,
	}
}

func mkSignatureByVariableTypes(ts []*VariableType) string {
	s := ""
	for _, v := range ts {
		s += v.Descriptor()
	}
	return s
}

func mkBuildFunction(name string, IsAnyNumberParameter bool, args []*VariableDefinition, rs []*VariableDefinition) *Function {
	f := &Function{}
	f.Isbuildin = true
	f.Used = true
	f.IsAnyNumberParameter = IsAnyNumberParameter
	f.Typ = &FunctionType{}
	f.Typ.Parameters = args
	f.Typ.Returns = rs
	f.mkVariableType()
	return f
}
