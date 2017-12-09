from src.cmd.lucy import command
import sys
import os


class Run(command.Command):
    def __init__(self):
        self.__files = []
        self.__packagename = ""
        self.__lucypath = []
        pass
    def __parsrArgs(self,args):
        files = []
        packagename = ""
        for v in args:
            if v.endswith(".lucy"):
                files.append(v)
            else:
                packagename = v
        if len(files) > 0 and packagename != "":
            print("run command mixed with package and lucy files")
            return -1
        if packagename != "":
            os.getenv("LUCYPATH")


        return 0

    def runCommand(self,args):
        lucypath = os.getenv("LUCYPATH")
        if lucypath == "" or lucypath == None:
            print("LUCYPATH is not set current directory as LUCYPATH")
            self.__lucypath = ["."]
            return
        else:
            t = lucypath.split(":") # lucy path is separated by :
            for v in t:
                if v != "":
                    self.__lucypath.append(v)

        print(self.__lucypath)
        args = args[1:] # skip run
        if 0 != self.__parsrArgs(args):
            Run.static_usage()
            sys.exit(1)



    def static_usage():
        print("run lucy files or package")