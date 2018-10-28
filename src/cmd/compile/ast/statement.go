package ast

import (
	"fmt"
)

type StatementTypeKind int

const (
	_ StatementTypeKind = iota
	StatementTypeExpression
	StatementTypeIf
	StatementTypeBlock
	StatementTypeFor
	StatementTypeContinue
	StatementTypeReturn
	StatementTypeBreak
	StatementTypeSwitch
	StatementTypeWhen
	StatementTypeLabel
	StatementTypeGoTo
	StatementTypeDefer
	StatementTypeClass
	StatementTypeEnum
	StatementTypeNop
	StatementTypeImport
	StatementTypeTypeAlias
)

type Statement struct {
	Type                      StatementTypeKind
	Checked                   bool // if checked
	Pos                       *Pos
	StatementIf               *StatementIf
	Expression                *Expression
	TypeAlias                 *TypeAlias
	StatementFor              *StatementFor
	StatementReturn           *StatementReturn
	StatementSwitch           *StatementSwitch
	StatementWhen             *StatementWhen
	StatementBreak            *StatementBreak
	Block                     *Block
	StatementContinue         *StatementContinue
	StatementLabel            *StatementLabel
	StatementGoTo             *StatementGoTo
	Defer                     *StatementDefer
	Class                     *Class
	Enum                      *Enum
	Import                    *Import
	isStaticFieldDefaultValue bool
	/*
		this.super()
		special case
	*/
	IsCallFatherConstructionStatement bool
}

func (s *Statement) isVariableDefinition() bool {
	if s.Type != StatementTypeExpression {
		return false
	}
	return s.Expression.Type == ExpressionTypeVarAssign ||
		s.Expression.Type == ExpressionTypeVar
}

func (s *Statement) simplifyIf() {
	if len(s.StatementIf.ElseIfList) > 0 {
		return
	}
	if len(s.StatementIf.PrefixExpressions) > 0 {
		return
	}
	if s.StatementIf.Condition.Type != ExpressionTypeBool {
		return
	}
	c := s.StatementIf.Condition.Data.(bool)
	if c {
		s.Type = StatementTypeBlock
		s.Block = &s.StatementIf.Block
	} else {
		if s.StatementIf.Else != nil {
			s.Type = StatementTypeBlock
			s.Block = s.StatementIf.Else
		} else {
			s.Type = StatementTypeNop
		}
	}
}

func (s *Statement) simplifyFor() {
	if s.StatementFor.Init == nil &&
		s.StatementFor.Increment == nil &&
		s.StatementFor.Condition != nil &&
		s.StatementFor.Condition.Type == ExpressionTypeBool &&
		s.StatementFor.Condition.Data.(bool) == false {
		s.Type = StatementTypeNop
		s.StatementFor = nil
	}
}
func (s *Statement) check(block *Block) []error {
	defer func() {
		s.Checked = true
	}()
	errs := []error{}
	switch s.Type {
	case StatementTypeExpression:
		return s.checkStatementExpression(block)
	case StatementTypeIf:
		es := s.StatementIf.check(block)
		s.simplifyIf()
		return es
	case StatementTypeFor:
		es := s.StatementFor.check(block)
		s.simplifyFor()
		return es
	case StatementTypeSwitch:
		return s.StatementSwitch.check(block)
	case StatementTypeBreak:
		return s.StatementBreak.check(block)
	case StatementTypeContinue:
		return s.StatementContinue.check(block)
	case StatementTypeReturn:
		return s.StatementReturn.check(block)
	case StatementTypeGoTo:
		err := s.StatementGoTo.checkStatementGoTo(block)
		if err != nil {
			return []error{err}
		}
	case StatementTypeDefer:
		block.InheritedAttribute.Function.HasDefer = true
		s.Defer.Block.inherit(block)
		s.Defer.Block.InheritedAttribute.Defer = s.Defer
		es := s.Defer.Block.check()
		block.Defers = append(block.Defers, s.Defer)
		return es
	case StatementTypeBlock:
		s.Block.inherit(block)
		return s.Block.check()
	case StatementTypeLabel:
		if block.InheritedAttribute.Defer != nil {
			block.InheritedAttribute.Defer.Labels =
				append(block.InheritedAttribute.Defer.Labels, s.StatementLabel)
		}
	case StatementTypeClass:
		PackageBeenCompile.statementLevelClass =
			append(PackageBeenCompile.statementLevelClass, s.Class)
		err := block.Insert(s.Class.Name, s.Class.Pos, s.Class)
		if err != nil {
			errs = append(errs, err)
		}
		return append(errs, s.Class.check(block)...)
	case StatementTypeEnum:
		es := s.Enum.check()
		err := block.Insert(s.Enum.Name, s.Enum.Pos, s.Enum)
		if err != nil {
			es = append(es, err)
		}
		return es
	case StatementTypeNop:
		//nop , should be never execute to here
		//
	case StatementTypeWhen:
		return s.StatementWhen.check(block, s)
	case StatementTypeImport:
		if block.InheritedAttribute.Function.TemplateClonedFunction == false {
			errs = append(errs, fmt.Errorf("%s cannot have 'import' at this scope , non-template function",
				errMsgPrefix(s.Import.Pos)))
			return errs
		}
		err := s.Import.MkAccessName()
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		if s.Import.Alias == UnderScore {
			errs = append(errs, fmt.Errorf("%s import at block scope , must be used",
				errMsgPrefix(s.Import.Pos)))
			return nil
		}
		if err := PackageBeenCompile.insertImport(s.Import); err != nil {
			errs = append(errs, err)
		}
	case StatementTypeTypeAlias:
		err := s.TypeAlias.Type.resolve(block)
		if err != nil {
			return []error{err}
		}
		err = block.Insert(s.TypeAlias.Name, s.TypeAlias.Pos, s.TypeAlias.Type)
		if err != nil {
			return []error{err}
		}
		return nil
	}
	return nil
}
