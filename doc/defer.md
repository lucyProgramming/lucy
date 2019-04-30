### **defer语句**
defer是block-level的延迟函数，defer语句会在这个block执行结束后运行
~~~
fn main (args []string) {
	for i := 0 ;i < 2 ; i++ {
		defer {
			print("defer:" , i )
		}
		print(i)
	}
	print("endFor")
}
运行结果：
0
defer: 0
1
defer: 1
endFor

~~~
同时defer也用于异常处理

~~~
fn main (args []string) {
	for i := 0 ;i < 2 ; i++ {
		defer {
			print("defer:" , i )
            x := catch()
			if x != null {
				print(x)
			}
		}
		print(i)
		[1][1] = 123
	}
	print("endFor")
}
运行结果：
0
defer: 0
java.lang.ArrayIndexOutOfBoundsException: index out range
1
defer: 1
java.lang.ArrayIndexOutOfBoundsException: index out range
endFor

~~~



