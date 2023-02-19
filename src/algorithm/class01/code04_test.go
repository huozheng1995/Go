package class01

import "testing"

func Test_04(t *testing.T) {

}

//二分法：输入数组有序，找是否存在某个数
func bsExist(sortedArr []int, num int) bool {
	if sortedArr == nil || len(sortedArr) == 0 {
		return false
	}

	left := 0
	right := len(sortedArr) - 1
	var mid int
	for left < right {
		mid = left + (right-left)>>1
		if sortedArr[mid] < num {
			left = mid + 1
		} else if num < sortedArr[mid] {
			right = mid - 1
		} else {
			return true
		}
	}

	return sortedArr[left] == num
}

//二分法：一个数组，值可以为[+,-,0]，相邻的两个数不相等，求局部最小值
func bsPartialMinValue(arr []int, num int) int {
	if arr == nil || len(arr) == 0 {
		return -1
	} else if len(arr) == 1 {
		return 0
	}

	if arr[0] < arr[1] {
		return 0
	}
	if arr[len(arr)-1] < arr[len(arr)-2] {
		return len(arr) - 1
	}

	left := 1
	right := len(arr) - 2
	var mid int
	for left < right {
		mid = left + (right-left)>>1
		if arr[mid-1] < arr[mid] {
			right = mid - 1
		} else if arr[mid] > arr[mid+1] {
			left = mid + 1
		} else {
			return mid
		}
	}

	return left
}
