package geo

// [中国的经纬度范围](https://www.cnblogs.com/inteliot/archive/2012/09/14/2684471.html)
func IsValidPointInChina(p Point) bool {
	if p.Lng < 73.66 || p.Lng > 135.05 {
		return false
	}
	if p.Lat < 3.86 || p.Lat > 53.55 {
		return false
	}

	return true
}
