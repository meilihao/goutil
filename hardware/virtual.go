package hardware

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var (
	ErrRunWithoutRoot = errors.New("权限必须以root运行")
)

func VirtualInfo() (string, error) {
	if os.Geteuid() != 0 {
		return "", ErrRunWithoutRoot
	}

	detectors := []struct {
		cmd       string
		args      []string
		pruneFunc func(string) string
	}{
		{
			cmd:  "systemd-detect-virt",
			args: []string{"--vm"}, // systemd-detect-virt can check container
			pruneFunc: nil
		},
		{
			cmd:       "virt-what",
			args:      nil,
			pruneFunc: func(result string) string {
				if result == "" {
					result = "none"
				}

				return result
			},
		},
	}

	for i := range detectors {
		data, err := RunCMD(detectors[i].cmd, detectors[i].args)
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

func RunCMD(cpath string, args []string) ([]byte, error) {
	if tmp, _ := exec.LookPath(cpath); tmp != "" {
		cpath = tmp
	} else {
		return nil, fmt.Errorf("not found cmd : %s", cpath)
	}

	var cmd *exec.Cmd
	if len(args) > 0 {
		cmd = exec.Command(cpath, args...)
	} else {
		cmd = exec.Command(cpath)
	}

	data, err := cmd.CombinedOutput()
	if err != nil {
		if len(data) > 0 {
			fmt.Printf("run %s err: %s\n", cpath, string(data))
		}

		return nil, err
	}

	return bytes.ToLower(bytes.TrimSpace(data)), nil
}
