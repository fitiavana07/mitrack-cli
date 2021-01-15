package encoding

import (
	"encoding/binary"
	"io"
)

// EncoderV3 is the encoder used in v3.
type EncoderV3 struct {
}

// DecoderV3 is the decoder used in v3.
type DecoderV3 struct {
}

// NewEncoderV3 returns a new v3 encoder.
func NewEncoderV3() Encoder {
	return &EncoderV3{}
}

// NewDecoderV3 returns a new v3 decoder.
func NewDecoderV3() Decoder {
	return &DecoderV3{}
}

// WriteEncoded writes encoded bytes into w.
func (e *EncoderV3) WriteEncoded(w io.Writer, data interface{}) error {
	switch data.(type) {
	case string:
		str := data.(string)
		if err := e.writeString(w, str); err != nil {
			return err
		}
	// TODO support marshal interface
	default:
		if err := e.writeNumeric(w, data); err != nil {
			return err
		}
	}
	return nil
}

// writeString writes an encoded string into w.
// It uses the format: len value, where len is an uint16.
// Thus, it only supports strings up to 2^16 long.
func (e *EncoderV3) writeString(w io.Writer, str string) error {
	if err := e.writeNumeric(w, uint16(len(str))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(str)); err != nil {
		return err
	}
	return nil
}

// writeNumeric is used to encode numeric values in v3.
// It does not check whether the provided value is a number.
func (e *EncoderV3) writeNumeric(w io.Writer, num interface{}) error {
	return binary.Write(w, binary.LittleEndian, num)
}

// ReadDecoded decodes the given data using v3 encoding,
// and writes the value into data.
func (d *DecoderV3) ReadDecoded(r io.Reader, data interface{}) error {
	switch data.(type) {
	case *string:
		ptr := data.(*string)
		if err := d.readString(r, ptr); err != nil {
			return err
		}
	// TODO support unmarshal interface
	default:
		if err := d.readNumeric(r, data); err != nil {
			return err
		}
	}

	return nil
}

// readString reads string from the next bytes of r.
func (d *DecoderV3) readString(r io.Reader, ptr *string) error {
	var length uint16
	if err := d.readNumeric(r, &length); err != nil {
		return err
	}

	b := make([]byte, length)
	if _, err := r.Read(b); err != nil {
		return err
	}

	*ptr = string(b)
	return nil
}

// readNumeric reads numeric from the next bytes of r.
func (d *DecoderV3) readNumeric(r io.Reader, num interface{}) error {
	return binary.Read(r, binary.LittleEndian, num)
}
