package sorting

import (
	"reflect"
	"testing"
)

func TestNumericSelectionSort(t *testing.T) {
	unsorted := []int{10, 32, 8, 2, 54}
	sortedExpected := []int{2, 8, 10, 32, 54}

	sorted := NumericSelectionSort(unsorted)

	if !reflect.DeepEqual(sorted, sortedExpected) {
		t.Errorf("sorted: %v, sortedExpected: %v", sorted, sortedExpected)
	}
}

func TestNumericSelectionSort_AlreadySorted(t *testing.T) {
	sortedPreviously := []int{2, 8, 10, 32, 54}

	sorted := NumericSelectionSort(sortedPreviously)

	if !reflect.DeepEqual(sorted, sortedPreviously) {
		t.Errorf("Should keep same sorting. sorted: %v, sortedExpected: %v", sorted, sortedPreviously)
	}
}
