package sorting

import (
	"algorithms/pkg/interfaces"
)

func NumericSelectionSort[T interfaces.Numeric](arr []T) []T {
	ordered := make([]T, len(arr), len(arr))
	arrCopy := make([]T, len(arr), len(arr))
	copy(arrCopy, arr)

	for orderedIndex := range ordered {
		orderValue := arrCopy[0]
		orderIndex := 0

		for index, value := range arrCopy {
			if value < orderValue {
				orderValue = value
				orderIndex = index
			}
		}

		ordered[orderedIndex] = orderValue
		arrCopy = append(arrCopy[0:orderIndex], arrCopy[orderIndex+1:]...)
	}

	return ordered
}
