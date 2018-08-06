package ast

import (
	"fmt"
)

func (e *Expression) checkNewExpression(block *Block, errs *[]error) *Type {
	no := e.Data.(*ExpressionNew)
	err := no.Type.resolve(block)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	if no.Type.Type == VariableTypeMap {
		return e.checkNewMapExpression(block, no, errs)
	}
	if no.Type.Type == VariableTypeArray {
		return e.checkNewArrayExpression(block, no, errs)
	}
	if no.Type.Type == VariableTypeJavaArray {
		return e.checkNewJavaArrayExpression(block, no, errs)
	}
	// new object
	if no.Type.Type != VariableTypeObject {
		*errs = append(*errs, fmt.Errorf("%s cannot have new on type '%s'",
			errMsgPrefix(e.Pos), no.Type.TypeString()))
		return nil
	}
	err = no.Type.Class.loadSelf()
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %v",
			errMsgPrefix(no.Type.Pos), err))
		return nil
	}
	if no.Type.Class.IsInterface() {
		*errs = append(*errs, fmt.Errorf("%s '%s' is interface",
			errMsgPrefix(no.Type.Pos), no.Type.Class.Name))
		return nil
	}
	if no.Type.Class.IsAbstract() {
		*errs = append(*errs, fmt.Errorf("%s '%s' is abstract",
			errMsgPrefix(no.Type.Pos), no.Type.Class.Name))
		return nil
	}
	ret := &Type{}
	*ret = *no.Type
	ret.Type = VariableTypeObject
	ret.Pos = e.Pos
	callArgTypes := checkExpressions(block, no.Args, errs, true)
	ms, matched, err := no.Type.Class.matchConstructionFunction(e.Pos, errs, no, nil, callArgTypes)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(no.Type.Pos), err))
		return ret
	}
	if matched {
		m := ms[0]
		if block.InheritedAttribute.Class != ret.Class {
			if (ret.Class.LoadFromOutSide && m.IsPublic() == false) ||
				(ret.Class.LoadFromOutSide == false && m.IsPrivate() == true) {
				*errs = append(*errs, fmt.Errorf("%s constuction cannot access from here", errMsgPrefix(no.Type.Pos)))
			}
		}
		no.Construction = m
		return ret
	}
	if len(ms) == 0 {
		*errs = append(*errs, fmt.Errorf("%s  'construction' not found",
			errMsgPrefix(e.Pos)))
	} else {
		*errs = append(*errs, msNotMatchError(no.Type.Pos, "constructor", ms, callArgTypes))
	}
	return ret
}

func (e *Expression) checkNewMapExpression(block *Block, newMap *ExpressionNew,
	errs *[]error) *Type {
	if len(newMap.Args) > 0 {
		*errs = append(*errs, fmt.Errorf("%s new 'map' expect no arguments",
			errMsgPrefix(newMap.Args[0].Pos)))
	}
	ret := newMap.Type.Clone()
	ret.Pos = e.Pos
	return ret
}

func (e *Expression) checkNewJavaArrayExpression(block *Block, newArray *ExpressionNew,
	errs *[]error) *Type {
	return e.checkNewArrayExpression(block, newArray, errs)
}

func (e *Expression) checkNewArrayExpression(block *Block, newArray *ExpressionNew,
	errs *[]error) *Type {
	ret := newArray.Type.Clone() // clone the type
	ret.Pos = e.Pos
	if len(newArray.Args) != 1 { // 0 and 1 is accept
		*errs = append(*errs,
			fmt.Errorf("%s new array expect at least 1 argument ",
				errMsgPrefix(e.Pos)))
		newArray.Args = []*Expression{} // reset to 0,continue to analyse
	}
	amount, es := newArray.Args[0].checkSingleValueContextExpression(block)
	if es != nil {
		*errs = append(*errs, es...)
	}
	if amount == nil {
		return ret
	}
	if amount.Type != VariableTypeInt {
		*errs = append(*errs, fmt.Errorf("%s argument must be 'int',but '%s'",
			errMsgPrefix(amount.Pos), amount.TypeString()))
	}
	//no further checks
	return ret
}
