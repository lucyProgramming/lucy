package ast

import (
	"fmt"
)

func errMsgPrefix(pos *Position) string {
	return fmt.Sprintf("%s:%d:%d", pos.Filename, pos.StartLine, pos.StartColumn)
}

func errorsNotEmpty(es []error) bool {
	return es != nil && len(es) > 0
}

func checkEnum(enums []*Enum) []error {
	ret := make([]error, 0)
	for _, v := range enums {
		if len(v.Enums) == 0 {
			continue
		}
		err := v.check()
		if err != nil {
			ret = append(ret, err)
		}
	}
	return ret
}

func divisionByZeroErr(pos *Position) error {
	return fmt.Errorf("%s division by zero", errMsgPrefix(pos))
}

func checkExpressions(block *Block, es []*Expression, errs *[]error) []*Type {
	ret := []*Type{}
	for _, v := range es {
		ts, e := v.check(block)
		if errorsNotEmpty(e) {
			*errs = append(*errs, e...)
		}
		if ts != nil {
			ret = append(ret, ts...)
		} else {
			ret = append(ret, nil)
		}
	}
	return ret
}

func mkVoidType(pos *Position) *Type {
	t := &Type{}
	t.Type = VARIABLE_TYPE_VOID // means no return;
	t.Pos = pos
	return t
}

func checkRightValuesValid(ts []*Type, errs *[]error) (ret []*Type) {
	ret = []*Type{}
	for _, v := range ts {
		ret = append(ret, v)
		if v == nil {
			continue
		}
		if false == v.RightValueValid() {
			*errs = append(*errs, fmt.Errorf("%s '%s' cannot used as right value",
				errMsgPrefix(v.Pos), v.TypeString()))
		}
	}
	return ret
}

/*
	when access from global,should check if access from package
*/
func shouldAccessFromImports(name string, from *Position, have *Position) (*Import, bool) {
	if have == nil { // in case buildIn types
		return nil, false
	}
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

func msNotMatchError(pos *Position, name string, ms []*ClassMethod, want []*Type) error {
	errMsg := fmt.Sprintf("%s method named '%s' have no suitable match:\n",
		errMsgPrefix(pos), name)
	errMsg += "\twant " + ms[0].Function.badParameterMsg(name, want) + "\n"
	for _, m := range ms {
		errMsg += "\thave " + m.Function.readableMsg(name) + "\n"
	}
	return fmt.Errorf(errMsg)
}

func searchBuildIns(name string) interface{} {
	var t interface{}
	var ok bool
	t, ok = buildInFunctionsMap[name]
	if ok {
		return t
	}
	if lucyBuildInPackage != nil {
		t, _ = lucyBuildInPackage.Block.NameExists(name)
		return t
	}
	return nil
}

func checkConst(block *Block, c *Constant, errs *[]error) error {
	if c.Type != nil {
		c.mkDefaultValue()
	}
	if c.Expression == nil {
		err := fmt.Errorf("%s const have no expression", errMsgPrefix(c.Pos))
		*errs = append(*errs, err)
		return err
	}
	is, err := c.Expression.constantFold()
	if err != nil {
		*errs = append(*errs, err)
		return err
	}
	if is == false {
		err := fmt.Errorf("%s const named '%s' is not defined by const value",
			errMsgPrefix(c.Pos), c.Name)
		*errs = append(*errs, err)
		return err
	}
	c.Value = c.Expression.Data
	tt, _ := c.Expression.check(block)
	if c.Type != nil {
		if c.Type.Equal(errs, tt[0]) == false {
			err := fmt.Errorf("%s cannot use '%s' as '%s' for initialization value",
				errMsgPrefix(c.Pos), c.Type.TypeString(), tt[0].TypeString())
			*errs = append(*errs, err)
			return err
		}
	} else { // means use old typec
		c.Type = tt[0]
	}
	return nil
}

func functionPointerCallWant(ts ParameterList) string {
	s := "("
	for k, v := range ts {
		s += " " + v.Name + " "
		s += v.Type.TypeString()
		s += " "
		if k != len(ts)-1 {
			s += ","
		}
	}
	s += ")"
	return s
}

func functionPointerCallHave(ts []*Type) string {
	s := "("
	for k, v := range ts {
		s += " " + v.Name + " "
		s += v.TypeString()
		if k != len(ts)-1 {
			s += ","
		}
	}
	s += ")"
	return s
}
func convertLiteralExpressionsToNeeds(es []*Expression, needs []*Type, checked []*Type) []error {
	errs := []error{}
	if len(es) == 0 {
		return errs
	}
	if len(es) != len(checked) || len(es) != len(needs) { // means multi return
		return errs
	}
	for k, e := range es {
		if e.IsLiteral() == false {
			continue
		}
		if needs[k] == nil {
			continue
		}
		if checked[k] == nil {
			continue
		}
		if needs[k].Equal(&errs, checked[k]) {
			continue // no need
		}
		if (needs[k].IsInteger() && checked[k].IsInteger()) ||
			(needs[k].IsFloat() && checked[k].IsFloat()) {
			pos := checked[k].Pos // keep pos
			e.convertNumberLiteralTo(needs[k].Type)
			*checked[k] = *needs[k]
			checked[k].Pos = pos
		}
	}
	return errs
}
