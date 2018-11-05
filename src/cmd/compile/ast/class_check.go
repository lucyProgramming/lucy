package ast

import (
	"errors"
	"fmt"
	"path/filepath"
)

func (this *Class) check(father *Block) []error {
	this.Block.inherit(father)
	this.Block.InheritedAttribute.Class = this
	errs := this.checkPhase1()
	errs = append(errs, this.checkPhase2()...)
	return errs
}

func (this *Class) checkPhase1() []error {
	if this.Block.InheritedAttribute.ClassAndFunctionNames == "" {
		this.Block.InheritedAttribute.ClassAndFunctionNames = filepath.Base(this.Name)
	} else {
		this.Block.InheritedAttribute.ClassAndFunctionNames += "$" + filepath.Base(this.Name)
	}
	this.mkDefaultConstruction()
	errs := this.Block.checkConstants()
	err := this.resolveFather()
	if err != nil {
		errs = append(errs, err)
	} else {
		err = this.checkIfClassHierarchyErr()
		if err != nil {
			errs = append(errs, err)
		}
	}
	errs = append(errs, this.checkModifierOk()...)
	errs = append(errs, this.resolveFieldsAndMethodsType()...)
	return errs
}

func (this *Class) checkPhase2() []error {
	errs := []error{}
	errs = append(errs, this.checkFields()...)
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	this.mkClassInitMethod()
	for name, ms := range this.Methods {
		if this.Fields != nil && this.Fields[name] != nil {
			f := this.Fields[name]
			if f.Pos.Line < ms[0].Function.Pos.Line {
				errMsg := fmt.Sprintf("%s method named '%s' already declared as field,at:\n",
					errMsgPrefix(ms[0].Function.Pos), name)
				errMsg += fmt.Sprintf("\t%s", errMsgPrefix(f.Pos))
				errs = append(errs, errors.New(errMsg))
			} else {
				errMsg := fmt.Sprintf("%s field named '%s' already declared as method,at:\n",
					errMsgPrefix(f.Pos), name)
				errMsg += fmt.Sprintf("\t%s", errMsgPrefix(ms[0].Function.Pos))
				errs = append(errs, errors.New(errMsg))
			}
			continue
		}
		if len(ms) > 1 {
			errMsg := fmt.Sprintf("%s class method named '%s' has declared %d times,which are:\n",
				errMsgPrefix(ms[0].Function.Pos),
				ms[0].Function.Name, len(ms))
			for _, v := range ms {
				errMsg += fmt.Sprintf("\t%s\n", errMsgPrefix(v.Function.Pos))
			}
			errs = append(errs, errors.New(errMsg))
		}
	}
	errs = append(errs, this.checkMethods()...)
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	errs = append(errs, this.checkIfOverrideFinalMethod()...)
	errs = append(errs, this.resolveInterfaces()...)
	if this.IsInterface() {
		errs = append(errs, this.checkOverrideInterfaceMethod()...)
	}
	if this.IsAbstract() {
		errs = append(errs, this.checkOverrideAbstractMethod()...)
	}
	errs = append(errs, this.suitableForInterfaces()...)
	if this.SuperClass != nil {
		errs = append(errs, this.suitableSubClassForAbstract(this.SuperClass)...)
	}
	return errs
}

func (this *Class) suitableSubClassForAbstract(super *Class) []error {
	errs := []error{}
	if super.Name != JavaRootClass {
		err := super.loadSuperClass(this.Pos)
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		if super.SuperClass == nil {
			return errs
		}
		length := len(errs)
		errs = append(errs, this.suitableSubClassForAbstract(super.SuperClass)...)
		if len(errs) > length {
			return errs
		}
	}
	if super.IsAbstract() {
		for _, v := range super.Methods {
			m := v[0]
			if m.IsAbstract() == false {
				continue
			}
			var nameMatch *ClassMethod
			implementation := this.implementMethod(this.Pos, m, &nameMatch, false, &errs)
			if implementation != nil {
				if err := m.implementationMethodIsOk(this.Pos, implementation); err != nil {
					errs = append(errs, err)
				}
			} else {
				pos := this.Pos
				if nameMatch != nil && nameMatch.Function.Pos != nil {
					pos = nameMatch.Function.Pos
				}
				if nameMatch != nil {
					errMsg := fmt.Sprintf("%s method is suitable for abstract super class\n", errMsgPrefix(pos))
					errMsg += fmt.Sprintf("\t have %s\n", nameMatch.Function.readableMsg())
					errMsg += fmt.Sprintf("\t want %s\n", m.Function.readableMsg())
					errs = append(errs, errors.New(errMsg))
				} else {
					errs = append(errs,
						fmt.Errorf("%s missing implementation method '%s' define on abstract class '%s'",
							pos.ErrMsgPrefix(), m.Function.readableMsg(), super.Name))
				}
			}
		}
	}
	return errs
}

func (this *Class) interfaceMethodExists(name string) *Class {
	if this.IsInterface() == false {
		panic("not a interface")
	}
	if this.Methods != nil && len(this.Methods[name]) > 0 {
		return this
	}
	for _, v := range this.Interfaces {
		if v.interfaceMethodExists(name) != nil {
			return v
		}
	}
	return nil
}

func (this *Class) abstractMethodExists(pos *Pos, name string) (*Class, error) {
	if this.IsAbstract() {
		if this.Methods != nil && len(this.Methods[name]) > 0 {
			method := this.Methods[name][0]
			if method.IsAbstract() {
				return this, nil
			}
		}
	}
	if this.Name == JavaRootClass {
		return nil, nil
	}
	err := this.loadSuperClass(pos)
	if err != nil {
		return nil, err
	}
	if this.SuperClass == nil {
		return nil, nil
	}
	return this.SuperClass.abstractMethodExists(pos, name)
}

func (this *Class) checkOverrideAbstractMethod() []error {
	errs := []error{}
	err := this.loadSuperClass(this.Pos)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	if this.SuperClass == nil {
		return errs
	}
	for _, v := range this.Methods {
		m := v[0]
		name := m.Function.Name
		if m.IsAbstract() == false {
			continue
		}
		exist, err := this.SuperClass.abstractMethodExists(m.Function.Pos, name)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if exist != nil {
			errs = append(errs, fmt.Errorf("%s method '%s' override '%s'",
				errMsgPrefix(v[0].Function.Pos), name, exist.Name))
		}
	}
	return errs
}

func (this *Class) checkOverrideInterfaceMethod() []error {
	errs := []error{}
	for name, v := range this.Methods {
		var exist *Class
		for _, vv := range this.Interfaces {
			exist = vv.interfaceMethodExists(name)
			if exist != nil {
				break
			}
		}
		if exist != nil {
			errs = append(errs, fmt.Errorf("%s method '%s' override '%s'",
				errMsgPrefix(v[0].Function.Pos), name, exist.Name))
		}
	}
	return errs
}

func (this *Class) checkIfClassHierarchyErr() error {
	m := make(map[string]struct{})
	arr := []string{}
	is := false
	class := this
	pos := this.Pos
	if err := this.loadSuperClass(pos); err != nil {
		return err
	}
	if this.SuperClass == nil {
		return nil
	}

	if this.SuperClass.IsFinal() {
		return fmt.Errorf("%s class name '%s' have super class  named '%s' that is final",
			this.Pos.ErrMsgPrefix(), this.Name, this.SuperClass.Name)
	}
	for class.Name != JavaRootClass {
		_, ok := m[class.Name]
		if ok {
			arr = append(arr, class.Name)
			is = true
			break
		}
		m[class.Name] = struct{}{}
		arr = append(arr, class.Name)
		err := class.loadSuperClass(pos)
		if err != nil {
			return err
		}
		if this.SuperClass == nil {
			return nil
		}
		class = class.SuperClass
	}
	if is == false {
		return nil
	}
	errMsg := fmt.Sprintf("%s class named '%s' detects a circularity in class hierarchy",
		this.Pos.ErrMsgPrefix(), this.Name)
	tab := "\t"
	index := len(arr) - 1
	for index >= 0 {
		errMsg += tab + arr[index] + "\n"
		tab += " "
		index--
	}
	return fmt.Errorf(errMsg)
}

func (this *Class) checkIfOverrideFinalMethod() []error {
	errs := []error{}
	if this.SuperClass != nil {
		for name, v := range this.Methods {
			if name == SpecialMethodInit {
				continue
			}
			if len(v) == 0 {
				continue
			}
			if len(this.SuperClass.Methods[name]) == 0 {
				// this class not found at super
				continue
			}
			m := v[0]
			for _, v := range this.SuperClass.Methods[name] {
				if v.IsFinal() == false {
					continue
				}
				f1 := &Type{
					Type:         VariableTypeFunction,
					FunctionType: &m.Function.Type,
				}
				f2 := &Type{
					Type:         VariableTypeFunction,
					FunctionType: &v.Function.Type,
				}
				if f1.Equal(f2) {
					errs = append(errs, fmt.Errorf("%s override final method",
						errMsgPrefix(m.Function.Pos)))
				}
			}
		}
	}
	return errs
}

func (this *Class) suitableForInterfaces() []error {
	errs := []error{}
	if this.IsInterface() {
		return errs
	}
	for _, i := range this.Interfaces {
		errs = append(errs, this.suitableForInterface(i)...)
	}
	return errs
}

func (this *Class) suitableForInterface(inter *Class) []error {
	errs := []error{}
	err := inter.loadSelf(this.Pos)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	for _, v := range inter.Methods {
		m := v[0]
		var nameMatch *ClassMethod
		implementation := this.implementMethod(this.Pos, m, &nameMatch, false, &errs)
		if implementation != nil {
			if err := m.implementationMethodIsOk(this.Pos, implementation); err != nil {
				errs = append(errs, err)
			}
		} else {
			pos := this.Pos
			if nameMatch != nil && nameMatch.Function.Pos != nil {
				pos = nameMatch.Function.Pos
			}
			errs = append(errs, fmt.Errorf("%s missing implementation method '%s' define on interface '%s'",
				pos.ErrMsgPrefix(), m.Function.readableMsg(), inter.Name))
		}
	}
	for _, v := range inter.Interfaces {
		es := this.suitableForInterface(v)
		errs = append(errs, es...)
	}
	return errs
}

func (this *Class) checkFields() []error {
	errs := []error{}
	if this.IsInterface() {
		for _, v := range this.Fields {
			errs = append(errs, fmt.Errorf("%s interface '%s' expect no field named '%s'",
				errMsgPrefix(v.Pos), this.Name, v.Name))
		}
		return errs
	}
	staticFieldAssignStatements := []*Statement{}
	for _, v := range this.Fields {
		if v.DefaultValueExpression != nil {
			assignment, es := v.DefaultValueExpression.
				checkSingleValueContextExpression(&this.Methods[SpecialMethodInit][0].Function.Block)
			errs = append(errs, es...)
			if assignment == nil {
				continue
			}
			if v.Type.assignAble(&errs, assignment) == false {
				errs = append(errs, fmt.Errorf("%s cannot assign '%s' as '%s' for default value",
					errMsgPrefix(v.Pos), assignment.TypeString(), v.Type.TypeString()))
				continue
			}
			if assignment.Type == VariableTypeNull {
				errs = append(errs, fmt.Errorf("%s pointer types default value is '%s' already",
					v.Pos.ErrMsgPrefix(), assignment.TypeString()))
				continue
			}
			if v.IsStatic() &&
				v.DefaultValueExpression.isLiteral() {
				v.DefaultValue = v.DefaultValueExpression.Data
				continue
			}
			if v.IsStatic() == false {
				// nothing to do
				continue
			}
			bin := &ExpressionBinary{}
			bin.Right = &Expression{
				Type: ExpressionTypeList,
				Op:   "list",
				Data: []*Expression{v.DefaultValueExpression},
			}
			{
				selection := &ExpressionSelection{}
				selection.Expression = &Expression{}
				selection.Expression.Op = "selection"
				selection.Expression.Value = &Type{
					Type:  VariableTypeClass,
					Class: this,
				}
				selection.Name = v.Name
				selection.Field = v
				left := &Expression{
					Type: ExpressionTypeSelection,
					Data: selection,
					Op:   "selection",
				}
				left.Value = v.Type
				bin.Left = &Expression{
					Type: ExpressionTypeList,
					Data: []*Expression{left},
				}
			}
			e := &Expression{
				Type: ExpressionTypeAssign,
				Data: bin,
				IsStatementExpression: true,
				Op: "assign",
			}
			staticFieldAssignStatements = append(staticFieldAssignStatements, &Statement{
				Type:                      StatementTypeExpression,
				Expression:                e,
				isStaticFieldDefaultValue: true,
			})
		}
	}
	if len(staticFieldAssignStatements) > 0 {
		b := &Block{}
		b.Statements = staticFieldAssignStatements
		if this.StaticBlocks != nil {
			this.StaticBlocks = append([]*Block{b}, this.StaticBlocks...)
		} else {
			this.StaticBlocks = []*Block{b}
		}
	}
	return errs
}

func (this *Class) checkMethods() []error {
	errs := []error{}
	if this.IsInterface() {
		return errs
	}
	for name, methods := range this.Methods {
		for _, method := range methods {
			errs = append(errs, method.checkModifierOk()...)
			if method.IsAbstract() {
				//nothing
			} else {
				if this.IsInterface() {
					errs = append(errs, fmt.Errorf("%s interface method cannot have implementation",
						errMsgPrefix(method.Function.Pos)))
					continue
				}
				errs = append(errs, method.Function.checkReturnVarExpression()...)
				isConstruction := name == SpecialMethodInit
				if isConstruction {
					if method.IsFirstStatementCallFatherConstruction() == false {
						errs = append(errs, fmt.Errorf("%s construction method should call father construction method first",
							errMsgPrefix(method.Function.Pos)))
					}
				}
				if isConstruction && method.Function.Type.VoidReturn() == false {
					errs = append(errs, fmt.Errorf("%s construction method expect no return values",
						errMsgPrefix(method.Function.Type.ParameterList[0].Pos)))
				}
				method.Function.Block.InheritedAttribute.IsConstructionMethod = isConstruction
				method.Function.checkBlock(&errs)
			}
		}
	}
	return errs
}

func (this *Class) checkModifierOk() []error {
	errs := []error{}
	if this.IsInterface() && this.IsFinal() {
		errs = append(errs, fmt.Errorf("%s interface '%s' cannot be final",
			errMsgPrefix(this.FinalPos), this.Name))
	}
	if this.IsAbstract() && this.IsFinal() {
		errs = append(errs, fmt.Errorf("%s abstract class '%s' cannot be final",
			errMsgPrefix(this.FinalPos), this.Name))
	}
	return errs
}
