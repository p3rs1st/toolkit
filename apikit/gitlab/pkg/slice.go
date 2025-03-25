package pkg

func MapFunc[S ~[]E, E, R any](s S, mapper func(E) R) []R {
	size := len(s)

	newslice := make([]R, size)
	for i := range s {
		newslice[i] = mapper(s[i])
	}

	return newslice
}

func FilterFunc[S ~[]E, E any](slice S, filter func(E) bool) []E {
	size := 0

	for i := range slice {
		if filter(slice[i]) {
			size++
		}
	}

	idx := 0
	newslice := make(S, size)

	for i := range slice {
		if filter(slice[i]) {
			newslice[idx] = slice[i]
			idx++
		}
	}

	return newslice
}
