package coding

import (
	"fmt"
	"math"
	"testing"
)

// Examples

func ExampleDecoder_Decompress() {
	e := NewEncoder()
	e.EncodeString("{ \"name\": \"pi\" }")
	e.EncodeFloat64(math.Pi)
	e.EncodeString("{ \"name\": \"phi\" }")
	e.EncodeFloat64(math.Phi)
	e.EncodeString("{ \"name\": \"e\" }")
	e.EncodeFloat64(math.E)
	e.EncodeString("{ \"name\": \"ln(2)\" }")
	e.EncodeFloat64(math.Ln2)

	fmt.Printf("Bytes: %d\n", len(e.Data()))

	cd, err := e.Compress()
	if err != nil {
		fmt.Printf("Error compressing data: %s\n", err)
		return
	}

	fmt.Printf("Compressed bytes: %d\n", len(cd))
	fmt.Printf("Change: %0.f%%\n", 100.0*float64(len(cd)-len(e.data))/float64(len(e.data)))

	d := NewDecoder(cd)
	if err = d.Decompress(); err != nil {
		fmt.Printf("Error decompressing data: %s\n", err)
		return
	}

	if err = d.Validate(); err != nil {
		fmt.Printf("Error validating data: %s\n", err)
		return
	}

	var s string
	if s, err = d.DecodeString(); err != nil {
		fmt.Printf("Error decoding string: %s\n", err)
		return
	}
	fmt.Printf("String: %s\n", s)

	var f float64
	if f, err = d.DecodeFloat64(); err != nil {
		fmt.Printf("Error decoding float: %s\n", err)
		return
	}
	fmt.Printf("Float: %f\n", f)

	// Output:
	// Bytes: 153
	// Compressed bytes: 106
	// Change: -28%
	// String: { "name": "pi" }
	// Float: 3.141593
}

// Encoder

func TestFlush(t *testing.T) {
	e := NewEncoder()
	e.EncodeString("Hello, World!")

	if len(e.data) == 0 {
		t.Fatal("Expected data.")
	}

	e.Flush()
	if len(e.data) > 0 {
		t.Errorf("Expected no data but received %d bytes.\n", len(e.data))
	}
}

// Decoder

func TestEndOfBufferError(t *testing.T) {
	e := NewEncoder()
	e.EncodeString("Hello, World!")
	e.data = e.data[:len(e.data)-2]

	d := NewDecoder(e.data)
	if s, err := d.DecodeString(); err != nil {
		if err != ErrEOB {
			t.Fatalf("Expected end of buffer error but received: %s\n", err)
		}
	} else {
		t.Errorf("Expected end of buffer error but received string %s.\n", s)
	}
}

func TestTypeMismatch(t *testing.T) {
	e := NewEncoder()
	e.EncodeBool(true)

	d := NewDecoder(e.data)
	if s, err := d.DecodeString(); err != nil {
		if err != ErrType {
			t.Fatalf("Expected a type mismatch error but received: %s\n", err)
		}
	} else {
		t.Errorf("Expected a type mismatch error but received string: %s.\n", s)
	}
}

// Bool

func TestEncodeDecodeBool_1(t *testing.T) {
	var i bool = true
	testEncodeDecode(i, t)
}

func TestEncodeDecodeBool_2(t *testing.T) {
	var i bool = false
	testEncodeDecode(i, t)
}

// Int

func TestEncodeDecodeInt_1(t *testing.T) {
	var i int = 0
	testEncodeDecode(i, t)
}

func TestEncodeDecodeInt_2(t *testing.T) {
	var i int = -1
	testEncodeDecode(i, t)
}

func TestEncodeDecodeInt_3(t *testing.T) {
	var i int = 10
	testEncodeDecode(i, t)
}

// Int64

func TestEncodeDecodeInt64_1(t *testing.T) {
	var i int64 = 0
	testEncodeDecode(i, t)
}

func TestEncodeDecodeInt64_2(t *testing.T) {
	var i int64 = -1
	testEncodeDecode(i, t)
}

func TestEncodeDecodeInt64_3(t *testing.T) {
	var i int64 = 10
	testEncodeDecode(i, t)
}

// Int32

func TestEncodeDecodeInt32_1(t *testing.T) {
	var i int32 = 0
	testEncodeDecode(i, t)
}

func TestEncodeDecodeInt32_2(t *testing.T) {
	var i int32 = -1
	testEncodeDecode(i, t)
}

func TestEncodeDecodeInt32_3(t *testing.T) {
	var i int32 = 10
	testEncodeDecode(i, t)
}

// Int16

func TestEncodeDecodeInt16_1(t *testing.T) {
	var i int16 = 0
	testEncodeDecode(i, t)
}

func TestEncodeDecodeInt16_2(t *testing.T) {
	var i int16 = -1
	testEncodeDecode(i, t)
}

func TestEncodeDecodeInt16_3(t *testing.T) {
	var i int16 = 10
	testEncodeDecode(i, t)
}

// Int8

func TestEncodeDecodeInt8_1(t *testing.T) {
	var i int8 = 0
	testEncodeDecode(i, t)
}

func TestEncodeDecodeInt8_2(t *testing.T) {
	var i int8 = -1
	testEncodeDecode(i, t)
}

func TestEncodeDecodeInt8_3(t *testing.T) {
	var i int8 = 10
	testEncodeDecode(i, t)
}

// Uint

func TestEncodeDecodeUint_1(t *testing.T) {
	var i uint = 0
	testEncodeDecode(i, t)
}

func TestEncodeDecodeUint_2(t *testing.T) {
	var i uint = 10
	testEncodeDecode(i, t)
}

// Uint64

func TestEncodeDecodeUint64_1(t *testing.T) {
	var i uint64 = 0
	testEncodeDecode(i, t)
}

func TestEncodeDecodeUint64_2(t *testing.T) {
	var i uint64 = 10
	testEncodeDecode(i, t)
}

// Uint32

func TestEncodeDecodeUint32_1(t *testing.T) {
	var i uint32 = 0
	testEncodeDecode(i, t)
}

func TestEncodeDecodeUint32_2(t *testing.T) {
	var i uint32 = 10
	testEncodeDecode(i, t)
}

// Uint16

func TestEncodeDecodeUint16_1(t *testing.T) {
	var i uint16 = 0
	testEncodeDecode(i, t)
}

func TestEncodeDecodeUint16_2(t *testing.T) {
	var i uint16 = 10
	testEncodeDecode(i, t)
}

// Uint8

func TestEncodeDecodeUint8_1(t *testing.T) {
	var i uint8 = 0
	testEncodeDecode(i, t)
}

func TestEncodeDecodeUint8_2(t *testing.T) {
	var i uint8 = 10
	testEncodeDecode(i, t)
}

// Float64

func TestEncodeDecodeFloat64_1(t *testing.T) {
	var i float64 = 0.0
	testEncodeDecode(i, t)
}

func TestEncodeDecodeFloat64_2(t *testing.T) {
	var i float64 = -math.Pi
	testEncodeDecode(i, t)
}

func TestEncodeDecodeFloat64_3(t *testing.T) {
	var i float64 = math.Pi
	testEncodeDecode(i, t)
}

// Float32

func TestEncodeDecodeFloat32_1(t *testing.T) {
	var i float32 = 0.0
	testEncodeDecode(i, t)
}

func TestEncodeDecodeFloat32_2(t *testing.T) {
	var i float32 = -math.Pi
	testEncodeDecode(i, t)
}

func TestEncodeDecodeFloat32_3(t *testing.T) {
	var i float32 = math.Pi
	testEncodeDecode(i, t)
}

// String

func TestEncodeDecodeString_1(t *testing.T) {
	var i string = ""
	testEncodeDecode(i, t)
}

func TestEncodeDecodeString_2(t *testing.T) {
	var i string = "Hello, World!"
	testEncodeDecode(i, t)
}

// Data

func TestEncodeDecodeData_1(t *testing.T) {
	var i []byte
	testEncodeDecode(i, t)
}

func TestEncodeDecodeData_2(t *testing.T) {
	var i []byte = []byte{0x00, 0x01, 0x02, 0x03}
	testEncodeDecode(i, t)
}

// CRC

func TestValidCRC_1(t *testing.T) {
	e := NewEncoder()
	e.EncodeString("")

	d := NewDecoder(e.Data())
	if err := d.Validate(); err != nil {
		t.Errorf("CRC check failed: %s\n", err)
	}
}

func TestValidCRC_2(t *testing.T) {
	e := NewEncoder()
	e.EncodeString("Hello, World!")

	d := NewDecoder(e.Data())
	if err := d.Validate(); err != nil {
		t.Errorf("CRC check failed: %s\n", err)
	}
}

func TestInalidCRC(t *testing.T) {
	e := NewEncoder()
	e.EncodeString("")

	d := NewDecoder(e.data)
	if err := d.Validate(); err == nil {
		t.Error("CRC check should have failed.\n", err)
	}
}

// Compression

func TestCompressDecompress_1(t *testing.T) {
	testCompressDecompress("", t)
}

func TestCompressDecompress_2(t *testing.T) {
	testCompressDecompress("Hello, World!", t)
}

// Non-exported functions

// testEncodeDecode attempts to encode the input value and then decode it.
func testEncodeDecode(i interface{}, t *testing.T) {
	e := NewEncoder()

	switch i.(type) {
	case bool:
		e.EncodeBool(i.(bool))
	case int:
		e.EncodeInt(i.(int))
	case int64:
		e.EncodeInt64(i.(int64))
	case int32:
		e.EncodeInt32(i.(int32))
	case int16:
		e.EncodeInt16(i.(int16))
	case int8:
		e.EncodeInt8(i.(int8))
	case uint:
		e.EncodeUint(i.(uint))
	case uint64:
		e.EncodeUint64(i.(uint64))
	case uint32:
		e.EncodeUint32(i.(uint32))
	case uint16:
		e.EncodeUint16(i.(uint16))
	case uint8:
		e.EncodeUint8(i.(uint8))
	case float64:
		e.EncodeFloat64(i.(float64))
	case float32:
		e.EncodeFloat32(i.(float32))
	case string:
		e.EncodeString(i.(string))
	case []byte:
		e.EncodeData(i.([]byte))
	}

	d := NewDecoder(e.Data())

	var o interface{}
	var err error

	switch i.(type) {
	case bool:
		o, err = d.DecodeBool()
	case int:
		o, err = d.DecodeInt()
	case int64:
		o, err = d.DecodeInt64()
	case int32:
		o, err = d.DecodeInt32()
	case int16:
		o, err = d.DecodeInt16()
	case int8:
		o, err = d.DecodeInt8()
	case uint:
		o, err = d.DecodeUint()
	case uint64:
		o, err = d.DecodeUint64()
	case uint32:
		o, err = d.DecodeUint32()
	case uint16:
		o, err = d.DecodeUint16()
	case uint8:
		o, err = d.DecodeUint8()
	case float64:
		o, err = d.DecodeFloat64()
	case float32:
		o, err = d.DecodeFloat32()
	case string:
		o, err = d.DecodeString()
	case []byte:
		o, err = d.DecodeData()
	}

	if err != nil {
		t.Fatalf("Error decoding type: %s\n", err)
	}

	switch i.(type) {
	case []byte:
		ib := i.([]byte)
		ob := o.([]byte)

		if len(ib) == len(ob) {
			for j, v := range ib {
				if v != ob[j] {
					t.Fatalf("Mismatched values %d and %d.\n", v, ob[j])
				}
			}
		} else {
			t.Fatalf("Data has unequal length %d != %d.\n", len(ib), len(ob))
		}
	default:
		if i != o {
			t.Fatalf("Expected output %v to match input %v.\n", o, i)
		}
	}
}

func testCompressDecompress(s string, t *testing.T) {
	e := NewEncoder()
	e.EncodeString(s)
	dc := append([]byte{}, e.Data()...)
	cd, err := e.Compress()
	if err != nil {
		t.Fatalf("Unable to compress data: %s\n", err)
	}

	d := NewDecoder(cd)
	if err = d.Decompress(); err != nil {
		t.Fatalf("Unable to decompress data: %s\n", err)
	}

	if len(dc) != len(d.data) {
		t.Fatalf("Unequal data lengths %d and %d.\n", len(cd), len(d.data))
	}

	for i, v := range dc {
		if v != d.data[i] {
			t.Fatalf("Mismatched data values %d and %d.\n", v, d.data[i])
		}
	}
}
