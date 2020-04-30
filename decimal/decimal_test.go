package decimal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToStringSmallInt(t *testing.T) {
	dec := FromInt32(177)
	require.Equal(t, "177", dec.String())
}

func TestToStringNegSmallInt(t *testing.T) {
	dec := FromInt32(-177)
	require.Equal(t, "-177", dec.String())
}

func TestToStringSmallFloat(t *testing.T) {
	dec, err := FromFloat64(177.92)
	require.NoError(t, err)
	require.Equal(t, "177.92", dec.String())
}

func TestToStringNegSmallFloat(t *testing.T) {
	dec, err := FromFloat64(-177.92)
	require.NoError(t, err)
	require.Equal(t, "-177.92", dec.String())
}

func TestToStringBigInt(t *testing.T) {
	dec := FromInt64(54654652135835)
	require.Equal(t, "54654652135835", dec.String())
}

func TestToStringNegBigInt(t *testing.T) {
	dec := FromInt64(-1651613254346)
	require.Equal(t, "-1651613254346", dec.String())
}

func TestToStringBigFloat(t *testing.T) {
	dec, err := FromFloat64(54656512546.5125465125)
	require.NoError(t, err)
	require.Equal(t, "54656512546.5125465125", dec.String())
}

func TestToStringNegBigFloat(t *testing.T) {
	dec, err := FromFloat64(-116767465441.67465441675)
	require.NoError(t, err)
	require.Equal(t, "-116767465441.67465441675", dec.String())
}

func TestToInt64(t *testing.T) {
	dec, err := FromFloat64(4545238794676.82347)
	require.NoError(t, err)
	require.Equal(t, int64(4545238794676), dec.Int64())
}

func TestToInt64NonFrac(t *testing.T) {
	dec := FromInt64(4545238794676)
	require.Equal(t, int64(4545238794676), dec.Int64())
}

func TestToInt64Neg(t *testing.T) {
	dec, err := FromFloat64(-4545238794676.82347)
	require.NoError(t, err)
	require.Equal(t, int64(-4545238794676), dec.Int64())
}

func TestToInt64NonFracNeg(t *testing.T) {
	dec := FromInt64(-4545238794676)
	require.Equal(t, int64(-4545238794676), dec.Int64())
}

func TestToInt32(t *testing.T) {
	dec, err := FromFloat64(238794676.82347)
	require.NoError(t, err)
	require.Equal(t, int32(238794676), dec.Int32())
}

func TestToInt32NonFrac(t *testing.T) {
	dec := FromInt64(238794676)
	require.Equal(t, int32(238794676), dec.Int32())
}

func TestToInt32Neg(t *testing.T) {
	dec, err := FromFloat64(-238794676.82347)
	require.NoError(t, err)
	require.Equal(t, int32(-238794676), dec.Int32())
}

func TestToInt32NonFracNeg(t *testing.T) {
	dec := FromInt64(-238794676)
	require.Equal(t, int32(-238794676), dec.Int32())
}
