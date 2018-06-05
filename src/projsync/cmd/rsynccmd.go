package cmd

///// IRsyncStep
type IRsyncStep interface {
	GetStepArgs() []string
}

///// RsyncCmd
type RsyncCmd struct {
	user     string
	password string
	host     string
	port     string

	// rsync step
	step IRsyncStep
}

func NewRsyncCmd() *RsyncCmd {
	return &RsyncCmd{}
}

func (cmd *RsyncCmd) GetCmdName() string {
	return "rsync"
}

func (cmd *RsyncCmd) GetCmdArgs() []string {
	args := make([]string, 0)
	args = append(args, "-v")
	args = append(args, "-t")
	args = append(args, "-e")
	args = append(args, cmd.genSSHArgs())
	args = append(args, cmd.genstepsCmd()...)
	return args
}

func (cmd *RsyncCmd) genSSHArgs() string {
	argtext := "ssh" + " -p " + cmd.port
	return argtext
}

func (cmd *RsyncCmd) genstepsCmd() []string {
	return cmd.step.GetStepArgs()
}

func (cmd *RsyncCmd) SetUser(user string) {
	cmd.user = user
}

func (cmd *RsyncCmd) SetPassword(password string) {
	cmd.password = password
}

func (cmd *RsyncCmd) SetHost(host string) {
	cmd.host = host
}

func (cmd *RsyncCmd) SetPort(port string) {
	cmd.port = port
}

func (cmd *RsyncCmd) SetRsyncStep(step IRsyncStep) {
	cmd.step = step
}

//// RsyncStepPutFile
type RsyncStepPutFile struct {
	user     string
	host     string

	localfile  string
	remotefile string
}

func NewRsyncStepPutFile() *RsyncStepPutFile {
	return &RsyncStepPutFile{}
}

func(step *RsyncStepPutFile) SetUser(user string) {
	step.user = user
}

func(step *RsyncStepPutFile) SetHost(host string) {
	step.host = host 
}
func (step *RsyncStepPutFile) SetLocalfile(file string) {
	step.localfile = file
}

func (step *RsyncStepPutFile) SetRemotefile(file string) {
	step.remotefile = file
}

func (step *RsyncStepPutFile) GetStepArgs() []string {
	args := make([]string, 0)
	args = append(args, step.localfile)
	args = append(args, step.genRemoteDirWithHost())
	return args
}

func (step *RsyncStepPutFile) genRemoteDirWithHost() string {
	remotefile := step.user + "@" + step.host + ":" + step.remotefile
	return remotefile
}
//// RsyncStepSync
type RsyncSyncDirection int

const (
	RSYNC_SYNC_DIRECTION_LOCAL_TO_REMOTE RsyncSyncDirection = 1
	RSYNC_SYNC_DIRECTION_REMOTE_TO_LOCAL RsyncSyncDirection = 2
)

type RsyncStepSync struct {
	user     string
	host     string

	localdir  string
	remotedir string

	// SyncDirect
	syncdirection RsyncSyncDirection

	// include file/dir, support file suffix
	include []string
	// exclude file/dir, support file suffix
	exclude []string
}

func NewRsyncStepSync() *RsyncStepSync {
	step := &RsyncStepSync{}
	step.SetDirection(RSYNC_SYNC_DIRECTION_LOCAL_TO_REMOTE)
	return step
}

func(step *RsyncStepSync) SetUser(user string) {
	step.user = user
}

func(step *RsyncStepSync) SetHost(host string) {
	step.host = host 
}

func (step *RsyncStepSync) SetLocalDir(dir string) {
	step.localdir = dir + "/"
}

func (step *RsyncStepSync) SetRemoteDir(dir string) {
	step.remotedir = dir + "/"
}

func (step *RsyncStepSync) SetDirection(direction RsyncSyncDirection) {
	step.syncdirection = direction
}

func (step *RsyncStepSync) GetStepArgs() []string {
	args := make([]string, 0)

	// include
	if len(step.include) > 0 {
		for _, include := range step.include {
			option := "--include=" + include
			args = append(args, option)
		}
	}

	// exclude
	if len(step.exclude) > 0 {
		for _, exclude := range step.exclude {
			option := "--exclude=" + exclude
			args = append(args, option)
		}
	}


	// recursive sync dir
	args = append(args, "-r")

	// sync
	if step.syncdirection == RSYNC_SYNC_DIRECTION_LOCAL_TO_REMOTE {
		args = append(args, step.localdir)
		args = append(args, step.genRemoteDirWithHost())
	} else if step.syncdirection == RSYNC_SYNC_DIRECTION_REMOTE_TO_LOCAL {
		args = append(args, step.genRemoteDirWithHost())
		args = append(args, step.localdir)
	}

	return args
}

func (step *RsyncStepSync) AddInclude(include string) {
	step.include = append(step.include, include)
}

func (step *RsyncStepSync) AddExclude(exclude string) {
	step.exclude = append(step.exclude, exclude)
}

func (step *RsyncStepSync) genRemoteDirWithHost() string {
	remotedir := step.user + "@" + step.host + ":" + step.remotedir
	return remotedir
}
