package jvm

import (
	"fmt"
	//	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

// a,b,c = 122,fdfd2232,"hello";
func (m *MakeExpression) buildAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
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
				mapDestination := false
				if classnames[0] == java_hashmap_class && targets[0].IsPointer() == false {
					m.prepareStackForMapAssignWhenValueIsNotPointer(class, code, targets[0])
					currentStack += 2
					mapDestination = true
				}
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
						m.pack2Object(class, code, targets[0])
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
				if mapDestination {
					currentStack -= 2
				}
			}
			continue
		}
		mapDestination := false
		if classnames[0] == java_hashmap_class && targets[0].IsPointer() == false {
			m.prepareStackForMapAssignWhenValueIsNotPointer(class, code, targets[0])
			currentStack += 2
			mapDestination = true
		}
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
				m.pack2Object(class, code, targets[0])
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
		if mapDestination {
			currentStack -= 2
		}
	}
	return
}

func (m *MakeExpression) prepareStackForMapAssignWhenValueIsNotPointer(class *cg.ClassHighLevel, code *cg.AttributeCode, t *ast.VariableType) (stack uint16) {
	t.IsPointer()
	stack = 2
	switch t.Typ {
	case ast.VARIABLE_TYPE_BOOL:
		code.Codes[code.CodeLength] = cg.OP_new
		class.InsertClassConst("java/lang/Boolean", code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
	case ast.VARIABLE_TYPE_BYTE:
		fallthrough
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		code.Codes[code.CodeLength] = cg.OP_new
		class.InsertClassConst("java/lang/Integer", code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
	case ast.VARIABLE_TYPE_LONG:
		code.Codes[code.CodeLength] = cg.OP_new
		class.InsertClassConst("java/lang/Long", code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
	case ast.VARIABLE_TYPE_FLOAT:
		code.Codes[code.CodeLength] = cg.OP_new
		class.InsertClassConst("java/lang/Float", code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
	case ast.VARIABLE_TYPE_DOUBLE:
		code.Codes[code.CodeLength] = cg.OP_new
		class.InsertClassConst("java/lang/Double", code.Codes[code.CodeLength+1:code.CodeLength+3])
		code.Codes[code.CodeLength+3] = cg.OP_dup
		code.CodeLength += 4
	}
	return
}

// ast will convert colon asssgin to assgin
func (m *MakeExpression) buildColonAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	panic(1)
	return
}

//func (m *MakeExpression) convertPrimitiveType2Object(class *cg.ClassHighLevel, code *cg.AttributeCode, t *ast.VariableType) (maxstack uint16) {
//	switch t.Typ {
//	case ast.VARIABLE_TYPE_BOOL:
//		code.Codes[code.CodeLength] = cg.OP_new
//		class.InsertClassConst("java/lang/Boolean", code.Codes[code.CodeLength+1:code.CodeLength+3])
//		code.CodeLength += 3
//		code.Codes[code.CodeLength] = cg.OP_dup_x1
//		code.Codes[code.CodeLength+1] = cg.OP_swap
//		code.Codes[code.CodeLength+2] = cg.OP_invokespecial
//		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
//			Class:      "java/lang/Boolean",
//			Name:       specail_method_init,
//			Descriptor: "(Z)V",
//		}, code.Codes[code.CodeLength+3:code.CodeLength+5])
//		code.CodeLength += 5
//		maxstack = 3
//	case ast.VARIABLE_TYPE_BYTE:
//		fallthrough
//	case ast.VARIABLE_TYPE_SHORT:
//		fallthrough
//	case ast.VARIABLE_TYPE_INT:
//		code.Codes[code.CodeLength] = cg.OP_new
//		class.InsertClassConst("java/lang/Integer", code.Codes[code.CodeLength+1:code.CodeLength+3])
//		code.CodeLength += 3
//		code.Codes[code.CodeLength] = cg.OP_dup_x1
//		code.Codes[code.CodeLength+1] = cg.OP_swap
//		code.Codes[code.CodeLength+2] = cg.OP_invokespecial
//		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
//			Class:      "java/lang/Integer",
//			Name:       specail_method_init,
//			Descriptor: "(I)V",
//		}, code.Codes[code.CodeLength+3:code.CodeLength+5])
//		code.CodeLength += 5
//		maxstack = 3
//	case ast.VARIABLE_TYPE_LONG:
//		code.Codes[code.CodeLength] = cg.OP_new
//		class.InsertClassConst("java/lang/Long", code.Codes[code.CodeLength+1:code.CodeLength+3])
//		code.CodeLength += 3
//		code.Codes[code.CodeLength] = cg.OP_dup_x2
//		code.Codes[code.CodeLength+1] = cg.OP_swap
//		code.Codes[code.CodeLength+2] = cg.OP_invokespecial
//		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
//			Class:      "java/lang/Long",
//			Name:       specail_method_init,
//			Descriptor: "(J)V",
//		}, code.Codes[code.CodeLength+3:code.CodeLength+5])
//		code.CodeLength += 5
//		maxstack = 4
//	case ast.VARIABLE_TYPE_FLOAT:
//		code.Codes[code.CodeLength] = cg.OP_new
//		class.InsertClassConst("java/lang/Float", code.Codes[code.CodeLength+1:code.CodeLength+3])
//		code.CodeLength += 3
//		code.Codes[code.CodeLength] = cg.OP_dup_x1
//		code.Codes[code.CodeLength+1] = cg.OP_swap
//		code.Codes[code.CodeLength+2] = cg.OP_invokespecial
//		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
//			Class:      "java/lang/Float",
//			Name:       specail_method_init,
//			Descriptor: "(F)V",
//		}, code.Codes[code.CodeLength+3:code.CodeLength+5])
//		code.CodeLength += 5
//		maxstack = 3
//	case ast.VARIABLE_TYPE_DOUBLE:
//		code.Codes[code.CodeLength] = cg.OP_new
//		class.InsertClassConst("java/lang/Double", code.Codes[code.CodeLength+1:code.CodeLength+3])
//		code.CodeLength += 3
//		code.Codes[code.CodeLength] = cg.OP_dup2_x1
//		code.Codes[code.CodeLength+1] = cg.OP_swap
//		code.Codes[code.CodeLength+2] = cg.OP_invokespecial
//		class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
//			Class:      "java/lang/Double",
//			Name:       specail_method_init,
//			Descriptor: "(J)V",
//		}, code.Codes[code.CodeLength+3:code.CodeLength+5])
//		code.CodeLength += 5
//		maxstack = 4
//	}
//	return
//}
