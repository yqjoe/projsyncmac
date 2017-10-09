package task

import (
	"net/rpc"
	"fmt"
	"projsync/proto"
)

type TaskPrinter struct {
	client *rpc.Client
}

func NewTaskPrinter(printersvraddr string) *TaskPrinter {
	tp := &TaskPrinter{nil}
	if len(printersvraddr) == 0 {
		// 不远程发送
		return tp
	}

	var err  error
	tp.client, err = rpc.Dial("tcp", printersvraddr)
	if err != nil {
		tp.client = nil
	}

	return tp
}


func (printer *TaskPrinter) Write(p []byte) (n int, err error) {
	//printer.client.Call("")
	if printer.client == nil {
		fmt.Printf("%v", string(p))
	} else {
		req := &proto.ReqPrintTaskInfo{}
		req.Info = string(p)
		var rsp proto.RspPrintTaskInfo
		printer.client.Go("RpcPrinterServer.PrintTaskInfo", req, &rsp, nil)
	}
	return len(p), nil
}

func (printer *TaskPrinter) Close() {
	if printer.client != nil {
		var req proto.ReqClosePrinterSvr = 0
		var rsp proto.RspClosePrinterSvr
		printer.client.Call("RpcPrinterServer.ClosePrinterSvr", req, &rsp)
		printer.client.Close()
	}
}