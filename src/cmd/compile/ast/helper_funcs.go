package ast

import (
	"fmt"
	"strings"
)

func divisionByZeroErr(pos *Pos) error {
	return fmt.Errorf("%s division by zero", pos.ErrMsgPrefix())
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
				if err := t.rightValueValid(); err != nil {
					*errs = append(*errs, err)
				}
			}
			ret = append(ret, ts...)
		}
	}
	return ret
}

func getExtraExpressionPos(args []*Expression, n int) *Pos {
	i := 0
	for _, a := range args {
		if a.HaveMultiValue() {
			i += len(a.MultiValues)
		} else {
			i++
		}
		if i >= n {
			return a.Pos
		}
	}
	return nil
}

func mkVoidType(pos *Pos) *Type {
	result := &Type{}
	result.Type = VariableTypeVoid // means no return;
	result.Pos = pos
	return result
}

/*
	when access from global,should check if access from package
*/
func shouldAccessFromImports(name string, from *Pos, alreadyHave *Pos) (*Import, bool) {
	//fmt.Println(name, from, alreadyHave)
	// different file
	// should access from import
	if from.Filename != alreadyHave.Filename {
		i := PackageBeenCompile.getImport(from.Filename, name)
		if i != nil {
			i.Used = true
			return i, true
		} else {
			return nil, false
		}
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
	should := from.Line < alreadyHave.Line
	if should {
		i.Used = true
	}
	return i, should
}

func methodsNotMatchError(pos *Pos, name string, ms []*ClassMethod, want []*Type) error {
	if len(ms) == 0 {
		return fmt.Errorf("%s method '%s' not found", pos.ErrMsgPrefix(), name)
	}
	var errMsg string
	if len(ms) == 1 {
		errMsg = fmt.Sprintf("%s cannot call method '%s':\n",
			pos.ErrMsgPrefix(), name)
	} else {
		errMsg = fmt.Sprintf("%s method named '%s' have no suitable match:\n",
			pos.ErrMsgPrefix(), name)
	}
	wantString := "fn " + name + " ("
	for k, v := range want {
		if v == nil {
			continue
		}
		wantString += v.TypeString()
		if k != len(want)-1 {
			wantString += ","
		}
	}
	wantString += ")"
	errMsg += "\twant " + wantString + "\n"
	for _, m := range ms {
		errMsg += "\thave fn " + name + " " + m.Function.Type.TypeString() + "\n"
	}
	return fmt.Errorf(errMsg)
}

func searchBuildIns(name string) interface{} {
	var t interface{}
	var ok bool
	t, ok = buildInFunctionsMap[name]
	if ok {
		//check
		if _, exists := lucyBuildInPackage.Block.NameExists(name); exists {
			panic(fmt.Sprintf("%s both exits in 'build function' and 'core package'",
				name))
		}
		return t
	}
	if lucyBuildInPackage != nil { // avoid lucy/lang package
		t, _ = lucyBuildInPackage.Block.NameExists(name)
		return t
	}
	return nil
}

func checkConst(block *Block, c *Constant) error {
	if c.Type != nil {
		c.mkDefaultValue()
	}
	if c.DefaultValueExpression == nil {
		err := fmt.Errorf("%s const have no expression", errMsgPrefix(c.Pos))
		return err
	}
	is, err := c.DefaultValueExpression.constantFold()
	if err != nil {
		return err
	}
	if is == false {
		err := fmt.Errorf("%s const named '%s' is not defined by const value",
			c.Pos.ErrMsgPrefix(), c.Name)
		return err
	}
	c.Value = c.DefaultValueExpression.Data
	t, _ := c.DefaultValueExpression.checkSingleValueContextExpression(block)
	if c.Type != nil {
		es := []error{}
		if c.Type.assignAble(&es, t) == false {
			if (c.Type.isInteger() && t.isInteger()) ||
				(c.Type.isFloat() && t.isFloat()) {
				fmt.Println(c.DefaultValueExpression.Data)
				c.DefaultValueExpression.convertLiteralToNumberType(c.Type.Type)
				c.Value = c.DefaultValueExpression.Data

			} else {
				err := fmt.Errorf("%s cannot use '%s' as '%s' for initialization value",
					c.Pos.ErrMsgPrefix(), c.Type.TypeString(), t.TypeString())
				return err
			}
		}
	} else { // means use old type
		c.Type = t
	}
	return nil
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
		if e.isLiteral() == false {
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

		if needs[k].assignAble(&errs, eval[k]) {
			continue // no need
		}
		if (needs[k].isInteger() && eval[k].isInteger()) ||
			(needs[k].isFloat() && eval[k].isFloat()) {
			pos := eval[k].Pos // keep pos
			e.convertToNumberType(needs[k].Type)
			eval[k] = e.Value
			eval[k].Pos = pos
		}
	}
	return
}

// check out package name is valid or not
func PackageNameIsValid(name string) bool {
	if strings.HasPrefix(name, `/`) || strings.HasSuffix(name, `/`) {
		return false
	}
	t := strings.Split(name, `/`)
	if len(t) == 1 {
		return true
	}
	for _, v := range t {
		allOK := true
		for _, vv := range []byte(v) {
			if (vv >= '0' && vv <= '9') ||
				(vv >= 'a' && vv <= 'z') ||
				(vv >= 'A' && vv <= 'Z') ||
				vv == '$' ||
				vv == '_' {
			} else {
				allOK = false
				break
			}
		}
		if allOK == false {
			return false
		}
	}
	return true
}
