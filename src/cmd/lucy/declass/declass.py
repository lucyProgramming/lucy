import src.cmd.lucy.command as command
import getopt
import sys
import os
from optparse import OptionParser

class Declass(command.Command):
    def __init__(self):
        self.__src = ""
        self.__dest = ""
        self.__help_msg = "declass jvm class files,command line args are -src and -dest"

    def __parseParameter(self,args):
        parser = OptionParser(usage=self.__help_msg)
        parser.add_option("--src",action="store",type="string",dest="src",help="source directory")
        parser.add_option("--dest", action="store", type="string", dest="dest", help="destination directory")
        opt,args = parser.parse_args(args)
        if opt.dest == "" or opt.src == "":
            Declass.static_usage()
            sys.exit(1)
        self.__src = opt.src
        self.__dest = opt.dest
        return 0

    def __parse(self):
        if os.path.exists(self.__src) == False:
            print("src %s directory is not exits" % (self.__src))
            return
        if os.path.exists(self.__dest) == False:
            print("dest %s directory is not exits" % (self.__dest))
            return
        return 0

    def static_usage():
        print("declass jvm class files,command line args are -src and -dest")

    def runCommand(self,args):
        args = args[1:] # skip run command
        if self.__parseParameter(args) != 0:
            sys.exit(1)

        if 0 != self.__parse():
            sys.exit(2)





class JvmClass:
    def __init__(self):
        pass


class JvmClassParser:
    def __init__(self,filepath,destfilepath):
        self.__filepath = filepath
        sefl.__descfilepath = destfilepath
    def parse(self):
        pass
