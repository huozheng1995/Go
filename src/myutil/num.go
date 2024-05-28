package myutil

type NumUtil[T any] interface {
	ToString(T) string
	GetDisplaySize() int
	ToNum(string2 string) T
}
