package file

type INumFile[T any] interface {
	Read(p []T) (int, error)
}
