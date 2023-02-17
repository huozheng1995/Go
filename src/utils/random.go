package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomArray(maxSize int32, maxValue int32) []int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	arr := make([]int32, r.Int31n(maxSize))
	for i := 0; i < len(arr); i++ {
		arr[i] = r.Int31n(maxValue+1) - r.Int31n(maxValue)
	}
	return arr
}
