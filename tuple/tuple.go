package tuple

type Pair[A, B any] struct {
	a A
	b B
}

func NewPair[A, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{a, b}
}

func (p Pair[A, B]) First() A {
	return p.a
}

func (p Pair[A, B]) Second() B {
	return p.b
}

type Triple[A, B, C any] struct {
	a A
	b B
	c C
}

func NewTriple[A, B, C any](a A, b B, c C) Triple[A, B, C] {
	return Triple[A, B, C]{a, b, c}
}

func (t Triple[A, B, C]) First() A {
	return t.a
}

func (t Triple[A, B, C]) Second() B {
	return t.b
}

func (t Triple[A, B, C]) Third() C {
	return t.c
}
