

public class ArrayInt extends Array {
	private int[] elements;
	public ArrayInt(int length,int[] values){
		this.start = 0;
		this.end = length;
		this.cap = length * 2;
		this.elements = new int[this.cap];
		if(values != null){
			for(int i = 0;i < values.length;i++){
				this.elements[i] = values[i];
			}
		}
	}
	private ArrayInt(){
		
	}
	public ArrayInt slice(int start,int end){
		ArrayInt result = new ArrayInt();
		if(start < 0 || start > end || end >= this.end){
			new ArrayIndexOutOfBoundsException(outOfRagneMsg);
		}
		result.elements = this.elements;
		result.start = start;
		result.end = end;
		result.cap = this.cap;
		return result;
	}
	public int get(int index){
		if(this.start + index >= this.end || index < 0){
			new ArrayIndexOutOfBoundsException(outOfRagneMsg);
		}
		return this.elements[this.start + index];
	}
	public void append(int e){
		if(this.end < this.cap){
			this.elements[this.end++] = e;
			return ;
		}
		this.expand(this.cap * 2);
		this.elements[this.end++] = e;
	}
	private void expand(int cap){
		int[] eles = new int[cap];
		int length = this.size();
		for(int i = 0;i < length;i++){
			eles[i] = this.elements[i + this.start];
		}
		this.start = 0;
		this.end = length;
		this.cap = cap;
		this.elements = eles;
	}
	public void append(int[] es){
		if(this.end + es.length < this.cap){
			for(int i = 0;i < es.length;i++){
				this.elements[this.end + i] = es[i];
			}
			this.end += es.length;
			return ;
		}
		this.expand(this.cap * 2);
		for(int i = 0;i < es.length;i++){
			this.elements[this.end + i] = es[i];
		}
		this.end += es.length;
	}
}