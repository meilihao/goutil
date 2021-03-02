// `lsscsi + nvme list -o json`
package hardware

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

const (
    SysFSRoot = "/sys"
    ClassNvme = "/class/nvme"
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
}

// func main() {
//     ns := ListNdevices()
//     fmt.Printf("found %d nvme devices.\n", len(ns))
//     for _, v := range ns {
//         fmt.Printf("%+v\n", *v)
//     }
// }

func ListNdevices() []*NvmeDevice {
    basePath := filepath.Join(SysFSRoot, ClassNvme)
    fs, _ := ioutil.ReadDir(basePath)
    if len(fs) == 0 {
        return []*NvmeDevice{}
    }

    nds := make([]*NvmeDevice, 0, len(fs))

    for _, v := range fs {
        if !v.IsDir() {
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
    fs, _ := ioutil.ReadDir(filepath.Join(SysFSRoot, ClassNvme, nc.Name))

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
    name = strings.TrimPrefix(name, nc.Name+"p")

    n, _ := strconv.Atoi(name)

    return n
}

func GetValue(filename string) string {
    data, _ := os.ReadFile(filename)

    return string(bytes.TrimSpace(data))
}

func GetNvmeSize(filename string) int {
    n, _ := strconv.Atoi(GetValue(filename))

    return n * 512 // block size is 512B
}
