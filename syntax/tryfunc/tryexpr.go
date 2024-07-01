//go:build !sugar_production

package tryfunc

func Try_(error) {
	panic("Try_() is unsupported at runtime")
}

func Try[A any](A, error) A {
	panic("Try() is unsupported at runtime")
}

func Try2[A, B any](A, B, error) (A, B) {
	panic("Try2() is unsupported at runtime")
}

func Try3[A, B, C any](A, B, C, error) (A, B, C) {
	panic("Try3() is unsupported at runtime")
}
