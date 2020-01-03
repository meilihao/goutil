package hardware

import (
	"testing"
)

func TestReaMACs(t *testing.T) {
	ls := RealMACs()
	if len(ls) == 0 {
		t.Errorf("no real mac")
	} else {
		t.Log(ls)
	}
}
