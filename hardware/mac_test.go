package hardware

import (
	"testing"
)

func TestReaMACs(t *testing.T) {
	ls, err := RealMACs()
	if err != nil {
		t.Fatal(err)
	}

	if len(ls) == 0 {
		t.Errorf("no real mac")
	} else {
		t.Log(ls)
	}
}
