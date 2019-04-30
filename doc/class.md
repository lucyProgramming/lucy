### **类(class)**
定义一个类

~~~
	class xxx {}
~~~
<br/>
定义一个继承之其他类的类

~~~
	class xxx {}
    class yyy extends xxx {}
~~~
<br/>

定义一个实现了接口的类

~~~
	interface inter {} 
	class xxx {}
    class yyy extends xxx  implements inter , ... {
    	
    }
~~~
<br/>

定义构造方法

~~~
	class xxx {
    	public fn xxx(){
        	
        }
    }
~~~
<br/>


调用父类的构造方法

~~~
	class xxx {
    	public fn xxx(){
        	this.super()
        }
    }
~~~
<br/>

定义字段 

~~~
	class xxx {
    	a,b int
    }
~~~

<br/>

定义包含默认值的字段

~~~
	class xxx {
    	a,b int = 1,2
    }
~~~

<br/>

定义方法
~~~
	class xxx {
    	public fn sayHai(){  // public 
        	 
        }
        private fn sayHa2i(){  // private 
        	 
        }
        protected fn sayHa2i(){  // protected 
        	 
        }
    }
~~~
<br/>

定义静态方法

~~~
	class xxx {
    	public static fn sayHai(){  // public 
        	 
        }
    }
~~~
<br/>


定义带返回值的方法

~~~
	class xxx {
    	public static fn sayHai(hai string = "hello world")->(a int = 1){  // public 
        	print(hai)
        }
    }
~~~






