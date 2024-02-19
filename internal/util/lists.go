package util

// Reverse reverses a list contents
func Reverse[X any](xs []X) []X {
	i, j := 0, len(xs)-1
	for i < j {
		xs[i], xs[j] = xs[j], xs[i]
		i++
		j--
	}
	return xs
}

// ZipApply will apply to a function to the a[i], b[i] and return a list
// of the result, where a and b are too independent lists. If the lists are
// different sizes, it will only zip until the length of the shorter list
func ZipApply[A, B, C any](f func(a A, b B) C, as []A, bs []B) []C {
	cs := make([]C, MinInt(len(as), len(bs)))
	for i := 0; i < len(cs); i++ {
		cs[i] = f(as[i], bs[i])
	}
	return cs
}

func Apply[A, B any](as []A, f func(a A) B) []B {
	bs := make([]B, len(as))
	for i, a := range as {
		bs[i] = f(a)
	}
	return bs
}

func Any[X any](xs []X, f func(X) bool) bool {
	for _, x := range xs {
		if f(x) {
			return true
		}
	}
	return false
}
