package util

import (
	"os"
	"os/exec"
	"syscall"
)

func Exec(name string, args ...string) {
	binary, lookErr := exec.LookPath(name)
	if lookErr != nil {
		panic(lookErr)
	}

	env := os.Environ()

	// Prepend the binary to the arguments slice
	args = append([]string{binary}, args...)

	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}
