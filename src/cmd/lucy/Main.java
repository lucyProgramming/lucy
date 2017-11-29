/**
 * Created by yuyang on 17-11-21.
 */





import cmd.lucy.command.*;
import cmd.lucy.run.*;
import java.util.HashMap;
import cmd.lucy.declass.*;

public class Main {
    public static void print_help(){
        String msg = "lucy\n";
        msg += "\t run       run a package or files\n";

        System.out.println(msg);
    }
    public  static HashMap<String, Command> handlers;

    public static  void main(String[] args){
        Main.handlers = new HashMap<String,Command>();
        Main.handlers.put("run",new Run());
        if(args.length == 0){
            Main.print_help();
            return ;
        }
        
        if (!Main.handlers.containsKey(args[0])){
            System.out.println("command " + args[0] + " is unkown");
            Main.print_help();
            return;
        }
        try {
            Main.handlers.get(args[0]).RunCommand(args);
        }catch (Exception e ){
            System.out.println("run command failed,err:" + e.toString());
        }
    }
}



