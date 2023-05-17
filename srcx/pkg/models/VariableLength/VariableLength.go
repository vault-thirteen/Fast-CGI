package vl

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/vault-thirteen/auxie/reader"
)

const (
	MaxVariableLength           = uint32(0x7FFF_FFFF) // = (1^32-1) >> 1.
	MaxSingleByteVariableLength = 127                 // = 255 >> 1.
	FourByteMarkerMask          = uint32(0x8000_0000) // = 1^31.
)

const (
	RawValueSizeOneByte   = 1
	RawValueSizeFourBytes = 4
)

const (
	ErrOverflow = "overflow"
)

// VariableLength is a variable length (size) in accordance to the FastCGI
// Specification Version 1.1.
type VariableLength struct {
	// Real computer integer value of the length.
	realValue uint32

	// Number of bytes used for storing the raw value.
	// Depending on this selector, the real raw value can be found in one or
	// another raw value fields: rawValue1B or rawValue4B.
	rawValueSize byte

	// Raw value when it is one byte long.
	rawValue1B byte

	// Raw value when it is four bytes long.
	rawValue4B uint32
}

// NewFromStream reads a variable length from the stream.
func NewFromStream(stream io.Reader) (vl *VariableLength, err error) {
	rdr := reader.New(stream)

	var b byte
	b, err = rdr.ReadByte()
	if err != nil {
		return nil, err
	}

	// When the greatest bit is 0, size must be a single byte.
	if (b >> 7) == 0 { // 127 or less is stored in a single byte.
		vl = &VariableLength{
			rawValueSize: RawValueSizeOneByte,
			rawValue1B:   b,
			realValue:    uint32(b),
		}
		return vl, nil
	}

	// When the greatest bit is 1, size must be four bytes.
	var buf = make([]byte, 4)
	buf[0] = b

	b, err = rdr.ReadByte()
	if err != nil {
		return nil, err
	}
	buf[1] = b

	b, err = rdr.ReadByte()
	if err != nil {
		return nil, err
	}
	buf[2] = b

	b, err = rdr.ReadByte()
	if err != nil {
		return nil, err
	}
	buf[3] = b

	vl = &VariableLength{
		rawValueSize: RawValueSizeFourBytes,
		rawValue4B:   binary.BigEndian.Uint32(buf),
	}
	vl.realValue = vl.rawValue4B & MaxVariableLength

	return vl, nil
}

// NewFromInteger creates a variable length from a specified integer value.
func NewFromInteger(n uint) (vl *VariableLength, err error) {
	if n > uint(MaxVariableLength) {
		return nil, errors.New(ErrOverflow)
	}

	if n <= uint(MaxSingleByteVariableLength) {
		vl = &VariableLength{
			realValue:    uint32(n),
			rawValueSize: RawValueSizeOneByte,
			rawValue1B:   byte(n),
		}
		return vl, nil
	}

	vl = &VariableLength{
		realValue:    uint32(n),
		rawValueSize: RawValueSizeFourBytes,
		rawValue4B:   uint32(n) | FourByteMarkerMask,
	}
	return vl, nil
}

// ToBytes returns the bytes storing the variable length (size) in accordance
// to the FastCGI Specification Version 1.1. In other words, it returns the raw
// value as an array (slice) of bytes.
func (vl *VariableLength) ToBytes() (ba []byte) {
	switch vl.rawValueSize {
	case RawValueSizeOneByte:
		return []byte{vl.rawValue1B}
	case RawValueSizeFourBytes:
		var buf4B = []byte{0, 0, 0, 0}
		binary.BigEndian.PutUint32(buf4B, vl.rawValue4B)
		return buf4B
	default:
		// This case is normally unreachable. It can be reached only if the
		// hardware's memory is corrupted by gremlins.
		panic(vl.rawValueSize)
	}
}

// RealValue return the real value of the length as an uint32.
func (vl *VariableLength) RealValue() (realValue uint32) {
	return vl.realValue
}

// Size returns the number of bytes required to store the raw value.
func (vl *VariableLength) Size() (bytesCount int) {
	return int(vl.rawValueSize)
}
