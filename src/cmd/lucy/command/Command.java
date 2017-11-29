package cmd.lucy.command;



/**
 * Created by yuyang on 17-11-21.
 */
public interface Command {
    void RunCommand(String[] args) throws  Exception;
}
