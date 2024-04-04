package binparsergen

import "fmt"

type StringParser struct {
	BaseParser
	Length uint64 `json:"length,omitempty"`
}

func (self StringParser) Prototype() string {
	return `
func ParseTerminatedString(reader io.ReaderAt, offset int64) string {
   var buf [1024]byte
   data := buf[:]
   n, err := reader.ReadAt(data, offset)
   if err != nil && err != io.EOF {
     return ""
   }
   idx := bytes.Index(data[:n], []byte{0})
   if idx < 0 {
      idx = n
   }
   return string(data[0:idx])
}

func ParseString(reader io.ReaderAt, offset int64, length int64) string {
   data := make([]byte, length)
   n, err := reader.ReadAt(data, offset)
   if err != nil && err != io.EOF {
      return ""
   }
   return string(data[:n])
}

`
}

func (self StringParser) PrototypeName() string {
	return "String"
}

func (self *StringParser) Compile(struct_name string, field_name string) string {
	if self.Length == 0 {
		return fmt.Sprintf(`

func (self *%[1]s) %[2]s() string {
  return ParseTerminatedString(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
}
`, struct_name, field_name)
	}

	return fmt.Sprintf(`

func (self *%[1]s) %[2]s() string {
  return ParseString(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset, %[3]v)
}
`, struct_name, field_name, self.Length)
}

func (self StringParser) GoType() string {
	return "string"
}

func (self StringParser) Size(value string) string {
	return "0"
}

type UTF16StringParser struct {
	BaseParser
	Length uint64 `json:"length,omitempty"`
}

func (self UTF16StringParser) Prototype() string {
	return `
func ParseTerminatedUTF16String(reader io.ReaderAt, offset int64) string {
   var buf [1024]byte
   data := buf[:]
   n, err := reader.ReadAt(data, offset)
   if err != nil && err != io.EOF {
     return ""
   }

   idx := bytes.Index(data[:n], []byte{0, 0})
   if idx < 0 {
      idx = n-1
   }
   if idx%2 != 0 {
      idx += 1
   }
   return UTF16BytesToUTF8(data[0:idx], binary.LittleEndian)
}

func ParseUTF16String(reader io.ReaderAt, offset int64, length int64) string {
   data := make([]byte, length)
   n, err := reader.ReadAt(data, offset)
   if err != nil && err != io.EOF {
     return ""
   }
   return UTF16BytesToUTF8(data[:n], binary.LittleEndian)
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

`
}

func (self UTF16StringParser) PrototypeName() string {
	return "UTF16String"
}

func (self *UTF16StringParser) Compile(struct_name string, field_name string) string {
	if self.Length == 0 {
		return fmt.Sprintf(`

func (self *%[1]s) %[2]s() string {
  return ParseTerminatedUTF16String(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
}
`, struct_name, field_name)
	}

	return fmt.Sprintf(`

func (self *%[1]s) %[2]s() string {
  return ParseUTF16String(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset, %[3]v)
}
`, struct_name, field_name, self.Length)
}

func (self UTF16StringParser) GoType() string {
	return "string"
}

func (self UTF16StringParser) Size(value string) string {
	return "0"
}
