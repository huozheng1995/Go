package myutil

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

func GenerateRandomArrayKM(maxKinds int, maxValue int, k int, m int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	numKinds := r.Intn(maxKinds) + 2
	arr := make([]int, k+(numKinds-1)*m)
	arrIndex := 0

	set := NewSet()
	a := r.Intn(maxValue+1) - r.Intn(maxValue)
	set.Add(a)
	for ; k > 0; k-- {
		arr[arrIndex] = a
		arrIndex++
	}
	numKinds--

	for ; numKinds > 0; numKinds-- {
		var b int
		for {
			b = r.Intn(maxValue+1) - r.Intn(maxValue)
			if !set.Contains(b) {
				break
			}
		}
		set.Add(b)
		for m2 := m; m2 > 0; m2-- {
			arr[arrIndex] = b
			arrIndex++
		}
	}

	for i := 0; i < len(arr); i++ {
		SwapInt(arr, i, r.Intn(len(arr)))
	}

	return arr
}
