package main

import (
	"fmt"
	"projsync/confmgr"
	"projsync/server"
)

func main() {
	fmt.Println("im projsync")
	if err := confmgr.Init(); err != nil {
		fmt.Println("confmgr.Init fail")
		return 
	}

	server.RunTaskServer()
}
