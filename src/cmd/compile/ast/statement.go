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
	StatementTypeGoTo
	StatementTypeDefer
	StatementTypeClass
	StatementTypeEnum
	StatementTypeNop
)

type Statement struct {
	Type int

	Checked                   bool // if checked
	Pos                       *Position
	StatementIf               *StatementIf
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
	isStaticFieldDefaultValue bool
	/*
		this.super()
		special case
	*/
	IsCallFatherConstructionStatement bool
}

func (s *Statement) isVariableDefinition() bool {
	return s.Type == StatementTypeExpression &&
		(s.Expression.Type == ExpressionTypeColonAssign || s.Expression.Type == ExpressionTypeVar)
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
		if block.InheritedAttribute.StatementFor == nil &&
			block.InheritedAttribute.StatementSwitch == nil &&
			block.InheritedAttribute.SwitchTemplateBlock == nil {
			return []error{fmt.Errorf("%s 'break' cannot in this scope", errMsgPrefix(s.Pos))}
		} else {
			if block.InheritedAttribute.Defer != nil {
				return []error{fmt.Errorf("%s cannot has 'break' in 'defer'",
					errMsgPrefix(s.Pos))}
			}
			if t, ok := block.InheritedAttribute.ForBreak.(*StatementFor); ok {
				s.StatementBreak.StatementFor = t
			} else if t, ok := block.InheritedAttribute.ForBreak.(*StatementSwitch); ok {
				s.StatementBreak.StatementSwitch = t
			} else {
				s.StatementBreak.SwitchTemplateBlock = block.InheritedAttribute.ForBreak.(*Block)
			}
			s.StatementBreak.mkDefers(block)
		}
	case StatementTypeContinue:
		if block.InheritedAttribute.StatementFor == nil {
			return []error{fmt.Errorf("%s 'continue' can`t in this scope",
				errMsgPrefix(s.Pos))}
		}
		if block.InheritedAttribute.Defer != nil {
			return []error{fmt.Errorf("%s cannot has 'continue' in 'defer'",
				errMsgPrefix(s.Pos))}
		}
		s.StatementContinue.StatementFor = block.InheritedAttribute.StatementFor
		s.StatementContinue.mkDefers(block)
	case StatementTypeReturn:
		if block.InheritedAttribute.Defer != nil {
			return []error{fmt.Errorf("%s cannot has 'return' in 'defer'",
				errMsgPrefix(s.Pos))}
		}
		es := s.StatementReturn.check(block)
		if len(s.StatementReturn.Defers) > 0 {
			block.InheritedAttribute.Function.MkAutoVarForReturnBecauseOfDefer()
		}
		return es
	case StatementTypeGoTo:
		err := s.checkStatementGoTo(block)
		if err != nil {
			return []error{err}
		}
	case StatementTypeDefer:
		block.InheritedAttribute.Function.mkAutoVarForException()
		s.Defer.Block.inherit(block)
		s.Defer.Block.InheritedAttribute.Defer = s.Defer
		es := s.Defer.Block.checkStatements()
		block.Defers = append(block.Defers, s.Defer)
		return es
	case StatementTypeBlock:
		s.Block.inherit(block)
		return s.Block.checkStatements()
	case StatementTypeLabel:
		// nothing to do
		//s.StatementLabel.InDefer = block.InheritedAttribute.Defer != nil
		if block.InheritedAttribute.Defer != nil {
			block.InheritedAttribute.Defer.Labels = append(block.InheritedAttribute.Defer.Labels, s.StatementLabel)
		}
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
	if s.Expression.Type == ExpressionTypeTypeAlias { // special case
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
	if esNotEmpty(es) {
		errs = append(errs, es...)
	}
	return errs
}
