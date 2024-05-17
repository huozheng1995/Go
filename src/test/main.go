package main

import (
	"fmt"
	"os"
)

func main() {
	bytes := make([]byte, 20)
	fmt.Println(bytes)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = byte(i)
	}
	fmt.Println(bytes)

	bytes2 := bytes[5:15]
	fmt.Println(bytes2)

	file, err := os.Open("go.mod")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes2 = bytes2[:cap(bytes2)]
	n, err := file.Read(bytes2)
	fmt.Println(n)
	fmt.Println(bytes2)
	fmt.Println(bytes)

}
