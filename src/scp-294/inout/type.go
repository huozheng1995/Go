package inout

type TypeCode int

const (
	Text TypeCode = iota
	File
)

type Type struct {
	Code TypeCode
	Name string
	Desc string
}

func CreateTypes() []Type {
	types := make([]Type, 0)
	text := Type{Text, "Text", "Input data in the textarea"}
	file := Type{File, "File", "Select a file to parse"}
	types = append(types, text)
	types = append(types, file)
	return types
}
