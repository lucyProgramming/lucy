/**
 * Created by yuyang on 17-11-21.
 */





import command.Command;
import run.*;
import java.util.HashMap;

public class Main {
    public static void print_help(){

        System.out.println("");
    }
    public  static HashMap<String,command.Command> handlers;




    public static  void main(String[] args){
        Main.handlers = new HashMap<String,command.Command>();
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



