package install_lucy_array

const (
	array_template = `
package lucy.deps;
import java.lang.reflect.* ; 
IMPORTS
public class ArrayTTT   {
	public int start;
	public int end; // not include
	public int cap;
	static String outOfRangeMsg = "index out range";
	public TTT[] elements;
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
	public ArrayTTT(TTT[] values){
		this.start = 0;
		this.end = values.length;
		this.cap = values.length;
		this.elements = values;
		DEFAULT_INIT
	}
	public ArrayTTT(){

	}
	public synchronized void set(int index , TTT value) {
		if (index < 0 ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		index += this.start ; 
		if (index >= this.end ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		this.elements[index] = value ; 
	}
	public synchronized TTT get(int index) {
		if (index < 0 ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		index += this.start ; 
		if (index >= this.end){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		return this.elements[index]  ; 
	}	


	public  synchronized ArrayTTT slice(int start,int end){
		if(start < 0 || start > end || end + this.start > this.end){
			throw new ArrayIndexOutOfBoundsException(outOfRangeMsg);
		}
		ArrayTTT result = new ArrayTTT();
		result.elements = this.elements;
		result.start = this.start + start;
		result.end = this.start + end;
		result.cap = this.cap;
		return result;
	}
	public synchronized void append(TTT e){
		if(this.end < this.cap){
		}else{
			this.expand(this.cap * 2);
		}
		this.elements[this.end++] = e;
	}
	public synchronized  void append(ArrayTTT es){
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
		TTT[] eles = (TTT[]) Array.newInstance(c.getComponentType() , cap );
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
	public TTT[] getJavaArray(){
		if(this.start == 0 && this.end == this.cap){
			return this.elements;
		}
		int length = this.end - this.start;
		TTT[] elements = new TTT[length];
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

`
)
