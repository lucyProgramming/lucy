package ast

import (
	"fmt"
)

func (this *Expression) checkNewExpression(block *Block, errs *[]error) *Type {
	no := this.Data.(*ExpressionNew)
	err := no.Type.resolve(block)
	if err != nil {
		*errs = append(*errs, err)
		return nil
	}
	if no.Type.Type == VariableTypeMap {
		return this.checkNewMapExpression(block, no, errs)
	}
	if no.Type.Type == VariableTypeArray {
		return this.checkNewArrayExpression(block, no, errs)
	}
	if no.Type.Type == VariableTypeJavaArray {
		return this.checkNewJavaArrayExpression(block, no, errs)
	}
	// new object
	if no.Type.Type != VariableTypeObject {
		*errs = append(*errs,
			fmt.Errorf("%s cannot have new on type '%s'",
				no.Type.Pos.ErrMsgPrefix(), no.Type.TypeString()))
		return nil
	}
	err = no.Type.Class.loadSelf(this.Pos)
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
	ret.Pos = this.Pos
	errsLength := len(*errs)
	callArgTypes := checkExpressions(block, no.Args, errs, true)
	if len(*errs) > errsLength {
		return ret
	}
	ms, matched, err := no.Type.Class.accessConstructionMethod(this.Pos, errs, no, nil, callArgTypes)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %v", this.Pos.ErrMsgPrefix(), err))
		return ret
	}
	if matched {
		m := ms[0]
		if err := no.Type.Class.constructionMethodAccessAble(this.Pos, m); err != nil {
			*errs = append(*errs, err)
		}
		no.Construction = m
		return ret
	}
	*errs = append(*errs, methodsNotMatchError(no.Type.Pos, no.Type.TypeString(), ms, callArgTypes))
	return ret
}

func (this *Expression) checkNewMapExpression(block *Block, newMap *ExpressionNew,
	errs *[]error) *Type {
	if len(newMap.Args) > 0 {
		*errs = append(*errs,
			fmt.Errorf("%s new 'map' expect no arguments",
				errMsgPrefix(newMap.Args[0].Pos)))
	}
	ret := newMap.Type.Clone()
	ret.Pos = this.Pos
	return ret
}

func (this *Expression) checkNewJavaArrayExpression(block *Block, newArray *ExpressionNew,
	errs *[]error) *Type {
	return this.checkNewArrayExpression(block, newArray, errs)
}

func (this *Expression) checkNewArrayExpression(block *Block, newArray *ExpressionNew,
	errs *[]error) *Type {
	ret := newArray.Type.Clone() // clone the type
	ret.Pos = this.Pos
	if len(newArray.Args) != 1 { //
		*errs = append(*errs,
			fmt.Errorf("%s new array expect at least 1 argument",
				errMsgPrefix(this.Pos)))
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
	} else {
		if amount.Type == VariableTypeLong {
			newArray.Args[0].convertToNumberType(VariableTypeLong)
		}
		if newArray.Args[0].isLiteral() {
			if a := newArray.Args[0].getLongValue(); a < 0 {
				*errs = append(*errs,
					fmt.Errorf("%s '%d' is negative ",
						errMsgPrefix(amount.Pos), a))
			}
		}
	}

	//no further checks
	return ret
}
