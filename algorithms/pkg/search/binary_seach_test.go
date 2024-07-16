package search

import "testing"

func TestBinarySearchNumericList(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	searchFor := 8
	expectPosition := 7

	position, err := BinarySearchNumericList(list, searchFor)

	if err != nil {
		t.Fatal("Error in BinarySearchNumberList")
	}

	if position != expectPosition {
		t.Errorf("BinarySearchNumberList: expect %d, got %d", expectPosition, position)
	}

	searchFor = 17
	expectPosition = 16

	position, err = BinarySearchNumericList(list, searchFor)

	if err != nil {
		t.Fatal("Error in BinarySearchNumberList")
	}

	if position != expectPosition {
		t.Errorf("BinarySearchNumberList: expect %d, got %d", expectPosition, position)
	}
}

func TestBinarySearchNumericList_ErrorOnEmpty(t *testing.T) {
	emptyList := make([]int, 0)
	searchFor := 8
	_, err := BinarySearchNumericList(emptyList, searchFor)

	if err == nil {
		t.Fatal("Should return an error when list is empty")
	}
}
