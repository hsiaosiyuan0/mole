package util

func Merge[T1 comparable, T2 any](src, stuff map[T1]T2) {
	for k, v := range stuff {
		if _, ok := src[k]; !ok {
			src[k] = v
		}
	}
}

func PickOne[T1 comparable, T2 any](sets map[T1]T2) (v T2) {
	for _, v = range sets {
		return
	}
	return
}

func PickOnePair[T1 comparable, T2 any](sets map[T1]T2) (k T1, v T2) {
	for k, v = range sets {
		return
	}
	return
}

func TakeOne[T1 comparable, T2 any](sets map[T1]T2) (v T2) {
	var k T1
	for k, v = range sets {
		delete(sets, k)
		return
	}
	return
}

func CloneMap[T1 comparable, T2 any](m map[T1]T2) map[T1]T2 {
	nm := make(map[T1]T2, 0)
	for k, v := range m {
		nm[k] = v
	}
	return nm
}
