
package lucy.deps;
import java.lang.reflect.* ; 

public class ArrayBool   {
	public int start;
	public int end; // not include
	public int cap;
	boolean readOnly ; 
	static String outOfRangeMsg = "index out range";
	static String readOnlyArrayException = "array is readOnly";
	public boolean[] elements;
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
	public ArrayBool(boolean[] values){
		this.start = 0;
		this.end = values.length;
		this.cap = values.length;
		this.elements = values;
	}
	public ArrayBool(){
		
	}

	public synchronized setReadOnly(boolean readOnly) {
		this.readOnly = readOnly
	}

	public synchronized void set(int index , boolean value) {
		if (index < 0 ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		if if(this.readOnly) {
			throw new Exception(readOnlyArrayException);
		}
		index += this.start ; 
		if (index >= this.end ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		this.elements[index] = value ; 
	}

	public synchronized boolean get(int index) {
		if (index < 0 ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		index += this.start ; 
		if (index >= this.end){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		return this.elements[index]  ; 
	}

	public  synchronized ArrayBool slice(int start,int end){
		int length = end - start ;
		if(start < 0 || length < 0 || (length + this.start + start) > this.end){
			throw new ArrayIndexOutOfBoundsException(outOfRangeMsg);
		}
		ArrayBool result = new ArrayBool();
		result.elements = this.elements;
		result.start = this.start + start;
		result.end = result.start + length;
		result.cap = this.cap;
		return result;
	}

	public synchronized void append(boolean e){
		if if(this.readOnly) {
			throw new Exception(readOnlyArrayException);
		}
		if(this.end < this.cap){
		}else{
			this.expand(this.cap * 2);
		}
		this.elements[this.end++] = e;
	}
	public synchronized  void append(ArrayBool es){
		if if(this.readOnly) {
			throw new Exception(readOnlyArrayException);
		}
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
		boolean[] eles = (boolean[]) Array.newInstance(c.getComponentType() , cap );
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

