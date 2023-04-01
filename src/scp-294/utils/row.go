package utils

import (
	"bytes"
)

func ByteArrayToRowBytes(byteToStr ByteToStr, arr []byte, off int, len2 int, withDetails bool) []byte {
	buffer := new(bytes.Buffer)
	buffer.Grow(len(arr) * (byteToStr.getWidth() + 2))
	count := 0
	for i := 0; i < len2; i++ {
		count++
		AppendStringWithSpace(buffer, byteToStr.toString(arr[off+i]), byteToStr.getWidth())
		if withDetails {
			if count&0x0F == 0 {
				buffer.WriteString(", ")
			} else {
				buffer.WriteString(" ")
			}
		} else {
			buffer.WriteString(", ")
		}
	}
	return buffer.Bytes()
}

func ByteArrayToRowDetailsBytes(arr []byte, off int, len int) []byte {
	buffer := new(bytes.Buffer)
	count := 0
	for i := off; i < off+len; i++ {
		count++
		if arr[i] >= 32 && arr[i] <= 126 {
			buffer.WriteByte(arr[i])
		} else {
			buffer.WriteString(CharNULL)
		}
		if count&0x0F == 0 {
			buffer.WriteString(", ")
		}
	}
	return buffer.Bytes()
}

func AppendStringWithChar(buffer *bytes.Buffer, str string, expectedLen int, char rune) {
	diff := expectedLen - len(str)
	if diff > 0 {
		for i := 0; i < diff; i++ {
			buffer.WriteString(string(char))
		}
	}

	buffer.WriteString(str)
}

func AppendStringWithSpace(buffer *bytes.Buffer, str string, expectedLen int) {
	AppendStringWithChar(buffer, str, expectedLen, ' ')
}

func AppendStringWith0(buffer *bytes.Buffer, str string, expectedLen int) {
	AppendStringWithChar(buffer, str, expectedLen, '0')
}
