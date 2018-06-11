package ast

func (e *Expression) checkConst(block *Block) []error {
	errs := []error{}
	cs := e.Data.([]*Const)
	for _, c := range cs {
		err := checkConst(block, c, &errs)
		if err != nil {
			continue
		}
		err = block.insert(c.Name, c.Pos, c)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
