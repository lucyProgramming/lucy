package ast

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type TemplateFunction struct {
	instances []*TemplateFunctionInstance
}

type TemplateFunctionInstance struct {
	parameterTypes []*Type
	Entrance       *cg.MethodHighLevel
	Function       *Function
}

func (t *TemplateFunction) instanceExists(parameterTypes []*Type) *TemplateFunctionInstance {
	equal := func(instance *TemplateFunctionInstance) bool {
		if len(instance.parameterTypes) != len(parameterTypes) {
			return false
		}
		for k, tType := range parameterTypes {
			if tType.Equal(instance.parameterTypes[k]) == false {
				//not equal
				return false
			}
		}
		return true
	}
	for _, v := range t.instances {
		if equal(v) {
			return v
		}
	}
	return nil
}

func (t *TemplateFunction) insert(parameterTypes []*Type) *TemplateFunctionInstance {
	if t := t.instanceExists(parameterTypes); t != nil {
		return t
	}
	ret := &TemplateFunctionInstance{
		parameterTypes: parameterTypes,
	}
	t.instances = append(t.instances, ret)
	return ret
}
