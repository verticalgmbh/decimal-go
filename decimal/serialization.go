package decimal

import (
	"io"
)

func toUInt32(data []byte) uint32 {
	return uint32(data[0]) | (uint32(data[1]) << 8) | (uint32(data[2]) << 16) | (uint32(data[3]) << 24)
}

// ReadDecimal reads a .NET decimal number type
func ReadDecimal(reader io.Reader) (*Decimal, error) {
	buffer := make([]byte, 16)
	_, err := io.ReadFull(reader, buffer)
	if err != nil {
		return nil, err
	}

	return &Decimal{
		lo:    toUInt32(buffer[0:4]),
		mid:   toUInt32(buffer[4:8]),
		hi:    toUInt32(buffer[8:12]),
		flags: toUInt32(buffer[12:16])}, nil
}

// WriteDecimal writes a decimal value in .Net decimal format to the specified writer
func WriteDecimal(writer io.Writer, dec *Decimal) error {
	var buffer []byte = []byte{
		byte(dec.lo & 0xFF), byte((dec.lo >> 8) & 0xFF), byte((dec.lo >> 16) & 0xFF), byte((dec.lo >> 24) & 0xFF),
		byte(dec.mid & 0xFF), byte((dec.mid >> 8) & 0xFF), byte((dec.mid >> 16) & 0xFF), byte((dec.mid >> 24) & 0xFF),
		byte(dec.hi & 0xFF), byte((dec.hi >> 8) & 0xFF), byte((dec.hi >> 16) & 0xFF), byte((dec.hi >> 24) & 0xFF),
		byte(dec.flags & 0xFF), byte((dec.flags >> 8) & 0xFF), byte((dec.flags >> 16) & 0xFF), byte((dec.flags >> 24) & 0xFF)}

	_, err := writer.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}
