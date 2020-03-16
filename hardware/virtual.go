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

// 建议使用高版本virt-what, 比如1.19, 因为1.14有未知原因的dump
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
