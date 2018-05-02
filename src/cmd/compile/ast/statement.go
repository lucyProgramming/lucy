package ast

import (
	"fmt"
)

const (
	_ = iota
	STATEMENT_TYPE_EXPRESSION
	STATEMENT_TYPE_IF
	STATEMENT_TYPE_BLOCK
	STATEMENT_TYPE_FOR
	STATEMENT_TYPE_CONTINUE
	STATEMENT_TYPE_RETURN
	STATEMENT_TYPE_BREAK
	STATEMENT_TYPE_SWITCH
	STATEMENT_TYPE_SKIP // skip this block
	STATEMENT_TYPE_LABLE
	STATEMENT_TYPE_GOTO
	STATEMENT_TYPE_DEFER
)

type Statement struct {
	Checked           bool
	Pos               *Pos
	Typ               int
	StatementIf       *StatementIF
	Expression        *Expression // expression statment like a=123
	StatementFor      *StatementFor
	StatementReturn   *StatementReturn
	StatementSwitch   *StatementSwitch
	StatementBreak    *StatementBreak
	Block             *Block
	StatementContinue *StatementContinue
	StatmentLable     *StatementLable
	StatementGoto     *StatementGoto
	Defer             *Defer
}

func (s *Statement) StatementName() string {
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
	case STATEMENT_TYPE_LABLE:
		return "lable statement"
	case STATEMENT_TYPE_GOTO:
		return "goto statement"
	case STATEMENT_TYPE_DEFER:
		return "defer statement"
	case STATEMENT_TYPE_BLOCK:
		return "block statement"
	case STATEMENT_TYPE_RETURN:
		return "return statement"
	default:
		panic(11)
	}
	return ""
}

func (s *Statement) isVariableDefinition() bool {
	return s.Typ == STATEMENT_TYPE_EXPRESSION &&
		(s.Expression.Typ == EXPRESSION_TYPE_COLON_ASSIGN || s.Expression.Typ == EXPRESSION_TYPE_VAR)
}

func (s *Statement) check(block *Block) []error { // b is father
	defer func() {
		s.Checked = true
	}()
	switch s.Typ {
	case STATEMENT_TYPE_EXPRESSION:
		return s.checkStatementExpression(block)
	case STATEMENT_TYPE_IF:
		return s.StatementIf.check(block)
	case STATEMENT_TYPE_FOR:
		return s.StatementFor.check(block)
	case STATEMENT_TYPE_SWITCH:
		return s.StatementSwitch.check(block)
	case STATEMENT_TYPE_BREAK:
		if block.InheritedAttribute.StatementFor == nil && block.InheritedAttribute.StatementSwitch == nil {
			return []error{fmt.Errorf("%s '%s' cannot in this scope", errMsgPrefix(s.Pos), s.StatementName())}
		} else {
			if block.InheritedAttribute.StatementFor != nil && block.InheritedAttribute.Defer != nil {
				return []error{fmt.Errorf("%s cannot has '%s' in both 'defer' and 'for'",
					errMsgPrefix(s.Pos), s.StatementName())}
			}
			s.StatementBreak = &StatementBreak{}
			if f, ok := block.InheritedAttribute.mostCloseIsForOrSwitch.(*StatementFor); ok {
				s.StatementBreak.StatementFor = f
			} else {
				s.StatementBreak.StatementSwitch = block.InheritedAttribute.mostCloseIsForOrSwitch.(*StatementSwitch)
			}
		}
	case STATEMENT_TYPE_CONTINUE:
		if block.InheritedAttribute.StatementFor == nil {
			return []error{fmt.Errorf("%s '%s' can`t in this scope",
				errMsgPrefix(s.Pos), s.StatementName())}
		}
		if block.InheritedAttribute.StatementFor != nil && block.InheritedAttribute.Defer != nil {
			return []error{fmt.Errorf("%s cannot has '%s' in both 'defer' and 'for'",
				errMsgPrefix(s.Pos), s.StatementName())}
		}
		if s.StatementContinue == nil {
			s.StatementContinue = &StatementContinue{}
		}
		if s.StatementContinue.StatementFor == nil { // for
			s.StatementContinue.StatementFor = block.InheritedAttribute.StatementFor
		}
	case STATEMENT_TYPE_RETURN:
		if block.InheritedAttribute.Defer != nil {
			return []error{fmt.Errorf("%s cannot has '%s' in 'defer'",
				errMsgPrefix(s.Pos), s.StatementName())}
		}
		es := s.StatementReturn.check(block)
		if len(s.StatementReturn.Defers) > 0 {
			block.InheritedAttribute.Function.MkAutoVarForReturnBecauseOfDefer()
		}
		return es
	case STATEMENT_TYPE_GOTO:
		err := s.checkStatementGoto(block)
		if err != nil {
			return []error{err}
		}
	case STATEMENT_TYPE_DEFER:
		block.InheritedAttribute.Function.mkAutoVarForException()
		s.Defer.Block.inherite(block)
		s.Defer.Block.InheritedAttribute.Defer = s.Defer
		s.Defer.allowCatch = block.IsFunctionTopBlock
		es := s.Defer.Block.check()
		block.Defers = append(block.Defers, s.Defer)
		return es
	case STATEMENT_TYPE_BLOCK:
		s.Block.inherite(block)
		return s.Block.check()
	case STATEMENT_TYPE_LABLE: // nothing to do
	case STATEMENT_TYPE_SKIP:
		if block.InheritedAttribute.Function.isPackageBlockFunction == false {
			return []error{fmt.Errorf("cannot have '%s' at this scope", s.StatementName())}
		}
		return nil
	}
	return nil
}

func (s *Statement) checkStatementExpression(b *Block) []error {
	errs := []error{}
	//
	if s.Expression.Typ == EXPRESSION_TYPE_TYPE_ALIAS { // special case
		t := s.Expression.Data.(*ExpressionTypeAlias)
		err := t.Typ.resolve(b)
		if err != nil {
			return []error{err}
		}
		err = b.insert(t.Name, t.Pos, t.Typ)
		if err != nil {
			return []error{err}
		}
		return nil
	}
	if s.Expression.canBeUsedAsStatemen() {
		if s.Expression.Typ == EXPRESSION_TYPE_FUNCTION {
			f := s.Expression.Data.(*Function)
			if f.Name == "" {
				err := fmt.Errorf("%s function must have a name",
					errMsgPrefix(s.Expression.Pos))
				errs = append(errs, err)
			} else {
				err := b.insert(f.Name, f.Pos, f)
				if err != nil {
					errs = append(errs, err)
				}
			}
			es := f.check(b)
			if errsNotEmpty(es) {
				errs = append(errs, es...)
			}
			f.IsClosureFunction = f.ClosureVars.NotEmpty(f)
			return errs
		}
	} else {
		err := fmt.Errorf("%s expression '%s' evaluate but not used",
			errMsgPrefix(s.Expression.Pos), s.Expression.OpName())
		errs = append(errs, err)
	}
	s.Expression.IsStatementExpression = true
	_, es := b.checkExpression(s.Expression)
	if errsNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}
