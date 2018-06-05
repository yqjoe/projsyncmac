package cmd

// SshCmd
type SshCmd struct {
	user     string
	host     string
	port     string

	calls []string
}

func NewSshCmd() *SshCmd {
	return &SshCmd{"", "", "", make([]string, 0)}
}

func (cmd *SshCmd) GetCmdName() string {
	return "ssh"
}

func (cmd *SshCmd) GetCmdArgs() []string {
	args := make([]string, 0)
	args = append(args, "-p")
	args = append(args, cmd.port)
	args = append(args, cmd.user + "@" + cmd.host)
	args = append(args, cmd.genCallArg())
	return args
}

func (cmd *SshCmd) AddCall(call string) {
	cmd.calls = append (cmd.calls, call)
}

func (cmd *SshCmd) SetUser(user string) {
	cmd.user = user
}

func (cmd *SshCmd) SetHost(host string) {
	cmd.host = host
}

func (cmd *SshCmd) SetPort(port string) {
	cmd.port = port
}

func (cmd *SshCmd) genCallArg() string {
	callarg := ""

	if len(cmd.calls) > 0 {
		for _, call := range cmd.calls {
			if len(callarg) > 0 {
				callarg += (";" + call)
			} else {
				callarg += call
			}
		}
	}

	return callarg
}