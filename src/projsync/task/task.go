package task

import (
	"fmt"
	"projsync/cmd"
	"projsync/confmgr"
	"strings"
	//"runtime"
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
	// 关闭printer
	defer task.taskprinter.Close()

	for _, icmd := range(task.cmdlist) {
		cmd.ExecCmd(icmd, task.taskprinter)
	}

	//fmt.Println("goroutine num:", runtime.NumGoroutine())
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
	case "svn":
		task.addTaskSvnCmd(cmdconf)
	case "xcopy":
		task.addTaskXcopyCmd(cmdconf)
	case "echo":
		task.addTaskEchoCmd(cmdconf)
	case "rsync":
		task.addTaskRsyncCmd(cmdconf)
	case "ssh":
		task.addTaskSshCmd(cmdconf)
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

func WinFormatRemotePath(path string) string {
	// char "\" to "/"
	return strings.Replace(path, "\\", "/", -1)
}

func WinFormatLocalPath(path string) string {
	// char "/" to "\"
	return strings.Replace(path, "/", "\\", -1)
}

func MacFormatPath(path string) string {
	// char "\" to "/"
	return strings.Replace(path, "\\", "/", -1)
}

func (task *Task) genRemotefileFromLocalfile(localfile string) string {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return ""
	}

	localdirlen := len(projconf.Localdir)
	remotefile := projconf.Remotedir + localfile[localdirlen:]
	return WinFormatRemotePath(remotefile)
}

func (task *Task) addTaskWinScpSyncStep(scpcmd *cmd.WinScpCmd, stepconf *confmgr.StepConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}

	syncstep := cmd.NewWinScpStepSync()
	syncstep.SetLocalDir(WinFormatLocalPath(projconf.Localdir + stepconf.Relativedir))
	syncstep.SetRemoteDir(WinFormatRemotePath(projconf.Remotedir + stepconf.Relativedir))
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

// svn task
func (task *Task) addTaskSvnCmd(cmdconf* confmgr.CmdConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}

	for _, stepconf := range(cmdconf.Step) {
		svncmd := cmd.NewSvnCmd()
		svncmd.SetOp(stepconf.StepName)
		svncmd.SetSvnDir(projconf.Localdir + stepconf.Relativedir)
		svncmd.SetUser(projconf.SvnUser)
		svncmd.SetPassword(projconf.SvnPassword)

		task.cmdlist = append(task.cmdlist, svncmd)
	}
}

// xcopy task
func (task *Task) addTaskXcopyCmd(cmdconf* confmgr.CmdConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}

	for _, stepconf := range(cmdconf.Step) {
		xcopycmd := cmd.NewXcopyCmd()
		xcopycmd.SetSrcFiles(projconf.Localdir + "\\" + stepconf.Relativedir + "\\" + stepconf.FileName)
		xcopycmd.SetDstDir(projconf.Localdir + "\\" + stepconf.DstRelativedir)

		task.cmdlist = append(task.cmdlist, xcopycmd)
	}
}

// echo task
func (task *Task) addTaskEchoCmd(cmdconf* confmgr.CmdConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}

	for _, stepconf := range(cmdconf.Step) {
		echocmd := cmd.NewEchoCmd()
		echocmd.AddEchoStr(stepconf.StepName)

		task.cmdlist = append(task.cmdlist, echocmd)
	}
}

// rsync task
func (task *Task) addTaskRsyncCmd(cmdconf* confmgr.CmdConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}
	
	for _, stepconf := range(cmdconf.Step) {
		rsynccmd := cmd.NewRsyncCmd()
		
		rsynccmd.SetUser(projconf.User)
		rsynccmd.SetPassword(projconf.Password)
		rsynccmd.SetHost(projconf.Host)
		rsynccmd.SetPort(projconf.Port)
		task.addTaskRsyncStep(rsynccmd, &stepconf)

		task.cmdlist = append(task.cmdlist, rsynccmd)
	}
}

func (task *Task) addTaskRsyncStep(rsynccmd *cmd.RsyncCmd, stepconf *confmgr.StepConf) {
	switch stepconf.StepName {
		case "put":
			task.addTaskRsyncPutStep(rsynccmd, stepconf)
		case "sync":
			task.addTaskRsyncSyncStep(rsynccmd, stepconf)
		case "syncdst":
			task.addTaskRsyncSyncDstStep(rsynccmd, stepconf)
		default:
			fmt.Println("Step Not Impl:", stepconf.StepName)
	}
}

func (task *Task) addTaskRsyncPutStep(rsynccmd *cmd.RsyncCmd, stepconf *confmgr.StepConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}

	putstep := cmd.NewRsyncStepPutFile()
	putstep.SetUser(projconf.User)
	putstep.SetHost(projconf.Host)
	putstep.SetLocalfile(task.putstepfile)
	putstep.SetRemotefile(task.genRemotefileFromLocalfile(task.putstepfile))
	rsynccmd.SetRsyncStep(putstep)
}

func (task *Task) addTaskRsyncSyncStep(rsynccmd *cmd.RsyncCmd, stepconf *confmgr.StepConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}

	syncstep := cmd.NewRsyncStepSync()
	syncstep.SetUser(projconf.User)
	syncstep.SetHost(projconf.Host)
	syncstep.SetLocalDir(MacFormatPath(projconf.Localdir + stepconf.Relativedir))
	syncstep.SetRemoteDir(MacFormatPath(projconf.Remotedir + stepconf.Relativedir))
	if stepconf.SyncDirection == "local2remote" {
		syncstep.SetDirection(cmd.RSYNC_SYNC_DIRECTION_LOCAL_TO_REMOTE)
	} else {
		syncstep.SetDirection(cmd.RSYNC_SYNC_DIRECTION_REMOTE_TO_LOCAL)
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
	
	rsynccmd.SetRsyncStep(syncstep)
}

func (task *Task) addTaskRsyncSyncDstStep(rsynccmd *cmd.RsyncCmd, stepconf *confmgr.StepConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}

	syncstep := cmd.NewRsyncStepSync()
	syncstep.SetUser(projconf.User)
	syncstep.SetHost(projconf.Host)
	syncstep.SetLocalDir(MacFormatPath(projconf.Localdir + stepconf.Relativedir))
	syncstep.SetRemoteDir(MacFormatPath(projconf.Remotedir + stepconf.DstRelativedir))
	if stepconf.SyncDirection == "local2remote" {
		syncstep.SetDirection(cmd.RSYNC_SYNC_DIRECTION_LOCAL_TO_REMOTE)
	} else {
		syncstep.SetDirection(cmd.RSYNC_SYNC_DIRECTION_REMOTE_TO_LOCAL)
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
	
	rsynccmd.SetRsyncStep(syncstep)
}

// sshtask
func (task *Task) addTaskSshCmd(cmdconf* confmgr.CmdConf) {
	projconf := confmgr.GetProjectConf(task.ProjectName)
	if nil == projconf {
		return
	}

	for _, stepconf := range(cmdconf.Step) {
		sshcmd := cmd.NewSshCmd()
		for _, shellcmd := range(stepconf.ShellCmd) {
			sshcmd.AddCall(shellcmd)
			sshcmd.SetUser(projconf.User)
			sshcmd.SetHost(projconf.Host)
			sshcmd.SetPort(projconf.Port)
		}

		task.cmdlist = append(task.cmdlist, sshcmd)
	}
}
