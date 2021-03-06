# Package coding

[![Go Tests](https://github.com/colinc86/coding/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/colinc86/coding/actions/workflows/go-test.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/colinc86/coding.svg)](https://pkg.go.dev/github.com/colinc86/coding)

Package coding contains structures for encoding and decoding values.

## Installing

Navigate to your module and execute the following.

```bash
$ go get github.com/colinc86/coding
```

## Using

Import the package:

```go
import "github.com/colinc86/coding"
```

### Encoding Data

The `Encoder` type is responsible for encoding values.

#### Creating an Encoder

Calling `NewEncoder` creates a new encoder with an empty internal buffer.

```go
e := coding.NewEncoder()
```

#### Encoding

The `Encoder` type supports encoding:

- [x] `bool`
- [x] `int`, `int64`, `int32`, `int16`, `int8`
- [x] `uint`, `uint64`, `uint32`, `uint16`, `uint8`
- [x] `float64`, `float32`
- [x] `string`
- [x] `[]byte`

 The order that you encode values is the order that they must be decoded with a `Decoder`.

```go
e.EncodeBool(true)
e.EncodeString("Hello, World!")
e.EncodeInt(42)

d, err := json.Marshal(someStruct)
e.EncodeData(d)
```

#### Flushing Data

If you need to start over, you can call `Flush` on the encoder to clear its internal buffer.

```go
e.Flush()
```

#### Getting Encoded Data

Use the `Data` function to get the encoder's encoded data. This method calculates the encoded data's CRC32 and appends the bytes to the end of the encoded data before returning it so it may be verified by a decoder.

```go
encodedData := e.Data()
```

#### Compressing

You can, optionally, call the `Compress` function which will return a compressed version of the data returned by calling `Data` on the encoder.

```go
compressedData, err := e.Compress()
```

### Decoding

The `Decoder` type is responsible for decoding values.

#### Creating a Decoder

Call `NewDecoder` and give it the data to decode.

```go
d := NewDecoder(compressedData)
```

#### Decompressing

If you obtained data from an encoder by calling `Compress`, then the first method you must call on the decoder is `Decompress`. Otherwise, this step is unnecessary and will return an error.

Calling `Decompress` decompresses the decoder's internal buffer so there is no need to keep track of new slices.

```go
err := d.Decompress()
```

#### Validating Data

It is advised that you validate the decoder's data. Validating checks the CRC bytes that were appened to the encoder's data when calling `Data`.

```go
if err := d.Validate; err != nil {
	fmt.Printf("Invalid data: %s\n", err)
}
```

#### Decoding Data

The `Decoder` type supports decoding all of the same types as encoders. As mentioned above, you must decode values in the order they were decoded.

```go
// You should catch errors, but for the sake of brevity...
boolValue, _ := d.DecodeBool()
stringValue, _ := d.DecodeString()
intValue, _ := d.DecodeInt()
jsonData, _ := d.DecodeData()

someStruct = new(SomeStruct)
_ := json.Unmarshal(jsonData, someStruct)
```

## Example

```go
package main

import (
	"fmt"

	"github.com/colinc86/coding"
)

func main() {
	e := coding.NewEncoder()
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

	d := coding.NewDecoder(cd)
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
}
```

Output:

```
Bytes: 153
Compressed bytes: 106
Change: -28%
String: { "name": "pi" }
Float: 3.141593
```
