package decimal

import (
	"math"
	"math/big"
)

// Decimal - number with arbitrary decimal precision
//           defined to work exactly like the decimal type in .NET
type Decimal struct {
	lo    uint32
	mid   uint32
	hi    uint32
	flags uint32
}

// FromBool - creates a decimal from a boolean value
func FromBool(value bool) *Decimal {
	if value {
		return FromUInt8(1)
	} else {
		return FromUInt8(0)
	}
}

// FromUInt8 - creates a decimal from a uint8 value
func FromUInt8(value uint8) *Decimal {
	dec := &Decimal{}
	dec.AssignUInt8(value)
	return dec
}

// FromUInt16 - creates a decimal from a uint16 value
func FromUInt16(value uint16) *Decimal {
	dec := &Decimal{}
	dec.AssignUInt16(value)
	return dec
}

// FromUInt32 - creates a decimal from a uint32 value
func FromUInt32(value uint32) *Decimal {
	dec := &Decimal{}
	dec.AssignUInt32(value)
	return dec
}

// FromUInt64 - creates a decimal from a uint64 value
func FromUInt64(value uint64) *Decimal {
	dec := &Decimal{}
	dec.AssignUInt64(value)
	return dec
}

// FromInt64Frac - creates a decimal from a mantisse and fractional information
func FromInt64Frac(value int64, fracdigits uint8) *Decimal {
	dec := &Decimal{}

	if value < 0 {
		value = -value
		dec.flags = 0x80000000
	}

	dec.lo = (uint32)(value & 0xFFFFFFFF)
	dec.mid = (uint32)((value >> 32) & 0xFFFFFFFF)
	dec.hi = 0
	dec.flags |= uint32(fracdigits) << 16
	return dec
}

// FromUInt64Frac - creates a decimal from a mantisse and fractional information
func FromUInt64Frac(value uint64, fracdigits uint8) *Decimal {
	dec := &Decimal{}

	dec.lo = (uint32)(value & 0xFFFFFFFF)
	dec.mid = (uint32)((value >> 32) & 0xFFFFFFFF)
	dec.hi = 0
	dec.flags = uint32(fracdigits) << 16
	return dec
}

// FromInt8 - creates a decimal from a int8 value
func FromInt8(value int8) *Decimal {
	dec := &Decimal{}
	dec.AssignInt8(value)
	return dec
}

// FromInt16 - creates a decimal from a int16 value
func FromInt16(value int16) *Decimal {
	dec := &Decimal{}
	dec.AssignInt16(value)
	return dec
}

// FromInt32 - creates a decimal from a int32 value
func FromInt32(value int32) *Decimal {
	dec := &Decimal{}
	dec.AssignInt32(value)
	return dec
}

// FromInt64 - creates a decimal from a int64 value
func FromInt64(value int64) *Decimal {
	dec := &Decimal{}
	dec.AssignInt64(value)
	return dec
}

// FromFloat32 - creates a decimal from a float32 value
func FromFloat32(value float32) (*Decimal, error) {
	return FromFloat64(float64(value))
}

// FromFloat64 - creates a decimal from a float64 value
func FromFloat64(value float64) (*Decimal, error) {
	dec := &Decimal{}
	err := dec.AssignFloat64(value)
	if err != nil {
		return nil, err
	}
	return dec, nil
}

// AssignFloat32 - assigns a float value to this decimal
func (dec *Decimal) AssignFloat32(value float32) error {
	return dec.AssignFloat64(float64(value))
}

// AssignFloat64 - assigns a float value to this decimal
func (dec *Decimal) AssignFloat64(value float64) error {
	lo, mid, hi, flags, err := varDecFromR8(value)
	if err != nil {
		return err
	}

	dec.lo = lo
	dec.mid = mid
	dec.hi = hi
	dec.flags = flags
	return nil
}

// AssignUInt8 - assigns a uint8 to this decimal
func (dec *Decimal) AssignUInt8(value uint8) {
	dec.AssignUInt32(uint32(value))
}

// AssignUInt16 - assigns a uint16 to this decimal
func (dec *Decimal) AssignUInt16(value uint16) {
	dec.AssignUInt32(uint32(value))
}

// AssignUInt32 - assigns a uint32 to this decimal
func (dec *Decimal) AssignUInt32(value uint32) {
	dec.lo = value
	dec.mid = 0
	dec.hi = 0
	dec.flags = 0
}

// AssignUInt64 - assigns a uint32 to this decimal
func (dec *Decimal) AssignUInt64(value uint64) {
	dec.lo = (uint32)(value & 0xFFFFFFFF)
	dec.mid = (uint32)((value >> 32) & 0xFFFFFFFF)
	dec.hi = 0
	dec.flags = 0
}

// AssignInt8 - assigns a uint8 to this decimal
func (dec *Decimal) AssignInt8(value int8) {
	dec.AssignInt32(int32(value))
}

// AssignInt16 - assigns a uint16 to this decimal
func (dec *Decimal) AssignInt16(value int16) {
	dec.AssignInt32(int32(value))
}

// AssignInt32 - assigns a uint32 to this decimal
func (dec *Decimal) AssignInt32(value int32) {
	if value < 0 {
		value = -value
		dec.flags = 0x80000000
	} else {
		dec.flags = 0
	}

	dec.lo = uint32(value)
	dec.mid = 0
	dec.hi = 0
}

// AssignInt64 - assigns a uint32 to this decimal
func (dec *Decimal) AssignInt64(value int64) {
	if value < 0 {
		value = -value
		dec.flags = 0x80000000
	} else {
		dec.flags = 0
	}

	dec.lo = (uint32)(value & 0xFFFFFFFF)
	dec.mid = (uint32)((value >> 32) & 0xFFFFFFFF)
	dec.hi = 0
}

// Float32 - converts the decimal to it's float32 representation
func (dec *Decimal) Float32() float32 {
	return float32(dec.Float64())
}

// Float64 - converts the decimal to it's float64 representation
func (dec *Decimal) Float64() float64 {
	scaleex := (dec.flags & 0xFF0000) >> 16
	scale := math.Pow(10, float64(scaleex))
	result := float64(dec.lo)/scale + float64(dec.mid)*4294967296/scale + float64(dec.hi)*4294967296*4294967296/scale

	if dec.flags&0x80000000 > 0 {
		return -result
	}

	return result
}

// IsNeg - returns whether the decimal is a negative number
func (dec *Decimal) IsNeg() bool {
	return dec.flags&0x80000000 > 0
}

// Exp - returns exponent of this decimal
func (dec *Decimal) Exp() uint8 {
	return uint8((dec.flags & 0xFF0000) >> 16)
}

// String - get string representation of the decimal
func (dec *Decimal) String() string {
	number := big.Int{}
	number.SetBytes([]byte{
		byte(dec.hi >> 24), byte(dec.hi >> 16), byte(dec.hi >> 8), byte(dec.hi),
		byte(dec.mid >> 24), byte(dec.mid >> 16), byte(dec.mid >> 8), byte(dec.mid),
		byte(dec.lo >> 24), byte(dec.lo >> 16), byte(dec.lo >> 8), byte(dec.lo)})

	numstr := number.String()
	if dec.IsNeg() {
		if dec.Exp() > 0 {
			return "-" + numstr[:len(numstr)-int(dec.Exp())] + "." + numstr[dec.Exp()+1:]
		}
		return "-" + numstr
	}

	if dec.Exp() > 0 {
		return numstr[:len(numstr)-int(dec.Exp())] + "." + numstr[dec.Exp()+1:]
	}

	return numstr
}
