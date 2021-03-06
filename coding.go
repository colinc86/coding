// Package coding contains structures for encoding and decoding values.
package coding

// A group of coding types.
const (
	codingTypeBool byte = 0x00

	codingTypeInt   byte = 0x01
	codingTypeInt64 byte = 0x02
	codingTypeInt32 byte = 0x03
	codingTypeInt16 byte = 0x04
	codingTypeInt8  byte = 0x05

	codingTypeUint   byte = 0x06
	codingTypeUint64 byte = 0x07
	codingTypeUint32 byte = 0x08
	codingTypeUint16 byte = 0x09
	codingTypeUint8  byte = 0x0A

	codingTypeFloat64 byte = 0x0B
	codingTypeFloat32 byte = 0x0C

	codingTypeString byte = 0x0D
	codingTypeData   byte = 0x0E
	codingTypeSlice  byte = 0x0F
)
