package utils

func SafeDeref[T any](ptr *T) T {
	var initial T
	if ptr != nil {
		return *ptr
	}
	return initial
}
