package utils

func Find[E any](collection []E, f func(*E) bool) *E {
	for _, v := range collection {
		if f(&v) {
			return &v
		}
	}
	return nil
}
