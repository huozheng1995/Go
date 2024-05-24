package myutil

// num to string

type NumToStr[T any] interface {
	ToString(T) string
	GetWidth() int
}

// string to num

type StrToNum[T any] interface {
	ToNum(string2 string) T
}
