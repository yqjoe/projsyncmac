package task

import (
	"fmt"
	"projsync/cmd"
	"projsync/confmgr"
	"strings"
)

type Task struct {
	ProjectName string
	TaskName string

	cmdlist []cmd.ICmd
	putstepfile string // winscp put 命令 同步的文件，本地文件的全路径

	// sync tool 打印服务器地址，如果没有则为空
	synctoolprintsvraddr string
	taskprinter *TaskPrinter
}

func NewTask(projname, taskname string) *Task {
	return &Task {
		projname,
		taskname,
		make([]cmd.ICmd, 0),
		"",
		"",
		nil}
}

func (task *Task) Run() {
	// 连接synctool print svr
	task.taskprinter = NewTaskPrinter(task.synctoolprintsvraddr)

	for _, icmd := range(task.cmdlist) {
		cmd.ExecCmd(icmd, task.taskprinter)
	}

	// 关闭printer
	task.taskprinter.Close()
}

func (task *Task) SetPutStepFile(file string) {
	task.putstepfile = file
}

func (task *Task) SetSyncToolPrintSvrAddr(addr string) {
	task.synctoolprintsvraddr = addr
}

func (task *Task) InitTaskFromConf() {
	taskconfobj := confmgr.GetTaskConf(task.ProjectName, task.TaskName)
	if taskconfobj == nil {
		fmt.Println("ProjectTaskConf Not found", taskconfobj.TaskName)
		return
	}

	for _, cmd := range(taskconfobj.Cmd) {
		task.addTaskCmd(&cmd)		
	}
}

func (task *Task) addTaskCmd(cmdconf *confmgr.CmdConf) {
	switch cmdconf.CmdName {
	case "winscp":
		task.addTaskWinScpCmd(cmdconf)
	default:
		fmt.Println("cmd not impl:", cmdconf.CmdName)
	}
}

func (task *Task) addTaskWinScpCmd(cmdconf *confmgr.CmdConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}

	scpcmd := cmd.NewWinScpCmd();
	scpcmd.SetUser(projconf.User)
	scpcmd.SetPassword(projconf.Password)
	scpcmd.SetHost(projconf.Host)
	scpcmd.SetPort(projconf.Port)

	for _, stepconf := range(cmdconf.Step) {
		task.addTaskWinScpStep(scpcmd, &stepconf)
	}

	task.cmdlist = append(task.cmdlist, scpcmd)
}

func (task *Task) addTaskWinScpStep(scpcmd *cmd.WinScpCmd, stepconf *confmgr.StepConf) {
	switch stepconf.StepName {
		case "put":
			task.addTaskWinScpPutStep(scpcmd, stepconf)
		case "sync":
			task.addTaskWinScpSyncStep(scpcmd, stepconf)
		case "call":
			task.addTaskWinScpCallStep(scpcmd, stepconf)
		default:
			fmt.Println("Step Not Impl:", stepconf.StepName)
	}
}

func (task *Task) addTaskWinScpPutStep(scpcmd *cmd.WinScpCmd, stepconf *confmgr.StepConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}

	putstep := cmd.NewWinScpStepPutFile()
	putstep.SetLocalfile(task.putstepfile)
	putstep.SetRemotefile(task.genRemotefileFromLocalfile(task.putstepfile))
	scpcmd.AddWinScpStep(putstep)
}

func FormatRemotePath(path string) string {
	// char "\" to "/"
	return strings.Replace(path, "\\", "/", -1)
}

func FormatLocalPath(path string) string {
	// char "/" to "\"
	return strings.Replace(path, "/", "\\", -1)
}

func (task *Task) genRemotefileFromLocalfile(localfile string) string {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return ""
	}

	localdirlen := len(projconf.Localdir)
	remotefile := projconf.Remotedir + localfile[localdirlen:]
	return FormatRemotePath(remotefile)
}

func (task *Task) addTaskWinScpSyncStep(scpcmd *cmd.WinScpCmd, stepconf *confmgr.StepConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}

	syncstep := cmd.NewWinScpStepSync()
	syncstep.SetLocalDir(FormatLocalPath(projconf.Localdir + stepconf.Relativedir))
	syncstep.SetRemoteDir(FormatRemotePath(projconf.Remotedir + stepconf.Relativedir))
	if stepconf.SyncDirection == "local2remote" {
		syncstep.SetDirection(cmd.WIN_SCP_SYNC_DIRECTION_LOCAL_TO_REMOTE)
	} else {
		syncstep.SetDirection(cmd.WIN_SCP_SYNC_DIRECTION_REMOTE_TO_LOCAL)
	}
	for _, include := range(stepconf.Include) {
		if len(include) > 0 {
			syncstep.AddInclude(include)
		}
	}
	for _, exclude := range(stepconf.Exclude) {
		if len(exclude) > 0 {
			syncstep.AddExclude(exclude)
		}
	}
	
	scpcmd.AddWinScpStep(syncstep)
}

func (task *Task) addTaskWinScpCallStep(scpcmd *cmd.WinScpCmd, stepconf *confmgr.StepConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}

	callstep := cmd.NewWinScpStepCall()
	for _, shellcmd := range(stepconf.ShellCmd) {
		callstep.AddShellCmd(shellcmd)
	}

	scpcmd.AddWinScpStep(callstep)
}