package cg

const (
	ACC_FIELD_PUBLIC    uint16 = 0x0001 //public，表示字段可以从任何包访问。
	ACC_FIELD_PRIVATE   uint16 = 0x0002 // private，表示字段仅能该类自身调用。
	ACC_FIELD_PROTECTED uint16 = 0x0004 //protected，表示字段可以被子类调用。
	ACC_FIELD_STATIC    uint16 = 0x0008 //static，表示静态字段。
	ACC_FIELD_FINAL     uint16 = 0x0010 //final，表示字段定义后值无法修改（JLS §17.5）
)

type FieldInfo struct {
	accessFlags     uint16
	nameIndex       uint16
	descriptorIndex uint16
	attributeCount  uint16
	attributes      []*AttributeInfo
}
