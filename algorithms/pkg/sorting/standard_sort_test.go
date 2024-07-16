package sorting

import (
	"fmt"
	"sort"
	"testing"
)

func TestIntSlice(t *testing.T) {
	unsorted := sort.IntSlice{10, 2, 5, 1, 0}

	fmt.Printf("unsorted: %v\n", unsorted)

	unsorted.Sort()

	fmt.Printf("sorted: %v\n", unsorted)

	index := sort.SearchInts(unsorted, 5)

	fmt.Printf("index: %v\n", index)
}
