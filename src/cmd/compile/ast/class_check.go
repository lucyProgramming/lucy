package ast

import (
	"errors"
	"fmt"
	"path/filepath"
)

func (c *Class) check(father *Block) []error {
	c.Block.inherit(father)
	c.Block.InheritedAttribute.Class = c
	errs := c.checkPhase1()
	es := c.checkPhase2()
	errs = append(errs, es...)
	return errs
}

func (c *Class) checkPhase1() []error {
	c.mkDefaultConstruction()
	errs := c.Block.checkConstants()
	err := c.resolveFather()
	if err != nil {
		errs = append(errs, err)
	} else {
		err = c.checkIfClassHierarchyErr()
		if err != nil {
			errs = append(errs, err)
		}
	}
	errs = append(errs, c.checkModifierOk()...)
	es := c.resolveFieldsAndMethodsType()

	errs = append(errs, es...)

	return errs
}

func (c *Class) checkPhase2() []error {
	errs := []error{}
	if c.Block.InheritedAttribute.ClassAndFunctionNames == "" {
		c.Block.InheritedAttribute.ClassAndFunctionNames = filepath.Base(c.Name)
	} else {
		c.Block.InheritedAttribute.ClassAndFunctionNames += "$" + filepath.Base(c.Name)
	}
	errs = append(errs, c.checkFields()...)
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	c.mkClassInitMethod()
	for name, ms := range c.Methods {
		if c.Fields != nil && c.Fields[name] != nil {
			f := c.Fields[name]
			if f.Pos.StartLine < ms[0].Function.Pos.StartLine {
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
	errs = append(errs, c.checkMethods()...)
	if PackageBeenCompile.shouldStop(errs) {
		return errs
	}
	errs = append(errs, c.checkIfOverrideFinalMethod()...)
	errs = append(errs, c.resolveInterfaces()...)
	if c.IsInterface() {
		errs = append(errs, c.checkOverrideInterfaceMethod()...)
	}
	if c.IsAbstract() {
		errs = append(errs, c.checkOverrideAbstractMethod()...)
	}
	errs = append(errs, c.suitableForInterfaces()...)
	err := c.loadSuperClass(c.Pos)
	if err != nil {
		errs = append(errs, err)
	} else {
		errs = append(errs, c.suitableSubClassForAbstract(c.SuperClass)...)
	}
	return errs
}

func (c *Class) suitableSubClassForAbstract(super *Class) []error {
	errs := []error{}
	if super.Name != JavaRootClass {
		err := super.loadSuperClass(c.Pos)
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		length := len(errs)
		errs = append(errs, c.suitableSubClassForAbstract(super.SuperClass)...)
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
			implementation := c.implementMethod(c.Pos, m, false, &errs)
			if implementation != nil {
				if err := m.implementationMethodIsOk(c.Pos, implementation); err != nil {
					errs = append(errs, err)
				}
			} else {
				errs = append(errs, fmt.Errorf("%s missing implementation method '%s' define on abstract class '%s'",
					errMsgPrefix(c.Pos), m.Function.readableMsg(), super.Name))
			}
		}
	}
	return errs
}

func (c *Class) interfaceMethodExists(name string) *Class {
	if c.IsInterface() == false {
		panic("not a interface")
	}
	if c.Methods != nil && len(c.Methods[name]) > 0 {
		return c
	}
	for _, v := range c.Interfaces {
		if v.interfaceMethodExists(name) != nil {
			return v
		}
	}
	return nil
}

func (c *Class) abstractMethodExists(pos *Pos, name string) (*Class, error) {
	if c.IsAbstract() {
		if c.Methods != nil && len(c.Methods[name]) > 0 {
			method := c.Methods[name][0]
			if method.IsAbstract() {
				return c, nil
			}
		}
	}
	if c.Name == JavaRootClass {
		return nil, nil
	}
	err := c.loadSuperClass(pos)
	if err != nil {
		return nil, err
	}
	return c.SuperClass.abstractMethodExists(pos, name)
}

func (c *Class) checkOverrideAbstractMethod() []error {
	errs := []error{}
	err := c.loadSuperClass(c.Pos)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	for _, v := range c.Methods {
		m := v[0]
		name := m.Function.Name
		if m.IsAbstract() == false {
			continue
		}
		exist, err := c.SuperClass.abstractMethodExists(m.Function.Pos, name)
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

func (c *Class) checkOverrideInterfaceMethod() []error {
	errs := []error{}
	for name, v := range c.Methods {
		var exist *Class
		for _, vv := range c.Interfaces {
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

func (c *Class) checkIfClassHierarchyErr() error {
	m := make(map[string]struct{})
	arr := []string{}
	is := false
	class := c
	pos := c.Pos
	if err := c.loadSuperClass(pos); err != nil {
		return err
	}
	if c.SuperClass.LoadFromOutSide && c.SuperClass.IsPublic() == false {
		return fmt.Errorf("%s class`s super-class named '%s' is not public",
			errMsgPrefix(c.Pos), c.SuperClass.Name)
	}
	if c.SuperClass.IsFinal() {
		return fmt.Errorf("%s class name '%s' have super class  named '%s' that is final",
			errMsgPrefix(c.Pos), c.Name, c.SuperClassName)
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
		class = class.SuperClass
	}
	if is == false {
		return nil
	}
	errMsg := fmt.Sprintf("%s class named '%s' detects a circularity in class hierarchy",
		errMsgPrefix(c.Pos), c.Name)
	tab := "\t"
	index := len(arr) - 1
	for index >= 0 {
		errMsg += tab + arr[index] + "\n"
		tab += " "
		index--
	}
	return fmt.Errorf(errMsg)
}

func (c *Class) checkIfOverrideFinalMethod() []error {
	err := c.loadSuperClass(c.Pos)
	if err != nil {
		return []error{err}
	}
	errs := []error{}
	for name, v := range c.Methods {
		if name == SpecialMethodInit {
			continue
		}
		if len(v) == 0 {
			continue
		}
		if len(c.SuperClass.Methods[name]) == 0 {
			// this class not found at super
			continue
		}
		m := v[0]
		for _, v := range c.SuperClass.Methods[name] {
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
			if f1.Equal(&errs, f2) == true {
				errs = append(errs, fmt.Errorf("%s override final method",
					errMsgPrefix(m.Function.Pos)))
			}
		}
	}
	return errs
}

func (c *Class) suitableForInterfaces() []error {
	errs := []error{}
	if c.IsInterface() {
		return errs
	}
	for _, i := range c.Interfaces {
		errs = append(errs, c.suitableForInterface(i)...)
	}
	return errs
}

func (c *Class) suitableForInterface(inter *Class) []error {
	errs := []error{}
	err := inter.loadSelf(c.Pos)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	for _, v := range inter.Methods {
		m := v[0]
		implementation := c.implementMethod(c.Pos, m, false, &errs)
		if implementation != nil {
			if err := m.implementationMethodIsOk(c.Pos, implementation); err != nil {
				errs = append(errs, err)
			}
		} else {
			errs = append(errs, fmt.Errorf("%s missing implementation method '%s' define on interface '%s'",
				errMsgPrefix(c.Pos), m.Function.readableMsg(), inter.Name))
		}
	}
	for _, v := range inter.Interfaces {
		es := c.suitableForInterface(v)
		errs = append(errs, es...)
	}
	return errs
}

func (c *Class) checkFields() []error {
	errs := []error{}
	if c.IsInterface() {
		for _, v := range c.Fields {
			errs = append(errs, fmt.Errorf("%s interface '%s' expect no field named '%s'",
				errMsgPrefix(v.Pos), c.Name, v.Name))
		}
		return errs
	}
	staticFieldAssignStatements := []*Statement{}
	for _, v := range c.Fields {
		if v.Expression != nil {
			t, es := v.Expression.checkSingleValueContextExpression(&c.Block)
			errs = append(errs, es...)
			if v.Type.Equal(&errs, t) == false {
				errs = append(errs, fmt.Errorf("%s cannot assign '%s' as '%s' for default value",
					errMsgPrefix(v.Pos), t.TypeString(), v.Type.TypeString()))
				continue
			}
			//TODO:: should check or not ???
			if t.Type == VariableTypeNull {
				errs = append(errs, fmt.Errorf("%s pointer types default value is '%s' already",
					errMsgPrefix(v.Pos), t.TypeString()))
				continue
			}
			if v.IsStatic() && v.Expression.IsLiteral() {
				v.DefaultValue = v.Expression.Data
				continue
			}
			if v.IsStatic() == false {
				// nothing to do
				continue
			}
			bin := &ExpressionBinary{}
			bin.Right = &Expression{
				Type:        ExpressionTypeList,
				Description: "list",
				Data:        []*Expression{v.Expression},
			}
			{
				selection := &ExpressionSelection{}
				selection.Expression = &Expression{}
				selection.Expression.Description = "selection"
				selection.Expression.Value = &Type{
					Type:  VariableTypeClass,
					Class: c,
				}
				selection.Name = v.Name
				selection.Field = v
				left := &Expression{
					Type:        ExpressionTypeSelection,
					Data:        selection,
					Description: "selection",
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
				Description:           "assign",
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
		if c.StaticBlocks != nil {
			c.StaticBlocks = append([]*Block{b}, c.StaticBlocks...)
		} else {
			c.StaticBlocks = []*Block{b}
		}
	}
	return errs
}

func (c *Class) checkMethods() []error {
	errs := []error{}
	if c.IsInterface() {
		return errs
	}
	for name, methods := range c.Methods {
		for _, method := range methods {
			errs = append(errs, method.checkModifierOk()...)
			if method.IsAbstract() {
				//nothing
			} else {
				if c.IsInterface() {
					errs = append(errs, fmt.Errorf("%s interface method cannot have implementation",
						errMsgPrefix(method.Function.Pos)))
					continue
				}
				isConstruction := name == SpecialMethodInit
				if isConstruction {
					if method.IsFirstStatementCallFatherConstruction() == false {
						errs = append(errs, fmt.Errorf("%s construction method should call father construction method first",
							errMsgPrefix(method.Function.Pos)))
					}
					if method.IsFinal() {
						errs = append(errs, fmt.Errorf("%s construction method cannot be final",
							errMsgPrefix(method.Function.Pos)))
					}
				}
				if method.IsStatic() == false { // bind this
					if method.Function.Block.Variables == nil {
						method.Function.Block.Variables = make(map[string]*Variable)
					}
					method.Function.Block.Variables[THIS] = &Variable{}
					method.Function.Block.Variables[THIS].Name = THIS
					method.Function.Block.Variables[THIS].Pos = method.Function.Pos
					method.Function.Block.Variables[THIS].Type = &Type{
						Type:  VariableTypeObject,
						Class: c,
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

func (c *Class) checkModifierOk() []error {
	errs := []error{}
	if c.IsInterface() && c.IsFinal() {
		errs = append(errs, fmt.Errorf("%s interface '%s' cannot be final",
			errMsgPrefix(c.FinalPos), c.Name))
	}
	if c.IsAbstract() && c.IsFinal() {
		errs = append(errs, fmt.Errorf("%s abstract class '%s' cannot be final",
			errMsgPrefix(c.FinalPos), c.Name))
	}
	return errs
}
