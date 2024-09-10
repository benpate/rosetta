package slice

import "testing"

func TestAt(t *testing.T) {

	var tests = []struct {
		slice  []int
		index  int
		result int
	}{
		{[]int{1, 2, 3, 4, 5}, 0, 1},
		{[]int{1, 2, 3, 4, 5}, 2, 3},
		{[]int{1, 2, 3, 4, 5}, 4, 5},
		{[]int{1, 2, 3, 4, 5}, 5, 0},
		{[]int{1, 2, 3, 4, 5}, -1, 0},
	}

	for _, test := range tests {
		result := At(test.slice, test.index)

		if result != test.result {
			t.Errorf("Expected At(%v, %v) to be %v, but got %v", test.slice, test.index, test.result, result)
		}
	}

}
