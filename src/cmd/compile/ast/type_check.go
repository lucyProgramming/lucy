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
	for _, v := range p.Funcs {
		errs = append(errs, v.check(nil)...)
	}
	return errs
}

func (p *Package) checkBlocks() []error {
	errs := []error{}
	for _, v := range p.Blocks {
		errs = append(errs, v.check(p)...)
	}
	return errs
}

func (p *Package) checkClass() []error {
	errs := []error{}
	for _, v := range p.Classes {
		errs = append(errs, v.check()...)
	}
	return errs
}

func (p *Package) checkConst() []error {
	errs := make([]error, 0)
	var err error
	for _, v := range p.Consts {
		if v.Expression == nil && v.Typ == nil {
			errs = append(errs, fmt.Errorf("%s %d:%d %s is has no type and no init value", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, v.Name))
		}
		if v.Expression != nil {
			is, t, value, err := v.Expression.getConstValue()
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if is == false {
				errs = append(errs, fmt.Errorf("%s %d:%d %s is not a const value", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, v.Name))
				continue
			}
			//rewrite
			v.Expression = &Expression{}
			v.Expression.Typ = t
			v.Expression.Data = value
		}
		if v.Typ != nil && v.Expression != nil {
			v.Data, err = v.Typ.assignExpression(p, v.Expression)
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

	var block Block

	for _, v := range p.Vars {
		if v.Expression == nil && v.Typ == nil {
			continue
		}
		if v.Expression != nil {
			err = v.Expression.constFold() //fold const error
			if err != nil {
				errs = append(errs, fmt.Errorf("%s %d:%d variable %s defined wrong,err:%v", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, err))
				continue
			}
		}
		if v.Typ == nil && v.Expression != nil { //means variable typed by assignment
			var es []error
			v.Typ, es = block.getTypeFromExpression(v.Expression)
			if es != nil {
				errs = append(errs, es...)
			}
			continue
		}
		if v.Typ != nil && v.Expression != nil { // if typ match
			t2, es := block.getTypeFromExpression(v.Expression)
			if len(es) > 0 {
				errs = append(errs, es...)
				continue
			}
			match := v.Typ.typeCompatible(t2)
			if !match {
				errs = append(errs, fmt.Errorf("%s %d:%d variable %s dose not matched by %s ", v.Pos.Filename, v.Pos.StartLine, v.Pos.StartColumn, v.Name, v.Expression.typeName()))
				continue
			}
		}
		panic("unhandled situation")
	}
	return errs
}
