package constants

type FormatCode int

const (
	Hex FormatCode = iota
	Dec
	Bin
	HexByte
	DecByte
	DecInt8
	FormattedHexByte
	FormattedDecByte
	FormattedDecInt8
	RawBytes
)

type Format struct {
	Code FormatCode
	Name string
	Desc string
}

type FormatMapping struct {
	InputFormat      Format
	OutputFormats    []Format
	SupportProcessor bool
}

func CreateFormatMappings() []FormatMapping {
	hex := Format{Hex, "Hex numbers", "ABAB5, 12EF1, 56, 75, CCCCC, 2CDD, DC11248, 05, 12, FE, FF, "}
	dec := Format{Dec, "Dec numbers", "703157, 77553, 86, 117, 838860, 11485, 230756936, 5, 18, 254, 255, "}
	bin := Format{Bin, "Bin numbers", "10101011101010110101, 1010110, 1110101, 11001100110011001100, 10110011011101, 101, 10010, 11111110, "}
	hexByte := Format{HexByte, "Hex Byte numbers", "AB, EF, 56, 75, CC, 2C, DC, BB, FE, FF, "}
	decByte := Format{DecByte, "Dec Byte numbers", "171, 239, 86, 117, 204, 44, 220, 187, 254, 255, "}
	decInt8 := Format{DecInt8, "Dec Int8 numbers", "-85, -17, 86, 117, -52, 44, -36, -69, -2, -1, "}
	formattedHexByte := Format{FormattedHexByte, "Formatted Hex Byte numbers", ""}
	formattedDecByte := Format{FormattedDecByte, "Formatted Dec Byte numbers", ""}
	formattedDecInt8 := Format{FormattedDecInt8, "Formatted Dec Int8 numbers", ""}
	rawBytes := Format{RawBytes, "Raw Bytes", "Raw Bytes in textarea or file"}
	mappings := make([]FormatMapping, 0)
	mappings = append(mappings, FormatMapping{hex, []Format{hex, dec, bin}, false})
	mappings = append(mappings, FormatMapping{dec, []Format{hex, dec, bin}, false})
	mappings = append(mappings, FormatMapping{bin, []Format{hex, dec, bin}, false})
	mappings = append(mappings, FormatMapping{hexByte,
		[]Format{hexByte, decByte, decInt8, rawBytes, formattedHexByte, formattedDecByte, formattedDecInt8}, true})
	mappings = append(mappings, FormatMapping{decByte,
		[]Format{hexByte, decByte, decInt8, rawBytes, formattedHexByte, formattedDecByte, formattedDecInt8}, true})
	mappings = append(mappings, FormatMapping{decInt8,
		[]Format{hexByte, decByte, decInt8, rawBytes, formattedHexByte, formattedDecByte, formattedDecInt8}, true})
	mappings = append(mappings, FormatMapping{rawBytes,
		[]Format{hexByte, decByte, decInt8, rawBytes, formattedHexByte, formattedDecByte, formattedDecInt8}, true})
	return mappings
}
