package cg

type AttributeInfo struct {
	NameIndex       uint16
	attributeLength uint32
	Info            []byte
}

type AttributeGroupedByName map[string][]*AttributeInfo

func (a AttributeGroupedByName) GetByName(name string) []*AttributeInfo {
	if a == nil {
		return nil
	}
	return a[name]
}
