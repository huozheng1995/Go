package myutil

// num to string

type NumToStr[T any] interface {
	ToString(T) string
	GetWidth() int
}
