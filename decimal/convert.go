package decimal

import (
	"errors"
	"math"
	"unsafe"
)

var sDoublePowers10 []float64 = []float64{
	1, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	1e20, 1e21, 1e22, 1e23, 1e24, 1e25, 1e26, 1e27, 1e28, 1e29,
	1e30, 1e31, 1e32, 1e33, 1e34, 1e35, 1e36, 1e37, 1e38, 1e39,
	1e40, 1e41, 1e42, 1e43, 1e44, 1e45, 1e46, 1e47, 1e48, 1e49,
	1e50, 1e51, 1e52, 1e53, 1e54, 1e55, 1e56, 1e57, 1e58, 1e59,
	1e60, 1e61, 1e62, 1e63, 1e64, 1e65, 1e66, 1e67, 1e68, 1e69,
	1e70, 1e71, 1e72, 1e73, 1e74, 1e75, 1e76, 1e77, 1e78, 1e79,
	1e80}

var sPowers10 []uint32 = []uint32{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000}

var sUlongPowers10 []uint64 = []uint64{
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
	1000000000000000000,
	10000000000000000000,
}

func getExponent(value float64) uint {
	return uint((*(*uint64)(unsafe.Pointer(&value)) >> 52) & 0x7FF)
}

func uint64x64To128(a uint64, b uint64) (uint32, uint32, uint32, uint32, error) {
	low := uint64(uint32(a)) * uint64(uint32(b))
	mid := uint64(uint32(a)) * uint64(uint32(b>>32))
	high := uint64(uint32(a>>32)) * uint64(uint32(b>>32))

	high += mid >> 32
	mid <<= 32
	low += mid
	if low < mid {
		high++
	}

	if high > math.MaxUint32 {
		return 0, 0, 0, 0, errors.New("Decimal Overflow")
	}

	return uint32(low), uint32(low >> 32), uint32(high), uint32(0), nil
}

func varDecFromR8(value float64) (uint32, uint32, uint32, uint32, error) {
	var low uint32 = 0
	var mid uint32 = 0
	var hi uint32 = 0
	var flags uint32 = 0
	var err error

	// The most we can scale by is 10^28, which is just slightly more
	// than 2^93.  So a float with an exponent of -94 could just
	// barely reach 0.5, but smaller exponents will always round to zero.
	//
	DBLBIAS := uint(1022)
	exp := int(getExponent(value) - DBLBIAS)
	if exp < -94 {
		return 0, 0, 0, 0, nil // result should be zeroed out
	}

	if exp > 96 {
		return 0, 0, 0, 0, errors.New("decimal overflow")
	}

	if value < 0 {
		value = -value
		flags = 0x80000000
	}

	// Round the input to a 15-digit integer.  The R8 format has
	// only 15 digits of precision, and we want to keep garbage digits
	// out of the Decimal were making.
	//
	// Calculate max power of 10 input value could have by multiplying
	// the exponent by log10(2).  Using scaled integer multiplcation,
	// log10(2) * 2 ^ 16 = .30103 * 65536 = 19728.3.
	//
	dbl := value
	power := 14 - ((exp * 19728) >> 16)

	// power is between -14 and 43
	if power >= 0 {
		// We have less than 15 digits, scale input up.
		//
		if power > 28 {
			power = int(28)
		}

		dbl *= sDoublePowers10[power]
	} else {
		if power != -1 || dbl >= 1E15 {
			dbl /= sDoublePowers10[-power]
		} else {
			power = 0 // didn't scale it
		}
	}

	if dbl < 1E14 && power < 28 {
		dbl *= 10
		power++
	}

	// Round to int64
	//
	var mant uint64
	mant = uint64(int64(dbl))
	dbl -= float64(int64(mant)) // difference between input & integer
	if dbl > 0.5 || dbl == 0.5 && (mant&1) != 0 {
		mant++
	}

	if mant == 0 {
		return 0, 0, 0, 0, nil // result should be zeroed out
	}

	if power < 0 {
		// Add -power factors of 10, -power <= (29 - 15) = 14.
		//
		power = -power
		if power < 10 {
			pow10 := sPowers10[power]
			low64 := uint64(uint64(mant) * uint64(pow10))
			hi64 := uint64(mant>>32) * uint64(pow10)
			low = uint32(low64)
			hi64 += low64 >> 32
			mid = uint32(hi64)
			hi64 >>= 32
			hi = uint32(hi64)
		} else {
			// Have a big power of 10.
			//
			low, mid, hi, _, err = uint64x64To128(mant, sUlongPowers10[power-1])
			if err != nil {
				return 0, 0, 0, 0, err
			}
		}
	} else {
		// Factor out powers of 10 to reduce the scale, if possible.
		// The maximum number we could factor out would be 14.  This
		// comes from the fact we have a 15-digit number, and the
		// MSD must be non-zero -- but the lower 14 digits could be
		// zero.  Note also the scale factor is never negative, so
		// we can't scale by any more than the power we used to
		// get the integer.
		//
		lmax := power
		if lmax > 14 {
			lmax = 14
		}

		if byte(mant) == 0 && lmax >= 8 {
			den := uint64(100000000)
			div := uint64(mant / den)
			if uint32(mant) == uint32(div*den) {
				mant = div
				power -= 8
				lmax -= 8
			}
		}

		if uint(mant&0xF) == 0 && lmax >= 4 {
			den := uint64(10000)
			div := uint64(mant / den)
			if uint32(mant) == uint32(div*den) {
				mant = div
				power -= 4
				lmax -= 4
			}
		}

		if uint32(mant&3) == 0 && lmax >= 2 {
			den := uint64(100)
			div := mant / den
			if uint32(mant) == uint32(div*den) {
				mant = div
				power -= 2
				lmax -= 2
			}
		}

		if uint32(mant&1) == 0 && lmax >= 1 {
			den := uint64(10)
			div := mant / den
			if uint32(mant) == uint32(div*den) {
				mant = div
				power--
			}
		}

		flags |= uint32(power) << 16
		low = uint32(mant)
		mid = uint32(mant >> 32)
	}

	return low, mid, hi, flags, nil
}
