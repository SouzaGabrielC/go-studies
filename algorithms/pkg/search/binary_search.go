package search

import (
	"algorithms/pkg/interfaces"
	"errors"
)

// BinarySearchNumericList find the numeric item in the list provided using binary search algorithm.
// list should be already in ascending order
// Example:
//
//		func Example() {
//		  list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
//		  searchFor := 8
//	      expectPosition := 7
//
//		  position, err := BinarySearchNumericList(list, searchFor)
//		  if err != nil { println(err) }
//			println("Position:", position)
//		  }
func BinarySearchNumericList[T interfaces.Numeric](list []T, item T) (int, error) {

	if len(list) == 0 {
		return -1, errors.New("list is empty")
	}

	var minIndex, maxIndex int

	minIndex = 0
	maxIndex = len(list) - 1

	for minIndex <= maxIndex {
		midIndex := (minIndex + maxIndex) / 2
		if list[midIndex] == item {
			return midIndex, nil
		}

		if item > list[midIndex] {
			minIndex = midIndex
		}

		if item < list[midIndex] {
			maxIndex = midIndex
		}
	}

	return -1, errors.New("item not found")
}
