package charset

type ICharset interface {
	GetName() string
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
	EncodeToStr([]byte) (string, error)
	DecodeToStr([]byte) (string, error)
}

type Charset struct {
	ICharset
}

func (c Charset) EncodeToStr(arr []byte) (string, error) {
	encoded, _ := c.Encode(arr)
	return string(encoded), nil
}

func (c Charset) DecodeToStr(arr []byte) (string, error) {
	decoded, _ := c.Decode(arr)
	return string(decoded), nil
}
