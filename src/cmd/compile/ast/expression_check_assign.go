package ast

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (e *Expression) checkColonAssignExpression(block *Block, errs *[]error) {
	bin := e.Data.(*ExpressionBinary)
	var names []*Expression
	if bin.Left.Type == ExpressionTypeList {
		names = bin.Left.Data.([]*Expression)
	} else {
		names = []*Expression{bin.Left}
	}
	noErr := true
	values := bin.Right.Data.([]*Expression)
	assignTypes := checkExpressions(block, values, errs)
	if len(names) > len(assignTypes) {
		pos := e.Pos
		getLastPosFromArgs(assignTypes, &pos)
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d values to %d destinations",
			errMsgPrefix(pos),
			len(assignTypes),
			len(names)))
		noErr = false
	} else if len(names) < len(assignTypes) {
		pos := e.Pos
		getFirstPosFromArgs(assignTypes[len(names):], &pos)
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d values to %d destinations",
			errMsgPrefix(pos),
			len(assignTypes),
			len(names)))
	}
	var err error
	noNewVariable := true
	declareVariableExpression := &ExpressionDeclareVariable{}
	declareVariableExpression.InitValues = values
	for k, v := range names {
		if v.Type != ExpressionTypeIdentifier {
			*errs = append(*errs, fmt.Errorf("%s not a name on the left,but '%s'",
				errMsgPrefix(v.Pos), v.OpName()))
			noErr = false
			continue
		}
		identifier := v.Data.(*ExpressionIdentifier)
		if identifier.Name == NoNameIdentifier {
			vd := &Variable{}
			vd.Name = identifier.Name
			declareVariableExpression.Variables = append(declareVariableExpression.Variables, vd)
			declareVariableExpression.IfDeclaredBefore = append(declareVariableExpression.IfDeclaredBefore, false)
			continue
		}
		var variableType *Type
		if k < len(assignTypes) {
			variableType = assignTypes[k]
		}
		if variable, ok := block.Variables[identifier.Name]; ok {
			if variableType != nil {
				if variable.Type.Equal(errs, assignTypes[k]) == false {
					*errs = append(*errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
						errMsgPrefix(assignTypes[k].Pos),
						variable.Type.TypeString(),
						assignTypes[k].TypeString()))
					noErr = false
				}
			}
			identifier.Variable = variable
			declareVariableExpression.Variables = append(declareVariableExpression.Variables, variable)
			declareVariableExpression.IfDeclaredBefore = append(declareVariableExpression.IfDeclaredBefore, true)
		} else { // should be no error
			noNewVariable = false
			vd := &Variable{}
			if k < len(assignTypes) {
				vd.Type = assignTypes[k]
			}
			vd.Name = identifier.Name
			vd.Pos = v.Pos
			if variableType != nil {
				vd.Type = variableType.Clone()
				vd.Type.Pos = e.Pos
			} else {
				vd.Type = &Type{}
				vd.Type.Type = VariableTypeVoid
				vd.Type.Pos = v.Pos
			}
			if vd.Type.isTyped() == false {
				*errs = append(*errs, fmt.Errorf("%s '%s' init value not typed",
					errMsgPrefix(v.Pos), identifier.Name))
			}
			err = block.Insert(vd.Name, v.Pos, vd)
			identifier.Variable = vd
			if err != nil {
				*errs = append(*errs, err)
				noErr = false
				continue
			}
			declareVariableExpression.Variables = append(declareVariableExpression.Variables, vd)
			declareVariableExpression.IfDeclaredBefore = append(declareVariableExpression.IfDeclaredBefore, false)
			if e.IsPublic { // only use when is is global
				vd.AccessFlags |= cg.ACC_FIELD_PUBLIC
			}
		}
	}
	if noNewVariable {
		*errs = append(*errs, fmt.Errorf("%s no new variables to create", errMsgPrefix(e.Pos)))
		noErr = false
	}
	if noErr == false {
		return
	}
	// no error,rewrite data
	e.Data = declareVariableExpression
}

func (e *Expression) checkOpAssignExpression(block *Block, errs *[]error) (t *Type) {
	bin := e.Data.(*ExpressionBinary)
	left := bin.Left.getLeftValue(block, errs)
	bin.Left.Value = left
	right, es := bin.Right.checkSingleValueContextExpression(block)
	if esNotEmpty(es) {
		*errs = append(*errs, es...)
	}
	if left == nil || right == nil {
		return
	}
	ret := left.Clone()
	ret.Pos = e.Pos
	if right.RightValueValid() == false {
		*errs = append(*errs, fmt.Errorf("%s '%s' is not right value valid",
			errMsgPrefix(bin.Right.Pos), right.TypeString()))
		return ret
	}
	/*
		var  s string;
		s += "11111111";
	*/
	if left.Type == VariableTypeString {
		if right.Type != VariableTypeString || (e.Type != ExpressionTypePlusAssign) {
			*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on string and '%s'",
				errMsgPrefix(e.Pos),
				e.OpName(),
				right.TypeString()))
		}
		return ret
	}
	//number
	if e.Type == ExpressionTypePlusAssign ||
		e.Type == ExpressionTypeMinusAssign ||
		e.Type == ExpressionTypeMulAssign ||
		e.Type == ExpressionTypeDivAssign ||
		e.Type == ExpressionTypeModAssign {
		if left.Equal(errs, right) {
			return ret
		}
		if left.IsInteger() && right.IsInteger() && bin.Right.IsLiteral() {
			bin.Right.ConvertToNumber(left.Type)
			return ret
		}
		if left.IsFloat() && right.IsFloat() && bin.Right.IsLiteral() {
			bin.Right.ConvertToNumber(left.Type)
			return ret
		}

	}
	if e.Type == ExpressionTypeAndAssign ||
		e.Type == ExpressionTypeOrAssign ||
		e.Type == ExpressionTypeXorAssign {
		if left.IsInteger() && left.Equal(errs, right) {
			return ret
		}
	}
	if e.Type == ExpressionTypeLshAssign ||
		e.Type == ExpressionTypeRshAssign {
		if left.IsInteger() && right.IsInteger() {
			if right.Type == VariableTypeLong {
				bin.Right.ConvertToNumber(VariableTypeInt)
			}
			return ret
		}
	}

	*errs = append(*errs, fmt.Errorf("%s cannot apply algorithm '%s' on '%s' and '%s'",
		errMsgPrefix(e.Pos),
		e.OpName(),
		left.TypeString(),
		right.TypeString()))

	return ret
}

func (e *Expression) checkAssignExpression(block *Block, errs *[]error) *Type {
	bin := e.Data.(*ExpressionBinary)
	lefts := make([]*Expression, 1)
	if bin.Left.Type == ExpressionTypeList {
		lefts = bin.Left.Data.([]*Expression)
	} else {
		lefts[0] = bin.Left
		bin.Left = &Expression{}
		bin.Left.Type = ExpressionTypeList
		bin.Left.Data = lefts // rewrite to list anyway
	}
	values := bin.Right.Data.([]*Expression)
	valueTypes := checkExpressions(block, values, errs)
	leftTypes := []*Type{}
	for _, v := range lefts {
		if v.Type == ExpressionTypeIdentifier {
			name := v.Data.(*ExpressionIdentifier)
			if name.Name == NoNameIdentifier { // skip "_"
				leftTypes = append(leftTypes, nil) // this is no assign situation
				continue
			}
		}
		t := v.getLeftValue(block, errs)
		v.Value = t
		leftTypes = append(leftTypes, t) // append even if it`s nil
	}
	convertLiteralExpressionsToNeeds(values, leftTypes, valueTypes)
	bin.Left.MultiValues = leftTypes
	if len(lefts) > len(valueTypes) { //expression length compare with value types is more appropriate
		pos := e.Pos
		getLastPosFromArgs(valueTypes, &pos)
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d value to %d detinations",
			errMsgPrefix(pos),
			len(valueTypes),
			len(lefts)))
	} else if len(lefts) < len(valueTypes) {
		pos := e.Pos
		getFirstPosFromArgs(valueTypes[len(lefts):], &pos)
		*errs = append(*errs, fmt.Errorf("%s cannot assign %d value to %d detinations",
			errMsgPrefix(pos),
			len(valueTypes),
			len(lefts)))
	}
	for k, v := range leftTypes {
		if v == nil { // get left value error or "_"
			continue
		}
		if k >= len(valueTypes) {
			continue
		}
		if valueTypes[k] == nil {
			continue
		}
		if false == leftTypes[k].Equal(errs, valueTypes[k]) {
			*errs = append(*errs, fmt.Errorf("%s cannot assign '%s' to '%s'",
				errMsgPrefix(e.Pos),
				valueTypes[k].TypeString(), leftTypes[k].TypeString()))
		}
	}
	voidReturn := mkVoidType(e.Pos)
	if len(leftTypes) > 1 || len(leftTypes) == 0 {
		return voidReturn
	}
	if leftTypes[0] == nil {
		return voidReturn
	}
	// here is safe
	tt := leftTypes[0].Clone()
	tt.Pos = e.Pos
	return tt
}
