package maps

func Values[M ~map[K]V, K comparable, V any](m M) []V {
	res := make([]V, 0, len(m))
	for _, v := range m {
		res = append(res, v)
	}
	return res
}
