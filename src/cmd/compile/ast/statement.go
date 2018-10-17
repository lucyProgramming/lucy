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
	StatementTypeSwitchTemplate
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
	StatementSwitchTemplate   *StatementSwitchTemplate
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

func (s *Statement) check(block *Block) []error {
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
		return s.StatementBreak.check(s, block)
	case StatementTypeContinue:
		return s.StatementContinue.check(s, block)
	case StatementTypeReturn:
		return s.StatementReturn.check(s, block)
	case StatementTypeGoTo:
		err := s.checkStatementGoTo(block)
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
		err := block.Insert(s.Class.Name, s.Pos, s.Class)
		if err != nil {
			errs = append(errs, err)
		}
		return append(errs, s.Class.check(block)...)
	case StatementTypeEnum:
		es := s.Enum.check()
		err := block.Insert(s.Enum.Name, s.Pos, s.Enum)
		if err != nil {
			es = append(es, err)
		}
		return es
	case StatementTypeNop:
		//nop , should be never execute to here
		//
	case StatementTypeSwitchTemplate:
		return s.StatementSwitchTemplate.check(block, s)
	case StatementTypeImport:
		if block.InheritedAttribute.Function.TemplateClonedFunction == false {
			errs = append(errs, fmt.Errorf("%s cannot have 'import' at this scope , non-template function",
				errMsgPrefix(s.Import.Pos)))
			return errs
		}
		err := s.Import.MkAccessName()
		if err != nil {
			errs = append(errs, fmt.Errorf("%s %v",
				errMsgPrefix(s.Import.Pos), err))
			return errs
		}
		if s.Import.AccessName == NoNameIdentifier {
			errs = append(errs, fmt.Errorf("%s import at block scope , must be used",
				errMsgPrefix(s.Import.Pos)))
			return nil
		}
		if PackageBeenCompile.Files == nil {
			PackageBeenCompile.Files = make(map[string]*SourceFile)
		}
		if PackageBeenCompile.Files[s.Pos.Filename] == nil {
			PackageBeenCompile.Files[s.Pos.Filename] = &SourceFile{}
		}
		if PackageBeenCompile.Files[s.Pos.Filename].Imports == nil {
			PackageBeenCompile.Files[s.Pos.Filename].Imports = make(map[string]*Import)
		}
		is := PackageBeenCompile.Files[s.Pos.Filename]
		if is.Imports == nil {
			is.Imports = make(map[string]*Import)
		}
		if is.ImportsByResources == nil {
			is.ImportsByResources = make(map[string]*Import)
		}
		_, ok := is.Imports[s.Import.AccessName]
		if ok {
			errs = append(errs, fmt.Errorf("%s package '%s' reimported",
				errMsgPrefix(s.Import.Pos), s.Import.AccessName))
			return nil
		}
		_, ok = is.ImportsByResources[s.Import.Import]
		if ok {
			errs = append(errs, fmt.Errorf("%s package '%s' reimported",
				errMsgPrefix(s.Import.Pos), s.Import.Import))
			return nil
		}
		is.Imports[s.Import.AccessName] = s.Import
		is.ImportsByResources[s.Import.Import] = s.Import

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
