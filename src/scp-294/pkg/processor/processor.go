package processor

type IProcessor interface {
	GetName() string
	GetType() ProcType
	Process(arr []byte) ([]byte, error)
}

// ProcName

const (
	IbmEBCDICEncoder = "IbmEBCDICEncoder"
	IbmEBCDICDecoder = "IbmEBCDICDecoder"
)

type ProcType int

const (
	Encode ProcType = iota
	Decode
	Compress
	Decompress
	Encrypt
	Decrypt
)
