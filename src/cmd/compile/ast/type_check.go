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
	for _, v := range p.Vars {
		if v.Typ == nil && v.Init != nil { //means variable typed by assignment
			v.Typ = p.getTypeFromExpression(v.Init)
			continue
		}
		if v.Typ != nil && v.Init != nil { // if typ match

		}
		panic("unhandled situation")
	}
	return errs
}

func (p *Package) getTypeFromExpression(e *Expression) *VariableType {
	switch e.Typ {
	case EXPRESSION_TYPE_BOOL:
		return &VariableType{
			Typ: VARIABLE_TYPE_BOOL,
		}
	case EXPRESSION_TYPE_BYTE:
		return &VariableType{
			Typ: VARIABLE_TYPE_BYTE,
		}
	case EXPRESSION_TYPE_INT:
		return &VariableType{
			Typ: VARIABLE_TYPE_INT,
		}
	case EXPRESSION_TYPE_FLOAT:
		return &VariableType{
			Typ: VARIABLE_TYPE_FLOAT,
		}
	case EXPRESSION_TYPE_STRING:
		return &VariableType{
			Typ: VARIABLE_TYPE_STRING,
		}
	default:
		panic("unhandled situation")
	}
}
