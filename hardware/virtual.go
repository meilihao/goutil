package hardware

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

var (
	ErrRunWithoutRoot = errors.New("权限必须以root运行")
	ErrNoVirtWhat     = errors.New("未找到virt-what命令")
)

func VirtualInfo() (string, error) {
	if os.Geteuid() != 0 {
		return "", ErrRunWithoutRoot
	}

	cmdPath := "/usr/sbin/virt-what"
	if !FileIsExist(cmdPath) {
		return "", ErrNoVirtWhat
	}

	cmd := exec.Command("cmdPath")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(data)), nil
}
