package main

import (
	"fmt"
	"github.com/edward/test/utils"
)

func main() {
	utils.ReadFileBytes("C:\\source\\binlog.000001")
}

func printASCII() {
	fmt.Printf("OK: string=%v, bytes=%v\n", "ABCDEFG", []byte("ABCEFG"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "HIJKLMN", []byte("HIJKLMN"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "OPQRST", []byte("OPQRST"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "UVWXYZ", []byte("UVWXYZ"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "abcdefg", []byte("abcdefg"))
}