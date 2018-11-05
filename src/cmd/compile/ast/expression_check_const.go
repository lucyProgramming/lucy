package ast

func (this *Expression) checkConstant(block *Block) []error {
	errs := []error{}
	cs := this.Data.([]*Constant)
	for _, c := range cs {
		err := checkConst(block, c)
		if err != nil {
			if err != nil {
				errs = append(errs, err)
			}
			continue
		}
		err = block.Insert(c.Name, c.Pos, c)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
