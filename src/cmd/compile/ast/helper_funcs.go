package ast

import (
	"fmt"
)

func errMsgPrefix(pos *Pos) string {
	return fmt.Sprintf("%s:%d:%d", pos.Filename, pos.StartLine, pos.StartColumn)
}
func ErrMsgPrefix(pos *Pos) string {
	return errMsgPrefix(pos)
}
func errsNotEmpty(es []error) bool {
	return es != nil && len(es) > 0
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

func mkBuildinFunction(name string, args []*VariableDefinition, rs []*VariableDefinition, checker CallChecker) *Function {
	f := &Function{}
	f.Name = name
	f.IsBuildin = true
	f.callchecker = checker
	f.Typ = &FunctionType{}
	f.Typ.ParameterList = args
	f.Typ.ReturnList = rs
	f.IsGlobal = true
	return f
}

func oneAnyTypeParameterChecker(block *Block, errs *[]error, args []*VariableType, pos *Pos) {
	if len(args) != 1 {
		*errs = append(*errs, fmt.Errorf("%s only expect one argument", errMsgPrefix(pos)))
	}
}

func devisionByZeroErr(pos *Pos) error {
	return fmt.Errorf("%s division by zero", errMsgPrefix(pos))
}

func checkExpressions(block *Block, es []*Expression, errs *[]error) []*VariableType {
	ret := []*VariableType{}
	for _, v := range es {
		ts, e := v.check(block)
		if errsNotEmpty(e) {
			*errs = append(*errs, e...)
		}
		if ts != nil {
			ret = append(ret, ts...)
		}
	}
	return ret
}

func mkVoidType(pos *Pos) *VariableType {
	t := &VariableType{}
	t.Typ = VARIABLE_TYPE_VOID // means no return;
	t.Pos = pos
	return t
}

func checkRightValuesValid(ts []*VariableType, errs *[]error) (ret []*VariableType) {
	ret = []*VariableType{}
	for _, v := range ts {
		if !v.rightValueValid() {
			*errs = append(*errs, fmt.Errorf("%s '%s' cannot used as right value", errMsgPrefix(v.Pos), v.TypeString()))
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

/*
	when access from global,should check if access from package
*/
func shouldAccessFromImports(name string, from *Pos, have *Pos) (*Import, bool) {
	// different file
	if from.Filename != have.Filename {
		i := PackageBeenCompile.getImport(from.Filename, name)
		return i, i != nil
	}
	i := PackageBeenCompile.getImport(from.Filename, name)
	if i == nil {
		return nil, false
	}
	return i, have.StartLine < from.StartLine
}

func moreClose(from, more, less *Pos) bool {
	if from.Filename == more.Filename && from.Filename != less.Filename {
		return true
	}
	return more.StartLine < less.StartLine
}

func msNotMatchError(pos *Pos, name string, ms []*ClassMethod, want []*VariableType) error {
	errmsg := fmt.Sprintf("%s method named '%s' have no suitable match:\n",
		errMsgPrefix(pos), name)
	errmsg += "\t want " + ms[0].Func.badParameterMsg(name, want) + "\n"
	for _, m := range ms {
		errmsg += "\t have " + m.Func.readableMsg(name) + "\n"
	}
	return fmt.Errorf(errmsg)
	//*errs = append(*errs, fmt.Errorf(errmsg))
}
