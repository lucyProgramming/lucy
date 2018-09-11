package ast

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type TemplateFunction struct {
	Pairs []*TemplateFunctionCallPair
}

type TemplateFunctionCallPair struct {
	parameterTypes map[string]*Type
	Generated      *cg.MethodHighLevel
	Function       *Function
	ClassName      string
}

func (t *TemplateFunction) callPairExists(parameterTypes map[string]*Type) *TemplateFunctionCallPair {
	equal := func(p *TemplateFunctionCallPair) bool {
		if len(p.parameterTypes) != len(parameterTypes) {
			return false
		}
		for tName, tType := range parameterTypes {
			t, ok := p.parameterTypes[tName]
			if ok == false {
				//not found
				return false
			}
			if tType.Equal(t) == false {
				//not equal
				return false
			}
		}
		return true
	}
	for _, v := range t.Pairs {
		if equal(v) {
			return v
		}
	}
	return nil
}

func (t *TemplateFunction) insert(parameterTypes map[string]*Type) *TemplateFunctionCallPair {
	if t := t.callPairExists(parameterTypes); t != nil {
		return t
	}
	ret := &TemplateFunctionCallPair{
		parameterTypes: parameterTypes,
	}
	t.Pairs = append(t.Pairs, ret)
	return ret
}
