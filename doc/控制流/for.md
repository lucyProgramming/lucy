### **switch语句**

如果您希望一遍又一遍地运行相同的代码，并且每次的值都不同，那么使用循环是很方便的。

* 死循环
~~~
	for {
    	print("hello world");
    }
~~~
<br/>

* 条件循环
~~~
	var i int = 0 ;
	for i < 10{
    	print("hello world");
        i++;
    }
~~~
<br/>

* 循环10次
~~~
	for i := 0 ; i < 10 ; i++{
    	print("hello world");
    }
~~~

<br/>

* range
~~~
	//arr
    arr := [1,2,3,4,5] ;
	for k,_ := range arr {
    	print(k);
    }
    for v := range arr {
    	print(v);
    }
    for k,v := range arr {
    	print(k,v);
    }
    // map 
    m := {1 -> "hello"};
    for k,_ := range m {
    	print(k);
    }
    for v := range m {
    	print(v);
    }
    for k,v := range m {
    	print(k,v);
    }
~~~

<br/><br/>

#### break语句
用于跳出循环
~~~
	for i := 0 ;; i ++{
    	print(i);
        if i == 9 {
        	break;
        }
    }
    
~~~
<br/>

#### continue语句
跳过本次循环

~~~
	for i := 0 ;; i ++{
    	if i % 3 == 0 {
        	continue ;
        }
    	print(i);
        if i == 9 {
        	break;
        }
    }
    
~~~

