package main

import (
	"fmt"
	"github.com/edward/test/utils"
)

func main() {
	//utils.PrintFileBytesHex("C:\\Users\\user\\downloads\\new 2.txt")
	//utils.PrintFileBytesHex("C:\\source\\binlog\\cdc\\binlog.000002")
	//utils.PrintFileBytesDec("C:\\source\\binlog\\driver\\binlog.000001")

	bytes := []int64{-40, -66}
	bytes2 := []byte{26, 101, 26, 32, 118, 101, 26, 101, 114, 110, byte(uint8(-19 + 256)), 32, 111, 98, 108, 111, 104, 121, 32, 115, 32, 26, 101, 114, 118, byte(uint8(-31 + 256)), 110, 107, 121}
	fmt.Println(string(bytes2))
	fmt.Println(utils.DecToHexArray(bytes))
	fmt.Println(utils.DecToHex(65536))
	fmt.Println(utils.HexToDec("F7E0"))
	fmt.Println(utils.DecToBin(-19))
}

func printASCII() {
	fmt.Printf("OK: string=%v, bytes=%v\n", "2021-12-05 14:55:55", []byte("2021-12-05 14:55:55"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "01ABCDEFG", []byte("01ABCDEFG"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "HIJKLMN", []byte("HIJKLMN"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "OPQRST", []byte("OPQRST"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "UVWXYZ", []byte("UVWXYZ"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "abcdefg", []byte("abcdefg"))
}
