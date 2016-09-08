package rplib

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	WritableImage = "writable_resized.e2fs"
)

func DD(input string, output string, args ...string) {
	args = append([]string{fmt.Sprintf("if=%s", input), fmt.Sprintf("of=%s", output)}, args...)
	Shellexec("dd", args...)
}

func Sync() {
	Shellexec("sync")
}

func Reboot() {
	Shellexec("reboot")
}

func FindDevice(blockDevice string) (device string, err error) {
	syspath := path.Dir(Realpath(filepath.Join("/sys/class/block", path.Base(blockDevice))))

	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/dev", syspath))
	if err != nil {
		return "", err
	}
	dat_str := strings.TrimSpace(string(dat))
	device = Realpath(fmt.Sprintf("/dev/block/%s", dat_str))
	return device, nil
}

func Findfs(arg string) string {
	return Shellexecoutput("findfs", arg)
}

func Realpath(path string) string {
	newPath, err := filepath.Abs(path)
	if err != nil {
		log.Panic(err)
	}

	newPath, err = filepath.EvalSymlinks(newPath)
	if err != nil {
		log.Panic(err)
	}
	return newPath
}

func SetPartitionFlag(device string, nr int, flag string) {
	Shellexec("parted", "-ms", device, "set", fmt.Sprintf("%v", nr), flag, "on")
}

func BlockSize(block string) (size int64) {
	// unit Byte
	sizeStr := Shellexecoutput("blockdev", "--getsize64", block)
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	Checkerr(err)
	return
}

func GetPartitionBeginEnd(device string, nr int) (begin, end int) {
	var err error
	line := Shellcmdoutput(fmt.Sprintf("parted -ms %s unit B print | grep \"^%d:\"", device, nr))
	log.Println("line:", line)
	fields := strings.Split(line, ":")
	begin, err = strconv.Atoi(strings.TrimRight(fields[1], "B"))
	Checkerr(err)
	end, err = strconv.Atoi(strings.TrimRight(fields[2], "B"))
	Checkerr(err)
	return
}

func GetPartitionBeginEnd64(device string, nr int) (begin, end int64) {
	var err error
	line := Shellcmdoutput(fmt.Sprintf("parted -ms %s unit B print | grep \"^%d:\"", device, nr))
	log.Println("line:", line)
	fields := strings.Split(line, ":")
	begin, err = strconv.ParseInt(strings.TrimRight(fields[1], "B"), 10, 64)
	Checkerr(err)
	end, err = strconv.ParseInt(strings.TrimRight(fields[2], "B"), 10, 64)
	Checkerr(err)
	return
}

func GetBootEntries(keyword string) (entries []string) {
	entryStr := Shellcmdoutput(fmt.Sprintf("efibootmgr -v | grep \"%s\" | cut -f 1 | sed 's/[^0-9]*//g'", keyword))
	log.Printf("entryStr: [%s]\n", entryStr)
	if "" == entryStr {
		entries = []string{}
	} else {
		entries = strings.Split(entryStr, "\n")
	}
	log.Println("entries:", entries)
	return
}

func CreateBootEntry(device string, partition int, loader string, label string) {
	Shellexec("efibootmgr", "-c", "-d", device, "-p", fmt.Sprintf("%v", partition), "-l", loader, "-L", label)
}

func ReadKernelCmdline() string {
	data, err := ioutil.ReadFile("/proc/cmdline")
	Checkerr(err)
	cmdline := string(data)
	return cmdline
}

func IsKernelCmdlineContains(substr string) bool {
	return strings.Contains(ReadKernelCmdline(), substr)
}
