package ast

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type TemplateFunction struct {
	Pairs []*TemplateFunctionCallPair
}

type TemplateFunctionCallPair struct {
	typeParameters map[string]*VariableType
	Generated      *cg.MethodHighLevel
	Function       *Function
	ClassName      string
}

func (t *TemplateFunction) callPairExists(typedParameters map[string]*VariableType, errs *[]error) *TemplateFunctionCallPair {
	f := func(p *TemplateFunctionCallPair) *TemplateFunctionCallPair {
		if len(p.typeParameters) != len(typedParameters) {
			return nil
		}
		for kk, vv := range typedParameters {
			t, ok := p.typeParameters[kk]
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

func (t *TemplateFunction) insert(typeParameters map[string]*VariableType, f *Function, errs *[]error) *TemplateFunctionCallPair {
	if t := t.callPairExists(typeParameters, errs); t != nil {
		return t
	}
	ret := &TemplateFunctionCallPair{
		typeParameters: typeParameters,
		Function:       f,
	}
	t.Pairs = append(t.Pairs, ret)
	return ret
}
