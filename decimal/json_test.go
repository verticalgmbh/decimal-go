package decimal

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// JsonObject - structure used for json tests
type JsonObject struct {
	Number Decimal
}

func TestJsonMarshal(t *testing.T) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.Encode(&JsonObject{
		Number: FromInt64Frac(7588, 3),
	})
	require.Equal(t, `{"Number":7.588}`+"\n", buffer.String())
}

func TestJsonUnmarshal(t *testing.T) {
	decoder := json.NewDecoder(strings.NewReader(`{"Number":7.588}`))
	result := &JsonObject{}
	decoder.Decode(result)
	require.Equal(t, 7.588, result.Number.Float64())
}
