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

func (this *Statement) isVariableDefinition() bool {
	if this.Type != StatementTypeExpression {
		return false
	}
	return this.Expression.Type == ExpressionTypeVarAssign ||
		this.Expression.Type == ExpressionTypeVar
}

//func (this *Statement) simplifyIf() {
//	if len(this.StatementIf.ElseIfList) > 0 {
//		return
//	}
//	if len(this.StatementIf.PrefixExpressions) > 0 {
//		return
//	}
//	if this.StatementIf.Condition.Type != ExpressionTypeBool {
//		return
//	}
//	c := this.StatementIf.Condition.Data.(bool)
//	if c {
//		this.Type = StatementTypeBlock
//		this.Block = &this.StatementIf.Block
//	} else {
//		if this.StatementIf.Else != nil {
//			this.Type = StatementTypeBlock
//			this.Block = this.StatementIf.Else
//		} else {
//			this.Type = StatementTypeNop
//		}
//	}
//}
//
//func (this *Statement) simplifyFor() {
//	if this.StatementFor.Init == nil &&
//		this.StatementFor.Increment == nil &&
//		this.StatementFor.Condition != nil &&
//		this.StatementFor.Condition.Type == ExpressionTypeBool &&
//		this.StatementFor.Condition.Data.(bool) == false {
//		this.Type = StatementTypeNop
//		this.StatementFor = nil
//	}
//}
func (this *Statement) check(block *Block) []error {
	defer func() {
		this.Checked = true
	}()
	errs := []error{}
	switch this.Type {
	case StatementTypeExpression:
		return this.checkStatementExpression(block)
	case StatementTypeIf:
		es := this.StatementIf.check(block)
		return es
	case StatementTypeFor:
		es := this.StatementFor.check(block)
		return es
	case StatementTypeSwitch:
		return this.StatementSwitch.check(block)
	case StatementTypeBreak:
		return this.StatementBreak.check(block)
	case StatementTypeContinue:
		return this.StatementContinue.check(block)
	case StatementTypeReturn:
		return this.StatementReturn.check(block)
	case StatementTypeGoTo:
		err := this.StatementGoTo.checkStatementGoTo(block)
		if err != nil {
			return []error{err}
		}
	case StatementTypeDefer:
		block.InheritedAttribute.Function.HasDefer = true
		this.Defer.Block.inherit(block)
		this.Defer.Block.InheritedAttribute.Defer = this.Defer
		es := this.Defer.Block.check()
		block.Defers = append(block.Defers, this.Defer)
		return es
	case StatementTypeBlock:
		this.Block.inherit(block)
		return this.Block.check()
	case StatementTypeLabel:
		if block.InheritedAttribute.Defer != nil {
			block.InheritedAttribute.Defer.Labels =
				append(block.InheritedAttribute.Defer.Labels, this.StatementLabel)
		}
	case StatementTypeClass:
		PackageBeenCompile.statementLevelClass =
			append(PackageBeenCompile.statementLevelClass, this.Class)
		err := block.Insert(this.Class.Name, this.Class.Pos, this.Class)
		if err != nil {
			errs = append(errs, err)
		}
		return append(errs, this.Class.check(block)...)
	case StatementTypeEnum:
		es := this.Enum.check()
		err := block.Insert(this.Enum.Name, this.Enum.Pos, this.Enum)
		if err != nil {
			es = append(es, err)
		}
		return es
	case StatementTypeNop:
		//nop , should be never execute to here
		//
	case StatementTypeWhen:
		return this.StatementWhen.check(block, this)
	case StatementTypeImport:
		if block.InheritedAttribute.Function.TemplateClonedFunction == false {
			errs = append(errs, fmt.Errorf("%s cannot have 'import' at this scope , non-template function",
				errMsgPrefix(this.Import.Pos)))
			return errs
		}
		err := this.Import.MkAccessName()
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		if this.Import.Alias == UnderScore {
			errs = append(errs, fmt.Errorf("%s import at block scope , must be used",
				errMsgPrefix(this.Import.Pos)))
			return nil
		}
		if err := PackageBeenCompile.insertImport(this.Import); err != nil {
			errs = append(errs, err)
		}
	case StatementTypeTypeAlias:
		err := this.TypeAlias.Type.resolve(block)
		if err != nil {
			return []error{err}
		}
		err = block.Insert(this.TypeAlias.Name, this.TypeAlias.Pos, this.TypeAlias.Type)
		if err != nil {
			return []error{err}
		}
		return nil
	}
	return nil
}
