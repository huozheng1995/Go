package processor

type IProcessor interface {
	GetName() string
	GetType() ProcType
	Process(arr []byte) ([]byte, error)
}

const (
	IbmEBCDICEncoder = "IbmEBCDICEncoder"
	IbmEBCDICDecoder = "IbmEBCDICDecoder"
)
