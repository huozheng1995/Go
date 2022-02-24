package main

import "fmt"

func main() {
	fmt.Println(ByteToBit(5))
	fmt.Println(BytesToBit([]byte{5, 6, 7}))
	fmt.Println(FloatToBit(1.0))
	fmt.Println(FloatToBit(1.5))
	fmt.Println(FloatToBit(8.25))
	fmt.Println(Float64ToBit(8.25))
}
