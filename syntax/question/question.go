//go:build !sugar_production

package question

type Question[T any] interface {
	Q() T
}

type QuestionImpl[T any] struct{}

func (QuestionImpl[T]) Q() T {
	panic("Question.Q() is unsupported at runtime")
}
