package main

import "fmt"

const (
	Arr1Size     = 4
	Arr2Size     = 5
	MergeArrSize = Arr1Size + Arr2Size

	SortArrSize = 6
)

func main() {
	// Task 1
	arr1 := [Arr1Size]int{1, 4, 8, 10}
	arr2 := [Arr2Size]int{2, 3, 5, 6, 7}

	mergedArray := mergeSortedArrays(arr1, arr2)
	fmt.Println("Merged Sorted Array:", mergedArray)

	// Task 2
	arr := [SortArrSize]int{64, 34, 25, 12, 22, 11} // Пример массива для сортировки
	fmt.Println("Unsorted Array:", arr)
	bubbleSort(&arr)
	fmt.Println("Sorted Array:", arr)
}

func mergeSortedArrays(arr1 [Arr1Size]int, arr2 [Arr2Size]int) [MergeArrSize]int {
	var merged [MergeArrSize]int
	i, j, k := 0, 0, 0

	// Слияние до тех пор, пока не закончатся элементы в одном из массивов
	for i < len(arr1) && j < len(arr2) {
		if arr1[i] < arr2[j] {
			merged[k] = arr1[i]
			i++
		} else {
			merged[k] = arr2[j]
			j++
		}
		k++
	}

	for i < Arr1Size {
		merged[k] = arr1[i]
		i++
		k++
	}

	for j < Arr2Size {
		merged[k] = arr2[j]
		j++
		k++
	}

	return merged
}

func bubbleSort(arr *[SortArrSize]int) {
	n := len(arr)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}
