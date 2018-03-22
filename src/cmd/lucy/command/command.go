package command

type RunCommand interface {
	RunCommand(command string, args []string)
}
