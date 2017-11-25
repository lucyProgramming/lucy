package cg

const (
	ACC_CLASS_PUBLIC     U2 = 0x0001 // 可以被包的类外访问。
	ACC_CLASS_FINAL      U2 = 0x0010 //不允许有子类。
	ACC_CLASS_SUPER      U2 = 0x0020 //当用到invokespecial指令时，需要特殊处理③的父类方法。
	ACC_CLASS_INTERFACE  U2 = 0x0200 // 标识定义的是接口而不是类。
	ACC_CLASS_ABSTRACT   U2 = 0x0400 //  不能被实例化。
	ACC_CLASS_SYNTHETIC  U2 = 0x1000 //标识并非Java源码生成的代码。
	ACC_CLASS_ANNOTATION U2 = 0x2000 // 标识注解类型
	ACC_CLASS_ENUM       U2 = 0x4000 // 标识枚举类型
	ACC_VOLATILE         U2 = 0x0040 //volatile，表示字段是易变的。
	ACC_TRANSIENT        U2 = 0x0080 //transient，表示字段不会被序列化。
	ACC_SYNTHETIC        U2 = 0x1000 //表示字段由编译器自动产生。
	ACC_ENUM             U2 = 0x4000 //enum，表示字段为枚举类型。
)

type Class struct {
	magic          uint32
	minorVersion   U2
	majorVersion   U2
	constPoolCount U2
	constPool      []*ConstPool
	accessFlag     U2
	thisClass      U2
	superClass     U2
	interfaceCount U2
	interfaces     []U2
	fieldCount     U2
	fields         []*FieldInfo
	methodCount    U2
	methods        []*MethodInfo
	attributeCount U2
	attributes     []*AttributeInfo
}
