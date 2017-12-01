import src.cmd.lucy.command as command


class Declass(command.Command):
    def parseParameter(self,args):
        return -1

    def help(self):
        print("help message")


    def runCommand(self,args):
        if self.parseParameter(args) != 0:
            self.help()
            return
        pass
