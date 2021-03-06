package helpers

import (
	"os"
	"os/exec"
)

func Command(name string, args ...string) (cmd *exec.Cmd) {
	cmd = exec.Command(name, args...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	return
}

func CommandOutput(name string, args ...string) (cmd *exec.Cmd) {
	cmd = Command(name, args...)
	cmd.Stdout = os.Stdout
	return
}

func System(command string) (err error) {
	cmd := CommandOutput("/bin/sh", "-c", command)
	return cmd.Run()
}

func Docker(command string, args ...string) (cmd *exec.Cmd) {
	return Command("docker", append([]string{command}, args...)...)
}

func Compose(command string, path string, args ...string) (cmd *exec.Cmd) {
	cmd = Command("docker-compose", append([]string{command}, args...)...)
	cmd.Dir = path
	return
}

func Git(args ...string) (cmd *exec.Cmd) {
	return Command("git", args...)
}
