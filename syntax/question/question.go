package question

type Question[T any] interface {
	Q() T
}

type QuestionImpl[T any] struct{}

func (QuestionImpl[T]) Q() T {
	panic("Question.Q() is unsupported at runtime")
}

func Try(error) {
	panic("Try() is unsupported at runtime")
}

func Try1[T any](T, error) T {
	panic("Try1() is unsupported at runtime")
}

func Try2[A, B any](A, B, error) (A, B) {
	panic("Try2() is unsupported at runtime")
}

func Try3[A, B, C any](A, B, C, error) (A, B, C) {
	panic("Try3() is unsupported at runtime")
}
