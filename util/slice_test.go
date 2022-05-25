package util

import (
	"testing"
)

func TestSwap(t *testing.T) {
	a := []int{1, 2, 3}
	Swap(a, 1, 2)

	AssertEqual(t, 3, a[1], "should be ok")
	AssertEqual(t, 2, a[2], "should be ok")
}

func TestRemoveAt(t *testing.T) {
	a := []int{1, 2, 3}
	RemoveAt(&a, 1)

	AssertEqual(t, 2, len(a), "should be ok")
	AssertEqual(t, 3, a[1], "should be ok")
}
