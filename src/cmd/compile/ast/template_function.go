package ast

type TemplateFunction struct {
	Pairs []*TemplateFunctionCallPair
}

type TemplateFunctionCallPair struct {
	Args      []*VariableType
	Returns   []*VariableType
	Generated bool
}

func (t *TemplateFunction) callPairExists(Args []*VariableType,
	Returns []*VariableType) *TemplateFunctionCallPair {
	f := func(p *TemplateFunctionCallPair) *TemplateFunctionCallPair {
		if len(p.Args) != len(Args) {
			return nil
		}
		if len(p.Returns) != len(Returns) {
			return nil
		}
		for k, v := range p.Args {
			if false == v.Equal(Args[k]) {
				return nil
			}
		}
		for k, v := range p.Returns {
			if false == v.Equal(Args[k]) {
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

func (t *TemplateFunction) insert(Args []*VariableType,
	Returns []*VariableType) *TemplateFunctionCallPair {
	if t := t.callPairExists(Args, Returns); t != nil {
		return t
	}
	ret := &TemplateFunctionCallPair{
		Args:    Args,
		Returns: Returns,
	}
	t.Pairs = append(t.Pairs, ret)
	return ret
}
