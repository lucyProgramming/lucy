package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type FunctionDefaultValueParse struct {
}

func (fd *FunctionDefaultValueParse) Encode(class *cg.ClassHighLevel, f *ast.Function) *cg.AttributeDefaultParameters {
	ret := &cg.AttributeDefaultParameters{}
	ret.Start = uint16(f.DefaultValueStartAt)
	for i := ret.Start; i < uint16(len(f.Type.ParameterList)); i++ {
		switch f.Type.ParameterList[i].Expression.Type {
		case ast.EXPRESSION_TYPE_BOOL:
			if f.Type.ParameterList[i].Expression.Data.(bool) {
				ret.Constants = append(ret.Constants, class.Class.InsertIntConst(1))
			} else {
				ret.Constants = append(ret.Constants, class.Class.InsertIntConst(0))
			}
		case ast.EXPRESSION_TYPE_BYTE:
			ret.Constants = append(ret.Constants, class.Class.InsertIntConst(
				int32(f.Type.ParameterList[i].Expression.Data.(byte))))
		case ast.EXPRESSION_TYPE_SHORT:
			fallthrough
		case ast.EXPRESSION_TYPE_INT:
			ret.Constants = append(ret.Constants, class.Class.InsertIntConst(
				f.Type.ParameterList[i].Expression.Data.(int32)))
		case ast.EXPRESSION_TYPE_LONG:
			ret.Constants = append(ret.Constants, class.Class.InsertLongConst(
				f.Type.ParameterList[i].Expression.Data.(int64)))
		case ast.EXPRESSION_TYPE_FLOAT:
			ret.Constants = append(ret.Constants, class.Class.InsertFloatConst(
				f.Type.ParameterList[i].Expression.Data.(float32)))
		case ast.EXPRESSION_TYPE_DOUBLE:
			ret.Constants = append(ret.Constants, class.Class.InsertDoubleConst(
				f.Type.ParameterList[i].Expression.Data.(float64)))
		case ast.EXPRESSION_TYPE_STRING:
			ret.Constants = append(ret.Constants, class.Class.InsertStringConst(
				f.Type.ParameterList[i].Expression.Data.(string)))
		}
	}
	return ret
}

func (fd *FunctionDefaultValueParse) Decode(class *cg.Class, f *ast.Function, dp *cg.AttributeDefaultParameters) {
	f.HaveDefaultValue = true
	f.DefaultValueStartAt = int(dp.Start)
	for i := uint16(0); i < uint16(len(dp.Constants)); i++ {
		v := f.Type.ParameterList[dp.Start+i]
		v.Expression = &ast.Expression{}
		v.Expression.ExpressionValue = v.Type
		switch v.Type.Type {
		case ast.VARIABLE_TYPE_BOOL:
			v.Expression.Type = ast.EXPRESSION_TYPE_BOOL
			v.Expression.Data = int32(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info)) != 0
		case ast.VARIABLE_TYPE_BYTE:
			v.Expression.Type = ast.EXPRESSION_TYPE_BYTE
			v.Expression.Data = byte(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VARIABLE_TYPE_SHORT:
			v.Expression.Type = ast.EXPRESSION_TYPE_SHORT
			v.Expression.Data = int32(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VARIABLE_TYPE_INT:
			v.Expression.Type = ast.EXPRESSION_TYPE_INT
			v.Expression.Data = int32(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VARIABLE_TYPE_LONG:
			v.Expression.Type = ast.EXPRESSION_TYPE_LONG
			v.Expression.Data = int64(binary.BigEndian.Uint64(class.ConstPool[dp.Constants[i]].Info))
		case ast.VARIABLE_TYPE_FLOAT:
			v.Expression.Type = ast.EXPRESSION_TYPE_FLOAT
			v.Expression.Data = float32(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VARIABLE_TYPE_DOUBLE:
			v.Expression.Type = ast.EXPRESSION_TYPE_DOUBLE
			v.Expression.Data = float64(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VARIABLE_TYPE_STRING:
			v.Expression.Type = ast.EXPRESSION_TYPE_STRING
			utf8Index := binary.BigEndian.Uint16(class.ConstPool[dp.Constants[i]].Info)
			v.Expression.Data = string(class.ConstPool[utf8Index].Info)
		}
	}
}
