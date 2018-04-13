public class test{




	public void sayhai(int i){
		String s = "123";
		boolean f = true;
		if (f){
			String xxx = "456";
			String xxx2 = "456";
			boolean dd = false;
			boolean gg = true;
			boolean zz = dd && gg;
		}
		boolean x = false;
		boolean g = x && f;
		System.out.println(g && x);
	}

}