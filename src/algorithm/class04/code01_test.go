package class04

import (
	"fmt"
	"myutil"
	"testing"
)

func Test_01(t *testing.T) {
	testTime := 10000
	maxSize := 10000
	maxValue := 1000

	for i := 0; i < testTime; i++ {
		arr0 := myutil.GenerateRandomArray(maxSize, maxValue)
		//fmt.Printf("random arr0: %#v \n", arr0)
		arr1 := myutil.CopyArray(arr0)
		arr2 := myutil.CopyArray(arr0)
		mergeSort(arr1)
		mergeSortComp(arr2)
		if !myutil.CompareArray(arr1, arr2) {
			fmt.Printf("Error! index: %d \n", i)
			fmt.Printf("arr0: %#v \n", arr0)
			fmt.Printf("arr1: %#v \n", arr1)
			fmt.Printf("arr2: %#v \n", arr2)
			return
		}
	}
}

func Test_02(t *testing.T) {
	arr1 := []int{3, 6, -5}
	fmt.Printf("arr1 1: %#v \n", arr1)
	mergeSortComp(arr1)
	fmt.Printf("arr1 2: %#v \n", arr1)

}

//归并排序
func mergeSort(arr []int) {
	if arr == nil || len(arr) == 0 {
		return
	}
	mergeSortImpl(arr, 0, len(arr)-1)
}

func mergeSortImpl(arr []int, begin int, end int) {
	if end == begin {
		return
	}
	if end == begin+1 {
		if arr[begin] > arr[end] {
			temp := arr[begin]
			arr[begin] = arr[end]
			arr[end] = temp
		}
		return
	}
	mid := (end-begin)>>1 + begin
	mergeSortImpl(arr, begin, mid)
	mergeSortImpl(arr, mid+1, end)
	merge(arr, begin, mid, end)
}

func merge(arr []int, begin int, mid int, end int) {
	tempArr := make([]int, end-begin+1)
	index := 0
	leftBegin := begin
	rightBegin := mid + 1
	for leftBegin <= mid && rightBegin <= end {
		if arr[leftBegin] <= arr[rightBegin] {
			tempArr[index] = arr[leftBegin]
			leftBegin++
		} else {
			tempArr[index] = arr[rightBegin]
			rightBegin++
		}
		index++
	}
	if leftBegin > mid {
		for rightBegin <= end {
			tempArr[index] = arr[rightBegin]
			rightBegin++
			index++
		}
	} else {
		for leftBegin <= mid {
			tempArr[index] = arr[leftBegin]
			leftBegin++
			index++
		}
	}
	for i := 0; i < len(tempArr); i++ {
		arr[begin+i] = tempArr[i]
	}
}

//归并排序非递归实现
func mergeSortComp(arr []int) {
	if arr == nil || len(arr) == 0 {
		return
	}

	stepSize := 1
	for {
		begin := 0
		for begin < len(arr) {
			mid := begin + stepSize - 1
			if mid >= len(arr)-1 {
				break
			}
			end := mid + stepSize
			if end > len(arr)-1 {
				end = len(arr) - 1
			}
			merge(arr, begin, mid, end)
			begin = end + 1
		}

		if stepSize < len(arr)/2+1 {
			stepSize <<= 1
		} else {
			break
		}
	}
}
