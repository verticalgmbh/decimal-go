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
func FromBool(value bool) Decimal {
	if value {
		return Decimal{
			lo: uint32(1)}
	}
	return Decimal{}
}

// FromUInt8 - creates a decimal from a uint8 value
func FromUInt8(value uint8) Decimal {
	return Decimal{
		lo: uint32(value)}
}

// FromUInt16 - creates a decimal from a uint16 value
func FromUInt16(value uint16) Decimal {
	return Decimal{
		lo: uint32(value)}
}

// FromUInt32 - creates a decimal from a uint32 value
func FromUInt32(value uint32) Decimal {
	return Decimal{
		lo: value}
}

// FromUInt64 - creates a decimal from a uint64 value
func FromUInt64(value uint64) Decimal {
	return Decimal{
		lo:  (uint32)(value & 0xFFFFFFFF),
		mid: (uint32)((value >> 32) & 0xFFFFFFFF)}
}

// FromInt64Frac - creates a decimal from a mantisse and fractional information
func FromInt64Frac(value int64, fracdigits uint8) Decimal {
	dec := FromInt64(value)
	dec.flags |= uint32(fracdigits) << 16
	return dec
}

// FromUInt64Frac - creates a decimal from a mantisse and fractional information
func FromUInt64Frac(value uint64, fracdigits uint8) Decimal {
	dec := FromUInt64(value)
	dec.flags |= uint32(fracdigits) << 16
	return dec
}

// FromInt8 - creates a decimal from a int8 value
func FromInt8(value int8) Decimal {
	return FromInt32(int32(value))
}

// FromInt16 - creates a decimal from a int16 value
func FromInt16(value int16) Decimal {
	return FromInt32(int32(value))
}

// FromInt32 - creates a decimal from a int32 value
func FromInt32(value int32) Decimal {
	dec := Decimal{}
	if value < 0 {
		value = -value
		dec.flags = 0x80000000
	} else {
		dec.flags = 0
	}

	dec.lo = uint32(value)
	dec.mid = 0
	dec.hi = 0
	return dec
}

// FromInt64 - creates a decimal from a int64 value
func FromInt64(value int64) Decimal {
	dec := Decimal{}
	if value < 0 {
		value = -value
		dec.flags = 0x80000000
	} else {
		dec.flags = 0
	}

	dec.lo = (uint32)(value & 0xFFFFFFFF)
	dec.mid = (uint32)((value >> 32) & 0xFFFFFFFF)
	dec.hi = 0
	return dec
}

// FromFloat32 - creates a decimal from a float32 value
func FromFloat32(value float32) (Decimal, error) {
	return FromFloat64(float64(value))
}

// FromFloat64 - creates a decimal from a float64 value
func FromFloat64(value float64) (Decimal, error) {
	lo, mid, hi, flags, err := varDecFromR8(value)
	if err != nil {
		return Decimal{}, err
	}

	return Decimal{
		lo:    lo,
		mid:   mid,
		hi:    hi,
		flags: flags}, nil
}

// Float32 - converts the decimal to it's float32 representation
func (dec Decimal) Float32() float32 {
	return float32(dec.Float64())
}

// Float64 - converts the decimal to it's float64 representation
func (dec Decimal) Float64() float64 {
	scaleex := (dec.flags & 0xFF0000) >> 16
	scale := math.Pow(10, float64(scaleex))
	result := float64(dec.lo)/scale + float64(dec.mid)*4294967296/scale + float64(dec.hi)*4294967296*4294967296/scale

	if dec.flags&0x80000000 > 0 {
		return -result
	}

	return result
}

// IsNeg - returns whether the decimal is a negative number
func (dec Decimal) IsNeg() bool {
	return dec.flags&0x80000000 > 0
}

// Exp - returns exponent of this decimal
func (dec Decimal) Exp() uint8 {
	return uint8((dec.flags & 0xFF0000) >> 16)
}

// Int64 - get int64 value stored in decimal
//         if value is a fractional value, only the non fractional
//         part is returned
func (dec Decimal) Int64() int64 {
	value := (int64(dec.mid) << 32) | int64(dec.lo)
	exp := dec.Exp()
	for exp > 0 {
		value /= 10
		exp--
	}

	if dec.IsNeg() {
		return -value
	}

	return value
}

// Int32 - get int32 value stored in decimal
//         if value is a fractional value, only the non fractional
//         part is returned
func (dec Decimal) Int32() int32 {
	value := (int64(dec.mid) << 32) | int64(dec.lo)
	exp := dec.Exp()
	for exp > 0 {
		value /= 10
		exp--
	}

	if dec.IsNeg() {
		return -int32(value)
	}

	return int32(value)
}

// String - get string representation of the decimal
func (dec Decimal) String() string {
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
