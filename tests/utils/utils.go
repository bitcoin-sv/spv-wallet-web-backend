package utils

// Find the first element in a collection that satisfies a specified condition.
func Find[E any](collection []E, predicate func(E) bool) *E {
	for _, v := range collection {
		if predicate(v) {
			return &v
		}
	}
	return nil
}
