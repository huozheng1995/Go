package utils

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

func DecToBin(n int64) string {
	if n < 0 {
		n += 256
	}
	if n == 0 {
		return "0"
	}
	s := ""
	for q := n; q > 0; q = q / 2 {
		m := q % 2
		s = fmt.Sprintf("%v%v", m, s)
	}
	return s
}

func DecToOct(d int64) int64 {
	if d == 0 {
		return 0
	}
	if d < 0 {
		d += 256
	}
	s := ""
	for q := d; q > 0; q = q / 8 {
		m := q % 8
		s = fmt.Sprintf("%v%v", m, s)
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Println("Decimal to octal error:", err.Error())
		return -1
	}
	return int64(n)
}

func DecToHex(n int64) string {
	if n < 0 {
		n = n + 256
	}
	if n == 0 {
		return "00"
	}
	hex := map[int64]int64{10: 65, 11: 66, 12: 67, 13: 68, 14: 69, 15: 70}
	s := ""
	for q := n; q > 0; q = q / 16 {
		m := q % 16
		if m > 9 && m < 16 {
			m = hex[m]
			s = fmt.Sprintf("%v%v", string(m), s)
			continue
		}
		s = fmt.Sprintf("%v%v", m, s)
	}
	if n/16 == 0 {
		return "0" + s
	}
	return s
}

func DecArrayToHexArray(intArr []int64) string {
	s := ""
	for _, val := range intArr {
		s += DecToHex(val) + ", "
	}
	return s
}

func BinToDec(b string) (n int64) {
	s := strings.Split(b, "")
	l := len(s)
	i := 0
	d := float64(0)
	for i = 0; i < l; i++ {
		f, err := strconv.ParseFloat(s[i], 10)
		if err != nil {
			log.Println("Binary to decimal error:", err.Error())
			return -1
		}
		d += f * math.Pow(2, float64(l-i-1))
	}
	return int64(d)
}

func OctToDec(o int64) (n int64) {
	s := strings.Split(strconv.Itoa(int(o)), "")
	l := len(s)
	i := 0
	d := float64(0)
	for i = 0; i < l; i++ {
		f, err := strconv.ParseFloat(s[i], 10)
		if err != nil {
			log.Println("Octal to decimal error:", err.Error())
			return -1
		}
		d += f * math.Pow(8, float64(l-i-1))
	}
	return int64(d)
}

func HexToInt64(h string) (n int64) {
	s := strings.Split(strings.ToUpper(h), "")
	l := len(s)
	i := 0
	d := float64(0)
	hex := map[string]string{"A": "10", "B": "11", "C": "12", "D": "13", "E": "14", "F": "15"}
	for i = 0; i < l; i++ {
		c := s[i]
		if v, ok := hex[c]; ok {
			c = v
		}
		f, err := strconv.ParseFloat(c, 10)
		if err != nil {
			log.Println("Hexadecimal to decimal error:", err.Error())
			return -1
		}
		d += f * math.Pow(16, float64(l-i-1))
	}
	return int64(d)
}

//func HexArrayToDecArray(str string) (intArr []int64) {
//	decArray := make([]int64, 0, 100)
//	var sLeft, sRight byte
//	var val byte
//	for i := 0; i < len(str)+1; i++ {
//		if i == len(str) {
//			val = 0
//		} else {
//			val = str[i]
//		}
//		if (val >= '0' && val <= '9') || (val >= 'a' && val <= 'z') || (val >= 'A' && val <= 'Z') {
//			if sLeft == 0 {
//				sLeft = val
//			} else if sRight == 0 {
//				sRight = val
//			}
//		} else {
//			if sLeft > 0 && sRight > 0 {
//				hexValue := BytesToString([]uint8{sLeft, sRight})
//				sLeft, sRight = 0, 0
//				decArray = append(decArray, HexToInt64(hexValue))
//			}
//		}
//	}
//
//	return decArray
//}

func DecArrayToByteArray(intArr []int64) []byte {
	byteArray := make([]byte, len(intArr))
	for i := 0; i < len(intArr); i++ {
		byteArray[i] = Int64ToBytes(intArr[i])[7]
	}
	return byteArray
}

func PrintByteArray(arr []byte) {
	rowCount := 32

	totalRow := len(arr) / rowCount
	lastRowCount := len(arr) % rowCount
	if lastRowCount > 0 {
		totalRow++
	}

	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		byteIndex := rowIndex * rowCount
		if rowIndex == totalRow-1 {
			rowCount = lastRowCount
		}
		fmt.Printf("row%s(%s, %s, %s, %s): %s\n", Fill0(strconv.Itoa(rowIndex), printLen),
			Fill0(strconv.Itoa(byteIndex), printLen), Fill0(strconv.Itoa(byteIndex+8), printLen),
			Fill0(strconv.Itoa(byteIndex+16), printLen), Fill0(strconv.Itoa(byteIndex+24), printLen),
			ByteArrayToLine(arr[byteIndex:byteIndex+rowCount]))
	}
}

func PrintInt8ArrayWithComma(arr []int8) {
	str := ""
	for _, val := range arr {
		str = str + strconv.Itoa(int(val)) + ", "
	}
	fmt.Println(str)
}

func PrintInt8Array(arr []int8) {
	rowCount := 32

	totalRow := len(arr) / rowCount
	lastRowCount := len(arr) % rowCount
	if lastRowCount > 0 {
		totalRow++
	}

	for rowIndex := 0; rowIndex < totalRow; rowIndex++ {
		byteIndex := rowIndex * rowCount
		if rowIndex == totalRow-1 {
			rowCount = lastRowCount
		}
		fmt.Printf("row%s(%s, %s, %s, %s): %s\n", Fill0(strconv.Itoa(rowIndex), printLen),
			Fill0(strconv.Itoa(byteIndex), printLen), Fill0(strconv.Itoa(byteIndex+8), printLen),
			Fill0(strconv.Itoa(byteIndex+16), printLen), Fill0(strconv.Itoa(byteIndex+24), printLen),
			Int8ArrayToLine(arr[byteIndex:byteIndex+rowCount]))
	}
}

func OctToBin(o int64) string {
	d := OctToDec(o)
	if d == -1 {
		return ""
	}
	return DecToBin(d)
}

func HexToBin(h string) string {
	d := HexToInt64(h)
	if d == -1 {
		return ""
	}
	return DecToBin(d)
}

func BinToOct(b string) int64 {
	d := BinToDec(b)
	if d == -1 {
		return -1
	}
	return DecToOct(d)
}

func BinToHex(b string) string {
	d := BinToDec(b)
	if d == -1 {
		return ""
	}
	return DecToHex(d)
}
