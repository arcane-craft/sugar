package tuple

type Pair[A, B any] func() (A, B)

func T2[A, B any](a A, b B) Pair[A, B] {
	return func() (A, B) {
		return a, b
	}
}

func First[A, B any](a A, _ B) A {
	return a
}

func Second[A, B any](_ A, b B) B {
	return b
}

type Tuple3[A, B, C any] func() (A, B, C)

func T3[A, B, C any](a A, b B, c C) Tuple3[A, B, C] {
	return func() (A, B, C) {
		return a, b, c
	}
}

func First3[A, B, C any](a A, _ B, _ C) A {
	return a
}

func Second3[A, B, C any](_ A, b B, _ C) B {
	return b
}

func Third[A, B, C any](_ A, _ B, c C) C {
	return c
}

type Tuple4[A, B, C, D any] func() (A, B, C, D)

func T4[A, B, C, D any](a A, b B, c C, d D) Tuple4[A, B, C, D] {
	return func() (A, B, C, D) {
		return a, b, c, d
	}
}

func First4[A, B, C, D any](a A, _ B, _ C, _ D) A {
	return a
}

func Second4[A, B, C, D any](_ A, b B, _ C, _ D) B {
	return b
}

func Third4[A, B, C, D any](_ A, _ B, c C, _ D) C {
	return c
}

func Fourth[A, B, C, D any](_ A, _ B, _ C, d D) D {
	return d
}

type Tuple5[A, B, C, D, E any] func() (A, B, C, D, E)

func T5[A, B, C, D, E any](a A, b B, c C, d D, e E) Tuple5[A, B, C, D, E] {
	return func() (A, B, C, D, E) {
		return a, b, c, d, e
	}
}

func First5[A, B, C, D, E any](a A, _ B, _ C, _ D, _ E) A {
	return a
}

func Second5[A, B, C, D, E any](_ A, b B, _ C, _ D, _ E) B {
	return b
}

func Third5[A, B, C, D, E any](_ A, _ B, c C, _ D, _ E) C {
	return c
}

func Fourth5[A, B, C, D, E any](_ A, _ B, _ C, d D, _ E) D {
	return d
}

func Fifth[A, B, C, D, E any](_ A, _ B, _ C, _ D, e E) E {
	return e
}
