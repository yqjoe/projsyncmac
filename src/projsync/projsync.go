package main

import (
	"fmt"
	"projsync/confmgr"
	"projsync/server"
)

func main() {
	if err := confmgr.Init(); err != nil {
		fmt.Println("confmgr.Init fail")
		return 
	}

	fmt.Println("projsync svr running")

	server.RunTaskServer()
}
