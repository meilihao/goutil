package time

import (
	"fmt"
	"log"
	"testing"
)

func TestTimeXBTrim10(t *testing.T) {
	x, err := NewTimeXB(32, "hxYYu2Gka5BGqym17C12id8rXOPf1vVI", 2, 10)
	if err != nil {
		log.Fatal(err)
	}

	src := "ABC"

	out := x.Generate(src)
	fmt.Println(out)
	tmp, now, err := x.Parse(out)
	if err != nil {
		t.Errorf("TimeXB.Parse failed: %v", err)
	}
	if tmp != src {
		t.Errorf("TimeXB.Parse failed: not equal src(%s != %s)", src, tmp)
	}
	fmt.Println(tmp, now)
}

func TestTimeXBNoTrim(t *testing.T) {
	x, err := NewTimeXB(32, "hxYYu2Gka5BGqym17C12id8rXOPf1vVI", 2, 0)
	if err != nil {
		log.Fatal(err)
	}

	src := "ABC"

	out := x.Generate(src)
	fmt.Println(out)
	tmp, now, err := x.Parse(out)
	if err != nil {
		t.Errorf("TimeXB.Parse failed: %v", err)
	}
	if tmp != src {
		t.Errorf("TimeXB.Parse failed: not equal src(%s != %s)", src, tmp)
	}
	fmt.Println(tmp, now)
}
