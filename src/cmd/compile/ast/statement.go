package ast

import (
	"fmt"
)

const (
	STATEMENT_TYPE_EXPRESSION = iota
	STATEMENT_TYPE_IF
	STATEMENT_TYPE_FOR
	STATEMENT_TYPE_CONTINUE
	STATEMENT_TYPE_RETURN
	STATEMENT_TYPE_BREAK
	STATEMENT_TYPE_SWITCH
	STATEMENT_TYPE_SKIP // skip this block
)

type Statement struct {
	Pos               Pos
	Typ               int
	StatementIf       *StatementIF
	Expression        *Expression // expression statment like a=123
	StatementFor      *StatementFor
	StatementReturn   *StatementReturn
	StatementTryCatch *StatementTryCatch
	StatmentSwitch    *StatmentSwitch
}

func (s *Statement) stateName() string {
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

func (s *Statement) check(b *Block) []error {
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
		s.StatementIf.Block.inherite(b)
		errs = append(errs, s.StatementIf.check()...)
	case STATEMENT_TYPE_FOR:
		s.StatementIf.Block.inherite(b)
		errs = append(errs, s.StatementIf.check()...)
	case STATEMENT_TYPE_SWITCH:
		s.StatementIf.Block.inherite(b)
		errs = append(errs, s.StatementFor.check()...)
	case STATEMENT_TYPE_BREAK:
	case STATEMENT_TYPE_CONTINUE:
		if b.InheritedAttribute.infor {
			errs = append(errs, fmt.Errorf("%s %d:%d %s can`t in this scope", s.Pos.Filename, s.Pos.StartLine, s.Pos.StartColumn, s.stateName()))
		}
	default:
		panic("unkown type statement")
	}
	return errs
}

func notFoundError(pos *Pos, typ, name string) error {
	return fmt.Errorf("%s %d:%d %s:%s not found", pos.Filename, pos.StartLine, pos.StartColumn, typ, name)
}

//func tooManyArgsTocall(p *Pos) error {
//	return fmt.Errorf("%s %d:%d too many args to call", p.Filename, p.StartLine, p.StartColumn)
//}
//func tooewFArgsTocall(p *Pos) error {
//	return fmt.Errorf("%s %d:%d too fw args to call", p.Filename, p.StartLine, p.StartColumn)
//}

func checkFunctionCall(f *Function, call *ExpressionFunctionCall, p *Pos) []error {
	errs := make([]error, 0)
	if len(call.Args) == 0 {
		return nil
	}
	if len(call.Args) != len(f.Typ.Parameters) {
		if len(call.Args) > len(f.Typ.Parameters) {
			errs = append(errs, fmt.Errorf("%s %d:%d too many args to call", p.Filename, p.StartLine, p.StartColumn))
		} else {
			errs = append(errs, fmt.Errorf("%s %d:%d too fw args to call", p.Filename, p.StartLine, p.StartColumn))
		}
		return errs
	}
	length := len(call.Args)
	for i := 0; i < length; i++ {
		t := getTypeFromExpression(call.Args[i])
		if !f.Typ.Parameters[i].TypedName.Typ.typeCompatible(t) {
			typstring1 := ""
			typstring2 := ""
			f.Typ.Parameters[i].TypedName.Typ.TypeString(&typstring1)
			t.TypeString(&typstring2)
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
	errs := []error{}
	//func1()
	if EXPRESSION_TYPE_FUNCTION_CALL == s.Expression.Typ {
		call := s.Expression.Data.(*ExpressionFunctionCall)
		f := b.searchFunction(call.Name)
		if f == nil {
			errs = append(errs, fmt.Errorf("%s %d:%d", notFoundError(&s.Pos, "function", call.Name)))
		} else {
			errs = append(errs, checkFunctionCall(f, call, &s.Pos)...)
		}
		return errs
	}
	//System.log("hello world")
	//if EXPRESSION_TYPE_METHOD_CALL == s.Expression.Typ {
	//	return errs
	//}

	// i++ i-- ++i --i
	if EXPRESSION_TYPE_INCREMENT == s.Expression.Typ ||
		EXPRESSION_TYPE_DECREMENT == s.Expression.Typ ||
		EXPRESSION_TYPE_PRE_INCREMENT == s.Expression.Typ ||
		EXPRESSION_TYPE_PRE_DECREMENT == s.Expression.Typ {
		left := s.Expression.Data.(*Expression)     // left means left value
		if left.Typ == EXPRESSION_TYPE_IDENTIFIER { //naming
			name := left.Data.(string)
			item := b.searchSymbolicItemAlsoGlobalVar(name)
			if item == nil {
				errs = append(errs, notFoundError(&s.Pos, "variable", name))
				return errs
			}
			return errs
		}
		errs = append(errs, fmt.Errorf("%s %d:%d cannot apply ++ or -- on %s", s.Pos.Filename, s.Pos.StartLine, s.Pos.StartColumn, left.humanReadableString()))
		return errs
	}
	if EXPRESSION_TYPE_COLON_ASSIGN == s.Expression.Typ { //declare variable
		binary := s.Expression.Data.(*ExpressionBinary)
		if binary.Left.Typ == EXPRESSION_TYPE_IDENTIFIER {
			errs = append(errs, fmt.Errorf("%s %d:%d no name on the left", s.Pos.Filename, s.Pos.StartLine, s.Pos.StartColumn))
			return errs
		}
		name := binary.Left.Data.(string)
		if _, ok := b.SymbolicTable.itemsMap[name]; ok {
			errs = append(errs, fmt.Errorf("%s %d:%d variable %s is already declared", s.Pos.Filename, s.Pos.StartLine, s.Pos.StartColumn))
			return errs
		}
		item := &SymbolicItem{}
		t := getTypeFromExpression(binary.Right)
		item.Typ = t
		item.Name = name
		b.SymbolicTable.itemsMap[name] = item
		return errs
	}
	//
	if s.Expression.Typ == EXPRESSION_TYPE_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_PLUS_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_MINUS_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_MUL_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_DIV_ASSIGN ||
		s.Expression.Typ == EXPRESSION_TYPE_MOD_ASSIGN {
		binary := s.Expression.Data.(*ExpressionBinary)
		if binary.Left.Typ == EXPRESSION_TYPE_IDENTIFIER {
			name := s.Expression.Data.(string)
			item := b.searchSymbolicItemAlsoGlobalVar(name)
			if item == nil {
				errs = append(errs, notFoundError(&s.Pos, "variable", name))
				return errs
			}
			return errs
		}
	}

	errs = append(errs, fmt.Errorf("%s %d:%d expression(%s) evaluated,but not used", s.Pos.Filename, s.Pos.StartLine, s.Pos.StartColumn, s.Expression.humanReadableString()))
	return errs
}

type StatementTryCatch struct {
	TryBlock     *Block
	CatchBlock   *Block
	FinallyBlock *Block
}

type StatmentSwitch struct {
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
	Expression []*Expression
}

func (s *StatementReturn) check() []error {
	errs := []error{}
	return errs
}

type StatementFor struct {
	Init      *Expression
	Condition *Expression
	Post      *Expression
	Block     *Block
}

func (s *StatementFor) check() []error {
	errs := []error{}
	return errs
}

type StatementIF struct {
	Condition  *Expression
	Block      *Block
	ElseBlock  *Block
	ElseIfList []*StatementElseIf
}

func (s *StatementIF) check() []error {
	errs := []error{}
	return errs
}

type StatementElseIf struct {
	Condition *Expression
	Block     *Block
}

func (s *StatementElseIf) check() []error {
	errs := []error{}
	return errs
}
