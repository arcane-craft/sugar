//go:build !sugar_production

package exception

type ExceptHandler func(error) error

func Try[T any](func() T, ...ExceptHandler) T {
	panic("")
}

func Catch(error, func(error) error) ExceptHandler {
	panic("")
}

func CatchT[T any](func(error) error) ExceptHandler {
	panic("")
}
