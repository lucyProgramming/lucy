package ast

import (
	"fmt"
)

const (
	STATEMENT_TYPE_EXPRESSION = iota
	STATEMENT_TYPE_IF
	STATEMENT_TYPE_BLOCK
	STATEMENT_TYPE_FOR
	STATEMENT_TYPE_CONTINUE
	STATEMENT_TYPE_RETURN
	STATEMENT_TYPE_BREAK
	STATEMENT_TYPE_SWITCH
	STATEMENT_TYPE_SKIP // skip this block

)

type Statement struct {
	Pos             *Pos
	Typ             int
	StatementIf     *StatementIF
	Expression      *Expression // expression statment like a=123
	StatementFor    *StatementFor
	StatementReturn *StatementReturn
	StatementSwitch *StatementSwitch
	Block           *Block
}

func (s *Statement) statementName() string {
	switch s.Typ {
	case STATEMENT_TYPE_EXPRESSION:
		return "expression statement"
	case STATEMENT_TYPE_IF:
		return "if statement"
	case STATEMENT_TYPE_FOR:
		return "for statement"
	case STATEMENT_TYPE_CONTINUE:
		return "continue statement"
	case STATEMENT_TYPE_BREAK:
		return "break statement"
	case STATEMENT_TYPE_SWITCH:
		return "switch statement"
	case STATEMENT_TYPE_SKIP:
		return "skip statement"
	}
	return ""
}

func (s *Statement) check(b *Block) []error { // b is father
	errs := []error{}
	if b.InheritedAttribute.istop {
		if s.Typ == STATEMENT_TYPE_SKIP { //special case
			return errs // 0 length error
		}
	}
	switch s.Typ {
	case STATEMENT_TYPE_EXPRESSION:
		errs = append(errs, s.checkStatementExpression(b)...)
	case STATEMENT_TYPE_IF:
		t, es := s.StatementIf.check(b)
		if len(es) > 0 {
			errs = append(errs, es...)
		}
		if t != nil {
			s.Typ = STATEMENT_TYPE_BLOCK
		}
		s.Block = t
	case STATEMENT_TYPE_FOR:

	case STATEMENT_TYPE_SWITCH:
		s.StatementIf.Block.inherite(b)
		errs = append(errs, s.StatementFor.check()...)
	case STATEMENT_TYPE_BREAK:
	case STATEMENT_TYPE_CONTINUE:
		if b.InheritedAttribute.infor {
			errs = append(errs, fmt.Errorf("%s %s can`t in this scope", errMsgPrefix(s.Pos), s.statementName()))
		}
	case STATEMENT_TYPE_RETURN:
		if b.InheritedAttribute.function == nil {
			errs = append(errs, fmt.Errorf("%s%s can`t in this scope", errMsgPrefix(s.Pos), s.statementName()))
			return errs
		}

	default:
		panic("unkown type statement" + s.statementName())
	}
	return errs
}

func checkFunctionCall(b *Block, f *Function, call *ExpressionFunctionCall, p *Pos) []error {
	errs := make([]error, 0)
	if len(call.Args) == 0 {
		return nil
	}
	if len(call.Args) != len(f.Typ.Parameters) {
		if len(call.Args) > len(f.Typ.Parameters) {
			errs = append(errs, fmt.Errorf("%s %d:%d too many args to call", p.Filename, p.StartLine, p.StartColumn))
		} else {
			errs = append(errs, fmt.Errorf("%s %d:%d too few args to call", p.Filename, p.StartLine, p.StartColumn))
		}
		return errs
	}
	length := len(call.Args)
	for i := 0; i < length; i++ {
		t, es := b.checkExpression(call.Args[i])
		if errsNotEmpty(es) {
			errs = append(errs, es...)
			continue
		}
		if !f.Typ.Parameters[i].Typ.typeCompatible(t) {
			typstring1 := f.Typ.Parameters[i].Typ.TypeString()
			typstring2 := t.TypeString()
			errs = append(errs,
				fmt.Errorf("%s %d:%d %s not match %s,cannot call function",
					p.Filename,
					p.StartLine,
					p.StartColumn,
					typstring1,
					typstring2,
				))
		}
	}
	return errs
}

func (s *Statement) checkStatementExpression(b *Block) []error {
	//fmt.Println("##############", s.Expression.OpName())
	//errs := []error{}
	//func1()
	//	if EXPRESSION_TYPE_FUNCTION_CALL == s.Expression.Typ {
	//		call := s.Expression.Data.(*ExpressionFunctionCall)
	//		f := b.searchFunction("")
	//		if f == nil {
	//			errs = append(errs, notFoundError(s.Pos, "function", call.Expression.OpName()))
	//		} else {
	//			errs = append(errs, checkFunctionCall(b, f, call, s.Pos)...)
	//		}
	//		return errs
	//	}
	//System.log("hello world")
	//if EXPRESSION_TYPE_METHOD_CALL == s.Expression.Typ {
	//	return errs
	//}
	// i++ i-- ++i --i
	//if EXPRESSION_TYPE_INCREMENT == s.Expression.Typ ||
	//	EXPRESSION_TYPE_DECREMENT == s.Expression.Typ ||
	//	EXPRESSION_TYPE_PRE_INCREMENT == s.Expression.Typ ||
	//	EXPRESSION_TYPE_PRE_DECREMENT == s.Expression.Typ {
	//	left := s.Expression.Data.(*Expression)     // left means left value
	//	if left.Typ != EXPRESSION_TYPE_IDENTIFIER { //naming
	//		name := left.Data.(string)
	//		item, _ := b.searchByName(name)
	//		if item == nil {
	//			errs = append(errs, notFoundError(s.Pos, "variable", name))
	//			return errs
	//		}
	//		return errs
	//	}
	//	errs = append(errs, fmt.Errorf("%s %d:%d cannot apply ++ or -- on %s", s.Pos.Filename, s.Pos.StartLine, s.Pos.StartColumn, left.OpName()))
	//	return errs
	//}
	//if EXPRESSION_TYPE_COLON_ASSIGN == s.Expression.Typ { //declare variable
	//	binary := s.Expression.Data.(*ExpressionBinary)
	//	if binary.Left.Typ != EXPRESSION_TYPE_IDENTIFIER && binary.Left.Typ != EXPRESSION_TYPE_LIST {
	//		errs = append(errs, fmt.Errorf("%s no name on the left,but %s", errMsgPrefix(binary.Left.Pos), binary.Left.OpName()))
	//		return errs
	//	}
	//	var names []*Expression
	//	if binary.Left.Typ == EXPRESSION_TYPE_IDENTIFIER {
	//		names = append(names, binary.Left)
	//	} else {
	//		names = binary.Left.Data.([]*Expression)
	//	}
	//	values := binary.Right.Data.([]*Expression)
	//	for k, v := range names {
	//		if v.Typ != EXPRESSION_TYPE_IDENTIFIER {
	//			errs = append(errs, fmt.Errorf("%s expression is not a name,but %s", errMsgPrefix(v.Pos), v.OpName()))
	//			continue
	//		}
	//		if v.Data.(string) == "_" { // not receive
	//			continue
	//		}
	//		vd := &VariableDefinition{}
	//		vd.Name = v.Data.(string)
	//		vd.Expression = values[k]
	//		es := b.checkVar(vd)
	//		if errsNotEmpty(es) {
	//			errs = append(errs, es...)
	//			continue
	//		}
	//		err := b.insert(vd.Name, s.Expression.Pos, vd)
	//		if err != nil {
	//			errs = append(errs, err)
	//		}
	//	}
	//	return errs
	//}

	_, es := b.checkExpression(s.Expression)
	return es
}

type StatementSwitch struct {
	Outter              *Block
	Condition           *Expression //switch
	StatmentSwitchCases []*StatmentSwitchCase
	Default             *Block
}

type StatmentSwitchCase struct {
	Match *Expression
	Block *Block
}

func (s *StatmentSwitchCase) check() []error {
	errs := []error{}
	return errs
}

type StatementReturn struct {
	Pos         *Pos // use some time
	Expressions []*Expression
}

func (s *StatementReturn) check(b *Block) []error {
	errs := make([]error, 0)
	if len(s.Expressions) > len(b.InheritedAttribute.returns) {
		errs = append(errs, fmt.Errorf("%s %d:%d too many to return", s.Pos.Filename, s.Pos.StartLine, s.Pos.StartColumn))
		return errs
	} else if len(s.Expressions) < len(b.InheritedAttribute.returns) {
		errs = append(errs, fmt.Errorf("%s %d:%d too few to return", s.Pos.Filename, s.Pos.StartLine, s.Pos.StartColumn))
		return errs
	} else {
		return errs // 0 length errs,that is ok
	}
	for k, v := range s.Expressions {
		t, es := b.checkExpression(v)
		if es != nil && len(es) > 0 {
			errs = append(errs, es...)
			continue
		}
		if !b.InheritedAttribute.returns[k].Typ.typeCompatible(t) {
			errs = append(errs, typeNotMatchError(v.Pos, b.InheritedAttribute.returns[k].Typ, t))
		}
	}
	return errs
}

func typeNotMatchError(pos *Pos, t1, t2 *VariableType) error {
	typestring1 := t1.TypeString()
	typestring2 := t1.TypeString()
	return fmt.Errorf("%s %d:%d type not match (%s!=%s)", pos.Filename, pos.StartLine, pos.StartColumn, typestring1, typestring2)
}

type StatementFor struct {
	Pos       *Pos
	Init      *Expression
	Condition *Expression
	Post      *Expression
	Block     *Block
}

func (s *StatementFor) check() []error {
	errs := []error{}
	return errs
}

type ElseIfList []*StatementElseIf

func (e ElseIfList) check() []error {
	errs := make([]error, 0)
	var err error
	for _, v := range e {
		t, es := v.Block.checkExpression(v.Condition)
		if errsNotEmpty(es) {
			errs = append(errs, es...)
			continue
		}
		if t.Typ != VARIABLE_TYPE_BOOL {
			errs = append(errs, fmt.Errorf("%s not a bool expression", errMsgPrefix(v.Condition.Pos), err))
			continue
		}
		errs = append(errs, v.Block.check(nil)...)
	}
	return errs
}

type StatementIF struct {
	Condition  *Expression
	Block      *Block
	ElseBlock  *Block
	ElseIfList ElseIfList
}

func (s *StatementIF) check(father *Block) (*Block, []error) {
	errs := []error{}
	//inherite
	s.Block.inherite(father)
	if s.ElseIfList != nil && len(s.ElseIfList) > 0 {
		for _, v := range s.ElseIfList {
			v.Block.inherite(father)
		}
	}
	if s.ElseBlock != nil {
		s.ElseBlock.inherite(father)
	}
	conditionType, es := s.Block.checkExpression(s.Condition)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	if conditionType != nil {
		if conditionType.Typ != VARIABLE_TYPE_BOOL {
			errs = append(errs, fmt.Errorf("%s condition is not a bool expression", errMsgPrefix(s.Condition.Pos)))
		}
	}
	// condition should a bool expression
	if s.Condition.Typ == EXPRESSION_TYPE_BOOL && s.Condition.Data.(bool) == true { // if(true){...}
		errs = append(errs, s.Block.check(nil)...)
		return s.Block, errs
	}
	if s.Condition.Typ == EXPRESSION_TYPE_BOOL && s.Condition.Data.(bool) == false { // if(false){}else{...}
		errs = append(errs, s.ElseBlock.check(nil)...)
		return s.Block, errs
	}

	errs = append(errs, s.Block.check(nil)...)
	if s.ElseIfList != nil && len(s.ElseIfList) > 0 {
		errs = append(errs, s.ElseIfList.check()...)
	}
	if s.ElseBlock != nil {
		errs = append(errs, s.ElseBlock.check(nil)...)
	}
	return nil, errs
}

type StatementElseIf struct {
	Condition *Expression
	Block     *Block
}

func (s *StatementElseIf) check() []error {
	errs := []error{}
	return errs
}
