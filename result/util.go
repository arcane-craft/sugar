package result

import (
	"github.com/arcane-craft/sugar"
	"github.com/arcane-craft/sugar/tuple"
)

func FromErr(e error) Result[sugar.Unit] {
	if e != nil {
		return Err[sugar.Unit](e)
	}
	return Ok(sugar.Unit{})
}

func FromR1Err[T any](v T, e error) Result[T] {
	if e != nil {
		return Err[T](e)
	}
	return Ok(v)
}

func FromR2Err[A, B any](r1 A, r2 B, e error) Result[tuple.Pair[A, B]] {
	if e != nil {
		return Err[tuple.Pair[A, B]](e)
	}
	return Ok(tuple.T2(r1, r2))
}

func FromR3Err[A, B, C any](r1 A, r2 B, r3 C, e error) Result[tuple.Tuple3[A, B, C]] {
	if e != nil {
		return Err[tuple.Tuple3[A, B, C]](e)
	}
	return Ok(tuple.T3(r1, r2, r3))
}
