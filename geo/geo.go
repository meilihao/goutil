package geo

import (
	"encoding/json"
	"fmt"
)

type Point struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

type Polygon struct {
	points []Point
	box
}

type box struct {
	LngMin float64
	LngMax float64

	LatMin float64
	LatMax float64
}

func NewPolygon(points []Point) *Polygon {
	p := &Polygon{points: points, box: box{}}
	p.SetBox()

	return p
}

func NewPolygonByString(s string) *Polygon {
	var a []Point

	if s == "" {
		fmt.Println("empty geo.")
		return nil
	}

	err := json.Unmarshal([]byte(s), &a)
	if err != nil || len(a) < 3 {
		fmt.Println("invalid geo : ", err)
		return nil
	}

	return NewPolygon(a)
}

func (p *Polygon) SetBox() {
	p.box.LngMin = p.points[0].Lng
	p.box.LngMax = p.points[0].Lng
	p.box.LatMin = p.points[0].Lat
	p.box.LatMax = p.points[0].Lat

	for _, v := range p.points {
		if v.Lng > p.box.LngMax {
			p.box.LngMax = v.Lng
		}
		if v.Lng < p.box.LngMin {
			p.box.LngMin = v.Lng
		}
		if v.Lat > p.box.LatMax {
			p.box.LatMax = v.Lat
		}
		if v.Lat < p.box.LatMin {
			p.box.LatMin = v.Lat
		}
	}
}

// http://alienryderflex.com/polygon/
// http://www.cnblogs.com/luxiaoxun/p/3722358.html
// http://blog.chinaunix.net/uid-30332431-id-5140349.html
// PNPoly算法
func (p *Polygon) IsInside(point Point) bool {
	if point.Lng > p.box.LngMax || point.Lng < p.box.LngMin ||
		point.Lat > p.box.LatMax || point.Lat < p.box.LatMin {
		return false
	}

	isIn := false
	n := len(p.points)

	for i, j := 0, n-1; i < n; i++ {
		if (p.points[i].Lat < point.Lat && p.points[j].Lat >= point.Lat ||
			p.points[j].Lat < point.Lat && p.points[i].Lat >= point.Lat) &&
			(p.points[i].Lng <= point.Lng || p.points[j].Lng <= point.Lng) { // 从待测点引出一条水平向左的射线
			// 有交点,必定应 : 交点.Lng<point.Lng
			if p.points[i].Lng+(point.Lat-p.points[i].Lat)/(p.points[j].Lat-p.points[i].Lat)*(p.points[j].Lng-p.points[i].Lng) < point.Lng {
				isIn = !isIn
			}
		}
		j = i
	}

	return isIn
}

func (p *Polygon) IsInsideLngLat(lng, lat float64) bool {
	return p.IsInside(Point{Lng: lng, Lat: lat})
}
