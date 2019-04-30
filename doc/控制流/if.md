### **if语句**
通常在写代码时，您总是需要为不同的决定来执行不同的动作。您可以在代码中使用条件语句来完成该任务。

我们可使用以下条件语句：

if 语句 - 只有当指定条件为 true 时，使用该语句来执行代码
if...else 语句 - 当条件为 true 时执行代码，当条件为 false 时执行其他代码
if...else if....else 语句 - 使用该语句来选择多个代码块之一来执行

#### if 语句
只有当指定条件为 true 时，该语句才会执行代码。
~~~
	if true {
    	print("hello")
    }
~~~
<br/>

#### if...else 语句
使用 if....else 语句在条件为 true 时执行代码，在条件为 false 时执行其他代码。
~~~
        condition := true;
        if  condition  {
			print("conditon is true")
         } else  {
			print("condition is false")
      	}
~~~
<br/>

#### if...else if....else  语句 
使用 if....else if...else 语句来选择多个代码块之一来执行。
~~~
	if  condition1 {
			print("condition1 is true")
	} else if  condition2 {
  			print("condtion2 is true")
	} else {
          	print("conditon1 and condtion2 is false")
	}
~~~
<br/>


#### 使用多个表达式，最后一个表达式是条件
~~~
	fn test() {
		print("hello world")
	}
	//可以很方便的创建临时变量
	if a:= true ; a = false ; test() ; a || true  {
		print(11)
	}

~~~
