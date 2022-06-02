package util

// exchange the elems at position `i` and `j`
func Swap[T any](arr []T, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

// move the targe elem to the tail of the slice then do the pop
// so this method will change the order of the origin slide
func RemoveAt[T any](arr *[]T, i int) {
	Swap(*arr, i, len(*arr)-1)
	*arr = (*arr)[0 : len(*arr)-1]
}
