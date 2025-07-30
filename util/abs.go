package util

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func Abs[T Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
