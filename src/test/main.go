package main

import "fmt"

func main() {
	fmt.Printf("OK: string=%v, bytes=%v\n", "ABCEFG", []byte("ABCEFG"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "HIJKLMN", []byte("HIJKLMN"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "OPQRST", []byte("OPQRST"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "UVWXYZ", []byte("UVWXYZ"))
}
