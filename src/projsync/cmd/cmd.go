package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"io"
)

//// ICmd interface
type ICmd interface {
	GetCmdName() string
	GetCmdArgs() []string
}

func ExecCmd(cmd ICmd, printer io.Writer) {
	fmt.Println("ExecCmd")
	execCmd(printer, cmd.GetCmdName(), cmd.GetCmdArgs()...)
}

func execCmd(printer io.Writer, cmdname string, args ...string) {
	fmt.Println(args)
	c := exec.Command(cmdname, args...)
	var cmderr bytes.Buffer
	c.Stdout = printer 
	c.Stderr = &cmderr

	if err := c.Run(); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("ErrInfo:", cmderr.String())
	}
}