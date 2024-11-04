package processor

type ProcType int

const (
	Encode ProcType = iota
	Decode
	Compress
	Decompress
	Encrypt
	Decrypt
)
