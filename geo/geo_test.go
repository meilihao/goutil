package geo

import (
	"testing"
)

func TestNewPolygon(t *testing.T) {
	points := []Point{Point{1, 1}}

	p := NewPolygon(points)
	if p == nil {
		t.Error("invalid polygon")
	}
}

func TestNewPolygonByString(t *testing.T) {
	var tests = []struct {
		s      string
		expect bool
	}{
		{"", false},
		{`[{"lng":0,"lat":0},{"lng":1,"lat":0}]`, false},
		{`[{"lng":0,"lat":0},{"lng":1,"lat":0},{"lng":1,"lat":1},{"lng":0,"lat":1}]`, true},
		{`[{"lng":-1,"lat":0},{"lng":1,"lat":1},{"lng":1,"lat":-1},{"lng":-1,"lat":-1}]`, true},
	}

	for _, tc := range tests {
		if p := NewPolygonByString(tc.s); (p != nil) != tc.expect {
			t.Errorf("create polygon by %s,get %t(%v), want %t", tc.s, p != nil, p, tc.expect)
		} else {
			t.Logf("%v\n", p)
		}
	}
}

func TestIsInside(t *testing.T) {
	var tests = []struct {
		s      string
		p      Point
		expect bool
	}{
		{`[{"lng":-1,"lat":1},{"lng":1,"lat":1},{"lng":1,"lat":-1},{"lng":-1,"lat":-1}]`, Point{0, 0}, true},
		//{`[{"lng":0,"lat":0},{"lng":1,"lat":0},{"lng":1,"lat":1},{"lng":0,"lat":1}]`, Point{0, 0}, false},      // 点与多边形顶点重合
		//{`[{"lng":-1,"lat":0},{"lng":1,"lat":1},{"lng":1,"lat":-1},{"lng":-1,"lat":-1}]`, Point{0, -1}, false}, // 点与多边形的边重合
	}

	for _, tc := range tests {
		if isIn := NewPolygonByString(tc.s).IsInside(tc.p); isIn != tc.expect {
			t.Errorf("p(%v) is inside in %s,get %t, want %t", tc.p, tc.s, isIn, tc.expect)
		} else {
			t.Logf("%v\n", isIn)
		}
	}
}
