### **goto语句**
当你想要跳出多重循环时非常有用
~~~
	for i := 0 ;; i ++{
    	print(i)
        if i == 9 {
        	goto end
        }
    }
    
    end:
~~~

<br/>

goto语句不能跳过变量定义，我们来看一个错误的例子

~~~
	for i := 0 ;; i ++{
    	print(i)
        if i == 9 {
        	goto end
        }
    }
    c := false  // cannot jump over variable definition
    
    end:
	
~~~