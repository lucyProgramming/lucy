package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type DefaultValueParse struct {
}

func (fd *DefaultValueParse) Encode(
	class *cg.ClassHighLevel,
	f *ast.Function) *cg.AttributeDefaultParameters {
	ret := &cg.AttributeDefaultParameters{}
	ret.Start = uint16(f.DefaultValueStartAt)
	for i := ret.Start; i < uint16(len(f.Type.ParameterList)); i++ {
		switch f.Type.ParameterList[i].DefaultValueExpression.Type {
		case ast.ExpressionTypeBool:
			if f.Type.ParameterList[i].DefaultValueExpression.Data.(bool) {
				ret.Constants = append(ret.Constants, class.Class.InsertIntConst(1))
			} else {
				ret.Constants = append(ret.Constants, class.Class.InsertIntConst(0))
			}
		case ast.ExpressionTypeByte:
			fallthrough
		case ast.ExpressionTypeShort:
			fallthrough
		case ast.ExpressionTypeInt:
			t := f.Type.ParameterList[i].DefaultValueExpression.Data.(int64)
			ret.Constants = append(ret.Constants, class.Class.InsertIntConst(int32(t)))
		case ast.ExpressionTypeLong:
			ret.Constants = append(ret.Constants, class.Class.InsertLongConst(
				f.Type.ParameterList[i].DefaultValueExpression.Data.(int64)))
		case ast.ExpressionTypeFloat:
			ret.Constants = append(ret.Constants, class.Class.InsertFloatConst(
				f.Type.ParameterList[i].DefaultValueExpression.Data.(float32)))
		case ast.ExpressionTypeDouble:
			ret.Constants = append(ret.Constants, class.Class.InsertDoubleConst(
				f.Type.ParameterList[i].DefaultValueExpression.Data.(float64)))
		case ast.ExpressionTypeString:
			ret.Constants = append(ret.Constants, class.Class.InsertStringConst(
				f.Type.ParameterList[i].DefaultValueExpression.Data.(string)))
		}
	}
	return ret
}

func (fd *DefaultValueParse) Decode(class *cg.Class, f *ast.Function, dp *cg.AttributeDefaultParameters) {
	f.HaveDefaultValue = true
	f.DefaultValueStartAt = int(dp.Start)
	for i := uint16(0); i < uint16(len(dp.Constants)); i++ {
		v := f.Type.ParameterList[dp.Start+i]
		v.DefaultValueExpression = &ast.Expression{}
		v.DefaultValueExpression.Value = v.Type
		switch v.Type.Type {
		case ast.VariableTypeBool:
			v.DefaultValueExpression.Type = ast.ExpressionTypeBool
			v.DefaultValueExpression.Data =
				binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info) != 0
		case ast.VariableTypeByte:
			fallthrough
		case ast.VariableTypeShort:
			fallthrough
		case ast.VariableTypeChar:
			fallthrough
		case ast.VariableTypeInt:
			v.DefaultValueExpression.Type = ast.ExpressionTypeInt
			v.DefaultValueExpression.Data =
				int64(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VariableTypeLong:
			v.DefaultValueExpression.Type = ast.ExpressionTypeLong
			v.DefaultValueExpression.Data =
				int64(binary.BigEndian.Uint64(class.ConstPool[dp.Constants[i]].Info))
		case ast.VariableTypeFloat:
			v.DefaultValueExpression.Type = ast.ExpressionTypeFloat
			v.DefaultValueExpression.Data =
				float32(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VariableTypeDouble:
			v.DefaultValueExpression.Type = ast.ExpressionTypeDouble
			v.DefaultValueExpression.Data =
				float64(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VariableTypeString:
			v.DefaultValueExpression.Type = ast.ExpressionTypeString
			utf8Index := binary.BigEndian.Uint16(class.ConstPool[dp.Constants[i]].Info)
			v.DefaultValueExpression.Data =
				string(class.ConstPool[utf8Index].Info)
		}
	}
}
