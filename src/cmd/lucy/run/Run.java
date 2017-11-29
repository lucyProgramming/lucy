package cmd.lucy.run;

import cmd.lucy.command.Command;


/**
 * Created by yuyang on 17-11-21.
 */
public class Run implements  Command {
    String[] files ;
    String packagename;
    private void print_help(){
        String helpStr = "command run at lease except one argument\n";
        System.out.println(helpStr);
    }
    // non-0 means some error
    private int parseArgs(String[] args){
        this.files = new String[args.length-1];
        for (int i = 1;i<args.length;i++){
            if(args[i].endsWith(".lucy")){ //
                this.files[i-1] = args[i];
            }else{
                this.packagename = args[i];
            }
        }
        if ( this.packagename != "" && args.length > 2){
            System.out.println("when run a package only except one argument\n");
            return 1;  // not 0
        }
        return 0;
    }

    public  void RunCommand(String[] args) throws  Exception{
        if(args.length == 1){
            this.print_help();
            return ;
        }
        if(0 != this.parseArgs(args)){  // not 0 means something is error
            return ;
        }



    }
}
