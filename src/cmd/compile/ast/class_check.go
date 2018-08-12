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
	if esNotEmpty(es) {
		errs = append(errs, es...)
	}
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
	if esNotEmpty(es) {
		errs = append(errs, es...)
	}

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
			errMsg := fmt.Sprintf("%s class method named '%s' already declared as field,at:\n",
				errMsgPrefix(ms[0].Function.Pos), ms[0].Function.Name)

			errMsg += fmt.Sprintf("\t%s", errMsgPrefix(f.Pos))
			errs = append(errs, errors.New(errMsg))
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
	errs = append(errs, c.suitableForInterfaces()...)
	return errs
}

func (c *Class) checkIfClassHierarchyErr() error {
	m := make(map[string]struct{})
	arr := []string{}
	is := false
	class := c
	if err := c.loadSuperClass(); err != nil {
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
		err := class.loadSuperClass()
		if err != nil {
			return err
		}
		if class.SuperClass == nil {
			panic("class is nil")
		}
		class = class.SuperClass
	}
	if is == false {
		return nil
	}
	errMsg := fmt.Sprintf("%s class named '%s' detects a circularity in class hierarchy\n",
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
	err := c.loadSuperClass()
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
	// c is class
	for _, i := range c.Interfaces {
		errs = append(errs, c.suitableForInterface(i, false)...)
	}
	return errs
}

func (c *Class) suitableForInterface(inter *Class, fromSub bool) []error {
	errs := []error{}
	for name, v := range inter.Methods {
		m := v[0]
		im, match := c.implementMethod(m, false, &errs, c.Pos)
		if match {
			if im.IsPublic() == false {
				err := fmt.Errorf("%s method '%s' is not public",
					errMsgPrefix(c.Pos), name)
				errs = append(errs, err)
			}
			if im.IsStatic() {
				err := fmt.Errorf("%s method '%s' is static",
					errMsgPrefix(c.Pos), name)
				errs = append(errs, err)
			}
			return errs
		} else {
			errs = append(errs, fmt.Errorf("%s missing implements method '%s' define on interface '%s'",
				errMsgPrefix(c.Pos), m.Function.readableMsg(), inter.Name))
		}
	}
	for _, vv := range inter.Interfaces {
		err := vv.loadSelf()
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		es := c.suitableForInterface(vv, true)
		if esNotEmpty(es) {
			errs = append(errs, es...)
		}
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
			if esNotEmpty(es) {
				errs = append(errs, es...)
			}
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
				Type: ExpressionTypeList,
				Data: []*Expression{v.Expression},
			}
			{
				selection := &ExpressionSelection{}
				selection.Expression = &Expression{}
				selection.Expression.Value = &Type{
					Type:  VariableTypeClass,
					Class: c,
				}
				selection.Name = v.Name
				selection.Field = v
				left := &Expression{
					Type: ExpressionTypeSelection,
					Data: selection,
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
			if name == SpecialMethodInit {
				if c.IsAbstract() {
					errs = append(errs, fmt.Errorf("%s abstract class cannot have construction method",
						errMsgPrefix(method.Function.Pos)))
					continue
				}
			}
			errs = append(errs, method.checkModifierOk()...)
			if method.IsAbstract() {
				//nothing
			} else {
				if c.IsInterface() {
					errs = append(errs, fmt.Errorf("%s interface method cannot have implementation",
						errMsgPrefix(method.Function.Pos)))
					continue
				}
				isConstruction := (name == SpecialMethodInit)
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
				if isConstruction && method.Function.NoReturnValue() == false {
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
