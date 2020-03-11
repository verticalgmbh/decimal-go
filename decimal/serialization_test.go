package decimal

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WriteAndReadDecimal(t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	err := WriteDecimal(writer, FromInt64Frac(99299, 3))
	require.NoError(t, err)
	writer.Flush()

	reader := bufio.NewReader(&buffer)
	value, err := ReadDecimal(reader)
	require.NoError(t, err)
	require.Equal(t, 99.299, value.Float64())
}

func Test_WriteAndReadNegativeDecimal(t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	err := WriteDecimal(writer, FromInt64Frac(-99299, 3))
	require.NoError(t, err)
	writer.Flush()

	reader := bufio.NewReader(&buffer)
	value, err := ReadDecimal(reader)
	require.NoError(t, err)
	require.Equal(t, -99.299, value.Float64())
}

func Test_WriteAndReadIntAsDecimal(t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	err := WriteDecimal(writer, FromInt64(167556422498554))
	require.NoError(t, err)
	writer.Flush()

	reader := bufio.NewReader(&buffer)
	value, err := ReadDecimal(reader)
	require.NoError(t, err)
	require.Equal(t, float64(167556422498554), value.Float64())
}

func Test_WriteAndReadNegativeIntAsDecimal(t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	err := WriteDecimal(writer, FromInt64(-167556422498554))
	require.NoError(t, err)
	writer.Flush()

	reader := bufio.NewReader(&buffer)
	value, err := ReadDecimal(reader)

	require.NoError(t, err)
	require.Equal(t, float64(-167556422498554), value.Float64())
}

func Test_WriteAndReadUIntAsDecimal(t *testing.T) {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	err := WriteDecimal(writer, FromUInt64(656543214884244))
	require.NoError(t, err)
	writer.Flush()

	reader := bufio.NewReader(&buffer)
	value, err := ReadDecimal(reader)

	require.NoError(t, err)
	require.Equal(t, float64(656543214884244), value.Float64())
}
