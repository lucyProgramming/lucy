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
	values := []*ast.VariableType{}
	ops := [][]byte{}
	classnames := []string{}
	fieldnames := []string{}
	fieldDescriptors := []string{}
	remainstacks := []uint16{}
	// put left value one the stack
	index := len(lefts) - 1
	for index >= 0 { //
		stack, remainstack, op, value, classname, fieldname, fieldDescriptor := m.getLeftValue(class, code, e, context)
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		values = append(values, value)
		ops = append(ops, op)
		classnames = append(classnames, classname)
		fieldnames = append(fieldnames, fieldname)
		fieldDescriptors = append(fieldDescriptors, fieldDescriptor)
		currentStack += remainstack
		remainstacks = append(remainstacks, remainstack)
		index--
	}
	//
	rights := bin.Right.Data.([]*ast.Expression)
	slice := func() {
		for k, v := range ops[0] {
			code.Codes[code.CodeLength+uint16(k)] = v
		}
		code.CodeLength += uint16(len(ops[0]))
		if classnames[0] != "" {
			class.InsertFieldRef(cg.CONSTANT_Fieldref_info_high_level{}, code.Codes[code.CodeLength:code.CodeLength+2])
		}
		values = values[1:] // slice
		ops = ops[1:]
		classnames = classnames[1:]
		fieldnames = fieldnames[1:]
		fieldDescriptors = fieldDescriptors[1:]
		remainstacks = remainstacks[1:]
	}
	for _, v := range rights {
		if (v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || v.Typ == ast.EXPRESSION_TYPE_METHOD_CALL) &&
			len(e.VariableTypes) > 0 {
			// stack top is arrayListObject object
			stack, _ := m.build(class, code, e, context)
			if t := currentStack + stack; t > maxstack {
				maxstack = t
			}
			if t := 2 + currentStack; t > maxstack {
				maxstack = t
			}
			for k, v := range v.VariableTypes {
				m.buildStoreArrayListAutoVar(code, context)
				switch k {
				case 0:
					code.Codes[code.CodeLength] = cg.OP_iconst_0
					code.CodeLength++
				case 1:
					code.Codes[code.CodeLength] = cg.OP_iconst_1
					code.CodeLength++
				case 2:
					code.Codes[code.CodeLength] = cg.OP_iconst_2
					code.CodeLength++
				case 3:
					code.Codes[code.CodeLength] = cg.OP_iconst_3
					code.CodeLength++
				case 4:
					code.Codes[code.CodeLength] = cg.OP_iconst_4
					code.CodeLength++
				case 5:
					code.Codes[code.CodeLength] = cg.OP_iconst_5
					code.CodeLength++
				default:
					if k > 255 {
						panic("over 255")
					}
					code.Codes[code.CodeLength] = cg.OP_bipush
					code.Codes[code.CodeLength+1] = byte(k)
					code.CodeLength += 2
				}
				code.Codes[code.CodeLength] = cg.OP_invokevirtual
				class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
					Class:      "",
					Name:       "",
					Descriptor: "",
				}, code.Codes[code.CodeLength+1:code.CodeLength+3])
				code.CodeLength += 3
				//
				switch v.Typ {
				case ast.VARIABLE_TYPE_BOOL:
					fallthrough
				case ast.VARIABLE_TYPE_BYTE:
					fallthrough
				case ast.VARIABLE_TYPE_SHORT:
					fallthrough
				case ast.VARIABLE_TYPE_INT:
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
						Class:      "",
						Name:       "",
						Descriptor: "",
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
				case ast.VARIABLE_TYPE_LONG:
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
						Class:      "",
						Name:       "",
						Descriptor: "",
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
				case ast.VARIABLE_TYPE_FLOAT:
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
						Class:      "",
						Name:       "",
						Descriptor: "",
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
				case ast.VARIABLE_TYPE_DOUBLE:
					code.Codes[code.CodeLength] = cg.OP_invokevirtual
					class.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
						Class:      "",
						Name:       "",
						Descriptor: "",
					}, code.Codes[code.CodeLength+1:code.CodeLength+3])
					code.CodeLength += 3
				}
				if values[0].IsNumber() && values[0].Typ != v.Typ { // value is number 2
					m.numberTypeConverter(code, v.Typ, values[0].Typ)
				}
				slice()
			}
			currentStack -= remainstacks[0]
			continue
		}
		stack, es := m.build(class, code, e, context)
		backPatchEs(es, code) // true or false need to backpatch
		if t := currentStack + stack; t > maxstack {
			maxstack = t
		}
		if t := currentStack + values[0].JvmSlotSize(); t > maxstack { // incase int convert to double or long
			maxstack = t
		}
		// convert number mostly
		variableType := e.VariableType
		if v.Typ == ast.EXPRESSION_TYPE_FUNCTION_CALL || v.Typ == ast.EXPRESSION_TYPE_METHOD_CALL {
			variableType = e.VariableTypes[0]
		}
		if values[0].IsNumber() && values[0].Typ != variableType.Typ { // value is number 2
			m.numberTypeConverter(code, variableType.Typ, values[0].Typ)
		}
		currentStack -= remainstacks[0]
		slice()
	}
	return
}
func (m *MakeExpression) buildColonAssign(class *cg.ClassHighLevel, code *cg.AttributeCode, e *ast.Expression, context *Context) (maxstack uint16) {
	return
}
