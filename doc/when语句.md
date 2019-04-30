### **when语句**
在模板函数中，根据参数类型的不同来执行不同的代码块。
~~~
fn p <T> (a T) {
	when T{
		case int:
			print("int")
		case string:
			print("string")
	}
}

fn main (args []string) {
	p(1)
	p("")
}

运行结果：
int
string
~~~

除了上述的用法，还可以根据变量类型来执行不同的代码块。
~~~
class xxx {}
class yyy{}
var x = new xxx()
when x.(type){
case xxx:
    print("x is instance of xxx")
case yyy:
    print("x is instance of yyy")
}
~~~


