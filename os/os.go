package os

// OSBit return 32/64
// - 对于32位系统:`^unit(0): 2^32 − 1, (2^32 − 1)>>63 = 0`, 得到32
// - 对于64位系统:`^unit(0): 2^64 − 1, (2^64 − 1)>>63 = 1`, 得到64
func OSBit() int {
	return 32 << (^uint(0) >> 63)
}
