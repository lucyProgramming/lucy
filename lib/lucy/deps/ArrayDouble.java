
package lucy.deps;
import java.lang.reflect.* ; 

public class ArrayDouble   {
	public int start;
	public int end; // not include
	public int cap;
	static String outOfRangeMsg = "index out range";
	public double[] elements;
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
	public ArrayDouble(double[] values){
		this.start = 0;
		this.end = values.length;
		this.cap = values.length;
		this.elements = values;
		
	}
	private ArrayDouble(){

	}
	public synchronized void set(int index , double value) {
		if (index < 0 ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		index += this.start ; 
		if (index >= this.end ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		this.elements[index] = value ; 
	}
	public synchronized double get(int index) {
		if (index < 0 ){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		index += this.start ; 
		if (index >= this.end){
			throw new ArrayIndexOutOfBoundsException (outOfRangeMsg);
		}
		return this.elements[index]  ; 
	}	


	public  synchronized ArrayDouble slice(int start,int end){
		if(start < 0 || start > end || end + this.start > this.end){
			throw new ArrayIndexOutOfBoundsException(outOfRangeMsg);
		}
		ArrayDouble result = new ArrayDouble();
		result.elements = this.elements;
		result.start = this.start + start;
		result.end = this.start + end;
		result.cap = this.cap;
		return result;
	}
	public synchronized ArrayDouble append(double e){
		if(this.end < this.cap){
		}else{
			this.expand(this.cap * 2);
		}
		this.elements[this.end++] = e;
		return this;
	}
	private void expand(int cap){
		if(cap <= 0){
		    cap = 10;
		}
		Class c = this.elements.getClass();
		double[] eles = (double[]) Array.newInstance(c.getComponentType() , cap );
		int length = this.size();
		for(int i = 0;i < length;i++){
			eles[i] = this.elements[i + this.start];
		}
		this.start = 0;
		this.end = length;
		this.cap = cap;
		this.elements = eles;
	}
	public synchronized  ArrayDouble append(ArrayDouble es){
		if(this.end + es.size() < this.cap){
		}else {
			this.expand((this.cap + es.size()) * 2);
		}
		for(int i = 0;i < es.size();i++){
			this.elements[this.end + i] = es.elements[es.start + i ];
		}
		this.end += es.size();
		return this;
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
	public double[] getJavaArray(){
		if(this.start == 0 && this.end == this.cap){
			return this.elements;
		}
		int length = this.end - this.start;
		double[] elements = new double[length];
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

