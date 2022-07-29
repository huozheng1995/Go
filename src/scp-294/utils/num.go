package utils

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

func HexToDec(h string) (n int64) {
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
