package jvm

import (
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
	fieldnames := make([]string, len(lefts))
	fieldDescriptors := make([]string, len(lefts))
	remainstacks := make([]uint16, len(lefts))
	noDestinations := make([]bool, len(lefts))
	// put left value one the stack
	index := len(lefts) - 1
	for index >= 0 { //
		stack, remainstack, op, target, classname, fieldname, fieldDescriptor := m.getLeftValue(class, code, lefts[index], context)
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		noDestination := false
		if lefts[index].Typ == ast.EXPRESSION_TYPE_IDENTIFIER {
			if lefts[index].Data.(*ast.ExpressionIdentifer).Name == ast.NO_NAME_IDENTIFIER {
				noDestination = true
			}
		}
		targets[index] = target
		ops[index] = op
		classnames[index] = classname
		fieldnames[index] = fieldname
		fieldDescriptors[index] = fieldDescriptor
		remainstacks[index] = remainstack
		noDestinations[index] = noDestination
		currentStack += remainstack
		index--
	}
	//
	rights := bin.Right.Data.([]*ast.Expression)
	slice := func() {
		currentStack -= remainstacks[0]
		if noDestinations[0] == false {
			if ops[0] != nil {
				for k, v := range ops[0] { // write op
					code.Codes[code.CodeLength+uint16(k)] = v
				}
				code.CodeLength += uint16(len(ops[0]))
			}
			if classnames[0] != "" {
				class.InsertFieldRef(cg.CONSTANT_Fieldref_info_high_level{
					Class:      classnames[0],
					Name:       fieldnames[0],
					Descriptor: fieldDescriptors[0],
				}, code.Codes[code.CodeLength:code.CodeLength+2])
			}
		}
		targets = targets[1:] // slice
		ops = ops[1:]
		classnames = classnames[1:]
		fieldnames = fieldnames[1:]
		fieldDescriptors = fieldDescriptors[1:]
		remainstacks = remainstacks[1:]
	}
	for _, v := range rights {
		if (v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || v.Typ == ast.EXPRESSION_TYPE_METHOD_CALL) && len(v.VariableTypes) > 1 {
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
		stack, es := m.build(class, code, v, context)
		backPatchEs(es, code) // true or false need to backpatch
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		variableType := v.VariableType
		if v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || v.Typ == ast.EXPRESSION_TYPE_METHOD_CALL {
			variableType = v.VariableTypes[0]
		}
		if noDestinations[0] == false {
			if t := currentStack + targets[0].JvmSlotSize(); t > maxstack { // incase int convert to double or long
				maxstack = t
			}
			if targets[0].IsNumber() && targets[0].Typ != variableType.Typ { // value is number 2
				m.numberTypeConverter(code, variableType.Typ, targets[0].Typ)
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

// ast will convert colon asssgin to assgin
func (m *MakeExpression) buildColonAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	panic(1)
	return
}
