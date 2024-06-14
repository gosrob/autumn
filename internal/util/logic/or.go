package logic

func OrGet[T any](is bool, t1, t2 T) T {
	if is {
		return t1
	}

	return t2
}
