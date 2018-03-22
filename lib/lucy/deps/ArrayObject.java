
	

package lucy.deps;
public class ArrayObject   {
	public int start;
	public int end; // not include
	public int cap;
	static String outOfRagneMsg = "index out range";
	public Object[] elements;
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
	public ArrayObject(Object[] values,int end){
		this.start = 0;
		this.end = end;
		this.cap = values.length;
		this.elements = values;
	}
	private ArrayObject(){
		
	}
	public ArrayObject slice(int start,int end){
		if(end  < 0 ){
		      end = this.end - this.start;  // whole length
		}
		ArrayObject result = new ArrayObject();
		if(start < 0 || start > end || end + this.start > this.end){
			throw new ArrayIndexOutOfBoundsException(outOfRagneMsg);
		}
		result.elements = this.elements;
		result.start = this.start + start;
		result.end = this.start + end;
		result.cap = this.cap;
		return result;
	}
	public Object get(int index){
		if(this.start + index >= this.end || index < 0){
			throw new ArrayIndexOutOfBoundsException(outOfRagneMsg);
		}
		return this.elements[this.start + index];
	}
	public void set(int index,Object v){
		if(this.start + index >= this.end || index < 0){
			new ArrayIndexOutOfBoundsException(outOfRagneMsg);
		}
		this.elements[this.start + index] = v;
	}
	public ArrayObject append(Object e){
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
		Object[] eles = new Object[cap];
		int length = this.size();
		for(int i = 0;i < length;i++){
			eles[i] = this.elements[i + this.start];
		}
		this.start = 0;
		this.end = length;
		this.cap = cap;
		this.elements = eles;
	}
	public ArrayObject append(Object[] es){
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

}

