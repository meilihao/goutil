package hardware

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var (
	ErrRunWithoutRoot = errors.New("权限必须以root运行")
)

// 建议使用高版本virt-what, 比如1.19, 因为1.14有未知原因的dump
func VirtualInfo() (string, error) {
	if data, _ := ioutil.ReadFile("/sys/class/dmi/id/sys_vendor"); len(data) > 0 { // phytium FT1500a(arm64)不支持dmi, 因此没有这个目录
		return string(bytes.TrimSpace(data)), nil
	}

	if os.Geteuid() != 0 {
		return "", ErrRunWithoutRoot
	}

	detectors := []struct {
		cmd          string
		args         []string
		pruneFunc    func(string) string
		resetErrFunc func(error) error
	}{
		{
			cmd:  "systemd-detect-virt",
			args: []string{"--vm"}, // systemd-detect-virt can check container
			pruneFunc: func(result string) string {
				if result == "qemu" {
					result = "kvm"
				}

				return result
			},
			resetErrFunc: func(err error) error { // systemd-detect-virt: exitcode of not virtaul is 1.
				var eerr *exec.ExitError
				if errors.As(err, &eerr) {
					return nil
				}

				return err
			},
		},
		{
			cmd:  "virt-what",
			args: nil,
			pruneFunc: func(result string) string {
				if result == "" {
					result = "none"
				}

				return result
			},
			resetErrFunc: nil,
		},
	}

	for i := range detectors {
		if tmp, _ := exec.LookPath(detectors[i].cmd); tmp == "" {
			continue
		}

		data, err := RunCMD(detectors[i].resetErrFunc, detectors[i].cmd, detectors[i].args)
		if err != nil {
			continue
		}
		if detectors[i].pruneFunc != nil {
			return detectors[i].pruneFunc(string(data)), nil
		}

		return string(data), nil
	}

	return "", errors.New("no get virtual tiype")
}

func RunCMD(resetErrFunc func(error) error, cpath string, args []string) ([]byte, error) {
	var cmd *exec.Cmd
	if len(args) > 0 {
		cmd = exec.Command(cpath, args...)
	} else {
		cmd = exec.Command(cpath)
	}

	data, err := cmd.CombinedOutput()
	if resetErrFunc != nil {
		err = resetErrFunc(err)
	}
	if err != nil {
		if len(data) > 0 {
			fmt.Printf("run %s err: %s\n", cpath, string(data))
		}

		return nil, err
	}

	return bytes.ToLower(bytes.TrimSpace(data)), nil
}
