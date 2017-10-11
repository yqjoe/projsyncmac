package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/rpc"
	"os"
	"projsync/confmgr"
	"projsync/proto"
	"projsynctool/server"
	"strconv"
	"time"
)

func main() {
	if err := confmgr.Init(); err != nil {
		fmt.Println("confmgr.Init fail")
		return
	}

	client, err := rpc.Dial("tcp", ":6547")
	if err != nil {
		fmt.Println("Dial fail")
		return
	}

	fmt.Println("tool start:", os.Args)

	if len(os.Args) < 3 {
		fmt.Printf("Usage:%v projectname taskname", os.Args[0])
		return
	}

	projectname := os.Args[1]
	taskname := os.Args[2]

	conf := confmgr.GetTaskConf(projectname, taskname)
	if conf == nil {
		fmt.Println("ProjectName or TaskName not impl")
		return
	}

	// Add Task
	req := &proto.ReqAddTask{}
	err = initReqAddTask(req, projectname, taskname, os.Args)
	if err != nil {
		fmt.Println("err:", err.Error())
	}

	var rsp proto.RspAddTask
	if conf.TaskPrinter == "yes" {
		req.SyncToolPrintSvrAddr = genLocalAttr()
		client.Go("TaskServer.AddTask", req, &rsp, nil)
		closechan := make(chan int, 1)
		printersvr := server.NewPrinterServer(req.SyncToolPrintSvrAddr, closechan)
		go printersvr.Serve()
		<-closechan
	} else { // no
		req.SyncToolPrintSvrAddr = ""
		client.Call("TaskServer.AddTask", req, &rsp)
	}
}

func genLocalAttr() string {
	rand.Seed(time.Now().UnixNano())
	return (":" + strconv.Itoa(10000+rand.Intn(9999)))
}

func initReqAddTask(req *proto.ReqAddTask, projectname, taskname string, args []string) error {
	req.ProjectName = projectname
	req.TaskName = taskname
	switch taskname {
	case "savefile":
		if len(args) < 4 {
			return errors.New("Less Args")
		}
		req.Putstepfile = args[3]
	}
	return nil
}
