package ast

import (
	"fmt"
)

func (p *Package) TypeCheck() []error {
	if p.NErros <= 2 {
		p.NErros = 10
	}
	errs := []error{}
	errs = append(errs, p.checkConst()...)
	if len(errs) > p.NErros {
		return errs
	}
	errs = append(errs, p.checkGlobalVariables()...)
	if len(errs) > p.NErros {
		return errs
	}
	errs = append(errs, p.checkFunctions()...)
	if len(errs) > p.NErros {
		return errs
	}
	errs = append(errs, p.checkBlocks()...)
	if len(errs) > p.NErros {
		return errs
	}
	errs = append(errs, p.checkClass()...)
	if len(errs) > p.NErros {
		return errs
	}
	return errs
}

func (p *Package) checkFunctions() []error {
	errs := []error{}
	return errs
}

func (p *Package) checkBlocks() []error {
	errs := []error{}
	for _, v := range p.Blocks {
		errs = append(errs, v.check()...)
	}
	return errs
}

func (p *Package) checkClass() []error {
	errs := []error{}
	return errs
}

func (p *Package) checkConst() []error {
	errs := make([]error, 0)
	var err error
	for _, v := range p.Consts {
		if v.Init == nil && v.Typ == nil {
			errs = append(errs, fmt.Errorf("%s %d:%d %s is has no type and no init value", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, v.Name))
		}
		if v.Init != nil {
			is, t, value, err := v.Init.getConstValue()
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if is == false {
				errs = append(errs, fmt.Errorf("%s %d:%d %s is not a const value", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, v.Name))
				continue
			}
			//rewrite
			v.Init = &Expression{}
			v.Init.Typ = t
			v.Init.Data = value
		}
		if v.Typ != nil && v.Init != nil {
			v.Data, err = v.Typ.assignExpression(p, v.Init)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s %d:%d can`t assign value to %s", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, v.Name))
				continue
			}
		}
	}
	return errs
}

func (p *Package) checkGlobalVariables() []error {
	errs := make([]error, 0)
	var err error
	for _, v := range p.Vars {
		if v.Init == nil && v.Typ == nil {
			continue
		}
		if v.Init != nil {
			err = v.Init.constFold() //fold const error
			if err != nil {
				errs = append(errs, fmt.Errorf("%s %d:%d variable %s defined wrong,err:%v", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, err))
				continue
			}
		}
		if v.Typ == nil && v.Init != nil { //means variable typed by assignment
			v.Typ, err = p.getTypeFromExpression(v.Init)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s %d:%d variable %s can`t assigned by %s ", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, v.Name, v.Init.typeName()))
				continue
			}
			continue
		}
		if v.Typ != nil && v.Init != nil { // if typ match
			match := v.Typ.typeCompatible(p.getTypeFromExpression(v.Init))
			if !match {
				errs = append(errs, fmt.Errorf("%s %d:%d variable %s dose not matched by %s ", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, v.Name, v.Init.typeName()))
				continue
			}
		}
		panic("unhandled situation")
	}
	return errs
}

func (p *Package) getTypeFromExpression(e *Expression) (t *VariableType, err error) {
	switch e.Typ {
	case EXPRESSION_TYPE_BOOL:
		t = &VariableType{
			Typ: VARIABLE_TYPE_BOOL,
		}
	case EXPRESSION_TYPE_BYTE:
		t = &VariableType{
			Typ: VARIABLE_TYPE_BYTE,
		}
	case EXPRESSION_TYPE_INT:
		t = &VariableType{
			Typ: VARIABLE_TYPE_INT,
		}
	case EXPRESSION_TYPE_FLOAT:
		t = &VariableType{
			Typ: VARIABLE_TYPE_FLOAT,
		}
	case EXPRESSION_TYPE_STRING:
		t = &VariableType{
			Typ: VARIABLE_TYPE_STRING,
		}
	default:
		panic("unhandled situation")
	}
	err = fmt.Errorf("can`t assign")
	return
}
