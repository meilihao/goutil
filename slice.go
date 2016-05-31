package goutil

// InSlice checks given string in string slice or not.
func InSlice(v string, sl []string) bool {
	for _, value := range sl {
		if value == v {
			return true
		}
	}
	return false
}

func SliceUnique(sl []string) (uniqueslice []string) {
	for _, v := range sl {
		if !InSlice(v, uniqueslice) {
			uniqueslice = append(uniqueslice, v)
		}
	}
	return
}

// InSliceIOne checks given interface in interface slice.
func InSliceIOne(v interface{}, sl []interface{}) bool {
	for _, value := range sl {
		if value == v {
			return true
		}
	}
	return false
}

// SliceIAllUnique cleans repeated values in slice.
func SliceIAllUnique(sl []interface{}) (uniqueslice []interface{}) {
	for _, v := range sl {
		if !InSliceIOne(v, uniqueslice) {
			uniqueslice = append(uniqueslice, v)
		}
	}
	return
}
