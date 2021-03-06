package coding

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"hash/crc32"
	"math"
)

// Encoder types encode encode values to binary data.
type Encoder struct {

	// The encoder's data.
	data []byte
}

// Initializers

// NewEncoder creates a new encoder.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Boolean

// EncodeBool encodes a boolean.
func (e *Encoder) EncodeBool(b bool) {
	e.appendByte(codingTypeBool)
	if b {
		e.appendByte(1)
	} else {
		e.appendByte(0)
	}
}

// Integer

// EncodeInt encodes an integer.
func (e *Encoder) EncodeInt(n int) {
	e.appendByte(codingTypeInt)
	b := make([]byte, 8, 8)
	_ = binary.PutVarint(b, int64(n))
	e.appendBytes(b)
}

// EncodeInt64 encodes an integer.
func (e *Encoder) EncodeInt64(n int64) {
	e.appendByte(codingTypeInt64)
	b := make([]byte, 8, 8)
	_ = binary.PutVarint(b, n)
	e.appendBytes(b)
}

// EncodeInt32 encodes an integer.
func (e *Encoder) EncodeInt32(n int32) {
	e.appendByte(codingTypeInt32)
	b := make([]byte, 4, 4)
	_ = binary.PutVarint(b, int64(n))
	e.appendBytes(b)
}

// EncodeInt16 encodes an integer.
func (e *Encoder) EncodeInt16(n int16) {
	e.appendByte(codingTypeInt16)
	b := make([]byte, 2, 2)
	_ = binary.PutVarint(b, int64(n))
	e.appendBytes(b)
}

// EncodeInt8 encodes an integer.
func (e *Encoder) EncodeInt8(n int8) {
	e.appendByte(codingTypeInt8)
	b := make([]byte, 1, 1)
	_ = binary.PutVarint(b, int64(n))
	e.appendBytes(b)
}

// Unsigned integer

// EncodeUint encodes an integer.
func (e *Encoder) EncodeUint(n uint) {
	e.appendByte(codingTypeUint)
	b := make([]byte, 8, 8)
	_ = binary.PutUvarint(b, uint64(n))
	e.appendBytes(b)
}

// EncodeUint64 encodes an integer.
func (e *Encoder) EncodeUint64(n uint64) {
	e.appendByte(codingTypeUint64)
	b := make([]byte, 8, 8)
	_ = binary.PutUvarint(b, n)
	e.appendBytes(b)
}

// EncodeUint32 encodes an integer.
func (e *Encoder) EncodeUint32(n uint32) {
	e.appendByte(codingTypeUint32)
	b := make([]byte, 4, 4)
	_ = binary.PutUvarint(b, uint64(n))
	e.appendBytes(b)
}

// EncodeUint16 encodes an integer.
func (e *Encoder) EncodeUint16(n uint16) {
	e.appendByte(codingTypeUint16)
	b := make([]byte, 2, 2)
	_ = binary.PutUvarint(b, uint64(n))
	e.appendBytes(b)
}

// EncodeUint8 encodes an integer.
func (e *Encoder) EncodeUint8(n uint8) {
	e.appendByte(codingTypeUint8)
	b := make([]byte, 1, 1)
	_ = binary.PutUvarint(b, uint64(n))
	e.appendBytes(b)
}

// Floating point

// EncodeFloat64 encodes a float.
func (e *Encoder) EncodeFloat64(f float64) {
	// Float type
	e.appendByte(codingTypeFloat64)

	// Get bits
	bits := math.Float64bits(f)

	// Encode and get byte length
	b := make([]byte, 16, 16)
	n := binary.PutUvarint(b, bits)

	// Encode byte length
	e.appendByte(byte(n))

	// Append float bytes
	e.appendBytes(b[:n])
}

// EncodeFloat32 encodes a float.
func (e *Encoder) EncodeFloat32(f float32) {
	// Float type
	e.appendByte(codingTypeFloat32)

	// Get bits
	bits := math.Float32bits(f)

	// Encode and get byte length
	b := make([]byte, 16, 16)
	n := binary.PutUvarint(b, uint64(bits))

	// Encode byte length
	e.appendByte(byte(n))

	// Append float bytes
	e.appendBytes(b[:n])
}

// Data

// EncodeString encodes the string.
func (e *Encoder) EncodeString(s string) {
	e.appendByte(codingTypeString)

	b := make([]byte, 8, 8)
	_ = binary.PutVarint(b, int64(len(s)))
	e.appendBytes(b)

	e.appendBytes([]byte(s))
}

// EncodeData encodes the data.
func (e *Encoder) EncodeData(b []byte) {
	e.appendByte(codingTypeData)

	bytes := make([]byte, 8, 8)
	_ = binary.PutVarint(bytes, int64(len(b)))
	e.appendBytes(bytes)

	e.appendBytes(b)
}

// Exported methods

// Data returns the encoder's data along with trailing CRC data.
func (e Encoder) Data() []byte {
	return append(e.data, e.crcBytes()...)
}

// Flush clears the encoder's data.
func (e *Encoder) Flush() {
	e.data = nil
}

// Compress compresses the encoder's data and returns the result.
//
// Compress calls the encoder's Data function so that its data's CRC is included
// in the compressed bytes.
func (e *Encoder) Compress() ([]byte, error) {
	var cmb bytes.Buffer
	w := zlib.NewWriter(&cmb)

	if _, err := w.Write(e.Data()); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return cmb.Bytes(), nil
}

// Non-exported methods

// crcBytes calculates the CRC32 of the encoder's data and returns the bytes
// that should be added to the data.
func (e *Encoder) crcBytes() []byte {
	crc := crc32.ChecksumIEEE(e.data)
	b := make([]byte, 16, 16)
	n := binary.PutUvarint(b, uint64(crc))

	var o []byte
	o = append(o, b[:n]...)
	o = append(o, byte(n))
	return o
}

// appendByte appends a single byte to the encoder's data.
func (e *Encoder) appendByte(b byte) {
	e.data = append(e.data, b)
}

// appendBytes appends a slice of bytes to the encoder's data.
func (e *Encoder) appendBytes(b []byte) {
	e.data = append(e.data, b...)
}
