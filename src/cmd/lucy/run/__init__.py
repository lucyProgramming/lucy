from src.cmd.lucy import command
import sys
import os
from src.cmd.lucy import declass
import subprocess



class Run(command.Command):
    def __init__(self):
        self.__files = []
        self.__path = ""
        self.__lucypath = []
        self.__classfiles = []
        pass

    def runCommand(self,args):
        lucypath = os.getenv("LUCYPATH")
        if lucypath == "" or lucypath == None:
            print("LUCYPATH is not set current directory as LUCYPATH")
            self.__lucypath = ["."]
            return
        else:
            t = lucypath.split(";") # lucy path is separated by ;
            for v in t:
                if v != "":
                    self.__lucypath.append(v)

        if 0 != self.__parsrArgs(args):
            Run.static_usage()
            return
        # compile.exe is program write by go
        ret = self.__build()
        if ret != 0:
            print(ret)
            return



    def __build(self):

        return 1








    #parse arguments and read dir
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
            found = False
            for v in self.__lucypath:
                if v.endswith("/"):
                    t = "%s%s" % (v, packagename)
                else:
                    t = "%s/%s" % (v, packagename)
                if os.path.exists(t): # found
                    packagename = t
                    self.__path = packagename
                    found = True
                    break
            if found == False:
                print("package %s not found" % (packagename))
                return -1
            fis = os.listdir(packagename)
            for v in fis:
                if v.endswith(".class"):
                    self.__classfiles.append(v)
                    continue
                if v.endswith(".lucy"):
                    self.__files.append(v)
            if len(self.__files) == 0 :
                print("no lucy files found in %s" % (packagename))
                return -1
        return 0

    # at this stage ,always rebuild project
    def __check_if_need_rebuild(self,lucyfiles,classfiles):
        return True
        if len(classfiles) == 0:
            return True
        class_file_latest_timestamp = 0
        for v in classfiles:
            t = os.path.getmtime("%s/%s" % (self.__path,v))
            t = int(t)
            if t > class_file_latest_timestamp:
                class_file_latest_timestamp = t
        lucy_file_latest_timestamp = 0
        for v in lucyfiles:
            t = os.path.getmtime("%s/%s" % (self.__path, v))
            t = int(t)
            if t > lucy_file_latest_timestamp:
                lucy_file_latest_timestamp = t
        if lucy_file_latest_timestamp > class_file_latest_timestamp:
            return True











    def static_usage():
        print("run lucy files or package")