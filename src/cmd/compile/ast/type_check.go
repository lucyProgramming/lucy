package ast

func (p *Package) TypeCheck() []error {
	if p.NErros <= 2 {
		p.NErros = 10
	}
	errs := []error{}
	errs = append(errs, p.Block.check(p)...)
	for _, v := range p.Blocks {
		errs = append(errs, v.check(p)...)
	}
	return errs
}
