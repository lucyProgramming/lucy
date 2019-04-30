### **return语句**
return 语句会终止函数的执行并返回函数的值。
~~~
public fn encode()->(bs []byte ){
		if this.x == null {
			return []byte("null")
		}
		if this.c.isPrimitive(){
			return []byte(this.x.toString())
		} 
		if isMap(this.c) {
			return this.encodeMap()
		}
		if this.c.isArray() {
			return this.encodeArray()
		}
		return this.encodeObject()
	}

~~~

