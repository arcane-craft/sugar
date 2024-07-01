package errorx

import "errors"

func IsIn(err error, targets ...error) bool {
	var ret bool
	for _, t := range targets {
		if errors.Is(err, t) {
			ret = true
			break
		}
	}
	return ret
}

func AsIn(err error, targets ...any) bool {
	var ret bool
	for _, t := range targets {
		if errors.As(err, t) {
			ret = true
			break
		}
	}
	return ret
}
