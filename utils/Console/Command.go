package Console

import (
	"fmt"
	"os/exec"
	"syscall"
)

func RunCommand(name string, args []string) {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Start(); err != nil { // 运行命令
		fmt.Println(err)
	}
}

func RunBatFile(fileName string) {
	RunCommand("cmd", []string{"/c", "start", fileName})
}
