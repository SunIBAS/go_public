package Console

import (
	"os"
	"os/exec"
)

func Clear() {
	cmd := exec.Command("cmd.exe", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
