package ast

import (
	"fmt"
)

func errMsgPrefix(pos *Pos) string {
	return fmt.Sprintf("%s:%d:%d", pos.Filename, pos.StartLine, pos.StartColumn)
}

func esNotEmpty(es []error) bool {
	return len(es) > 0
}

func divisionByZeroErr(pos *Pos) error {
	return fmt.Errorf("%s division by zero", errMsgPrefix(pos))
}

func checkExpressions(block *Block, es []*Expression, errs *[]error, singleValueContext bool) []*Type {
	ret := []*Type{}
	for _, v := range es {
		ts, es := v.check(block)
		*errs = append(*errs, es...)
		if ts == nil {
			ret = append(ret, nil)
		} else {
			if len(ts) > 1 && singleValueContext {
				*errs = append(*errs, fmt.Errorf("%s multi value in single value context",
					errMsgPrefix(v.Pos)))
			}
			for _, t := range ts {
				if t == nil {
					continue
				}
				if false == t.RightValueValid() {
					*errs = append(*errs, fmt.Errorf("%s '%s' cannot used as right value",
						errMsgPrefix(t.Pos), t.TypeString()))
				}
			}
			ret = append(ret, ts...)
		}
	}
	return ret
}

func getFirstPosFromArgs(args []*Type, pos **Pos) {
	for _, a := range args {
		if a != nil {
			*pos = a.Pos
			break
		}
	}
}

func getLastPosFromArgs(args []*Type, pos **Pos) {
	index := len(args) - 1
	for index >= 0 {
		if args[index] != nil {
			*pos = args[index].Pos
			break
		}
		index--
	}
}

func mkVoidType(pos *Pos) *Type {
	t := &Type{}
	t.Type = VariableTypeVoid // means no return;
	t.Pos = pos
	return t
}

/*
	when access from global,should check if access from package
*/
func shouldAccessFromImports(name string, from *Pos, alreadyHave *Pos) (*Import, bool) {
	//fmt.Println(name, from, alreadyHave)
	// different file
	if from.Filename != alreadyHave.Filename {
		i := PackageBeenCompile.getImport(from.Filename, name)
		should := i != nil
		if should {
			i.Used = true
		}
		return i, should
	}
	i := PackageBeenCompile.getImport(from.Filename, name)
	if i == nil {
		return nil, false
	}
	// this is should
	/*
		import
		from
		alreadyHave
	*/
	should := from.StartLine < alreadyHave.StartLine
	if should {
		i.Used = true
	}
	return i, should
}

func msNotMatchError(pos *Pos, name string, ms []*ClassMethod, want []*Type) error {
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
	t, _ := c.Expression.checkSingleValueContextExpression(block)
	if c.Type != nil {
		if c.Type.Equal(errs, t) == false {
			err := fmt.Errorf("%s cannot use '%s' as '%s' for initialization value",
				errMsgPrefix(c.Pos), c.Type.TypeString(), t.TypeString())
			*errs = append(*errs, err)
			return err
		}
	} else { // means use old type
		c.Type = t
	}
	return nil
}

func callWant(ft *FunctionType) string {
	s := "("
	for k, v := range ft.ParameterList {
		s += v.Name + " "
		s += v.Type.TypeString()
		if k != len(ft.ParameterList)-1 {
			s += ","
		}
	}
	if ft.VArgs != nil {
		if len(ft.ParameterList) > 0 {
			s += ","
		}
		s += ft.VArgs.Name + " "
		s += ft.VArgs.Type.TypeString()
	}
	s += ")"
	return s
}

func callHave(ts []*Type) string {
	s := "("
	for k, v := range ts {
		if v == nil {
			continue
		}
		if v.Name != "" {
			s += v.Name + " "
		}
		s += v.TypeString()
		if k != len(ts)-1 {
			s += ","
		}
	}
	s += ")"
	return s
}
func convertExpressionToNeed(e *Expression, need *Type, eval *Type) {
	convertExpressionsToNeeds([]*Expression{e}, []*Type{need}, []*Type{eval})
}
func convertExpressionsToNeeds(es []*Expression, needs []*Type, eval []*Type) {
	errs := []error{} // no return
	if len(es) == 0 {
		return
	}
	for k, e := range es {
		if e.IsLiteral() == false {
			continue
		}
		if k >= len(needs) {
			break
		}
		if needs[k] == nil {
			continue
		}
		if eval[k] == nil {
			continue
		}
		if needs[k].Equal(&errs, eval[k]) {
			continue // no need
		}
		if (needs[k].IsInteger() && eval[k].IsInteger()) ||
			(needs[k].IsFloat() && eval[k].IsFloat()) {
			pos := eval[k].Pos // keep pos
			e.convertNumberLiteralTo(needs[k].Type)
			*eval[k] = *needs[k]
			eval[k].Pos = pos
		}
	}
	return
}
