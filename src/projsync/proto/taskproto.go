package proto

// 加任务协议
type ReqAddTask struct {
	ProjectName string
	TaskName string

	Putstepfile string // winscp put 命令 同步的文件，本地文件的全路径

	// 同步工具打印日志服务器的地址，如果不开启，则为空
	SyncToolPrintSvrAddr string
}

type RspAddTask struct {
	Ret int
	Err string
}

// 输出任务执行信息协议
type ReqPrintTaskInfo struct {
	Info string
}

type RspPrintTaskInfo int

// 关闭printersvr
type ReqClosePrinterSvr int
type RspClosePrinterSvr int
