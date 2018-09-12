package jvm

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type DefaultValueParse struct {
}

func (fd *DefaultValueParse) Encode(class *cg.ClassHighLevel, f *ast.Function) *cg.AttributeDefaultParameters {
	ret := &cg.AttributeDefaultParameters{}
	ret.Start = uint16(f.DefaultValueStartAt)
	for i := ret.Start; i < uint16(len(f.Type.ParameterList)); i++ {
		switch f.Type.ParameterList[i].Expression.Type {
		case ast.ExpressionTypeBool:
			if f.Type.ParameterList[i].Expression.Data.(bool) {
				ret.Constants = append(ret.Constants, class.Class.InsertIntConst(1))
			} else {
				ret.Constants = append(ret.Constants, class.Class.InsertIntConst(0))
			}
		case ast.ExpressionTypeByte:
			ret.Constants = append(ret.Constants, class.Class.InsertIntConst(
				int32(f.Type.ParameterList[i].Expression.Data.(byte))))
		case ast.ExpressionTypeShort:
			fallthrough
		case ast.ExpressionTypeInt:
			ret.Constants = append(ret.Constants, class.Class.InsertIntConst(
				f.Type.ParameterList[i].Expression.Data.(int32)))
		case ast.ExpressionTypeLong:
			ret.Constants = append(ret.Constants, class.Class.InsertLongConst(
				f.Type.ParameterList[i].Expression.Data.(int64)))
		case ast.ExpressionTypeFloat:
			ret.Constants = append(ret.Constants, class.Class.InsertFloatConst(
				f.Type.ParameterList[i].Expression.Data.(float32)))
		case ast.ExpressionTypeDouble:
			ret.Constants = append(ret.Constants, class.Class.InsertDoubleConst(
				f.Type.ParameterList[i].Expression.Data.(float64)))
		case ast.ExpressionTypeString:
			ret.Constants = append(ret.Constants, class.Class.InsertStringConst(
				f.Type.ParameterList[i].Expression.Data.(string)))
		}
	}
	return ret
}

func (fd *DefaultValueParse) Decode(class *cg.Class, f *ast.Function, dp *cg.AttributeDefaultParameters) {
	f.HaveDefaultValue = true
	f.DefaultValueStartAt = int(dp.Start)
	for i := uint16(0); i < uint16(len(dp.Constants)); i++ {
		v := f.Type.ParameterList[dp.Start+i]
		v.Expression = &ast.Expression{}
		v.Expression.Value = v.Type
		switch v.Type.Type {
		case ast.VariableTypeBool:
			v.Expression.Type = ast.ExpressionTypeBool
			v.Expression.Data = binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info) != 0
		case ast.VariableTypeByte:
			v.Expression.Type = ast.ExpressionTypeByte
			v.Expression.Data = byte(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VariableTypeShort:
			v.Expression.Type = ast.ExpressionTypeShort
			v.Expression.Data = int32(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VariableTypeInt:
			v.Expression.Type = ast.ExpressionTypeInt
			v.Expression.Data = int32(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VariableTypeLong:
			v.Expression.Type = ast.ExpressionTypeLong
			v.Expression.Data = int64(binary.BigEndian.Uint64(class.ConstPool[dp.Constants[i]].Info))
		case ast.VariableTypeFloat:
			v.Expression.Type = ast.ExpressionTypeFloat
			v.Expression.Data = float32(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VariableTypeDouble:
			v.Expression.Type = ast.ExpressionTypeDouble
			v.Expression.Data = float64(binary.BigEndian.Uint32(class.ConstPool[dp.Constants[i]].Info))
		case ast.VariableTypeString:
			v.Expression.Type = ast.ExpressionTypeString
			utf8Index := binary.BigEndian.Uint16(class.ConstPool[dp.Constants[i]].Info)
			v.Expression.Data = string(class.ConstPool[utf8Index].Info)
		}
	}
}
