package time

import (
	"crypto/sha1"
	"fmt"
	"testing"
)

func TestTimeXTrim10(t *testing.T) {
	x := NewTimeX(sha1.New, "123", 8)
	src := "ABC"

	out := x.Generate(src)
	fmt.Println(out)
	tmp, now, err := x.Parse(out)
	if err != nil {
		t.Errorf("TimeX.Parse failed: %v", err)
	}
	if tmp != src {
		t.Errorf("TimeX.Parse failed: not equal src(%s != %s)", src, tmp)
	}
	fmt.Println(tmp, now)
}


func TestTimeXNoTrim(t *testing.T) {
	x := NewTimeX(sha1.New, "123", 2)
	src := "ABC"

	out := x.Generate(src)
	fmt.Println(out)
	tmp, now, err := x.Parse(out)
	if err != nil {
		t.Errorf("TimeX.Parse failed: %v", err)
	}
	if tmp != src {
		t.Errorf("TimeX.Parse failed: not equal src(%s != %s)", src, tmp)
	}
	fmt.Println(tmp, now)
}
