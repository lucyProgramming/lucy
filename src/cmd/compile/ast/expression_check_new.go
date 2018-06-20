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
	if no.Type.Type == VARIABLE_TYPE_MAP {
		return e.checkNewMapExpression(block, no, errs)
	}
	if no.Type.Type == VARIABLE_TYPE_ARRAY {
		return e.checkNewArrayExpression(block, no, errs)
	}
	if no.Type.Type == VARIABLE_TYPE_JAVA_ARRAY {
		return e.checkNewJavaArrayExpression(block, no, errs)
	}

	// new object
	if no.Type.Type != VARIABLE_TYPE_OBJECT {
		*errs = append(*errs, fmt.Errorf("%s cannot have new on type '%s'",
			errMsgPrefix(e.Pos), no.Type.TypeString()))
		return nil
	}
	if no.Type.Class.IsInterface() {
		*errs = append(*errs, fmt.Errorf("%s '%s' is interface",
			errMsgPrefix(e.Pos), no.Type.Class.Name))
		return nil
	}
	ret := &Type{}
	*ret = *no.Type
	ret.Type = VARIABLE_TYPE_OBJECT
	ret.Pos = e.Pos
	args := checkExpressions(block, no.Args, errs)
	ms, matched, err := no.Type.Class.matchConstructionFunction(e.Pos, errs, args, &no.Args)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s %v", errMsgPrefix(e.Pos), err))
		return ret
	}
	if matched {
		no.Construction = ms[0]
		return ret
	}
	if len(ms) == 0 {
		*errs = append(*errs, fmt.Errorf("%s  'construction' not found",
			errMsgPrefix(e.Pos)))
	} else {
		*errs = append(*errs, msNotMatchError(e.Pos, "constructor", ms, args))
	}
	return ret
}

func (e *Expression) checkNewMapExpression(block *Block, newMap *ExpressionNew,
	errs *[]error) *Type {
	if len(newMap.Args) > 0 {
		*errs = append(*errs, fmt.Errorf("%s new map expect no arguments",
			errMsgPrefix(newMap.Args[0].Pos)))
	}
	tt := newMap.Type.Clone()
	tt.Pos = e.Pos
	return tt
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
				errMsgPrefix(newArray.Args[0].Pos)))
		newArray.Args = []*Expression{} // reset to 0,continue to analyse
	}
	ts := checkRightValuesValid(checkExpressions(block, newArray.Args, errs), errs)
	amount, err := e.mustBeOneValueContext(ts)
	if err != nil {
		*errs = append(*errs, err)
	}
	if amount == nil {
		return ret
	}
	if amount.Type != VARIABLE_TYPE_INT {
		//if amount.Typ == VARIABLE_TYPE_JAVA_ARRAY && newArray.Typ.ArrayType.Equal(errs, amount.ArrayType) {
		//	//convert java array to array
		//	newArray.IsConvertJavaArray2Array = true
		//} else {
		*errs = append(*errs, fmt.Errorf("%s argument must be 'int',but '%s'",
			errMsgPrefix(amount.Pos), amount.TypeString()))
		//}
	}
	//no further checks
	return ret
}
