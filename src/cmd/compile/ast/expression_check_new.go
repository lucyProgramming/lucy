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
		*errs = append(*errs,
			fmt.Errorf("%s cannot have new on type '%s'",
				no.Type.Pos.ErrMsgPrefix(), no.Type.TypeString()))
		return nil
	}
	err = no.Type.Class.loadSelf(e.Pos)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %v",
			no.Type.Pos.ErrMsgPrefix(), err))
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
	ms, matched, err := no.Type.Class.accessConstructionFunction(e.Pos, errs, no, nil, callArgTypes)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %v", e.Pos.ErrMsgPrefix(), err))
		return ret
	}
	if matched {
		m := ms[0]
		if err := no.Type.Class.constructionMethodAccessAble(e.Pos, m); err != nil {
			*errs = append(*errs, err)
		}
		no.Construction = m
		return ret
	}
	*errs = append(*errs, methodsNotMatchError(no.Type.Pos, no.Type.TypeString(), ms, callArgTypes))
	return ret
}

func (e *Expression) checkNewMapExpression(block *Block, newMap *ExpressionNew,
	errs *[]error) *Type {
	if len(newMap.Args) > 0 {
		*errs = append(*errs,
			fmt.Errorf("%s new 'map' expect no arguments",
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
	if len(newArray.Args) != 1 { //
		*errs = append(*errs,
			fmt.Errorf("%s new array expect at least 1 argument",
				errMsgPrefix(e.Pos)))
		return ret
	}
	amount, es := newArray.Args[0].checkSingleValueContextExpression(block)
	if es != nil {
		*errs = append(*errs, es...)
	}
	if amount == nil {
		return ret
	}
	if amount.isInteger() == false {
		*errs = append(*errs,
			fmt.Errorf("%s argument must be 'int',but '%s'",
				errMsgPrefix(amount.Pos), amount.TypeString()))
	}
	if amount.Type == VariableTypeLong {
		newArray.Args[0].convertToNumber(VariableTypeLong)
	}
	//no further checks
	return ret
}
