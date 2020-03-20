package hardware

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os/exec"
)

var (
	ErrNoMachineID = errors.New("no machine id")
)

func MachineID() (string, error) {
	file := "/etc/machine-id"

	if !FileIsExist(file) {
		// http://www.jinbuguo.com/systemd/machine-id.html
		// 没用dbus-uuidgen, 因为部分系统会报错"dbus-uuidgen: /lib/x86_64-linux-gnu/libdbus-1.so.3: version `LIBDBUS_PRIVATE_1.10.6' not found (required by dbus-uuidgen)"
		// 也可尝试systemd-machine-id-setup命令
		result, err := exec.Command("uuidgen").CombinedOutput()
		if err != nil {
			return "", err
		}

		ioutil.WriteFile(file, bytes.ReplaceAll(result, []byte{'-'}, nil), 0644)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return "", ErrNoMachineID
	}

	return string(data), err
}
