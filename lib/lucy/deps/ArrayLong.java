
package lucy.deps;
import java.lang.reflect.* ; 

public class ArrayLong   {
	public int start;
	public int end; // not include
	public int cap;
	static String outOfRangeMsg = "index out range";
	public long[] elements;
	public int size(){
		return this.end - this.start;
	}
	public synchronized int start(){
        return this.start;
	}
	public int end(){
         return this.end;
	}
	public int cap(){
         return this.cap;
	}
	public ArrayLong(long[] values){
		this.start = 0;
		this.end = values.length;
		this.cap = values.length;
		this.elements = values;
		
	}
	public ArrayLong(){

	}
	public synchronized void set(int index , long value) {
		if (index < 0 ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		index += this.start ; 
		if (index >= this.end ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		this.elements[index] = value ; 
	}
	public synchronized long get(int index) {
		if (index < 0 ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		index += this.start ; 
		if (index >= this.end){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		return this.elements[index]  ; 
	}	


	public  synchronized ArrayLong slice(int start,int end){
		if(start < 0 || start > end || end + this.start > this.end){
			throw new ArrayIndexOutOfBoundsException(outOfRangeMsg);
		}
		ArrayLong result = new ArrayLong();
		result.elements = this.elements;
		result.start = this.start + start;
		result.end = this.start + end;
		result.cap = this.cap;
		return result;
	}
	public synchronized void append(long e){
		if(this.end < this.cap){
		}else{
			this.expand(this.cap * 2);
		}
		this.elements[this.end++] = e;
	}
	public synchronized  void append(ArrayLong es){
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
	private void expand(int cap){
		if(cap <= 0){
		    cap = 10;
		}
		Class c = this.elements.getClass();
		long[] eles = (long[]) Array.newInstance(c.getComponentType() , cap );
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
	public long[] getJavaArray(){
		if(this.start == 0 && this.end == this.cap){
			return this.elements;
		}
		int length = this.end - this.start;
		long[] elements = new long[length];
		for(int i = 0; i < length; i ++){
			elements[i] = this.elements[i + this.start];
		}
		this.start = 0;
		this.end = length;
		this.elements = elements;
		this.cap = length;
		return elements;
	}
}

