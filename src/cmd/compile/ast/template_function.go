package ast

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type TemplateFunction struct {
	Pairs []*TemplateFunctionCallPair
}

type TemplateFunctionCallPair struct {
	parameterTypes map[string]*VariableType
	Generated      *cg.MethodHighLevel
	Function       *Function
	ClassName      string
}

func (t *TemplateFunction) callPairExists(parameterTypes map[string]*VariableType, errs *[]error) *TemplateFunctionCallPair {
	f := func(p *TemplateFunctionCallPair) *TemplateFunctionCallPair {
		if len(p.parameterTypes) != len(parameterTypes) {
			return nil
		}
		for kk, vv := range parameterTypes {
			t, ok := p.parameterTypes[kk]
			if ok == false {
				// not found
				return nil
			}
			if vv.Equal(errs, t) == false {
				// not equal
				return nil
			}
		}
		return p
	}
	for _, v := range t.Pairs {
		if p := f(v); p != nil {
			return p
		}
	}
	return nil
}

func (t *TemplateFunction) insert(parameterTypes map[string]*VariableType, f *Function, errs *[]error) *TemplateFunctionCallPair {
	if t := t.callPairExists(parameterTypes, errs); t != nil {
		return t
	}
	ret := &TemplateFunctionCallPair{
		parameterTypes: parameterTypes,
		Function:       f,
	}
	t.Pairs = append(t.Pairs, ret)
	return ret
}
