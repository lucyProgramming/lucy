package cg

const (
	AccFieldPublic    uint16 = 0x0001 //public，表示字段可以从任何包访问。
	AccFieldPrivate   uint16 = 0x0002 // private，表示字段仅能该类自身调用。
	AccFieldProtected uint16 = 0x0004 //protected，表示字段可以被子类调用。
	AccFieldStatic    uint16 = 0x0008 //static，表示静态字段。
	AccFieldFinal     uint16 = 0x0010 //final，表示字段定义后值无法修改（JLS §17.5）
	AccFieldVolatile  uint16 = 0x0040 //volatile，表示字段是易变的。
	AccFieldTransient uint16 = 0x0080 //transient，表示字段不会被序列化。
	AccFieldSynthetic uint16 = 0x1000 //表示字段由编译器自动产生。
	AccFieldEnum      uint16 = 0x4000 //enum，表示字段为枚举类型。
)

type FieldInfo struct {
	AccessFlags            uint16
	NameIndex              uint16
	DescriptorIndex        uint16
	Attributes             []*AttributeInfo
	AttributeGroupedByName AttributeGroupedByName
}
