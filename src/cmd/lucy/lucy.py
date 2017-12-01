
from  src.cmd.lucy import  run
from  src.cmd.lucy import command
from  src.cmd.lucy import declass
import sys


class Help(command.Command):
    def runCommand(self):
        str = '''
            lucy command tool
                run         run lucy file or package
                declass     declass jvm class files
        
        '''
        print("show help")






commands = {
    "run":run.Run(),
    "help":Help(),
    "declass":declass.Declass(),
}


if sys.argv[1] in commands.keys():
    commands[sys.argv[1]].runCommand(sys.argv)
else:
    commands["help"].runCommand()












