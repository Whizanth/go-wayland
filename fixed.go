package wayland

import "math"

type Fixed uint32

func (f Fixed) Float64() float64 {
	return float64(int32(f)) / 256.0
}

func (f Fixed) Int32() int32 {
	return int32(f) / 256
}

func (f Fixed) Int() int {
	return int(f.Int32())
}

func ParseFixed[T int | int32 | float64](val T) Fixed {
	return Fixed(uint32(int32(math.Round(float64(val) * 256.0))))
}
