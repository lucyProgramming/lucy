
package lucy.deps;
import java.lang.reflect.* ; 

		import java.lang.Object;

public class ArrayObject   {
	public int start;
	public int end; // not include
	public int cap;
	static String outOfRangeMsg = "index out range";
	public Object[] elements;
	public int size(){
		return this.end - this.start;
	}
	public synchronized int start(){
        return this.start;
	}
	public synchronized int end(){
         return this.end;
	}
	public synchronized int cap(){
         return this.cap;
	}
	public ArrayObject(Object[] values){
		this.start = 0;
		this.end = values.length;
		this.cap = values.length;
		this.elements = values;
		
	}
	public ArrayObject(){
		
	}
	public synchronized void set(int index , Object value) {
		if (index < 0 ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		index += this.start ; 
		if (index >= this.end ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		this.elements[index] = value ; 
	}
	public synchronized Object get(int index) {
		if (index < 0 ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		index += this.start ; 
		if (index >= this.end){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		return this.elements[index]  ; 
	}	
	

	public  synchronized ArrayObject slice(int start,int end){
		int length = end - start ;
		if(start < 0 || length < 0 || (length + this.start + start) > this.end){
			throw new ArrayIndexOutOfBoundsException(outOfRangeMsg);
		}
		ArrayObject result = new ArrayObject();
		result.elements = this.elements;
		result.start = this.start + start;
		result.end = result.start + length;
		result.cap = this.cap;
		return result;
	}
	public synchronized void append(Object e){
		if(this.end < this.cap){
		}else{
			this.expand(this.cap * 2);
		}
		this.elements[this.end++] = e;
	}
	public synchronized  void append(ArrayObject es){
		if (es == null) { //no need 
			return  ;
		}
		if(this.end + es.size() < this.cap){
		}else {
			this.expand((this.cap + es.size()) * 2);
		}
		for(int i = 0;i < es.size();i++){
			this.elements[this.end + i] = es.elements[es.start + i ];
		}
		this.end += es.size();
		 
	}
	private synchronized void expand(int cap){
		if(cap <= 0){
		    cap = 10;
		}
		Class c = this.elements.getClass();
		Object[] eles = (Object[]) Array.newInstance(c.getComponentType() , cap );
		int length = this.size();
		for(int i = 0;i < length;i++){
			eles[i] = this.elements[i + this.start];
		}
		this.start = 0;
		this.end = length;
		this.cap = cap;
		this.elements = eles;
	}
	
	public synchronized String toString(){
	    String s = "[";
	    int size = this.end - this.start;
	    for(int i= 0;i < size;i ++){
            s += this.elements[this.start + i ];
            if(i != size -1){
                s += " ";
            }
	    }
	    s += "]";
	    return s;
	}
}

