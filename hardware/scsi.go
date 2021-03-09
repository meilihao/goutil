// `lsscsi + nvme list -o json`
package hardware

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	SysFSRoot     = "/sys"
	ClassNvme     = "/sys/class/nvme"
	ClassBlock    = "/sys/class/block"
	BusScsiDevs   = "/sys/bus/scsi/devices"
	BusVirtioDevs = "/sys/bus/virtio/devices"
)

type NvmeController struct {
	Name        string
	Path        string
	Minor       int
	Model       string
	FirmwareRev string
	Serial      string // serial no
}

type NvmeDevice struct {
	Controller *NvmeController
	Name       string
	Path       string
	Namespace  int
	Size       int
	Disktype   string
}

// func main() {
// 	ns := NDevices()
// 	fmt.Printf("found %d nvme devices.\n", len(ns))
// 	for _, v := range ns {
// 		fmt.Printf("%+v\n", *v)
// 	}

// 	ss := SDevices()
// 	fmt.Printf("found %d scsi devices.\n", len(ss))
// 	for _, v := range ss {
// 		fmt.Printf("%+v\n", *v)

// 		if v.Typename == ScsiShortDeviceTypes[13] {
// 			for _, vv := range v.SubDevices {
//                if vv == nil {
//                    fmt.Printf("--->: %+v\n", vv)
//                } else {
//                    fmt.Printf("--->: %+v\n", *vv)
//                }
// 			}
// 		}
// 	}

// 	vs := VDevices()
// 	fmt.Printf("found %d virtio devices.\n", len(vs))
// 	for _, v := range vs {
// 		fmt.Printf("%v\n", *v)
// 	}
// }

// --- nvme from lsscsi
func NDevices() []*NvmeDevice {
	basePath := ClassNvme
	fs, _ := ioutil.ReadDir(basePath)
	if len(fs) == 0 {
		return []*NvmeDevice{}
	}

	nds := make([]*NvmeDevice, 0, len(fs))

	for _, v := range fs {
		if v.Mode()&os.ModeSymlink == 0 {
			continue
		}

		// example: nvme0
		nc := &NvmeController{
			Name:        v.Name(),
			Path:        filepath.Join(basePath, v.Name()),
			Minor:       NvmeControllerMinor(v.Name()),
			Model:       GetValue(filepath.Join(basePath, v.Name(), "model")),
			FirmwareRev: GetValue(filepath.Join(basePath, v.Name(), "firmware_rev")),
			Serial:      GetValue(filepath.Join(basePath, v.Name(), "serial")),
		}

		if ns := ListNNamespace(nc); len(ns) > 0 {
			nds = append(nds, ns...)
		}
	}

	return nds
}

// range nvme controller's namespaces
func ListNNamespace(nc *NvmeController) []*NvmeDevice {
	fs, _ := ioutil.ReadDir(filepath.Join(ClassNvme, nc.Name))

	nds := make([]*NvmeDevice, 0, len(fs))

	for _, v := range fs {
		if !v.IsDir() || !strings.HasPrefix(v.Name(), nc.Name) {
			continue
		}

		// example: nvme0n1
		nd := &NvmeDevice{
			Controller: nc,
			Name:       v.Name(),
			Path:       filepath.Join(nc.Path, v.Name()),
			Namespace:  NvmeDeviceNamespace(v.Name(), nc),
			Size:       GetNvmeSize(filepath.Join(nc.Path, v.Name(), "size")), // GB = Size /1000/1000/1000
			Disktype:   GetDiskType(v.Name()),
		}

		nds = append(nds, nd)
	}

	return nds
}

func NvmeControllerMinor(name string) int {
	name = strings.TrimPrefix(name, "nvme")

	n, _ := strconv.Atoi(name)

	return n
}

func NvmeDeviceNamespace(name string, nc *NvmeController) int {
	name = strings.TrimPrefix(name, nc.Name+"n")

	n, _ := strconv.Atoi(name)

	return n
}

func GetValue(filename string) string {
	data, _ := ioutil.ReadFile(filename)

	return string(bytes.TrimSpace(data))
}

func GetNvmeSize(filename string) int {
	n, _ := strconv.Atoi(GetValue(filename))

	return n * 512 // block size is 512B
}

// --- scsi from lsscsi
var (
	// from lsscsi.c scsi_short_device_types
	// [SES(SCSI Enclosure Services) device's Device Type is 13](https://www.systutorials.com/docs/linux/man/8-sg_ses/)
	ScsiShortDeviceTypes = []string{
		"disk",
		"tape",
		"printer",
		"process",
		"worm",
		"cd/dvd",
		"scanner",
		"optical",
		"mediumx",
		"comms",
		"(0xa)",
		"(0xb)",
		"storage",
		"enclosu",
		"sim dsk",
		"opti rd",
		"bridge",
		"osd",
		"adi",
		"sec man",
		"zbc",
		"(0x15)",
		"(0x16)",
		"(0x17)",
		"(0x18)",
		"(0x19)",
		"(0x1a)",
		"(0x1b)",
		"(0x1c)",
		"(0x1e)",
		"wlun",
		"nodev",
	}
)

// /sys/class/enclosure/0:0:23:0/id, 唯一标识enclosure
type ScsiDevice struct {
	Addr         string
	Type         string
	Typename     string
	Vendor       string
	Model        string
	Rev          string
	Sg           string // sg0
	Name         string // sda
	Serial       string
	Disktype     string
	Size         int
	SubDevices   []*ScsiDevice
	Parent       *ScsiDevice
	Slot         int // default(-1) not in enclosure
	Subenclosure int // default(-1) is not enclosure; 0, primary enclosure
}

func SDevices() []*ScsiDevice {
	basePath := BusScsiDevs
	fs, _ := ioutil.ReadDir(basePath)
	if len(fs) == 0 {
		return []*ScsiDevice{}
	}

	sds := make([]*ScsiDevice, 0, len(fs))

	filterFn := func(name string) bool {
		if strings.HasPrefix(name, "host") { // scsi host
			return false
		}

		if strings.HasPrefix(name, "target") { // scsi target
			return false
		}

		return true
	}

	sgMap := make(map[string]*ScsiDevice, len(fs))
	for _, v := range fs {
		if v.Mode()&os.ModeSymlink == 0 {
			continue
		}

		if !filterFn(v.Name()) {
			continue
		}

		// example: 0:0:1:0, [scsi addr](https://www.tldp.org/HOWTO/SCSI-2.4-HOWTO/scsiaddr.html)
		sd := &ScsiDevice{
			Slot:         -1,
			Subenclosure: -1,
			Addr:         v.Name(),
			Type:         GetValue(filepath.Join(basePath, v.Name(), "type")),
			Vendor:       GetValue(filepath.Join(basePath, v.Name(), "vendor")),
			Model:        GetValue(filepath.Join(basePath, v.Name(), "model")),
			Rev:          GetValue(filepath.Join(basePath, v.Name(), "rev")),
		}
		sd.Typename = ScsiTypename(sd.Type)
		sd.Sg = ScsiSg(filepath.Join(basePath, v.Name()))

		if sd.Typename == ScsiShortDeviceTypes[0] {
			sd.Name = BlockName(filepath.Join(basePath, v.Name()))
			sd.Size = GetScsiSize(filepath.Join(basePath, v.Name(), "block", sd.Name)) // GB = Size /1000/1000/1000
			sd.Disktype = GetDiskType(sd.Name)
			sd.Serial = GetSDiskSerial(sd.Name)

			sgMap[sd.Sg] = sd
		}

		sds = append(sds, sd)
	}

	// 先处理disk后enclosure, 避免解析enclosure solt时sds中不存在该disk
	for _, sd := range sds {
		if sd.Typename != ScsiShortDeviceTypes[13] {
			continue
		}

		sd.SubDevices = ParseEnclosure(sd, sgMap, filepath.Join(basePath, sd.Addr, "enclosure", sd.Addr)) // = /sys/class/enclosure/ + v.Name()
		sd.Subenclosure = GetSubenclosure(sd.Sg)                                                          // 0, primary enclosure
	}

	return sds
}

func GetDiskType(name string) string {
	subsystem, _ := os.Readlink(filepath.Join(ClassBlock, name, "device", "subsystem"))
	if strings.Index(subsystem, "virtio") != -1 {
		return "virtio/virtio"
	}
	if strings.Index(subsystem, "nvme") != -1 {
		return "ssd/nvme"
	}

	isSSD := false
	rotational := GetValue(filepath.Join(ClassBlock, name, "queue", "rotational"))
	if rotational == "0" {
		isSSD = true
	}

	isSAS := false
	sas, _ := os.Stat(filepath.Join(ClassBlock, name, "device", "sas_address"))
	if sas != nil {
		isSAS = true
	}

	if isSSD {
		if isSAS {
			return "ssd/sas"
		} else {
			return "ssd/sata"
		}
	} else {
		if isSAS {
			return "hdd/sas"
		} else {
			return "hdd/sata"
		}
	}
}

var (
	subenclosureReg        = regexp.MustCompile(`subenclosure id: (\d) \[`)
	scsiDiskSerialShortReg = regexp.MustCompile(`ID_SERIAL_SHORT=(\w+)`)
	scsiDiskSerialReg      = regexp.MustCompile(`ID_SERIAL=(\w+)`)
)

func GetSubenclosure(sg string) int {
	out, err := exec.Command("sg_ses", "-p", "0xa", "/dev/"+sg).CombinedOutput()
	if err != nil {
		return -1
	}

	ls := subenclosureReg.FindStringSubmatch(string(out))
	if len(ls) >= 2 {
		n, _ := strconv.Atoi(ls[1])

		return n
	}

	return -1
}

func GetSDiskSerial(name string) string {
	out, err := exec.Command("udevadm", "info", "--query=property", "--name=/dev/"+name).CombinedOutput() // --query=property
	if err != nil {
		return ""
	}

	ls := scsiDiskSerialReg.FindStringSubmatch(string(out))
	if len(ls) >= 2 {
		return strings.TrimSpace(ls[1])
	}

	ls = scsiDiskSerialShortReg.FindStringSubmatch(string(out))
	if len(ls) >= 2 {
		return strings.TrimSpace(ls[1])
	}

	return ""
}

// 当slot没有插盘时: /sys/class/enclosure/0:0:22:0/SlotX下不存在名为`device`的symlink
func ParseEnclosure(parent *ScsiDevice, m map[string]*ScsiDevice, base string) []*ScsiDevice {
	count := EnclosureComponentsCount(base)
	sds := make([]*ScsiDevice, 0, count)

	fs, _ := ioutil.ReadDir(base)
	if len(fs) == 0 {
		return []*ScsiDevice{}
	}

	filter := func(name string) bool {
		if strings.HasPrefix(name, "Slot") {
			return true
		}

		return false
	}

	var solt int
	var sg string
	var tmpSd *ScsiDevice
	var ok bool

	for _, v := range fs {
		if !v.IsDir() {
			continue
		}

		if !filter(v.Name()) {
			continue
		}

		solt, _ = strconv.Atoi(GetValue(filepath.Join(base, v.Name(), "slot"))) // strconv.Atoi(strings.TrimPrefix(v.Name(), "Slot"))
		sg = ScsiSg(filepath.Join(base, v.Name(), "device")) // no device in solt, so no `device` symlink
		if tmpSd, ok = m[sg]; ok {
			tmpSd.Slot = solt
			tmpSd.Parent = parent

			sds = append(sds, tmpSd)
		} else {
			sds = append(sds, nil)
		}
	}

	return sds
}

func EnclosureComponentsCount(base string) int {
	n, _ := strconv.Atoi(GetValue(filepath.Join(base, "components")))

	return n
}

func ScsiTypename(index string) string {
	i, err := strconv.Atoi(index)
	if err == nil && i >= 0 && i < len(ScsiShortDeviceTypes) {
		return ScsiShortDeviceTypes[i]
	}

	return "unknown"
}

// == genericDev := Getvalue("/sys/bus/scsi/devices/1:0:0:0/generic/dev")) // for disk
// == genericDev := Getvalue("/sys/bus/scsi/devices/1:0:0:0/enclosure/1:0:0:0/generic/dev")) // for enclosure
// minor := strings.Split(genericDev, ":")[1]
// sgDev := "/dev/sg" + minor
func ScsiSg(base string) string {
	p, err := os.Readlink(filepath.Join(base, "generic"))
	if err != nil {
		return ""
	}

	return filepath.Base(p)
}

func BlockName(base string) string {
	fs, _ := ioutil.ReadDir(filepath.Join(base, "block"))
	for _, v := range fs {
		if v.IsDir() {
			return v.Name()
		}
	}

	return ""
}

func GetScsiSize(base string) int {
	n, _ := strconv.Atoi(GetValue(filepath.Join(base, "size")))

	return n * 512 // block size is 512B
}

// --- virtio
type VirtioDevice struct {
	Addr     string
	Vendor   string
	Model    string
	Rev      string
	Name     string // sda
	Serial   string // virtio_blk serial由qemu -device serial=xxx属性提供, [qemu scsi-block devices不支持-device serial属性](https://libvirt.org/formatdomain.html).
	Disktype string
	Size     int
}

// for aliyun ecs
func VDevices() []*VirtioDevice {
	basePath := BusVirtioDevs
	fs, _ := ioutil.ReadDir(basePath)
	if len(fs) == 0 {
		return []*VirtioDevice{}
	}

	vds := make([]*VirtioDevice, 0, len(fs))

	filterFn := func(name string) bool {
		s, _ := os.Stat(filepath.Join(basePath, name, "block"))
		if s != nil {
			return true
		}

		return false
	}

	for _, v := range fs {
		if v.Mode()&os.ModeSymlink == 0 {
			continue
		}

		if !filterFn(v.Name()) {
			continue
		}

		// example: 0:0:1:0, [scsi addr](https://www.tldp.org/HOWTO/SCSI-2.4-HOWTO/scsiaddr.html)
		vd := &VirtioDevice{
			Name:   BlockName(filepath.Join(basePath, v.Name())),
			Addr:   v.Name(),
			Vendor: GetValue(filepath.Join(basePath, v.Name(), "vendor")),
			Model:  GetValue(filepath.Join(basePath, v.Name(), "model")),
			Rev:    GetValue(filepath.Join(basePath, v.Name(), "rev")),
		}

		vd.Size = GetScsiSize(filepath.Join(basePath, v.Name(), "block", vd.Name)) // GB = Size /1024/1024/1024
		vd.Disktype = GetDiskType(vd.Name)
		vd.Serial = GetSDiskSerial(vd.Name)

		vds = append(vds, vd)
	}

	return vds
}
