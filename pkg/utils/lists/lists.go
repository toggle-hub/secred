package lists

func Map[T any, U any](list []T, fn func(elem T, index int, list []T) U) []U {
	newList := make([]U, len(list))

	for i, elem := range list {
		newList[i] = fn(elem, i, list)
	}

	return newList
}

func Reduce[T any, U any](list []T, accumulator U, fn func(U, T) U) U {
	acc := accumulator

	for _, elem := range list {
		acc = fn(acc, elem)
	}

	return acc
}
