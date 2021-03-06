package coding

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"math"
)

var (
	// ErrEOB is an end of buffer error.
	ErrEOB error = errors.New("end of buffer")

	// ErrType is an incorrect type error.
	ErrType error = errors.New("incorrect type")

	// ErrByteLength is an incorrect byte length error.
	ErrByteLength error = errors.New("incorrect byte length")

	// ErrCRC is a CRC check error.
	ErrCRC = errors.New("crc check failed")
)

// Decoder types decode bytes and keep track of an offset.
type Decoder struct {
	data   []byte
	offset int
}

// NewDecoder creates and returns a new decoder with the given data.
func NewDecoder(data []byte) *Decoder {
	return &Decoder{
		data: data,
	}
}

// Boolean

// DecodeBool decodes the next value as a boolean.
func (d *Decoder) DecodeBool() (bool, error) {
	if err := d.checkType(codingTypeBool); err != nil {
		return false, err
	}

	if !d.checkLength(1) {
		return false, ErrEOB
	}

	bb := d.getByte()
	return bb == 1, nil
}

// Integer

// DecodeInt decodes the next value as an integer.
func (d *Decoder) DecodeInt() (int, error) {
	if err := d.checkType(codingTypeInt); err != nil {
		return 0, err
	}

	i, err := d.decodeInt64(8)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

// DecodeInt64 decodes the next value as an integer.
func (d *Decoder) DecodeInt64() (int64, error) {
	if err := d.checkType(codingTypeInt64); err != nil {
		return 0, err
	}

	i, err := d.decodeInt64(8)
	if err != nil {
		return 0, err
	}

	return i, nil
}

// DecodeInt32 decodes the next value as an integer.
func (d *Decoder) DecodeInt32() (int32, error) {
	if err := d.checkType(codingTypeInt32); err != nil {
		return 0, err
	}

	i, err := d.decodeInt64(4)
	if err != nil {
		return 0, err
	}

	return int32(i), nil
}

// DecodeInt16 decodes the next value as an integer.
func (d *Decoder) DecodeInt16() (int16, error) {
	if err := d.checkType(codingTypeInt16); err != nil {
		return 0, err
	}

	i, err := d.decodeInt64(2)
	if err != nil {
		return 0, err
	}

	return int16(i), nil
}

// DecodeInt8 decodes the next value as an integer.
func (d *Decoder) DecodeInt8() (int8, error) {
	if err := d.checkType(codingTypeInt8); err != nil {
		return 0, err
	}

	i, err := d.decodeInt64(1)
	if err != nil {
		return 0, err
	}

	return int8(i), nil
}

// Unsigned integer

// DecodeUint decodes the next value as an integer.
func (d *Decoder) DecodeUint() (uint, error) {
	if err := d.checkType(codingTypeUint); err != nil {
		return 0, err
	}

	i, err := d.decodeUint64(8)
	if err != nil {
		return 0, err
	}

	return uint(i), nil
}

// DecodeUint64 decodes the next value as an integer.
func (d *Decoder) DecodeUint64() (uint64, error) {
	if err := d.checkType(codingTypeUint64); err != nil {
		return 0, err
	}

	i, err := d.decodeUint64(8)
	if err != nil {
		return 0, err
	}

	return uint64(i), nil
}

// DecodeUint32 decodes the next value as an integer.
func (d *Decoder) DecodeUint32() (uint32, error) {
	if err := d.checkType(codingTypeUint32); err != nil {
		return 0, err
	}

	i, err := d.decodeUint64(4)
	if err != nil {
		return 0, err
	}

	return uint32(i), nil
}

// DecodeUint16 decodes the next value as an integer.
func (d *Decoder) DecodeUint16() (uint16, error) {
	if err := d.checkType(codingTypeUint16); err != nil {
		return 0, err
	}

	i, err := d.decodeUint64(2)
	if err != nil {
		return 0, err
	}

	return uint16(i), nil
}

// DecodeUint8 decodes the next value as an integer.
func (d *Decoder) DecodeUint8() (uint8, error) {
	if err := d.checkType(codingTypeUint8); err != nil {
		return 0, err
	}

	i, err := d.decodeUint64(1)
	if err != nil {
		return 0, err
	}

	return uint8(i), nil
}

// Floating point

// DecodeFloat64 decodes the next value as a floating point number.
func (d *Decoder) DecodeFloat64() (float64, error) {
	if err := d.checkType(codingTypeFloat64); err != nil {
		return 0, err
	}

	if !d.checkLength(1) {
		return 0.0, ErrEOB
	}

	n, err := d.decodeUint64(1)
	if err != nil {
		return 0.0, err
	}

	i, err := d.decodeUint64(int(n))
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(i), nil
}

// DecodeFloat32 decodes the next value as a floating point number.
func (d *Decoder) DecodeFloat32() (float32, error) {
	if err := d.checkType(codingTypeFloat32); err != nil {
		return 0, err
	}

	if !d.checkLength(1) {
		return 0.0, ErrEOB
	}

	n, err := d.decodeUint64(1)
	if err != nil {
		return 0.0, err
	}

	i, err := d.decodeUint64(int(n))
	if err != nil {
		return 0, err
	}

	return math.Float32frombits(uint32(i)), nil
}

// Data

// DecodeString decodes the next value as a string.
func (d *Decoder) DecodeString() (string, error) {
	if err := d.checkType(codingTypeString); err != nil {
		return "", err
	}

	if !d.checkLength(8) {
		return "", ErrEOB
	}

	l, err := d.decodeInt64(8)
	if err != nil {
		return "", err
	}

	if !d.checkLength(int(l)) {
		return "", ErrEOB
	}

	if l == 0 {
		return "", nil
	}
	return string(d.getBytes(int(l))), nil
}

// DecodeData decodes the next value as a byte array.
func (d *Decoder) DecodeData() ([]byte, error) {
	if err := d.checkType(codingTypeData); err != nil {
		return nil, err
	}

	l, err := d.decodeInt64(8)
	if err != nil {
		return nil, err
	}

	if !d.checkLength(int(l)) {
		return nil, ErrEOB
	}

	if l == 0 {
		return nil, nil
	}
	return d.getBytes(int(l)), nil
}

// Exported methods

// Decompress decompresses the decoder's data and places the result in data.
func (d *Decoder) Decompress() error {
	b := bytes.NewBuffer(d.data)
	r, err := zlib.NewReader(b)
	if err != nil {
		return err
	}

	ob := new(bytes.Buffer)
	if _, err := ob.ReadFrom(r); err != nil {
		return err
	}

	if err := r.Close(); err != nil {
		return err
	}

	d.data = ob.Bytes()
	return nil
}

// Validate validates the decoder's data by calculating its CRC32 and comparing.
func (d Decoder) Validate() error {
	if len(d.data) < 1 {
		return ErrByteLength
	}

	l := int(d.data[len(d.data)-1:][0])

	r := bytes.NewReader(d.data[len(d.data)-l-1 : len(d.data)-1])
	i, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}

	crc := uint32(i)
	ccrc := crc32.ChecksumIEEE(d.data[:len(d.data)-l-1])
	if crc != ccrc {
		return ErrCRC
	}
	return nil
}

// Non-exported methods

// checkType checks the given type against the next type byte in the decoder's
// data.
//
// If the type given matches the type byte, then the current byte offset is
// increased by one, otherwise it is not changed.
func (d *Decoder) checkType(t byte) error {
	// Can we get the type byte?
	if !d.checkLength(1) {
		return ErrEOB
	}

	// Get it and check
	tb := d.getByte()
	if tb == t {
		return nil
	}

	// Since they weren't the same, undo the offset change from getByte
	d.decrementOffset(1)
	return ErrType
}

// decodeInt64 decodes the next byteLength bytes in to an int64 value.
func (d *Decoder) decodeInt64(byteLength int) (int64, error) {
	r, err := d.getIntByteReader(byteLength)
	if err != nil {
		return 0, err
	}
	return binary.ReadVarint(r)
}

// decodeUint64 decodes the next byteLength bytes in to an uint64 value.
func (d *Decoder) decodeUint64(byteLength int) (uint64, error) {
	r, err := d.getIntByteReader(byteLength)
	if err != nil {
		return 0, err
	}
	return binary.ReadUvarint(r)
}

// getIntByteReader creates a new byte reader with the next byteLength bytes.
func (d *Decoder) getIntByteReader(byteLength int) (*bytes.Reader, error) {
	if !d.checkLength(byteLength) {
		return nil, ErrEOB
	}

	ib := d.getBytes(byteLength)
	return bytes.NewReader(ib), nil
}

// checkLength checks that there are enough bytes in the buffer from the
// decoder's offset to satisfy the given length.
func (d *Decoder) checkLength(l int) bool {
	return d.offset+l-1 < len(d.data)
}

// getByte gets the next byte at offset and increments offset.
func (d *Decoder) getByte() byte {
	defer d.incrementOffset(1)
	return d.data[d.offset]
}

// getBytes gets the next n bytes at offset and increments offset by n.
func (d *Decoder) getBytes(n int) []byte {
	defer d.incrementOffset(n)
	return d.data[d.offset : d.offset+n]
}

// incrementOffset increments offset by n.
func (d *Decoder) incrementOffset(n int) {
	d.offset += n
}

// decrementOffset decrements offset by n.
func (d *Decoder) decrementOffset(n int) {
	d.offset -= n
}
