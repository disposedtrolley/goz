package helpers

func Bits(val uint, start uint, num uint) uint {
	return (((1 << num) - 1) & (val >> (start - 1)))
}
