package rplib

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func Shellexec(args ...string) {
	log.Println(strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	Checkerr(err)
}

func Shellexecoutput(args ...string) string {
	log.Println(strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.Output()
	Checkerr(err)

	return strings.TrimSpace(string(out[:]))
}

func Shellcmd(command string) {
	cmd := exec.Command("sh", "-c", command)
	log.Println(cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	Checkerr(err)
}

func Shellcmdoutput(command string) string {
	cmd := exec.Command("sh", "-c", command)
	log.Println(cmd.Args)
	out, err := cmd.Output()
	Checkerr(err)

	return strings.TrimSpace(string(out[:]))
}
