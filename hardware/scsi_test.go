package hardware

import (
	"fmt"
	"testing"
)

func TestSCSI(t *testing.T) {
	ns := NDevices()
	fmt.Printf("found %d nvme devices.\n", len(ns))
	for _, v := range ns {
		t.Logf("%+v\n", *v)
	}

	ss := SDevices()
	fmt.Printf("found %d scsi devices.\n", len(ss))
	for _, v := range ss {
		fmt.Printf("%+v\n", *v)

		if v.Typename == ScsiShortDeviceTypes[13] {
			for _, vv := range v.SubDevices {
				if vv == nil {
					t.Logf("--->: %+v\n", vv)
				} else {
					t.Logf("--->: %+v\n", *vv)
				}
			}
		}
	}

	vs := VDevices()
	t.Logf("found %d virtio devices.\n", len(vs))
	for _, v := range vs {
		t.Logf("%v\n", *v)
	}
}
