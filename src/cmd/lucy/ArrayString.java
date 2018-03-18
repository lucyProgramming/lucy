



package lucy.lang;
public class ArrayString   {
	public int start;
	public int end; // not include
	public int cap;
	static String outOfRagneMsg = "index out range";
	public String[] elements;
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
	public ArrayString(String[] values,int end){
		this.start = 0;
		this.end = end;
		this.cap = values.length;
		this.elements = values;
	}
	private ArrayString(){

	}
	public ArrayString slice(int start,int end){
		if(end  < 0 ){
		      end = this.end - this.start;  // whole length
		}
		ArrayString result = new ArrayString();
		if(start < 0 || start > end || end + this.start > this.end){
			throw new ArrayIndexOutOfBoundsException(outOfRagneMsg);
		}
		result.elements = this.elements;
		result.start = this.start + start;
		result.end = this.start + end;
		result.cap = this.cap;
		return result;
	}
	public String get(int index){
		if(this.start + index >= this.end || index < 0){
			throw new ArrayIndexOutOfBoundsException(outOfRagneMsg);
		}
		return this.elements[this.start + index];
	}
	public void set(int index,String v){
		if(this.start + index >= this.end || index < 0){
			new ArrayIndexOutOfBoundsException(outOfRagneMsg);
		}
		this.elements[this.start + index] = v;
	}
	public void append(String e){
		if(this.end < this.cap){
		}else{
			this.expand(this.cap * 2);
		}
		this.elements[this.end++] = e;
	}
	private void expand(int cap){
		String[] eles = new String[cap];
		int length = this.size();
		for(int i = 0;i < length;i++){
			eles[i] = this.elements[i + this.start];
		}
		this.start = 0;
		this.end = length;
		this.cap = cap;
		this.elements = eles;
	}
	public void append(String[] es){
		if(this.end + es.length < this.cap){
		}else {
			this.expand((this.cap + es.length) * 2);
		}
		for(int i = 0;i < es.length;i++){
			this.elements[this.end + i] = es[i];
		}
		this.end += es.length;
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











