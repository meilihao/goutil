package hardware

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	ClassEnclosure = "/sys/class/enclosure"
)

type Enclosure struct {
	Id         string
	ShelfNo    int
	Vendor     string
	Model      string
	Components int
	SgName     string
}

func (e *Enclosure) Name() string {
	return e.Vendor + "_" + e.Model
}

func Enclosures() ([]*Enclosure, error) {
	if !FileIsDir(ClassEnclosure) {
		return nil, nil
	}

	entries, err := ioutil.ReadDir(ClassEnclosure)
	if err != nil {
		return nil, err
	}

	enclosureList := make([]*Enclosure, 0, len(entries))
	mDup := make(map[string]bool, len(entries))
	var tmpId string

	entries = SortSCSIAddress(entries, false)

	for _, v := range entries {
		tpy := GetValue(filepath.Join(ClassEnclosure, v.Name(), "/device/type"))
		if tpy != "13" { // [SES(SCSI Enclosure Services) device's Device Type is 13](https://www.systutorials.com/docs/linux/man/8-sg_ses/)
			continue
		}

		tmpId = GetValue(filepath.Join(ClassEnclosure, v.Name(), "/id"))
		if mDup[tmpId] {
			continue
		}

		e := &Enclosure{}

		e.Vendor = GetValue(filepath.Join(ClassEnclosure, v.Name(), "/device/vendor"))
		e.Model = GetValue(filepath.Join(ClassEnclosure, v.Name(), "/device/model"))
		e.Id = tmpId
		mDup[e.Id] = true

		e.Components = GetValueInt(filepath.Join(ClassEnclosure, v.Name(), "components"))
		e.SgName = ScsiSg(filepath.Join(ClassEnclosure, v.Name(), "device"))
		e.ShelfNo = GetSubenclosure(e.SgName)

		enclosureList = append(enclosureList, e)
	}

	return enclosureList, nil
}

func SortSCSIAddress(list []os.FileInfo, withPath bool) []os.FileInfo {
	keys := make([]int, len(list))
	m := make(map[int]os.FileInfo, 4)
	for i := range list {
		if withPath {
			keys[i] = weightSCSIAddress(filepath.Base(list[i].Name()))
		} else {
			keys[i] = weightSCSIAddress(list[i].Name())
		}
		m[keys[i]] = list[i]
	}

	sort.Ints(keys)

	result := make([]os.FileInfo, 0, len(list))
	for _, v := range keys {
		result = append(result, m[v])
	}

	return result
}

// // [scsi addr](https://www.tldp.org/HOWTO/SCSI-2.4-HOWTO/scsiaddr.html)
// - SCSI adapter number [host]
// - channel number [bus]
// - id number [target]
// - lun [lun]
// example : 1:0:0:0
func weightSCSIAddress(addr string) int {
	ss := strings.Split(addr, ":")

	host, _ := strconv.Atoi(ss[0])
	bus, _ := strconv.Atoi(ss[1])
	target, _ := strconv.Atoi(ss[2])

	return host<<32 + bus<<16 + target<<8 + target
}
