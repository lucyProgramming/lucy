package jvm

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (m *MakeExpression) buildExpressionAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	//
	bin := e.Data.(*ast.ExpressionBinary)
	left := bin.Left.Data.([]*ast.Expression)[0]
	right := bin.Right.Data.([]*ast.Expression)[0]
	maxstack, remainStack, op, target, classname, name, descriptor := m.getLeftValue(class, code, left, context)
	stack, es := m.build(class, code, right, context)
	backPatchEs(es, code.CodeLength)
	if t := remainStack + stack; t > maxstack {
		maxstack = t
	}
	if target.IsNumber() && target.Typ != right.VariableType.Typ {
		m.numberTypeConverter(code, right.VariableType.Typ, target.Typ)
	}
	currentStack := remainStack + target.JvmSlotSize()
	if currentStack > maxstack {
		maxstack = currentStack
	}
	if t := currentStack + m.controlStack2FitAssign(code, op, classname, target); t > maxstack {
		maxstack = t
	}
	copyOPLeftValue(class, code, op, classname, name, descriptor)
	return
}

// a,b,c = 122,fdfd2232,"hello";
func (m *MakeExpression) buildAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	if e.IsStatementExpression == false {
		return m.buildExpressionAssign(class, code, e, context)
	}
	bin := e.Data.(*ast.ExpressionBinary)
	lefts := bin.Left.Data.([]*ast.Expression)
	currentStack := uint16(0)
	targets := make([]*ast.VariableType, len(lefts))
	ops := make([][]byte, len(lefts))
	classnames := make([]string, len(lefts))
	names := make([]string, len(lefts))
	descriptors := make([]string, len(lefts))
	remainstacks := make([]uint16, len(lefts))
	noDestinations := make([]bool, len(lefts))
	// put left value one the stack
	index := len(lefts) - 1
	for index >= 0 { //
		if lefts[index].Typ == ast.EXPRESSION_TYPE_IDENTIFIER &&
			lefts[index].Data.(*ast.ExpressionIdentifer).Name == ast.NO_NAME_IDENTIFIER {
			noDestinations[index] = true
		} else {
			stack, remainstack, op, target, classname, name, descriptor := m.getLeftValue(class, code, lefts[index], context)
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			targets[index] = target
			ops[index] = op
			classnames[index] = classname
			names[index] = name
			descriptors[index] = descriptor
			remainstacks[index] = remainstack
			currentStack += remainstack
		}
		index--
	}
	fmt.Println(noDestinations)
	//
	rights := bin.Right.Data.([]*ast.Expression)
	slice := func() {
		currentStack -= remainstacks[0]
		if noDestinations[0] == false {
			copyOPLeftValue(class, code, ops[0], classnames[0], names[0], descriptors[0])
		}
		targets = targets[1:] // slice
		ops = ops[1:]
		classnames = classnames[1:]
		names = names[1:]
		descriptors = descriptors[1:]
		remainstacks = remainstacks[1:]
		noDestinations = noDestinations[1:]
	}
	for _, v := range rights {
		if v.IsCall() && len(v.VariableTypes) > 1 {
			stack, _ := m.build(class, code, v, context)
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			// stack top is arrayListObject object
			if t := 1 + currentStack; t > maxstack {
				maxstack = t
			}
			m.buildStoreArrayListAutoVar(code, context) // store it into local
			for k, v := range v.VariableTypes {         // unpack
				mapDestination := (classnames[0] == java_hashmap_class && targets[0].IsPointer() == false)
				stack = m.unPackArraylist(class, code, k, v, context)
				if t := stack + currentStack; t > maxstack {
					maxstack = t
				}
				if noDestinations[0] == false {
					if t := currentStack + targets[0].JvmSlotSize(); t > maxstack { // incase int convert to double or long
						maxstack = t
					}
					if targets[0].IsNumber() && targets[0].Typ != v.Typ { // value is number 2
						m.numberTypeConverter(code, v.Typ, targets[0].Typ)
					}
					if mapDestination { // convert to primitive
						primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, targets[0])
					}
				} else { // pop fron stack
					if v.JvmSlotSize() == 1 {
						code.Codes[code.CodeLength] = cg.OP_pop
					} else {
						code.Codes[code.CodeLength] = cg.OP_pop2
					}
					code.CodeLength++
				}
				slice()
			}
			continue
		}
		mapDestination := (classnames[0] == java_hashmap_class && targets[0].IsPointer() == false)
		stack, es := m.build(class, code, v, context)
		backPatchEs(es, code.CodeLength) // true or false need to backpatch
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		variableType := v.VariableType
		if v.IsCall() {
			variableType = v.VariableTypes[0]
		}
		if noDestinations[0] == false {
			if t := currentStack + targets[0].JvmSlotSize(); t > maxstack { // incase int convert to double or long
				maxstack = t
			}
			if targets[0].IsNumber() && targets[0].Typ != variableType.Typ { // value is number 2
				m.numberTypeConverter(code, variableType.Typ, targets[0].Typ)
			}
			if mapDestination { // convert to primitive
				primitiveObjectConverter.putPrimitiveInObjectStaticWay(class, code, targets[0])
			}
		} else { // pop fron stack
			if variableType.JvmSlotSize() == 1 {
				code.Codes[code.CodeLength] = cg.OP_pop
			} else {
				code.Codes[code.CodeLength] = cg.OP_pop2
			}
			code.CodeLength++
		}
		slice()
	}
	return
}
