package utils

// Reverse slice
func ReverseSlice[T any](arr []T) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i] // literally swapping the first and the last numbers
	}
}

// Since slices in Go are references to underlying arrays, any modification to the slice elements inside the function will reflect in the original slice. There’s no need to pass a pointer to the slice or return it.