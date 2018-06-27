package ast

import (
	"fmt"
)

const (
	_ = iota
	StatementTypeExpression
	StatementTypeIf
	StatementTypeBlock
	StatementTypeFor
	StatementTypeContinue
	StatementTypeReturn
	StatementTypeBreak
	StatementTypeSwitch
	StatementTypeSwitchTemplate
	StatementTypeLabel
	StatementTypeGoto
	StatementTypeDefer
	StatementTypeClass
	StatementTypeEnum
	StatementTypeNop
)

type Statement struct {
	isStaticFieldDefaultValue bool
	Checked                   bool // if checked
	Pos                       *Position
	Type                      int
	StatementIf               *StatementIF
	Expression                *Expression
	StatementFor              *StatementFor
	StatementReturn           *StatementReturn
	StatementSwitch           *StatementSwitch
	StatementSwitchTemplate   *StatementSwitchTemplate
	StatementBreak            *StatementBreak
	Block                     *Block
	StatementContinue         *StatementContinue
	StatementLabel            *StatementLabel
	StatementGoTo             *StatementGoTo
	Defer                     *StatementDefer
	Class                     *Class
	Enum                      *Enum
	/*
		this.super()
		special case
	*/
	IsCallFatherConstructionStatement bool
}

func (s *Statement) StatementName() string {
	switch s.Type {
	case StatementTypeExpression:
		return "expression statement"
	case StatementTypeIf:
		return "if statement"
	case StatementTypeFor:
		return "for statement"
	case StatementTypeContinue:
		return "continue statement"
	case StatementTypeBreak:
		return "break statement"
	case StatementTypeSwitch:
		return "switch statement"
	case StatementTypeLabel:
		return "label statement"
	case StatementTypeGoto:
		return "goto statement"
	case StatementTypeDefer:
		return "defer statement"
	case StatementTypeBlock:
		return "block statement"
	case StatementTypeReturn:
		return "return statement"
	case StatementTypeClass:
		return "class"
	case StatementTypeEnum:
		return "enum"
	case StatementTypeNop:
		return "nop"
	case StatementTypeSwitchTemplate:
		return "switch template"
	}
	return ""
}

func (s *Statement) isVariableDefinition() bool {
	return s.Type == StatementTypeExpression &&
		(s.Expression.Type == EXPRESSION_TYPE_COLON_ASSIGN || s.Expression.Type == EXPRESSION_TYPE_VAR)
}

func (s *Statement) check(block *Block) []error { // block is father
	defer func() {
		s.Checked = true
	}()
	errs := []error{}
	switch s.Type {
	case StatementTypeExpression:
		return s.checkStatementExpression(block)
	case StatementTypeIf:
		return s.StatementIf.check(block)
	case StatementTypeFor:
		return s.StatementFor.check(block)
	case StatementTypeSwitch:
		return s.StatementSwitch.check(block)
	case StatementTypeBreak:
		if block.InheritedAttribute.StatementFor == nil && block.InheritedAttribute.StatementSwitch == nil {
			return []error{fmt.Errorf("%s '%s' cannot in this scope", errMsgPrefix(s.Pos), s.StatementName())}
		} else {
			if block.InheritedAttribute.Defer != nil {
				return []error{fmt.Errorf("%s cannot has '%s' in 'defer'",
					errMsgPrefix(s.Pos), s.StatementName())}
			}
			if f, ok := block.InheritedAttribute.ForBreak.(*StatementFor); ok {
				s.StatementBreak.StatementFor = f
			} else {
				s.StatementBreak.StatementSwitch = block.InheritedAttribute.ForBreak.(*StatementSwitch)
			}
			s.StatementBreak.mkDefers(block)
		}
	case StatementTypeContinue:
		if block.InheritedAttribute.StatementFor == nil {
			return []error{fmt.Errorf("%s '%s' can`t in this scope",
				errMsgPrefix(s.Pos), s.StatementName())}
		}
		if block.InheritedAttribute.Defer != nil {
			return []error{fmt.Errorf("%s cannot has '%s' in 'defer'",
				errMsgPrefix(s.Pos), s.StatementName())}
		}
		s.StatementContinue.StatementFor = block.InheritedAttribute.StatementFor
		s.StatementContinue.mkDefers(block)
	case StatementTypeReturn:
		if block.InheritedAttribute.Defer != nil {
			return []error{fmt.Errorf("%s cannot has '%s' in 'defer'",
				errMsgPrefix(s.Pos), s.StatementName())}
		}
		es := s.StatementReturn.check(block)
		if len(s.StatementReturn.Defers) > 0 {
			block.InheritedAttribute.Function.MkAutoVarForReturnBecauseOfDefer()
		}
		return es
	case StatementTypeGoto:
		err := s.checkStatementGoTo(block)
		if err != nil {
			return []error{err}
		}
	case StatementTypeDefer:
		block.InheritedAttribute.Function.mkAutoVarForException()
		s.Defer.Block.inherit(block)
		s.Defer.Block.InheritedAttribute.Defer = s.Defer
		//s.Defer.allowCatch = block.IsFunctionBlock
		es := s.Defer.Block.checkStatements()
		block.Defers = append(block.Defers, s.Defer)
		return es
	case StatementTypeBlock:
		s.Block.inherit(block)
		return s.Block.checkStatements()
	case StatementTypeLabel:
		// nothing to do
	case StatementTypeClass:
		err := block.Insert(s.Class.Name, s.Pos, s.Class)
		if err != nil {
			errs = append(errs, err)
		}
		return append(errs, s.Class.check(block)...)
	case StatementTypeEnum:
		err := s.Enum.check()
		if err != nil {
			return []error{err}
		}
		err = block.Insert(s.Enum.Name, s.Pos, s.Enum)
		if err != nil {
			return []error{err}
		} else {
			return nil
		}
	case StatementTypeNop:
		//nop , should be never execute to here
	case StatementTypeSwitchTemplate:
		return s.StatementSwitchTemplate.check(block, s)
	}
	return nil
}

func (s *Statement) checkStatementExpression(b *Block) []error {
	errs := []error{}
	//
	if s.Expression.Type == EXPRESSION_TYPE_TYPE_ALIAS { // special case
		t := s.Expression.Data.(*ExpressionTypeAlias)
		err := t.Type.resolve(b)
		if err != nil {
			return []error{err}
		}
		err = b.Insert(t.Name, t.Pos, t.Type)
		if err != nil {
			return []error{err}
		}
		return nil
	}
	if s.Expression.canBeUsedAsStatement() {
		s.Expression.IsStatementExpression = true
	} else {
		err := fmt.Errorf("%s expression '%s' evaluate but not used",
			errMsgPrefix(s.Expression.Pos), s.Expression.OpName())
		errs = append(errs, err)
	}
	_, es := s.Expression.check(b)
	if errorsNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}
