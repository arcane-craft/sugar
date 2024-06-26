//go:build !sugar_production

package exception

type target struct{}

func Try(func()) handler {
	panic("Try() is unsupported at runtime")
}

type handler struct{}

func (handler) Catch(target, func(error)) handler {
	panic("Catch() is unsupported at runtime")
}

func Return(...any) {
	panic("Return() is unsupported at runtime")
}

func Throw(err error) {
	panic("Throw() is unsupported at runtime")
}

func Error(...error) target {
	return target{}
}

func Type[E error]() target {
	return target{}
}

func Type2[A, B error]() target {
	return target{}
}

func Type3[A, B, C error]() target {
	return target{}
}

func Type4[A, B, C, D error]() target {
	return target{}
}

func Type5[A, B, C, D, E error]() target {
	return target{}
}

func Type6[A, B, C, D, E, F error]() target {
	return target{}
}

func Type7[A, B, C, D, E, F, G error]() target {
	return target{}
}

func Type8[A, B, C, D, E, F, G, H error]() target {
	return target{}
}

func Type9[A, B, C, D, E, F, G, H, I error]() target {
	return target{}
}

func Type10[A, B, C, D, E, F, G, H, I, J error]() target {
	return target{}
}
