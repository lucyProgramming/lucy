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
	for i := ret.Start; i < uint16(len(f.Typ.ParameterList)); i++ {
		switch f.Typ.ParameterList[i].Expression.Typ {
		case ast.EXPRESSION_TYPE_BOOL:
			if f.Typ.ParameterList[i].Expression.Data.(bool) {
				ret.Consts = append(ret.Consts, class.Class.InsertIntConst(1))
			} else {
				ret.Consts = append(ret.Consts, class.Class.InsertIntConst(0))
			}
		case ast.EXPRESSION_TYPE_BYTE:
			ret.Consts = append(ret.Consts, class.Class.InsertIntConst(int32(f.Typ.ParameterList[i].Expression.Data.(byte))))
		case ast.EXPRESSION_TYPE_SHORT:
			fallthrough
		case ast.EXPRESSION_TYPE_INT:
			ret.Consts = append(ret.Consts, class.Class.InsertIntConst(f.Typ.ParameterList[i].Expression.Data.(int32)))
		case ast.EXPRESSION_TYPE_LONG:
			ret.Consts = append(ret.Consts, class.Class.InsertLongConst(f.Typ.ParameterList[i].Expression.Data.(int64)))
		case ast.EXPRESSION_TYPE_FLOAT:
			ret.Consts = append(ret.Consts, class.Class.InsertFloatConst(f.Typ.ParameterList[i].Expression.Data.(float32)))
		case ast.EXPRESSION_TYPE_DOUBLE:
			ret.Consts = append(ret.Consts, class.Class.InsertDoubleConst(f.Typ.ParameterList[i].Expression.Data.(float64)))
		case ast.EXPRESSION_TYPE_STRING:
			ret.Consts = append(ret.Consts, class.Class.InsertStringConst(f.Typ.ParameterList[i].Expression.Data.(string)))
		}
	}
	return ret
}

func (fd *FunctionDefaultValueParse) Decode(class *cg.Class, f *ast.Function, dp *cg.AttributeDefaultParameters) {
	for i := uint16(0); i < uint16(len(dp.Consts)); i++ {
		v := f.Typ.ParameterList[dp.Start+i]
		v.Expression = &ast.Expression{}
		switch v.Typ.Typ {
		case ast.VARIABLE_TYPE_BOOL:
			v.Expression.Typ = ast.EXPRESSION_TYPE_BOOL
			v.Expression.Data = int32(binary.BigEndian.Uint32(class.ConstPool[dp.Consts[i]].Info)) != 0
		case ast.VARIABLE_TYPE_BYTE:
			v.Expression.Typ = ast.EXPRESSION_TYPE_BYTE
			v.Expression.Data = byte(binary.BigEndian.Uint32(class.ConstPool[dp.Consts[i]].Info))
		case ast.VARIABLE_TYPE_SHORT:
			v.Expression.Typ = ast.EXPRESSION_TYPE_SHORT
			v.Expression.Data = int32(binary.BigEndian.Uint32(class.ConstPool[dp.Consts[i]].Info))
		case ast.VARIABLE_TYPE_INT:
			v.Expression.Typ = ast.EXPRESSION_TYPE_INT
			v.Expression.Data = int32(binary.BigEndian.Uint32(class.ConstPool[dp.Consts[i]].Info))
		case ast.VARIABLE_TYPE_LONG:
			v.Expression.Typ = ast.EXPRESSION_TYPE_LONG
			v.Expression.Data = int64(binary.BigEndian.Uint64(class.ConstPool[dp.Consts[i]].Info))
		case ast.VARIABLE_TYPE_FLOAT:
			v.Expression.Typ = ast.EXPRESSION_TYPE_FLOAT
			v.Expression.Data = float32(binary.BigEndian.Uint32(class.ConstPool[dp.Consts[i]].Info))
		case ast.VARIABLE_TYPE_DOUBLE:
			v.Expression.Typ = ast.EXPRESSION_TYPE_DOUBLE
			v.Expression.Data = float64(binary.BigEndian.Uint32(class.ConstPool[dp.Consts[i]].Info))
		case ast.VARIABLE_TYPE_STRING:
			v.Expression.Typ = ast.EXPRESSION_TYPE_STRING
			utf8Index := binary.BigEndian.Uint16(class.ConstPool[dp.Consts[i]].Info)
			v.Expression.Data = string(class.ConstPool[utf8Index].Info)
		}
	}
}
