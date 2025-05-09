package binparsergen

import "fmt"

// A parser is an object which generates code to extract a specific
// object from binary data.
type Parser interface {

	// Generate a method on struct_name that extracts field field_name.
	Compile(struct_name string, field_name string) string

	// The name of the profile we are generating.
	ProfileName() string

	// Generate a free function which can parse this object from a
	// reader at a particular offset.
	Prototype() string

	// The name of the Prototype() method.
	PrototypeName() string

	// The GoType we will use to represent this object.
	GoType() string

	// If we use a pointer to represent the object this method
	// should return a "*"
	GoTypePointer() string

	// The size of the object.
	Size(value string) string

	// The list of PrototypeNames() that this types requires to be present
	Dependencies() []Parser
}

type BaseParser struct {
	Profile string
}

func (self BaseParser) Compile(struct_name string, field_name string) string {
	return ""
}

func (self BaseParser) Prototype() string {
	return ""
}

func (self BaseParser) PrototypeName() string {
	return ""
}

func (self BaseParser) ProfileName() string {
	return self.Profile
}

func (self BaseParser) GoType() string {
	return ""
}

func (self BaseParser) GoTypePointer() string {
	return ""
}

func (self BaseParser) Size(value string) string {
	return ""
}

func (self BaseParser) Dependencies() []Parser {
	return nil
}

type NullParser struct {
	BaseParser
}

type Uint64Parser struct {
	BaseParser
}

func (self Uint64Parser) Compile(struct_name string, field_name string) string {
	return fmt.Sprintf(`
func (self *%[1]s) %[2]s() uint64 {
    return ParseUint64(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
}
`, struct_name, field_name)
}

func (self Uint64Parser) Prototype() string {
	return fmt.Sprintf(`
func ParseUint64(reader io.ReaderAt, offset int64) uint64 {
	var buf [8]byte
	data := buf[:]
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return binary.LittleEndian.Uint64(data)
}
`)
}

func (self Uint64Parser) PrototypeName() string {
	return "ParseUint64"
}

func (self Uint64Parser) GoType() string {
	return "uint64"
}

func (self Uint64Parser) Size(value string) string {
	return "8"
}

type Int64Parser struct {
	BaseParser
}

func (self Int64Parser) Compile(struct_name string, field_name string) string {
	return fmt.Sprintf(`
func (self *%[1]s) %[2]s() int64 {
    return ParseInt64(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
}
`, struct_name, field_name)
}

func (self Int64Parser) Prototype() string {
	return fmt.Sprintf(`
func ParseInt64(reader io.ReaderAt, offset int64) int64 {
	var buf [8]byte
	data := buf[:]
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return int64(binary.LittleEndian.Uint64(data))
}
`)
}

func (self Int64Parser) PrototypeName() string {
	return "ParseInt64"
}

func (self Int64Parser) GoType() string {
	return "int64"
}

func (self Int64Parser) Size(value string) string {
	return "8"
}

type Uint32Parser struct {
	BaseParser
}

func (self Uint32Parser) Prototype() string {
	return `
func ParseUint32(reader io.ReaderAt, offset int64) uint32 {
	var buf [4]byte
	data := buf[:]
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return binary.LittleEndian.Uint32(data)
}
`
}

func (self Uint32Parser) PrototypeName() string {
	return "ParseUint32"
}

func (self Uint32Parser) Compile(struct_name string, field_name string) string {
	return fmt.Sprintf(`
func (self *%[1]s) %[2]s() uint32 {
   return ParseUint32(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
}
`, struct_name, field_name)
}
func (self Uint32Parser) GoType() string {
	return "uint32"
}
func (self Uint32Parser) Size(value string) string {
	return "4"
}

type Int32Parser struct {
	BaseParser
}

func (self Int32Parser) Prototype() string {
	return `
func ParseInt32(reader io.ReaderAt, offset int64) int32 {
	var buf [4]byte
	data := buf[:]
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return int32(binary.LittleEndian.Uint32(data))
}
`
}

func (self Int32Parser) PrototypeName() string {
	return "ParseInt32"
}

func (self Int32Parser) Compile(struct_name string, field_name string) string {
	return fmt.Sprintf(`
func (self *%[1]s) %[2]s() int32 {
   return ParseInt32(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
}
`, struct_name, field_name)
}
func (self Int32Parser) GoType() string {
	return "int32"
}
func (self Int32Parser) Size(value string) string {
	return "4"
}

type Uint16Parser struct {
	BaseParser
}

func (self Uint16Parser) Prototype() string {
	return `
func ParseUint16(reader io.ReaderAt, offset int64) uint16 {
	var buf [2]byte
	data := buf[:]
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return binary.LittleEndian.Uint16(data)
}
`
}

func (self Uint16Parser) PrototypeName() string {
	return "ParseUint16"
}

func (self Uint16Parser) Compile(struct_name string, field_name string) string {
	return fmt.Sprintf(`
func (self *%[1]s) %[2]s() uint16 {
   return ParseUint16(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
}
`, struct_name, field_name)
}
func (self Uint16Parser) GoType() string {
	return "uint16"
}
func (self Uint16Parser) Size(value string) string {
	return "2"
}

type Int16Parser struct {
	BaseParser
}

func (self Int16Parser) Prototype() string {
	return `
func ParseInt16(reader io.ReaderAt, offset int64) int16 {
	var buf [2]byte
	data := buf[:]
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return int16(binary.LittleEndian.Uint16(data))
}
`
}

func (self Int16Parser) PrototypeName() string {
	return "ParseInt16"
}

func (self Int16Parser) Compile(struct_name string, field_name string) string {
	return fmt.Sprintf(`
func (self *%[1]s) %[2]s() int16 {
   return ParseInt16(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
}
`, struct_name, field_name)
}
func (self Int16Parser) GoType() string {
	return "int16"
}
func (self Int16Parser) Size(value string) string {
	return "2"
}

type Uint8Parser struct {
	BaseParser
}

func (self Uint8Parser) Prototype() string {
	return `
func ParseUint8(reader io.ReaderAt, offset int64) byte {
	var buf [1]byte
	data := buf[:]
    _, err := reader.ReadAt(data, offset)
    if err != nil {
       return 0
    }
    return data[0]
}
`
}
func (self Uint8Parser) PrototypeName() string {
	return "ParseUint8"
}

func (self Uint8Parser) Compile(struct_name string, field_name string) string {
	return fmt.Sprintf(`
func (self *%[1]s) %[2]s() byte {
   return ParseUint8(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
}
`, struct_name, field_name)
}
func (self Uint8Parser) GoType() string {
	return "byte"
}
func (self Uint8Parser) Size(v string) string {
	return "1"
}

type Int8Parser struct {
	BaseParser
}

func (self Int8Parser) Prototype() string {
	return `
func ParseInt8(reader io.ReaderAt, offset int64) int8 {
	var buf [1]byte
	data := buf[:]
    _, err := reader.ReadAt(result, offset)
    if err != nil {
       return 0
    }
    return int8(result[0])
}
`
}
func (self Int8Parser) PrototypeName() string {
	return "ParseInt8"
}

func (self Int8Parser) Compile(struct_name string, field_name string) string {
	return fmt.Sprintf(`
func (self *%[1]s) %[2]s() int8 {
   return ParseInt8(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
}
`, struct_name, field_name)
}
func (self Int8Parser) GoType() string {
	return "int8"
}
func (self Int8Parser) Size(v string) string {
	return "1"
}

type ArrayParser struct {
	BaseParser
	Target       *FieldDefinition
	Count        int
	DynamicCount string
}

func (self ArrayParser) Prototype() string {
	parser := self.Target.GetParser()
	return fmt.Sprintf(`
func ParseArray_%[1]s(profile *%[2]s, reader io.ReaderAt, offset int64, count int) []%[3]s%[1]s {
    if count <= 0 {
      count = 0
    }
    if count > 4000000 {
       count = 4000000
    }
    result := make([]%[3]s%[1]s, 0, count)
    for i:=0; i<count; i++ {
      value := %[4]s(reader, offset)
      result = append(result, value)
      offset += int64(%[5]s)
    }
    return result
}
`, parser.GoType(), parser.ProfileName(), parser.GoTypePointer(),
		parser.PrototypeName(), parser.Size("value"))
}

func (self ArrayParser) PrototypeName() string {
	parser := self.Target.GetParser()
	return fmt.Sprintf("ParseArray_%s", parser.GoType())
}

func (self ArrayParser) Compile(struct_name string, field_name string) string {
	parser := self.Target.GetParser()

	if self.DynamicCount == "" {
		return fmt.Sprintf(`
func (self *%[1]s) %[2]s() []%[3]s%[4]s {
   return %[5]s(self.Profile, self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset, %[6]d)
}
`, struct_name, field_name, parser.GoTypePointer(),
			parser.GoType(), self.PrototypeName(),
			self.Count)
	} else {
		return fmt.Sprintf(`
func (self *%[1]s) %[2]s() []%[3]s%[4]s {
   return %[5]s(self.Profile, self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset, int(self.%[6]s()))
}
`, struct_name, field_name, parser.GoTypePointer(),
			parser.GoType(), self.PrototypeName(),
			self.DynamicCount)

	}
}

func (self ArrayParser) GoType() string {
	return "XXX"
}

func (self ArrayParser) Size(value string) string {
	parser := self.Target.GetParser()
	return fmt.Sprintf("%d * %s", self.Count, parser.Size(value))
}

type StructParser struct {
	BaseParser
	Target string
}

func (self StructParser) Prototype() string {
	return ""
	return fmt.Sprintf(`
func Parse_%[1]s(profile %[2]s, reader io.ReaderAt, offset int) *%[3]s {
   return profile.%[3]s(reader, offset)
}
`, self.Target, self.ProfileName(), self.Target)
}

func (self StructParser) PrototypeName() string {
	return fmt.Sprintf("profile.%s", self.Target)
}

func (self StructParser) Compile(struct_name string, field_name string) string {
	return fmt.Sprintf(`
func (self *%[1]s) %[2]s() *%[3]s {
    return self.Profile.%[3]s(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
}
`, struct_name, field_name, self.Target)
}

func (self StructParser) GoType() string {
	return self.Target
}

func (self StructParser) GoTypePointer() string {
	return "*"
}

func (self StructParser) Size(value string) string {
	return fmt.Sprintf("%s.Size()", value)
}

type Pointer struct {
	BaseParser
	Target *FieldDefinition
}

func (self *Pointer) Prototype() string {
	return ""
}

func (self Pointer) PrototypeName() string {
	return ""
}

func (self Pointer) Compile(struct_name string, field_name string) string {
	parser := self.Target.GetParser()

	return fmt.Sprintf(`
func (self *%[1]s) %[2]s() *%[3]s {
   deref := ParseUint64(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
   return self.Profile.%[3]s(self.Reader, int64(deref))
}
`, struct_name, field_name, parser.GoType())
}

func (self Pointer) GoType() string {
	parser := self.Target.GetParser()
	return "*" + parser.GoType()
}

func (self Pointer) Size(value string) string {
	return "8"
}

func (self Pointer) Dependencies() []Parser {
	return []Parser{&Uint64Parser{}}
}

type BitField struct {
	BaseParser
	StartBit uint64 `json:"start_bit,omitempty"`
	EndBit   uint64 `json:"end_bit,omitempty"`
	Target   string `json:"target,omitempty"`
}

func (self BitField) getParser() Parser {
	switch self.Target {
	case "unsigned long long":
		return &Uint64Parser{}
	case "unsigned long":
		return &Uint32Parser{}
	case "unsigned short":
		return &Uint16Parser{}
	case "unsigned char":
		return &Uint8Parser{}
	}
	return &Uint64Parser{}
}

func (self *BitField) Prototype() string {
	return ""
}

func (self BitField) PrototypeName() string {
	return ""
}

func (self BitField) Compile(struct_name string, field_name string) string {

	return fmt.Sprintf(`
func (self *%[1]s) %[2]s() uint64 {
   value := %[3]s(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
   return (uint64(value) & %#[4]x) >> %#[5]x
}
`, struct_name, field_name, self.getParser().PrototypeName(),
		(1<<uint64(self.EndBit))-1, self.StartBit)
}

func (self BitField) GoType() string {
	return "uint64"
}

func (self BitField) Size(value string) string {
	return "8"
}

func (self BitField) Dependencies() []Parser {
	return []Parser{self.getParser()}
}
