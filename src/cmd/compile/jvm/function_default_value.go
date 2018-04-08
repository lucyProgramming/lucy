package jvm

import (
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

}
