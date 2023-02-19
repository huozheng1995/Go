package utils

import (
	"math/rand"
	"time"
)

func SwapInt(arr []int, index1 int, index2 int) {
	if index1 == index2 {
		return
	}
	arr[index1] = arr[index1] ^ arr[index2]
	arr[index2] = arr[index1] ^ arr[index2]
	arr[index1] = arr[index1] ^ arr[index2]
}

func GenerateRandomArray(maxSize int, maxValue int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	arr := make([]int, r.Intn(maxSize))
	for i := 0; i < len(arr); i++ {
		arr[i] = r.Intn(maxValue+1) - r.Intn(maxValue)
	}
	return arr
}
