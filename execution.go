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
	debugPrint("execution started command:", cmd, ", in function:", cmdName)
	exechan := make(chan executionOut, 1)
	go exeTimeout(cmd, cmdName, exechan)
	select {
	case exeOut := <-exechan:
		if exeOut.err != nil {
			debugPrint("execution got Error command:", cmd, ", in function:", cmdName, "err:", exeOut.err)
			return "", exeOut.err
		}
		debugPrint("execution done command:", cmd, ", in function:", cmdName, "out:", exeOut.out)
		return exeOut.out, nil
	case <-time.After(timeout):
		debugPrint("execution has timeout command:", cmd, ", in function:", cmdName)
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

func debugPrint(out ...interface{}) {
	fmt.Println(out...)
}
