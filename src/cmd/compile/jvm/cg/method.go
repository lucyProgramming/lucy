package cg

const (
	ACC_METHOD_PUBLIC       U2 = 0x0001 //public，方法可以从包外访问
	ACC_METHOD_PRIVATE      U2 = 0x0002 //private，方法只能本类中访问
	ACC_METHOD_PROTECTED    U2 = 0x0004 //protected，方法在自身和子类可以访问
	ACC_METHOD_STATIC       U2 = 0x0008 //static，静态方法
	ACC_METHOD_FINAL        U2 = 0x0010 //final，方法不能被重写（覆盖）
	ACC_METHOD_SYNCHRONIZED U2 = 0x0020 //synchronized，方法由管程同步
	ACC_METHOD_BRIDGE       U2 = 0x0040 //bridge，方法由编译器产生
	ACC_METHOD_VARARGS      U2 = 0x0080 //表示方法带有变长参数
	ACC_METHOD_NATIVE       U2 = 0x0100 //native，方法引用非java语言的本地方法
	ACC_METHOD_ABSTRACT     U2 = 0x0400 //abstract，方法没有具体实现
	ACC_METHOD_STRICT       U2 = 0x0800 //strictfp，方法使用FP-strict浮点格式
	ACC_METHOD_SYNTHETIC    U2 = 0x1000 //方法在源文件中不出现，由编译器产生
)

type MethodInfo struct {
	accessFlags     U2
	nameIndex       U2
	descriptorIndex U2
	attributeCount  U2
	attributes      []*AttributeInfo
}
