
	

package lucy.deps;
public class ArrayDouble   {
	public int start;
	public int end; // not include
	public int cap;
	static String outOfRagneMsg = "index out range";
	public double[] elements;
	public int size(){
		return this.end - this.start;
	}
	public int start(){
        return this.start;
	}
	public int end(){
         return this.end;
	}
	public int cap(){
         return this.end;
	}
	public ArrayDouble(double[] values){
		this.start = 0;
		this.end = values.length;
		this.cap = values.length;
		this.elements = values;
	}
	private ArrayDouble(){
		
	}
	public ArrayDouble slice(int start,int end){
		if(end  < 0 ){
		      end = this.end - this.start;  // whole length
		}
		ArrayDouble result = new ArrayDouble();
		if(start < 0 || start > end || end + this.start > this.end){
			throw new ArrayIndexOutOfBoundsException(outOfRagneMsg);
		}
		result.elements = this.elements;
		result.start = this.start + start;
		result.end = this.start + end;
		result.cap = this.cap;
		return result;
	}
	public double get(int index){
		if(this.start + index >= this.end || index < 0){
			throw new ArrayIndexOutOfBoundsException(outOfRagneMsg);
		}
		return this.elements[this.start + index];
	}
	public void set(int index,double v){
		if(this.start + index >= this.end || index < 0){
			new ArrayIndexOutOfBoundsException(outOfRagneMsg);
		}
		this.elements[this.start + index] = v;
	}
	public ArrayDouble append(double e){
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
		double[] eles = new double[cap];
		int length = this.size();
		for(int i = 0;i < length;i++){
			eles[i] = this.elements[i + this.start];
		}
		this.start = 0;
		this.end = length;
		this.cap = cap;
		this.elements = eles;
	}
	public ArrayDouble append(double[] es){
		if(this.end + es.length < this.cap){
		}else {
			this.expand((this.cap + es.length) * 2);
		}
		for(int i = 0;i < es.length;i++){
			this.elements[this.end + i] = es[i];
		}
		this.end += es.length;
		return this;
	}
	public String toString(){
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

