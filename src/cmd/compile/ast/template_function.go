package ast

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type TemplateFunction struct {
	Pairs []*TemplateFunctionCallPair
}

type TemplateFunctionCallPair struct {
	typedParameters map[string]*VariableType
	Generated       *cg.MethodHighLevel
	Function        *Function
	ClassName       string
}

func (t *TemplateFunction) callPairExists(typedParameters map[string]*VariableType, errs *[]error) *TemplateFunctionCallPair {
	f := func(p *TemplateFunctionCallPair) *TemplateFunctionCallPair {
		if len(p.typedParameters) != len(typedParameters) {
			return nil
		}
		for kk, vv := range typedParameters {
			t, ok := p.typedParameters[kk]
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

func (t *TemplateFunction) insert(typedParameters map[string]*VariableType, f *Function, errs *[]error) *TemplateFunctionCallPair {
	if t := t.callPairExists(typedParameters, errs); t != nil {
		return t
	}
	ret := &TemplateFunctionCallPair{
		typedParameters: typedParameters,
		Function:        f,
	}
	t.Pairs = append(t.Pairs, ret)
	return ret
}
