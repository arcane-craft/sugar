package option

import "github.com/arcane-craft/sugar"

type Option[T any] interface {
	IsSome() bool
	IsNone() bool
	Expect(msg string) T
	Unwrap() T
	UnwrapOr(def T) T

	sugar.Question[T]
}

type oSome[T any] struct {
	v T

	sugar.QuestionImpl[T]
}

func (oSome[T]) IsSome() bool {
	return true
}

func (oSome[T]) IsNone() bool {
	return false
}

func (o oSome[T]) Expect(_ string) T {
	return o.v
}

func (o oSome[T]) Unwrap() T {
	return o.v
}

func (o oSome[T]) UnwrapOr(_ T) T {
	return o.v
}

type oNone[T any] struct {
	sugar.QuestionImpl[T]
}

func (oNone[T]) IsSome() bool {
	return false
}

func (oNone[T]) IsNone() bool {
	return true
}

func (oNone[T]) Expect(msg string) T {
	panic(msg)
}

func (oNone[T]) Unwrap() T {
	panic("unwrap None() is not allowed")
}

func (oNone[T]) UnwrapOr(def T) T {
	return def
}

func Some[T any](v T) Option[T] {
	return oSome[T]{v: v}
}

func None[T any]() Option[T] {
	return oNone[T]{}
}
