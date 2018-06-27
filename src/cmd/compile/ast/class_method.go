package ast

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type ClassMethod struct {
	isCompilerAuto  bool // compile auto method
	Function        *Function
	LoadFromOutSide bool
}

func (m *ClassMethod) IsPublic() bool {
	return (m.Function.AccessFlags & cg.ACC_METHOD_PUBLIC) != 0
}
func (m *ClassMethod) IsProtected() bool {
	return (m.Function.AccessFlags & cg.ACC_METHOD_PROTECTED) != 0
}

func (m *ClassMethod) IsStatic() bool {
	return (m.Function.AccessFlags & cg.ACC_METHOD_STATIC) != 0
}

func (m *ClassMethod) IsPrivate() bool {
	return (m.Function.AccessFlags & cg.ACC_METHOD_PRIVATE) != 0
}

func (m *ClassMethod) IsFirstStatementCallFatherConstruction() bool {
	if len(m.Function.Block.Statements) == 0 {
		return false
	}
	s := m.Function.Block.Statements[0]
	if s.Type != StatementTypeExpression {
		return false
	}
	e := s.Expression
	if e.Type != EXPRESSION_TYPE_METHOD_CALL {
		return false
	}
	call := s.Expression.Data.(*ExpressionMethodCall)
	if call.Expression.isThis() == false || call.Name != SUPER {
		return false
	}
	return true
}
