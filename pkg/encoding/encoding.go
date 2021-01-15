package encoding

import "io"

// Encoding versions.
const (
	FormatVersionCurrent uint32 = FormatVersionV3

	FormatVersionV3 uint32 = 3
	FormatVersionV2 uint32 = 2
	FormatVersionV1 uint32 = 1
)

// Encoder wraps the Encode method.
// Encode encodes the given data and writes its encoded value in dst.
type Encoder interface {
	WriteEncoded(w io.Writer, data interface{}) error
}

// Decoder wraps the Decode method.
// Decode decodes the given bytes and writes
type Decoder interface {
	ReadDecoded(r io.Reader, data interface{}) error
}
