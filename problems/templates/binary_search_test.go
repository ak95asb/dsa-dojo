package problems

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinarySearch(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   int
	}{
		{
			name:   "target found in middle",
			nums:   []int{-1, 0, 3, 5, 9, 12},
			target: 9,
			want:   4,
		},
		{
			name:   "target not found",
			nums:   []int{-1, 0, 3, 5, 9, 12},
			target: 2,
			want:   -1,
		},
		{
			name:   "single element found",
			nums:   []int{5},
			target: 5,
			want:   0,
		},
		{
			name:   "single element not found",
			nums:   []int{5},
			target: 3,
			want:   -1,
		},
		{
			name:   "target at beginning",
			nums:   []int{1, 2, 3, 4, 5},
			target: 1,
			want:   0,
		},
		{
			name:   "target at end",
			nums:   []int{1, 2, 3, 4, 5},
			target: 5,
			want:   4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BinarySearch(tt.nums, tt.target)
			assert.Equal(t, tt.want, got, "BinarySearch(%v, %d)", tt.nums, tt.target)
		})
	}
}
