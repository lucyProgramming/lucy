package ast

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type TemplateFunction struct {
	Instances []*TemplateFunctionInstance
}

type TemplateFunctionInstance struct {
	parameterTypes map[string]*Type
	Generated      *cg.MethodHighLevel
	Function       *Function
	ClassName      string
}

func (t *TemplateFunction) instanceExists(parameterTypes map[string]*Type) *TemplateFunctionInstance {
	equal := func(p *TemplateFunctionInstance) bool {
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
	for _, v := range t.Instances {
		if equal(v) {
			return v
		}
	}
	return nil
}

func (t *TemplateFunction) insert(parameterTypes map[string]*Type) *TemplateFunctionInstance {
	if t := t.instanceExists(parameterTypes); t != nil {
		return t
	}
	ret := &TemplateFunctionInstance{
		parameterTypes: parameterTypes,
	}
	t.Instances = append(t.Instances, ret)
	return ret
}
