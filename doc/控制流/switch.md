### **switch语句**
使用 switch 语句来选择要执行的多个代码块之一。

~~~
	var day int = 1;
	switch  day  {
		case 0:
			print("Today it's Sunday")
			break
		case 1:
			print("Today it's Monday")
			break
		case 2:
			print("Today it's Tuesday")
			break
		case 3:
			print("Today it's Wednesday")
			break
		case 4:
			print("Today it's Thursday")
			break
		case 5:
			print("Today it's Friday")
			break
		case 6:
			print("Today it's Saturday")
			break
		default:
			print("impossible")
	}
	
~~~

*break语句不是必须的，一个代码块执行完成后会自动跳出switch语句*


lucy支持一个case多个表达式示例如下，
~~~
	var day int = 1;
	switch  day  {
        case 0 ,1 :
            print("Today it's Sunday or Monday")
            break
        case 2,3:
            print("Today it's Tuesday or Wednesday")
            break
        case 4:
            print("Today it's Thursday")
            break
        case 5:
            print("Today it's Friday")
            break
        case 6:
            print("Today it's Saturday")
            break
        default:
            print("impossible")
}
~~~
<br/>

#### break语句
用于跳出switch。
<br/>


#### 枚举类型完备性
*要求枚举类型的case语句必须完备*
下面的例子定义了3个枚举项，虽然case只列举了2个但是有'default'，所有的分支都可以被执行到，所以可以通过编译。
~~~
	enum TT {
		XXX ,
		YYY,
		ZZZ 
	}
	
	var t TT;
	switch t{
		case XXX:
			print("xxx")
		case YYY:
			print("yyy")
		default:
	}
~~~
或者列举出所有的枚举项。

~~~
	enum TT {
		XXX ,
		YYY,
		ZZZ 
	}
	
	var t TT;
	switch t{
		case XXX:
			print("xxx")
		case YYY:
			print("yyy")
		case ZZZ:
	}
~~~
下面的语句无法通过编译。

~~~
	enum TT {
		XXX ,
		YYY,
		ZZZ 
	}
	
	var t TT;
	switch t{
		case XXX:
			print("xxx")
		case YYY:
			print("yyy")
		//misssing 'ZZZ' or default
	}
~~~

#### 和bool表达式混用
~~~

       private static fn hexByte2Int(b byte ) -> (v int ) {
               switch b {
                       case b >= 'a' && b <= 'f':
                               return int(b - 'a') + 10 
                       case b >= 'A' && b <= 'F':
                               return int(b - 'A') + 10 
                       default:
                               return int(b - '0')
               }
       }
//上面的语句看起来很像if语句，但是这种写法更简洁
//switch语句把每个case都当成一个条件
//当case表达式的值类型不是bool类型时去匹配和条件的值是否相等，相等则匹配
//当case表达式的值类型是bool时，case表达式的值为真是则匹配
~~~


