package time

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeSince(t *testing.T) {
	cases := []struct {
		then time.Time
		want string
	}{
		{time.Now().Add(1 * time.Second), ""},
		{time.Now().Add(-1 * time.Second), ""},
		{time.Now().Add(-61 * time.Second), ""},
		{time.Now().Add(-1 * time.Hour * 23).Add(-2 * time.Second), ""},
		{time.Now().AddDate(0, 0, -1), ""},
		{time.Now().AddDate(0, -1, -1), ""},
		{time.Now().AddDate(-1, -1, -1), ""},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("pager(%s)", c.then), func(t *testing.T) {
			d := TimeSince(c.then)
			fmt.Println(d)

			if d == "" {
				t.Errorf("TimeSince(%s) failed", c.then)
			}
		})
	}
}
