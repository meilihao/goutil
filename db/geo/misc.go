package geo

// [中国的经纬度范围](https://www.cnblogs.com/inteliot/archive/2012/09/14/2684471.html)
func IsValidPointInChina(p Point) bool {
	if p[0] < 73.66 || p[0] > 135.05 {
		return false
	}
	if p[1] < 3.86 || p[1] > 53.55 {
		return false
	}

	return true
}
