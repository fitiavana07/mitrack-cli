package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoderV3String(t *testing.T) {
	encoder := NewEncoderV3()
	decoder := NewDecoderV3()

	s := "string-to-encode"
	b := new(bytes.Buffer)
	err := encoder.WriteEncoded(b, s)

	assert.NoError(t, err)

	var s2 string
	err = decoder.ReadDecoded(b, &s2)

	assert.NoError(t, err)
	assert.Equal(t, s, s2)
}

func TestEncoderV3Numeric(t *testing.T) {
	encoder := NewEncoderV3()
	decoder := NewDecoderV3()

	var v int64 = 5123
	b := new(bytes.Buffer)
	err := encoder.WriteEncoded(b, v)

	assert.NoError(t, err)

	var v2 int64
	err = decoder.ReadDecoded(b, &v2)

	assert.NoError(t, err)
	assert.Equal(t, v, v2)
}
