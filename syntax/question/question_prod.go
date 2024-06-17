//go:build sugar_production

package question

type Question[T any] interface{}

type QuestionImpl[T any] struct{}
