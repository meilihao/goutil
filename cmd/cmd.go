package cmd

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"
)

var (
	regID = regexp.MustCompile(`(\d+\(\w+\))`)
)

// func main() {
// 	fmt.Println(UID("jr"))
// 	fmt.Println(UIDGID("jr", "adm"))
// }

// UID,GID不一定相等
// group同名的user可能不存在
// 查找gid: getent group ${group_name}
func UIDGID(user, group string) (string, string, error) {
	output, err := exec.Command("id", user).Output()
	if err != nil {
		return "", "", err
	}

	l := regID.FindAllString(string(output), -1)
	if len(l) < 3 {
		return "", "", err
	}

	ids := make([]string, 0, len(l))
	names := make([]string, 0, len(l))

	var i int
	for _, v := range l {
		if i = strings.Index(v, "("); i == -1 {
			continue
		}

		ids = append(ids, v[:i])
		names = append(names, strings.TrimSuffix(v[i+1:], ")"))
	}

	var uid, gid string
	for i, v := range names {
		if v == user {
			uid = ids[i]
			break
		}
	}
	for i, v := range names {
		if v == group {
			gid = ids[i]
			break
		}
	}

	if uid == "" {
		return "", "", errors.New(user + "不存在")
	}
	if gid == "" {
		return "", "", errors.New(group + "组不是" + user + "的支持组")
	}

	return uid, gid, nil
}

func UID(user string) (string, error) {
	if user == "" {
		return "", errors.New("用户名为空")
	}

	output, err := exec.Command("id", "-u", user).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
