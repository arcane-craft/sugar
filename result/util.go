package result

import (
	"fmt"

	"github.com/arcane-craft/sugar"
	"github.com/arcane-craft/sugar/tuple"
)

func From_(e error) Result[sugar.Unit] {
	if e != nil {
		return Err[sugar.Unit](e)
	}
	return Ok(sugar.Unit{})
}

func From[T any](v T, e error) Result[T] {
	if e != nil {
		return Err[T](e)
	}
	return Ok(v)
}

func From2[A, B any](r1 A, r2 B, e error) Result[tuple.Pair[A, B]] {
	if e != nil {
		return Err[tuple.Pair[A, B]](e)
	}
	return Ok(tuple.NewPair(r1, r2))
}

func From3[A, B, C any](r1 A, r2 B, r3 C, e error) Result[tuple.Triple[A, B, C]] {
	if e != nil {
		return Err[tuple.Triple[A, B, C]](e)
	}
	return Ok(tuple.NewTriple(r1, r2, r3))
}

func WrapErr[T any](desc string, r *Result[T]) {
	if r != nil && *r != nil {
		*r = (*r).MapErr(func(err error) error {
			return fmt.Errorf("%s %w", desc, err)
		})
	}
}
