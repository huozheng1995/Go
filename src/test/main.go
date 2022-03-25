package main

import (
	"fmt"
	"github.com/edward/test/utils"
)

func main() {
	//utils.PrintFileBytesHex2("C:\\source\\binlog\\binlog.000001", 0, 50)
	utils.PrintFileBytesHex2("C:\\source\\binlog\\test.000013", 0, 50)
	//utils.PrintFileBytes("C:\\source\\binlog\\binlog.000002")
	//utils.PrintFileBytesHex("C:\\source\\binlog\\binlog.01")
	//utils.PrintFileBytes("C:\\source\\binlog\\binlog.01")
	//utils.PrintFile("C:\\source\\binlog\\binlog.000001")
}

func printASCII() {
	fmt.Printf("OK: string=%v, bytes=%v\n", "ABCDEFG", []byte("ABCEFG"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "HIJKLMN", []byte("HIJKLMN"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "OPQRST", []byte("OPQRST"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "UVWXYZ", []byte("UVWXYZ"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "abcdefg", []byte("abcdefg"))
}
