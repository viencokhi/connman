package connman

import (
	"fmt"
	"os/exec"
	"time"
)

type executionOut struct {
	out string
	err error
}

var (
	timeout time.Duration = 15 * time.Second
	debug   bool          = true
)

func errTimeOut(cmd string) error {
	return fmt.Errorf("Execution timeout cmd:%v", cmd)
}

func exe(cmd, cmdName string) (string, error) {
	fmt.Println(">>", cmd)
	exechan := make(chan executionOut, 1)
	go exeTimeout(cmd, cmdName, exechan)
	select {
	case exeOut := <-exechan:
		if exeOut.err != nil {
			fmt.Println("error")
			return "", exeOut.err
		}
		fmt.Println("done")
		return exeOut.out, nil
	case <-time.After(timeout):
		fmt.Println("timeout")
		return "", errTimeOut(cmd)
	}
}

func exeTimeout(cmd, cmdName string, exechan chan<- executionOut) {
	out, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		exechan <- executionOut{out: string(out), err: err}
	}
	exechan <- executionOut{out: string(out), err: nil}
}
