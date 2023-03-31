package main

import "fmt"

func main() {
	slice1 := make([]byte, 0, 10)
	fmt.Println(slice1)

	slice2 := make([]byte, 10)
	fmt.Println(slice2)
}
