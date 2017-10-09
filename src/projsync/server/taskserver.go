package server

import (
	"net/rpc"
	"net"
	"fmt"
	"projsync/proto"
	"projsync/task"
	"projsync/confmgr"
)

type TaskServer struct {
}

func (server *TaskServer) AddTask(req *proto.ReqAddTask, rsp *proto.RspAddTask) error {
	onetask := task.NewTask(req.ProjectName, req.TaskName)
	switch req.TaskName {
		case "savefile":
			onetask.SetPutStepFile(req.Putstepfile)
	}
	if len(req.SyncToolPrintSvrAddr) > 0 {
		onetask.SetSyncToolPrintSvrAddr(req.SyncToolPrintSvrAddr)
	}
	onetask.InitTaskFromConf()

	onetask.Run()
	return nil
}

func RunTaskServer() {
	tasksvr := new(TaskServer)
	svr := rpc.NewServer()
	svr.Register(tasksvr)

	l, err := net.Listen("tcp", confmgr.GetTaskServerAddr())
	if err != nil {
		fmt.Println("Listen fail")
		return
	}
	svr.Accept(l)
}
