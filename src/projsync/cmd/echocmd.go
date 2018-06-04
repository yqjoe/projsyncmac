package cmd

// EchoCmd
type EchoCmd struct {
	echoStrs []string
}

func NewEchoCmd() *EchoCmd {
	return &EchoCmd{make([]string, 0)}
}

func (cmd *EchoCmd) GetCmdName() string {
	return "echo"
}

func (cmd *EchoCmd) GetCmdArgs() []string {
	args := make([]string, 0)
	args = append(args, cmd.echoStrs...)
	return args
}

func (cmd *EchoCmd) AddEchoStr(str string) {
	cmd.echoStrs = append (cmd.echoStrs, str)
}