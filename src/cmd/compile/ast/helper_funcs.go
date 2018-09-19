package ast

import (
	"fmt"
	"strings"
)

func errMsgPrefix(pos *Pos) string {
	return fmt.Sprintf("%s:%d:%d", pos.Filename, pos.StartLine, pos.StartColumn)
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
		}
		return i, i != nil
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

func methodsNotMatchError(pos *Pos, name string, ms []*ClassMethod, want []*Type) error {
	if len(ms) == 0 {
		return fmt.Errorf("%s method '%s' not found", errMsgPrefix(pos), name)
	}
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
		if c.Type.assignAble(errs, t) == false {
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
		if needs[k].assignAble(&errs, eval[k]) {
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

// check out package name is valid or not
func PackageNameIsValid(name string) bool {
	t := strings.Split(name, `/`)
	if len(t) == 1 {
		return true
	}
	if t[0] == "" || t[1] == "" {
		return false
	}
	for _, v := range t {
		allOK := true
		for _, vv := range []byte(v) {
			if (vv >= '0' && vv <= '9') ||
				(vv >= 'a' && vv <= 'z') ||
				(vv >= 'A' && vv <= 'Z') ||
				vv == '$' {
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
