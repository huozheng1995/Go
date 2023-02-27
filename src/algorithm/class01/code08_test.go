package class01

import (
	"fmt"
	"math/rand"
	"myutil"
	"testing"
	"time"
)

func Test_08(t *testing.T) {
	kinds := 10
	range1 := 200
	testTime := 10000
	max := 9
	fmt.Println("begin...")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < testTime; i++ {
		a := r.Intn(max) + 1
		b := r.Intn(max) + 1
		k := a
		m := b
		if a == b {
			m++
		} else if a > b {
			k = b
			m = a
		}

		arrayKM := myutil.GenerateRandomArrayKM(kinds, range1, k, m)
		a1 := int32(findKTimesNumComp(arrayKM, k, m))
		a2 := int32(findKTimesNum(arrayKM, k, m))
		if a1 != a2 {
			fmt.Printf("Error! a1: %d, a2: %d \n", a1, a2)
		}
	}
	fmt.Println("end...")
}

//异或运算：arr中有一种数a出现了K次，其它数出现了M次，
//其中M>1且K<M，要求找到a
//额外空间复杂度O(1), 时间复杂度0(N)
func findKTimesNum(arr []int, k int, m int) (a int) {
	//所有数的每一位按位相加
	tempArr := make([]int, 32)
	for _, val := range arr {
		for i := 0; i < 32; i++ {
			tempArr[i] += val >> i & 1
		}
		//fmt.Println(tempArr)
	}

	a = 0
	for i := 0; i < 32; i++ {
		if tempArr[i]%m > 0 {
			a |= 1 << i
		}
	}

	return a
}

func findKTimesNumComp(arr []int, k int, m int) (a int) {
	mapA := make(map[int]int)
	for _, val := range arr {
		if mapA[val] > 0 {
			mapA[val]++
		} else {
			mapA[val] = 1
		}
	}

	for key, val := range mapA {
		if val == k {
			return key
		}
	}

	return -1
}
