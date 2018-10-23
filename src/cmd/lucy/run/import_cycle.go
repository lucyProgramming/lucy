package run

import "fmt"

type ImportStack struct {
	Stacks []*PackageCompiled
	M      map[string]*PackageCompiled
}

func (i *ImportStack) fromLast(last *ImportStack) *ImportStack {
	i.Stacks = make([]*PackageCompiled, len(last.Stacks))
	copy(i.Stacks, last.Stacks)
	i.M = make(map[string]*PackageCompiled)
	for k, v := range last.M {
		i.M[k] = v
	}
	return i
}

/*
	check if import cycling
*/
func (i *ImportStack) insert(c *PackageCompiled) error {
	if _, ok := i.M[c.packageName]; ok {
		errMsg := fmt.Sprintf("package named '%s' import cycling\n", c.packageName)
		errMsg += "\t"
		for _, v := range i.Stacks {
			errMsg += fmt.Sprintf("'%s' -> ", v.packageName)
		}
		errMsg += c.packageName
		return fmt.Errorf(errMsg)
	}
	i.Stacks = append(i.Stacks, c)
	i.M[c.packageName] = c
	return nil
}
