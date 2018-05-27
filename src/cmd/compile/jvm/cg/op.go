package cg

const (
	OP_nop byte = 0x00 //什么都不做。
	//
	OP_aconst_null byte = 0x01 // 将 null 推送至栈顶。
	//
	OP_iconst_m1 byte = 0x02 //将 int 型-1 推送至栈顶。
	OP_iconst_0  byte = 0x03 // 将 int 型 0 推送至栈顶。
	OP_iconst_1  byte = 0x04 // 将 int 型 1 推送至栈顶。
	OP_iconst_2  byte = 0x05 //将 int 型 2 推送至栈顶。
	OP_iconst_3  byte = 0x06 // 将 int 型 3 推送至栈顶。
	OP_iconst_4  byte = 0x07 //将 int 型 4 推送至栈顶。
	OP_iconst_5  byte = 0x08 //将 int 型 5 推送至栈顶。
	//
	OP_lconst_0 byte = 0x09 //将 long 型 0 推送至栈顶。
	OP_lconst_1 byte = 0x0a // 将 long 型 1 推送至栈顶。
	//
	OP_fconst_0 byte = 0x0b // 将 float 型 0 推送至栈顶。
	OP_fconst_1 byte = 0x0c // 将 float 型 1 推送至栈顶。
	OP_fconst_2 byte = 0x0d //将 float 型 2 推送至栈顶。
	//
	OP_dconst_0 byte = 0x0e //将 double 型 0 推送至栈顶。
	OP_dconst_1 byte = 0x0f // 将 double 型 1 推送至栈顶。
	//
	OP_bipush byte = 0x10 //将单字节的常量值（-128~127）推送至栈顶。
	OP_sipush byte = 0x11 //将一个短整型常量值（-32768~32767）推送至栈顶。
	//
	OP_ldc    byte = 0x12 //将 int，float 或 String 型常量值从常量池中推送至栈顶。
	OP_ldc_w  byte = 0x13 //将 int，float 或 String 型常量值从常量池中推送至栈顶（宽 索引）。
	OP_ldc2_w byte = 0x14 //将 long 或 double 型常量值从常量池中推送至栈顶（宽索引）。
	//
	OP_iload byte = 0x15 //将指定的 int 型局部变量推送至栈顶。
	OP_lload byte = 0x16 //将指定的 long 型局部变量推送至栈顶。
	OP_fload byte = 0x17 //将指定的 float 型局部变量推送至栈顶。
	OP_dload byte = 0x18 //将指定的 double 型局部变量推送至栈顶。
	OP_aload byte = 0x19 //将指定的引用类型局部变量推送至栈顶。
	//
	OP_iload_0 byte = 0x1a // 将第一个 int 型局部变量推送至栈顶。
	OP_iload_1 byte = 0x1b //将第二个 int 型局部变量推送至栈顶。
	OP_iload_2 byte = 0x1c //将第三个 int 型局部变量推送至栈顶。
	OP_iload_3 byte = 0x1d //将第四个 int 型局部变量推送至栈顶。
	//
	OP_lload_0 byte = 0x1e //将第一个 long 型局部变量推送至栈顶。
	OP_lload_1 byte = 0x1f // 将第二个 long 型局部变量推送至栈顶。
	OP_lload_2 byte = 0x20 //将第三个 long 型局部变量推送至栈顶。
	OP_lload_3 byte = 0x21 // 将第四个 long 型局部变量推送至栈顶。
	//
	OP_fload_0 byte = 0x22 // 将第一个 float 型局部变量推送至栈顶。
	OP_fload_1 byte = 0x23 //将第二个 float 型局部变量推送至栈顶。
	OP_fload_2 byte = 0x24 // 将第三个 float 型局部变量推送至栈顶
	OP_fload_3 byte = 0x25 //将第四个 float 型局部变量推送至栈顶。
	//
	OP_dload_0 byte = 0x26 //将第一个 double 型局部变量推送至栈顶。
	OP_dload_1 byte = 0x27 //将第二个 double 型局部变量推送至栈顶。
	OP_dload_2 byte = 0x28 //将第三个 double 型局部变量推送至栈顶。
	OP_dload_3 byte = 0x29 //将第四个 double 型局部变量推送至栈顶。
	//
	OP_aload_0 byte = 0x2a //将第一个引用类型局部变量推送至栈顶。
	OP_aload_1 byte = 0x2b //将第二个引用类型局部变量推送至栈顶。
	OP_aload_2 byte = 0x2c //将第三个引用类型局部变量推送至栈顶。
	OP_aload_3 byte = 0x2d //将第四个引用类型局部变量推送至栈顶。
	//
	OP_iaload byte = 0x2e //将 int 型数组指定索引的值推送至栈顶。
	OP_laload byte = 0x2f //将 long 型数组指定索引的值推送至栈顶。
	OP_faload byte = 0x30 //将 float 型数组指定索引的值推送至栈顶。
	OP_daload byte = 0x31 // 将 double 型数组指定索引的值推送至栈顶。
	OP_aaload byte = 0x32 //将引用型数组指定索引的值推送至栈顶。
	OP_baload byte = 0x33 //将 boolean 或 byte 型数组指定索引的值推送至栈顶。
	OP_caload byte = 0x34 // 将 char 型数组指定索引的值推送至栈顶。
	OP_saload byte = 0x35 //将 short 型数组指定索引的值推送至栈顶。
	//
	OP_istore byte = 0x36 // 将栈顶 int 型数值存入指定局部变量。
	OP_lstore byte = 0x37 //将栈顶 long 型数值存入指定局部变量。
	OP_fstore byte = 0x38 //将栈顶 float 型数值存入指定局部变量。
	OP_dstore byte = 0x39 //将栈顶 double 型数值存入指定局部变量。
	OP_astore byte = 0x3a // 将栈顶引用型数值存入指定局部变量。
	//
	OP_istore_0 byte = 0x3b //将栈顶 int 型数值存入第一个局部变量。
	OP_istore_1 byte = 0x3c // 将栈顶 int 型数值存入第二个局部变量。
	OP_istore_2 byte = 0x3d //将栈顶 int 型数值存入第三个局部变量。
	OP_istore_3 byte = 0x3e // 将栈顶 int 型数值存入第四个局部变量。
	//
	OP_lstore_0 byte = 0x3f //将栈顶 long 型数值存入第一个局部变量。
	OP_lstore_1 byte = 0x40 // 将栈顶 long 型数值存入第二个局部变量。
	OP_lstore_2 byte = 0x41 //将栈顶 long 型数值存入第三个局部变量。
	OP_lstore_3 byte = 0x42 // 将栈顶 long 型数值存入第四个局部变量。
	//
	OP_fstore_0 byte = 0x43 //将栈顶 float 型数值存入第一个局部变量。
	OP_fstore_1 byte = 0x44 //将栈顶 float 型数值存入第二个局部变量。
	OP_fstore_2 byte = 0x45 //将栈顶 float 型数值存入第三个局部变量。
	OP_fstore_3 byte = 0x46 //将栈顶 float 型数值存入第四个局部变量。
	//
	OP_dstore_0 byte = 0x47 //将栈顶 double 型数值存入第一个局部变量。
	OP_dstore_1 byte = 0x48 //将栈顶 double 型数值存入第二个局部变量。
	OP_dstore_2 byte = 0x49 // 将栈顶 double 型数值存入第三个局部变量。
	OP_dstore_3 byte = 0x4a //将栈顶 double 型数值存入第四个局部变量。
	//
	OP_astore_0 byte = 0x4b //将栈顶引用型数值存入第一个局部变量。
	OP_astore_1 byte = 0x4c ///将栈顶引用型数值存入第二个局部变量。
	OP_astore_2 byte = 0x4d //将栈顶引用型数值存入第三个局部变量
	OP_astore_3 byte = 0x4e //将栈顶引用型数值存入第四个局部变量。
	//
	OP_iastore byte = 0x4f // 将栈顶 int 型数值存入指定数组的指定索引位置
	OP_lastore byte = 0x50 //将栈顶 long 型数值存入指定数组的指定索引位置。
	OP_fastore byte = 0x51 //将栈顶 float 型数值存入指定数组的指定索引位置。
	OP_dastore byte = 0x52 //将栈顶 double 型数值存入指定数组的指定索引位置。
	OP_aastore byte = 0x53 //将栈顶引用型数值存入指定数组的指定索引位置。
	OP_bastore byte = 0x54 //将栈顶 boolean 或 byte 型数值存入指定数组的指定索引位置。
	OP_castore byte = 0x55 //将栈顶 char 型数值存入指定数组的指定索引位置
	OP_sastore byte = 0x56 // 将栈顶 short 型数值存入指定数组的指定索引位置。
	//
	OP_pop             byte = 0x57 //将栈顶数值弹出（数值不能是 long 或 double 类型的）。
	OP_pop2            byte = 0x58 //将栈顶的一个（long 或 double 类型的）或两个数值弹出（其 它）。
	OP_dup             byte = 0x59 //复制栈顶数值并将复制值压入栈顶。
	OP_dup_x1          byte = 0x5a //复制栈顶数值并将两个复制值压入栈顶。
	OP_dup_x2          byte = 0x5b //复制栈顶数值并将三个（或两个）复制值压入栈顶。
	OP_dup2            byte = 0x5c //复制栈顶一个（long 或 double 类型的)或两个（其它）数值并 将复制值压入栈顶。
	OP_dup2_x1         byte = 0x5d //dup_x1 指令的双倍版本。
	OP_dup2_x2         byte = 0x5e // dup_x2 指令的双倍版本。
	OP_swap            byte = 0x5f //将栈最顶端的两个数值互换（数值不能是 long 或 double 类型 的）。
	OP_iadd            byte = 0x60 //将栈顶两 int 型数值相加并将结果压入栈顶。
	OP_ladd            byte = 0x61 //将栈顶两 long 型数值相加并将结果压入栈顶。
	OP_fadd            byte = 0x62 //将栈顶两 float 型数值相加并将结果压入栈顶。
	OP_dadd            byte = 0x63 //将栈顶两 double 型数值相加并将结果压入栈顶。
	OP_isub            byte = 0x64 //将栈顶两 int 型数值相减并将结果压入栈顶。
	OP_lsub            byte = 0x65 // 将栈顶两 long 型数值相减并将结果压入栈顶。
	OP_fsub            byte = 0x66 //将栈顶两 float 型数值相减并将结果压入栈顶。
	OP_dsub            byte = 0x67 //将栈顶两 double 型数值相减并将结果压入栈顶。
	OP_imul            byte = 0x68 //将栈顶两 int 型数值相乘并将结果压入栈顶。。
	OP_lmul            byte = 0x69 //将栈顶两 long 型数值相乘并将结果压入栈顶。
	OP_fmul            byte = 0x6a //将栈顶两 float 型数值相乘并将结果压入栈顶。
	OP_dmul            byte = 0x6b //将栈顶两 double 型数值相乘并将结果压入栈顶。
	OP_idiv            byte = 0x6c //将栈顶两 int 型数值相除并将结果压入栈顶。
	OP_ldiv            byte = 0x6d //将栈顶两 long 型数值相除并将结果压入栈顶。
	OP_fdiv            byte = 0x6e //将栈顶两 float 型数值相除并将结果压入栈顶。
	OP_ddiv            byte = 0x6f //将栈顶两 double 型数值相除并将结果压入栈顶。
	OP_irem            byte = 0x70 //将栈顶两 int 型数值作取模运算并将结果压入栈顶。
	OP_lrem            byte = 0x71 //将栈顶两 long 型数值作取模运算并将结果压入栈顶。
	OP_frem            byte = 0x72 //将栈顶两 float 型数值作取模运算并将结果压入栈顶。
	OP_drem            byte = 0x73 //将栈顶两 double 型数值作取模运算并将结果压入栈顶。
	OP_ineg            byte = 0x74 //将栈顶 int 型数值取负并将结果压入栈顶。
	OP_lneg            byte = 0x75 //将栈顶 long 型数值取负并将结果压入栈顶。
	OP_fneg            byte = 0x76 //将栈顶 float 型数值取负并将结果压入栈顶。
	OP_dneg            byte = 0x77 //将栈顶 double 型数值取负并将结果压入栈顶。
	OP_ishl            byte = 0x78 //将 int 型数值左移位指定位数并将结果压入栈顶。
	OP_lshl            byte = 0x79 //将 long 型数值左移位指定位数并将结果压入栈顶。
	OP_ishr            byte = 0x7a //将 int 型数值右（有符号）移位指定位数并将结果压入栈顶。
	OP_lshr            byte = 0x7b //将 long 型数值右（有符号）移位指定位数并将结果压入栈顶。
	OP_iushr           byte = 0x7c //将 int 型数值右（无符号）移位指定位数并将结果压入栈顶。
	OP_lushr           byte = 0x7d //将 long 型数值右（无符号）移位指定位数并将结果压入栈顶。
	OP_iand            byte = 0x7e //将栈顶两 int 型数值作“按位与”并将结果压入栈顶。
	OP_land            byte = 0x7f //将栈顶两 long 型数值作“按位与”并将结果压入栈顶。
	OP_ior             byte = 0x80 //将栈顶两 int 型数值作“按位或”并将结果压入栈顶。
	OP_lor             byte = 0x81 //将栈顶两 long 型数值作“按位或”并将结果压入栈顶。
	OP_ixor            byte = 0x82 //将栈顶两 int 型数值作“按位异或”并将结果压入栈顶。
	OP_lxor            byte = 0x83 // 将栈顶两 long 型数值作“按位异或”并将结果压入栈顶。
	OP_iinc            byte = 0x84 //将指定 int 型变量增加指定值。
	OP_i2l             byte = 0x85 //将栈顶 int 型数值强制转换成 long 型数值并将结果压入栈顶。
	OP_i2f             byte = 0x86 //将栈顶 int 型数值强制转换成 float 型数值并将结果压入栈顶。
	OP_i2d             byte = 0x87 //将栈顶 int 型数值强制转换成 double 型数值并将结果压入栈顶。
	OP_l2i             byte = 0x88 // 将栈顶 long 型数值强制转换成 int 型数值并将结果压入栈顶。
	OP_l2f             byte = 0x89 //将栈顶 long 型数值强制转换成 float 型数值并将结果压入栈顶。
	OP_l2d             byte = 0x8a //将栈顶 long 型数值强制转换成 double 型数值并将结果压入栈顶。
	OP_f2i             byte = 0x8b //将栈顶 float 型数值强制转换成 int 型数值并将结果压入栈顶。
	OP_f2l             byte = 0x8c //将栈顶 float 型数值强制转换成 long 型数值并将结果压入栈顶。
	OP_f2d             byte = 0x8d //将栈顶float型数值强制转换成double型数值并将结果压入栈顶。
	OP_d2i             byte = 0x8e //将栈顶 double 型数值强制转换成 int 型数值并将结果压入栈顶。
	OP_d2l             byte = 0x8f //将栈顶 double 型数值强制转换成 long 型数值并将结果压入栈顶。
	OP_d2f             byte = 0x90 //将栈顶double型数值强制转换成float型数值并将结果压入栈 顶。
	OP_i2b             byte = 0x91 //将栈顶 int 型数值强制转换成 byte 型数值并将结果压入栈顶。
	OP_i2c             byte = 0x92 //将栈顶 int 型数值强制转换成 char 型数值并将结果压入栈顶。
	OP_i2s             byte = 0x93 //将栈顶 int 型数值强制转换成 short 型数值并将结果压入栈顶。
	OP_lcmp            byte = 0x94 //比较栈顶两 long 型数值大小，并将结果（1，0，-1）压入栈顶
	OP_fcmpl           byte = 0x95 //比较栈顶两 float 型数值大小，并将结果（1，0，-1）压入栈 顶；当其中一个数值为“NaN”时，将-1压入栈顶。
	OP_fcmpg           byte = 0x96 //比较栈顶两 float 型数值大小，并将结果（1，0，-1）压入栈顶；当其中一个数值为“NaN”时，将 1压入栈顶。
	OP_dcmpl           byte = 0x97 //比较栈顶两 double 型数值大小，并将结果（1，0，-1）压入栈顶；当其中一个数值为“NaN”时，将-1压入栈顶。
	OP_dcmpg           byte = 0x98 //比较栈顶两 double 型数值大小，并将结果（1，0，-1）压入栈顶；当其中一个数值为“NaN”时，将 1压入栈顶。
	OP_ifeq            byte = 0x99 //当栈顶 int 型数值等于 0 时跳转。
	OP_ifne            byte = 0x9a //当栈顶 int 型数值不等于 0 时跳转。
	OP_iflt            byte = 0x9b //当栈顶 int 型数值小于 0 时跳转。
	OP_ifge            byte = 0x9c //当栈顶 int 型数值大于等于 0 时跳转。
	OP_ifgt            byte = 0x9d //当栈顶 int 型数值大于 0 时跳转。
	OP_ifle            byte = 0x9e //当栈顶 int 型数值小于等于 0 时跳转。
	OP_if_icmpeq       byte = 0x9f //比较栈顶两 int 型数值大小，当结果等于 0 时跳转。
	OP_if_icmpne       byte = 0xa0 //比较栈顶两 int 型数值大小，当结果不等于 0 时跳转。
	OP_if_icmplt       byte = 0xa1 //比较栈顶两 int 型数值大小，当结果小于 0 时跳转。
	OP_if_icmpge       byte = 0xa2 //比较栈顶两 int 型数值大小，当结果大于等于 0 时跳转。
	OP_if_icmpgt       byte = 0xa3 //比较栈顶两 int 型数值大小，当结果大于 0 时跳转
	OP_if_icmple       byte = 0xa4 //比较栈顶两 int 型数值大小，当结果小于等于 0 时跳转。
	OP_if_acmpeq       byte = 0xa5 //比较栈顶两引用型数值，当结果相等时跳转。
	OP_if_acmpne       byte = 0xa6 //比较栈顶两引用型数值，当结果不相等时跳转。
	OP_goto            byte = 0xa7 //无条件跳转。
	OP_jsr             byte = 0xa8 //跳转至指定 16 位 offset 位置，并将 jsr 下一条指令地址压入栈顶。
	OP_ret             byte = 0xa9 //返回至局部变量指定的 index 的指令位置（一般与 jsr，jsr_w联合使用）。
	OP_tableswitch     byte = 0xaa //用于 switch 条件跳转，case 值连续（可变长度指令）。
	OP_lookupswitch    byte = 0xab //用于 switch 条件跳转，case 值不连续（可变长度指令）。
	OP_ireturn         byte = 0xac //从当前方法返回 int。
	OP_lreturn         byte = 0xad //从当前方法返回 long。
	OP_freturn         byte = 0xae //从当前方法返回 float。
	OP_dreturn         byte = 0xaf //从当前方法返回 double。
	OP_areturn         byte = 0xb0 //从当前方法返回对象引用。
	OP_return          byte = 0xb1 //从当前方法返回 void。
	OP_getstatic       byte = 0xb2 //获取指定类的静态域，并将其值压入栈顶。
	OP_putstatic       byte = 0xb3 //为指定的类的静态域赋值。
	OP_getfield        byte = 0xb4 //获取指定类的实例域，并将其值压入栈顶。
	OP_putfield        byte = 0xb5 //为指定的类的实例域赋值。
	OP_invokevirtual   byte = 0xb6 //调用实例方法。
	OP_invokespecial   byte = 0xb7 //调用超类构造方法，实例初始化方法，私有方法。
	OP_invokestatic    byte = 0xb8 //调用静态方法。
	OP_invokeinterface byte = 0xb9 //调用接口方法。
	OP_invokedynamic   byte = 0xba //调用动态链接方法①。
	OP_new             byte = 0xbb //创建一个对象，并将其引用值压入栈顶。
	OP_newarray        byte = 0xbc //创建一个指定原始类型（如 int、float、char„„）的数组，并将其引用值压入栈顶。
	OP_anewarray       byte = 0xbd //创建一个引用型（如类，接口，数组）的数组，并将其引用值压入栈顶。
	OP_arraylength     byte = 0xbe //获得数组的长度值并压入栈顶。
	OP_athrow          byte = 0xbf //将栈顶的异常抛出。
	OP_checkcast       byte = 0xc0 //检验类型转换，检验未通过将抛出 ClassCastException。
	OP_instanceof      byte = 0xc1 //检验对象是否是指定的类的实例，如果是将 1 压入栈顶，否则将0 压入栈顶。
	OP_monitorenter    byte = 0xc2 //获得对象的 monitor，用于同步方法或同步块。
	OP_monitorexit     byte = 0xc3 //释放对象的 monitor，用于同步方法或同步块。
	OP_wide            byte = 0xc4 //扩展访问局部变量表的索引宽度。
	OP_multianewarray  byte = 0xc5 //创建指定类型和指定维度的多维数组（执行该指令时，操作栈中 必须包含各维度的长度值），并将其引用值压入栈顶。
	OP_ifnull          byte = 0xc6 //为 null 时跳转。
	OP_ifnonnull       byte = 0xc7 //不为 null 时跳转。
	OP_goto_w          byte = 0xc8 //无条件跳转（宽索引）。
	OP_jsr_w           byte = 0xc9 // 跳转至指定 32 位地址偏移量位置，并将 jsr_w 下一条指令地址压入栈顶。
	//保留指令
	breakpoint byte = 0xca //调试时的断点标志。
	//    byte = 0xfe //用于在特定硬件中使用的语言后门。
	//impdep1    byte = 0xff //用于在特定硬件中使用的语言后门。
)
