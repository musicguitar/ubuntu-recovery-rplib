package rplib

import (
"os"
"os/exec"
"strings"
)

func Shellexec(args ...string) {
        cmd := exec.Command(args[0], args[1:]...)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        err := cmd.Run()
        Checkerr(err)
}

func Shellexecoutput(args ...string) string {
        out, err := exec.Command(args[0], args[1:]...).Output()
        Checkerr(err)

        return strings.TrimSpace(string(out[:]))
}

func Shellcmd(command string) {
        cmd := exec.Command("sh", "-c", command)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        err := cmd.Run()
        Checkerr(err)
}

func Shellcmdoutput(command string) string {
        out, err := exec.Command("sh", "-c", command).Output()
        Checkerr(err)

        return strings.TrimSpace(string(out[:]))
}
