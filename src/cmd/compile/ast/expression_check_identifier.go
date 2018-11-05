package ast

import (
	"fmt"
	"time"
)

func (this *Expression) checkIdentifierExpression(block *Block) (*Type, error) {
	identifier := this.Data.(*ExpressionIdentifier)
	if identifier.Name == UnderScore {
		//_ is not valid
		return nil, fmt.Errorf("%s '%s' is not a valid name",
			this.Pos.ErrMsgPrefix(), identifier.Name)
	}
	//handle magic identifier
	switch identifier.Name {
	case magicIdentifierFile:
		this.Type = ExpressionTypeString
		this.Data = this.Pos.Filename
		result, _ := this.checkSingleValueContextExpression(block)
		return result, nil
	case magicIdentifierLine:
		this.Type = ExpressionTypeInt
		this.Data = int64(this.Pos.Line)
		result, _ := this.checkSingleValueContextExpression(block)
		return result, nil
	case magicIdentifierTime:
		this.Type = ExpressionTypeLong
		this.Data = int64(time.Now().UnixNano())
		result, _ := this.checkSingleValueContextExpression(block)
		return result, nil
	case magicIdentifierClass:
		if block.InheritedAttribute.Class == nil {
			return nil,
				fmt.Errorf("%s '%s' must in class scope", this.Pos.ErrMsgPrefix(), identifier.Name)
		}
		result := &Type{}
		result.Type = VariableTypeClass
		result.Pos = this.Pos
		result.Class = block.InheritedAttribute.Class
		return result, nil
	case magicIdentifierFunction:
		if block.InheritedAttribute.Function.isPackageInitBlockFunction {
			return nil, fmt.Errorf("%s '%s' must in function scope", this.Pos.ErrMsgPrefix(), identifier.Name)
		}
		result := &Type{}
		result.Type = VariableTypeMagicFunction
		result.Pos = this.Pos
		result.Function = block.InheritedAttribute.Function
		return result, nil
	}
	isCaptureVar := false
	d, err := block.searchIdentifier(this.Pos, identifier.Name, &isCaptureVar)
	if err != nil {
		return nil, err
	}
	if d == nil {
		i := PackageBeenCompile.getImport(this.Pos.Filename, identifier.Name)
		if i != nil {
			i.Used = true
			return this.checkIdentifierThroughImports(i)
		}
	}
	if d == nil {
		return nil, fmt.Errorf("%s '%s' not found", this.Pos.ErrMsgPrefix(), identifier.Name)
	}
	switch d.(type) {
	case *Function:
		f := d.(*Function)
		if f.IsGlobalMain() {
			// not allow
			return nil, fmt.Errorf("%s fucntion is global main", errMsgPrefix(this.Pos))
		}
		if f.IsBuildIn {
			return nil, fmt.Errorf("%s fucntion '%s' is buildin",
				this.Pos.ErrMsgPrefix(), f.Name)
		}
		if f.TemplateFunction != nil {
			return nil, fmt.Errorf("%s fucntion '%s' a template function",
				this.Pos.ErrMsgPrefix(), f.Name)
		}
		// try from import
		if f.IsBuildIn == false {
			i, should := shouldAccessFromImports(identifier.Name, this.Pos, f.Pos)
			if should {
				return this.checkIdentifierThroughImports(i)
			}
		}
		result := &Type{}
		result.Type = VariableTypeFunction
		result.FunctionType = &f.Type
		result.Pos = this.Pos
		identifier.Function = f
		return result, nil
	case *Variable:
		t := d.(*Variable)
		if t.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, this.Pos, t.Pos)
			if should {
				return this.checkIdentifierThroughImports(i)
			}
		}
		if isCaptureVar {
			t.BeenCapturedAsRightValue++
		}
		t.Used = true
		result := t.Type.Clone()
		result.Pos = this.Pos
		identifier.Variable = t
		return result, nil
	case *Constant:
		t := d.(*Constant)
		if t.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, this.Pos, t.Pos)
			if should {
				return this.checkIdentifierThroughImports(i)
			}
		}
		t.Used = true
		this.fromConst(t)
		result := t.Type.Clone()
		result.Pos = this.Pos
		return result, nil
	case *Class:
		c := d.(*Class)
		if c.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, this.Pos, c.Pos)
			if should {
				return this.checkIdentifierThroughImports(i)
			}
		}
		result := &Type{}
		result.Type = VariableTypeClass
		result.Pos = this.Pos
		result.Class = c
		return result, nil
	case *EnumName:
		enumName := d.(*EnumName)
		if enumName.Enum.IsBuildIn == false { // try from import
			i, should := shouldAccessFromImports(identifier.Name, this.Pos, enumName.Pos)
			if should {
				return this.checkIdentifierThroughImports(i)
			}
		}
		result := &Type{}
		result.Pos = enumName.Pos
		result.Type = VariableTypeEnum
		result.EnumName = enumName
		result.Enum = enumName.Enum
		identifier.EnumName = enumName
		return result, nil
	}
	return nil, fmt.Errorf("%s identifier '%s' is not a expression , but '%s'",
		this.Pos.ErrMsgPrefix(), identifier.Name, block.identifierIsWhat(d))
}

func (this *Expression) checkIdentifierThroughImports(it *Import) (*Type, error) {
	p, err := PackageBeenCompile.load(it.Import)
	if err != nil {
		return nil, fmt.Errorf("%s %v", this.Pos.ErrMsgPrefix(), err)
	}
	result := &Type{}
	result.Pos = this.Pos
	if pp, ok := p.(*Package); ok {
		result.Package = pp
		result.Type = VariableTypePackage
	} else {
		result.Class = p.(*Class)
		result.Type = VariableTypeClass
	}
	return result, nil
}
