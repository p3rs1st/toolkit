package pkg

func MapFunc[S ~[]E, E, R any](s S, mapper func(E) R) []R {
	size := len(s)
	newslice := make([]R, size)
	for i := range s {
		newslice[i] = mapper(s[i])
	}
	return newslice
}

func FilterFunc[S ~[]E, E any](s S, filter func(E) bool) S {
	size := 0
	for i := range s {
		if filter(s[i]) {
			size += 1
		}
	}

	j := 0
	newslice := make(S, size)
	for i := range s {
		if filter(s[i]) {
			newslice[j] = s[i]
			j += 1
		}
	}
	return newslice
}
