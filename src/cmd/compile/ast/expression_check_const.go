package ast

func (e *Expression) checkConstant(block *Block) []error {
	errs := []error{}
	cs := e.Data.([]*Constant)
	for _, c := range cs {
		err := checkConst(block, c, &errs)
		if err != nil {
			continue
		}
		err = block.Insert(c.Name, c.Pos, c)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
