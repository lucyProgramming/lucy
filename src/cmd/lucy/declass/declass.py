import src.cmd.lucy.command as command
import getopt
import sys
import os
from optparse import OptionParser
import struct


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
            os.mkdir(self.__dest)

        self.__parseDir(self.__src,self.__dest)

        return 0

    def __parseDir(self,src ,dest):
        print("read dir " + src)
        if os.path.isdir(src)  == False :
            return
        fis = os.listdir(src)
        for d in fis:
            if d.endswith(".class"):  # class file
                if d.find("$") != -1:  #name contains $ means a inner class
                    continue
                self.__parseFile("%s/%s" % (src,d),dest)
            else:
                self.__parseDir("%s/%s" % (src,d),"%s/%s" % (dest,d))

    def __parseFile(self,src,dest):
        p = JvmClassParser(src,dest)
        ret = p.parse()
        if "ok" not in ret:
            print("declass file %s failed,err:%s" % (src,ret.reason))


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
        self.__descfilepath = destfilepath
        self.__result = JvmClass() # hold result in this
    def parse(self):  # file is definitely exits
        fd = open(self.__filepath,"rb")
        try:
            self.__content = fd.read()
        finally:
            fd.close()

        ok = self.__parseMagicAndVersion()
        if 0 != ok:
            return {"reason": ok}

        ok = self.__parseConstPool()
        if 0 != ok:
            return {"reason": ok}

        return {"ok":True}

    def __parseConstPool(self):
        ret =  struct.unpack_from("!H",self.__content[0:])
        size = ret[0]
        for i in range(size):
            continue

        return 0






    def __parseMagicAndVersion(self):
        ret = struct.unpack_from("!I",self.__content[0:])
        self.__result.magic = ret[0]
        self.__content = self.__content[4:]
        ret = struct.unpack_from("!HH",self.__content[0:])
        self.__result.minorVersion = ret[0]
        self.__result.majorVersion = ret[1]
        self.__content = self.__content[4:]
        return 0



