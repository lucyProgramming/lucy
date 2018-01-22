package ast

import (
	"fmt"
)

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

func mkSignatureByVariableTypes(ts []*VariableType) string {
	s := ""
	for _, v := range ts {
		s += v.Descriptor()
	}
	return s
}

func mkBuildinFunction(name string, args []*VariableDefinition, rs []*VariableDefinition, checker CallChecker) *Function {
	f := &Function{}
	f.Isbuildin = true
	f.callchecker = checker
	f.Typ = &FunctionType{}
	f.Typ.ParameterList = args
	f.Typ.ReturnList = rs
	f.mkVariableType()
	return f
}

func oneAnyTypeParameterChecker(errs *[]error, args []*VariableType, pos *Pos) {
	if len(args) != 1 {
		*errs = append(*errs, fmt.Errorf("%s only expect one argument", errMsgPrefix(pos)))
	}

}

func devisionByZeroErr(pos *Pos) error {
	return fmt.Errorf("%s division by zero", errMsgPrefix(pos))
}
