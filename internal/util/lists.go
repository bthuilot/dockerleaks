package util

func Reverse[X any](xs []X) []X {
	i, j := 0, len(xs)-1
	for i < j {
		xs[i], xs[j] = xs[j], xs[i]
		i++
		j--
	}
	return xs
}

func ZipApply[A, B, C any](f func(a A, b B) C, as []A, bs []B) []C {
	cs := make([]C, MinInt(len(as), len(bs)))
	for i := 0; i < len(cs); i++ {
		cs[i] = f(as[i], bs[i])
	}
	return cs
}
