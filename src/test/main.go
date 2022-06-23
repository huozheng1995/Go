package main

import (
	"fmt"
	"github.com/edward/test/utils"
)

func main() {
	//utils.PrintFileBytesHex("C:\\Users\\user\\downloads\\new 2.txt")
	//utils.PrintFileBytesHex("C:\\source\\binlog\\cdc\\binlog.000002")
	//utils.PrintFileBytesDec("C:\\source\\binlog\\driver\\binlog.000001")

	//bytes := []int64{-40, -66}
	//bytes2 := []byte{26, 101, 26, 32, 118, 101, 26, 101, 114, 110, byte(uint8(-19 + 256)), 32, 111, 98, 108, 111, 104, 121, 32, 115, 32, 26, 101, 114, 118, byte(uint8(-31 + 256)), 110, 107, 121}
	//fmt.Println(string(bytes2))
	//fmt.Println(utils.DecToHexArray(bytes))
	//fmt.Println(utils.DecToHex(65536))
	//fmt.Println(utils.HexToDec("F7E0"))
	//fmt.Println(utils.DecToBin(-19))

	decArray1 := utils.HexArrayToDecArray("00 00 00 00 00 30 \n30 30 30 30 53 51 4C 31 31 30 31 33 00 00 00 00 \n00 00 00 00 00 01 00 00 00 01 00 00 00 22 FF FF \nFF 00 00 00 00 20 20 20 20 20 20 20 20 20 20 20 \n00 12 54 45 53 54 20 20 20 20 20 20 20 20 20 20 \n20 20 20 20 00 00 00 00 FF 00 00 00 00 00 00 00 \n00 00 07 00 00 00 00 00 00 00 00 00 00 00 00 00 \n00 00 01 00 00 00 00 00 80 00 00 00 00 00 00 00 \nC1 01 04 E4 00 00 00 00 00 00 00 00 08 00 00 00 \n00 00 00 00 0C 43 41 54 41 4C 4F 47 5F 4E 41 4D \n45 00 00 00 00 00 00 00 00 FF 00 00 00 00 00 00 \n00 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 \n09 53 51 4C 54 41 42 4C 45 53 00 00 00 08 53 59 \n53 49 42 4D 20 20 00 00 00 0C 43 41 54 41 4C 4F \n47 5F 4E 41 4D 45 00 00 00 00 FF ")
	byteArray1 := utils.DecArrayToByteArray(decArray1)
	int8Array1 := utils.ByteArrayToInt8Array(byteArray1)

	decArray := utils.HexArrayToDecArray("00 00 00 00 00 30 \n30 30 30 30 53 51 4C 31 31 30 35 37 00 00 00 00 \n00 00 00 00 00 01 00 00 00 01 00 00 00 22 FF FF \nFF 00 00 00 00 20 20 20 20 20 20 20 20 20 20 20 \n00 12 54 45 53 54 44 42 20 20 20 20 20 20 20 20 \n20 20 20 20 00 00 00 00 FF 00 00 00 00 00 00 00 \n00 00 07 00 00 00 00 00 00 00 00 00 00 00 00 00 \n00 00 01 00 00 00 00 00 80 00 00 00 00 00 00 00 \nC1 01 04 B8 00 00 00 00 00 00 00 00 08 00 00 00 \n00 00 0C 43 41 54 41 4C 4F 47 5F 4E 41 4D 45 00 \n00 00 00 00 00 00 00 00 00 FF 00 00 00 00 00 00 \n00 01 00 00 00 00 00 00 00 00 00 00 00 00 09 53 \n51 4C 54 41 42 4C 45 53 00 00 00 08 53 59 53 49 \n42 4D 20 20 00 00 00 0C 43 41 54 41 4C 4F 47 5F \n4E 41 4D 45 00 00 00 00 00 00 FF ")
	byteArray := utils.DecArrayToByteArray(decArray)
	int8Array := utils.ByteArrayToInt8Array(byteArray)

	fmt.Println(string(byteArray1))
	fmt.Println(string(byteArray))
	utils.PrintInt8Array(int8Array1)
	utils.PrintInt8Array(int8Array)

}

func printASCII() {
	fmt.Printf("OK: string=%v, bytes=%v\n", "2021-12-05 14:55:55", []byte("2021-12-05 14:55:55"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "01ABCDEFG", []byte("01ABCDEFG"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "HIJKLMN", []byte("HIJKLMN"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "OPQRST", []byte("OPQRST"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "UVWXYZ", []byte("UVWXYZ"))
	fmt.Printf("OK: string=%v, bytes=%v\n", "abcdefg", []byte("abcdefg"))
}
