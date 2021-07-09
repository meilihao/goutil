package targetcli

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	// This format is currently dictated by the iSCSI target backend,
	// specifically the rtslib-fb library.
	// A notable difference in this implementation (which also differs from
	// RFC3720, where the IQN format is defined) is that we require the
	// "unique" part after the colon to be present.
	//
	// See also the source code of rtslib-fb for the original regex:
	// https://github.com/open-iscsi/rtslib-fb/blob/master/rtslib/utils.py#L388

	// iqn.yyyy-mm.naming-authority:unique name
	regexWWN = regexp.MustCompile(`^iqn\.(\d{4}-\d{2})\.([^:]+)(:)([^,:\s']+)$`)
)

const (
	BackstoreTpyBlock  = "block"
	BackstoreTpyFileio = "fileio"
)

// # targetcli /backstores/fileio create test2 /iscsi_test2
// Created fileio test2 with size 33554432
func AddBackstoresObjectCmd(typ, name, filePath string) string {
	otps := []string{
		targetcliBinary,
		fmt.Sprintf("/backstores/%s", typ),
		"create",
		name,
		filePath,
	}

	return strings.Join(otps, " ")
}

// # targetcli /backstores/fileio delete test2
// Deleted storage object test2.
func DeleteBackstoresObjectCmd(typ, name string) string {
	otps := []string{
		targetcliBinary,
		fmt.Sprintf("/backstores/%s", typ),
		"delete",
		name,
	}

	return strings.Join(otps, " ")
}

// # targetcli ls /backstores/fileio/test2
// No such path /backstores/fileio/test2
func IsExistBackstoresObject(backstoreTpy, objName string) string {
	otps := []string{
		targetcliBinary,
		"ls",
		"/backstores/" + backstoreTpy + "/" + objName,
	}

	return strings.Join(otps, " ")
}

func CheckIQN(iqn string) error {
	if strings.ContainsAny(iqn, "_ ") {
		return errors.New("IQN cannot contain the characters '_' (underscore) or ' ' (space)")
	}

	if !regexWWN.MatchString(iqn) {
		return fmt.Errorf("Given IQN ('%s') does not match the regular expression '%s'", iqn, regexWWN.String())
	}

	return nil
}

// CreateIscsiTarget will create a iSCSI target using the name specified. If name is
// unspecified, a name will be generated. Notice the name must comply with iSCSI
// name format.
//
// # targetcli /iscsi create
// Created target iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.64cc17ed0de5.
// Created TPG 1.
// Global pref auto_add_default_portal=true
// Created default portal listening on all IPs (0.0.0.0), port 3260.
func CreateIscsiTargetCmd(name string) string {
	if name != "" {
		if err := CheckIQN(name); err != nil {
			panic(err)
		}
	}

	otps := []string{
		targetcliBinary,
		"/iscsi",
		fmt.Sprintf("create %s", name),
	}

	return strings.Join(otps, " ")
}

func ParseCreateIscsiTargetCmdResult(output string) (targetName string, tpgId int, err error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "Created target ") {
			targetName = CleanResultString("Created target ", scanner.Text())
		}
		if strings.HasPrefix(scanner.Text(), "Created TPG ") {
			tpgId, err = strconv.Atoi(CleanResultString("Created TPG ", scanner.Text()))
			if err != nil {
				return
			}
		}
	}

	return
}

// # targetcli ls /iscsi/iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.ca3c7dfe1233
// No such path /iscsi/iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.ca3c7dfe1233
func IsExistIscsiTarget(targetName string) string {
	otps := []string{
		targetcliBinary,
		"ls",
		"/iscsi/" + targetName,
	}

	return strings.Join(otps, " ")
}

// DeleteIscsiTargetCmd will remove a iSCSI target specified by target name
//
// # targetcli /iscsi delete iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.64cc17ed0de5
// Deleted Target iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.64cc17ed0de5.
func DeleteIscsiTargetCmd(targetName string) string {
	otps := []string{
		targetcliBinary,
		"/iscsi",
		"delete",
		targetName,
	}

	return strings.Join(otps, " ")
}

// AddLun will add a LUN in an existing target, which backing by
// specified file
//
// # targetcli /iscsi/iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.64cc17ed0de5/tpg1/luns create /backstores/fileio/test2 [lun]
// Created LUN 0.
func AddIscsiLunCmd(targetName string, tpgId int64, lunId int64, backstoreTpy, objName string) string {
	otps := []string{
		targetcliBinary,
		fmt.Sprintf("/iscsi/%s/tpg%d/luns", targetName, tpgId),
		"create",
		fmt.Sprintf("/backstores/%s/%s", backstoreTpy, objName),
	}
	if lunId >= 0 {
		otps = append(otps, fmt.Sprintf("%d", lunId))
	}

	return strings.Join(otps, " ")
}

func ParseAddIscsiLunCmdResult(output string) (lunId int64, err error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "Created LUN ") {
			lunId, err = strconv.ParseInt(CleanResultString("Created LUN ", scanner.Text()), 10, 64)
			if err != nil {
				return
			}
		}
	}

	return
}

// DeleteLun will remove a LUN from an target
//
// # targetcli /iscsi/iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.64cc17ed0de5/tpg1/luns delete 0
func DeleteIscsiLunCmd(targetName string, tpgId int64, lunId int64) string {
	otps := []string{
		targetcliBinary,
		fmt.Sprintf("/iscsi/%s/tpg%d/luns", targetName, tpgId),
		"delete",
		fmt.Sprintf("%d", lunId),
	}

	return strings.Join(otps, " ")
}

// BindInitiator will add permission to allow certain initiator(s) to connect to
// certain target
//
// # targetcli /iscsi/iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.0d690d398ec5/tpg1/acls create iqn.1993-08.org.debian:01:7ed7bee79b99
// Created mapped LUN 0.
func BindInitiatorCmd(targetName string, tpgId int64, initiator string) string {
	otps := []string{
		targetcliBinary,
		fmt.Sprintf("/iscsi/%s/tpg%d/acls", targetName, tpgId),
		"create",
		initiator,
	}

	return strings.Join(otps, " ")
}

// UnbindInitiator will remove permission to allow certain initiator(s) to connect to
// certain target.
// # targetcli /iscsi/iqn.2003-01.org.linux-iscsi.fyhdesktop29.x8664:sn.0d690d398ec5/tpg1/acls delete iqn.1993-08.org.debian:01:7ed7bee79b99
// Deleted Node ACL iqn.1993-08.org.debian:01:7ed7bee79b99.
func UnbindInitiatorCmd(targetName string, tpgId int64, initiator string) string {
	otps := []string{
		targetcliBinary,
		fmt.Sprintf("/iscsi/%s/tpg%d/acls", targetName, tpgId),
		"delete",
		initiator,
	}

	return strings.Join(otps, " ")
}

// status is 0 or 1
func SwitchIscsiTargetChapStatusCmd(targetName string, tpgId int64, status int32) string {
	otps := []string{
		targetcliBinary,
		fmt.Sprintf("/iscsi/%s/tpg%d", targetName, tpgId),
		"set attribute",
		fmt.Sprintf("authentication=%d", status),
	}

	return strings.Join(otps, " ")
}

func SetIscsiTargetChapCmd(targetName string, tpgId int64, initiator, user, password string) string {
	otps := []string{
		targetcliBinary,
		fmt.Sprintf("/iscsi/%s/tpg%d/acls/%s", targetName, tpgId, initiator),
		"set auth",
		fmt.Sprintf("userid=%s password=%s", user, password),
	}

	return strings.Join(otps, " ")
}

// cancel chap (ep: `set auth userid=`), must one by one, not allow `set auth userid= password=`, it will cause auth info confusion.
func CleanIscsiTargetChapUseridCmd(targetName string, tpgId int64, initiator string) string {
	otps := []string{
		targetcliBinary,
		fmt.Sprintf("/iscsi/%s/tpg%d/acls/%s", targetName, tpgId, initiator),
		"set auth",
		fmt.Sprintf("userid=%s", ""),
	}

	return strings.Join(otps, " ")
}

// cancel chap (ep: `set auth password=`), must one by one, not allow `set auth userid= password=`, it will cause auth info confusion.
func CleanIscsiTargetChapPasswordCmd(targetName string, tpgId int64, initiator string) string {
	otps := []string{
		targetcliBinary,
		fmt.Sprintf("/iscsi/%s/tpg%d/acls/%s", targetName, tpgId, initiator),
		"set auth",
		fmt.Sprintf("password=%s", ""),
	}

	return strings.Join(otps, " ")
}
