package result

import (
	"github.com/arcane-craft/sugar/option"
	"github.com/arcane-craft/sugar/syntax/question"
)

type Result[T any] interface {
	IsOk() bool
	IsErr() bool
	Ok() option.Option[T]
	Err() option.Option[error]
	Expect(string) T
	Unwrap() T
	ExpectErr(string) error
	UnwrapErr() error
	UnwrapOr(T) T
	Map(func(T) T) Result[T]
	MapErr(func(error) error) Result[T]
	Mutate(func(*T, *error))

	question.Question[T]
}

type rOk[T any] struct {
	v T

	question.QuestionImpl[T]
}

func (rOk[T]) IsOk() bool {
	return true
}

func (rOk[T]) IsErr() bool {
	return false
}

func (r rOk[T]) Ok() option.Option[T] {
	return option.Some(r.v)
}

func (rOk[T]) Err() option.Option[error] {
	return option.None[error]()
}

func (r rOk[T]) Expect(_ string) T {
	return r.v
}

func (r rOk[T]) Unwrap() T {
	return r.v
}

func (rOk[T]) ExpectErr(msg string) error {
	panic(msg)
}

func (rOk[T]) UnwrapErr() error {
	panic("unwrap_err Ok() is not allowed!")
}

func (r rOk[T]) UnwrapOr(_ T) T {
	return r.v
}

func (r rOk[T]) Map(f func(T) T) Result[T] {
	return Ok(f(r.v))
}

func (r rOk[T]) MapErr(_ func(error) error) Result[T] {
	return &r
}

func (r *rOk[T]) Mutate(f func(*T, *error)) {
	f(&r.v, nil)
}

type rErr[T any] struct {
	v error

	question.QuestionImpl[T]
}

func (rErr[T]) IsOk() bool {
	return false
}

func (rErr[T]) IsErr() bool {
	return true
}

func (rErr[T]) Ok() option.Option[T] {
	return option.None[T]()
}

func (r rErr[T]) Err() option.Option[error] {
	return option.Some(r.v)
}

func (rErr[T]) Expect(msg string) T {
	panic(msg)
}

func (r rErr[T]) Unwrap() T {
	panic("unwrap Err() is not allowed!")
}

func (r rErr[T]) ExpectErr(_ string) error {
	return r.v
}

func (r rErr[T]) UnwrapErr() error {
	return r.v
}

func (rErr[T]) UnwrapOr(def T) T {
	return def
}

func (r rErr[T]) Map(_ func(T) T) Result[T] {
	return &r
}

func (r rErr[T]) MapErr(f func(error) error) Result[T] {
	return Err[T](f(r.v))
}

func (r *rErr[T]) Mutate(f func(*T, *error)) {
	f(nil, &r.v)
}

func Ok[T any](v T) Result[T] {
	return &rOk[T]{v: v}
}

func Err[T any](e error) Result[T] {
	return &rErr[T]{v: e}
}
