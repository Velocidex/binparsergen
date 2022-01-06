package binparsergen

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
	"unicode/utf16"
	"unicode/utf8"

	"gotest.tools/assert"
)

func TestUTFParser(t *testing.T) {
	en := []uint8{72, 0, 101, 0, 108, 0, 108, 0, 111, 0, 87, 0, 111, 0, 114, 0, 108, 0, 100, 0, 0, 0, 74, 0, 0, 0}
	zh := []uint8{96, 79, 125, 89, 22, 78, 76, 117, 0, 0, 0, 0, 74, 0, 0, 0}
	zh2 := []uint8{96, 79, 125, 89, 22, 78, 76, 117, 12}

	assert.Equal(t, ParseTerminatedUTF16String(bytes.NewReader(en), 0), "HelloWorld")
	assert.Equal(t, ParseTerminatedUTF16String(bytes.NewReader(zh), 0), "你好世界")
	assert.Equal(t, ParseTerminatedUTF16String(bytes.NewReader(zh2), 0), "你好世界")
}

func ParseTerminatedUTF16String(reader io.ReaderAt, offset int64) string {
	data := make([]byte, 1024)
	n, err := reader.ReadAt(data, offset)
	if err != nil && err != io.EOF {
		return ""
	}

	idx := bytes.Index(data[:n], []byte{0, 0})
	if idx < 0 {
		idx = n - 1
	}
	if idx%2 != 0 {
		idx += 1
	}
	return UTF16BytesToUTF8(data[0:idx], binary.LittleEndian)
}

func UTF16BytesToUTF8(b []byte, o binary.ByteOrder) string {
	if len(b) < 2 {
		return ""
	}

	if b[0] == 0xff && b[1] == 0xfe {
		o = binary.BigEndian
		b = b[2:]
	} else if b[0] == 0xfe && b[1] == 0xff {
		o = binary.LittleEndian
		b = b[2:]
	}

	utf := make([]uint16, (len(b)+(2-1))/2)

	for i := 0; i+(2-1) < len(b); i += 2 {
		utf[i/2] = o.Uint16(b[i:])
	}
	if len(b)/2 < len(utf) {
		utf[len(utf)-1] = utf8.RuneError
	}

	return string(utf16.Decode(utf))
}
