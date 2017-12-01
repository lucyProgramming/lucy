from src.cmd.lucy import command
import sys



class Run(command.Command):
    def __init__(self):
        pass
    def __parsrArgs(self,args):

    def runCommand(self,args):
        args = args[1:] # skip run
        if 0 != self.__parsrArgs(self,args):
            Run.static_usage()
            sys.exit(1)
    def static_usage():
        print("run lucy files or package")