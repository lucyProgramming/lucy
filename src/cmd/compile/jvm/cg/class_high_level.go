package cg

type ClassHighLevel struct {
	IsClosureFunctionClass bool
	MainClass              *ClassHighLevel
	InnerClasss            []*ClassHighLevel
	AccessFlags            uint16
	IntConsts              map[int32][][]byte
	LongConsts             map[int64][][]byte
	StringConsts           map[string][][]byte
	FloatConsts            map[float32][][]byte
	DoubleConsts           map[float64][][]byte
	FieldRefs              map[CONSTANT_Fieldref_info_high_level][][]byte
	Name                   string
	SuperClass             string
	Interfaces             []string
	Fields                 map[string]*FiledHighLevel
	Methods                map[string][]*MethodHighLevel
}

type CONSTANT_Fieldref_info_high_level struct {
	Class string
	Field string
}

func (c *ClassHighLevel) InsertFieldRef(fr CONSTANT_Fieldref_info_high_level, location []byte) {
	if c.FieldRefs == nil {
		c.FieldRefs = make(map[CONSTANT_Fieldref_info_high_level][][]byte)
	}
	if x, ok := c.FieldRefs[fr]; ok {
		x = append(x, location)
	} else {
		c.FieldRefs[fr] = [][]byte{location}
	}
}

// func (c *ClassHighLevel) InsertClasses(s string, location []byte) {
// 	if c.Classes == nil {
// 		c.Classes = make(map[int32][][]byte)
// 	}
// 	if x, ok := c.Classes[i]; ok {
// 		x = append(x, location)
// 	} else {
// 		c.Classes[i] = [][]byte{location}
// 	}
// }

func (c *ClassHighLevel) InsertIntConst(i int32, location []byte) {
	if c.IntConsts == nil {
		c.IntConsts = make(map[int32][][]byte)
	}
	if x, ok := c.IntConsts[i]; ok {
		x = append(x, location)
	} else {
		c.IntConsts[i] = [][]byte{location}
	}
}
func (c *ClassHighLevel) InsertLongConst(i int64, location []byte) {
	if c.LongConsts == nil {
		c.LongConsts = make(map[int64][][]byte)
	}
	if x, ok := c.LongConsts[i]; ok {
		x = append(x, location)
	} else {
		c.LongConsts[i] = [][]byte{location}
	}
}
func (c *ClassHighLevel) InsertStringConst(s string, location []byte) {
	if c.StringConsts == nil {
		c.StringConsts = make(map[string][][]byte)
	}
	if x, ok := c.StringConsts[s]; ok {
		x = append(x, location)
	} else {
		c.StringConsts[s] = [][]byte{location}
	}
}

func (c *ClassHighLevel) InsertFloatConst(f float32, location []byte) {
	if c.FloatConsts == nil {
		c.FloatConsts = make(map[float32][][]byte)
	}
	if x, ok := c.FloatConsts[f]; ok {
		x = append(x, location)
	} else {
		c.FloatConsts[f] = [][]byte{location}
	}
}
func (c *ClassHighLevel) InsertDoubleConst(d float64, location []byte) {
	if c.DoubleConsts == nil {
		c.DoubleConsts = make(map[float64][][]byte)
	}
	if x, ok := c.DoubleConsts[d]; ok {
		x = append(x, location)
	} else {
		c.DoubleConsts[d] = [][]byte{location}
	}
}

type FiledHighLevel struct {
	BackPatchs [][]byte
	Name       string
	Descriptor string
	FieldInfo
}
type MethodHighLevel struct {
	BackPatchs     [][]byte
	ClassHighLevel *ClassHighLevel
	Name           string
	Descriptor     string
	MethodInfo
	Code AttributeCode
}
