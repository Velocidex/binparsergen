package binparsergen

import "fmt"

type SignatureParser struct {
	BaseParser
	Value string `json:"value,omitempty"`
}

func (self SignatureParser) Prototype() string {
	return `

type Signature struct {
    value, signature string
}

func (self Signature) IsValid() bool {
   return self.value == self.signature
}

func ParseSignature(reader io.ReaderAt, offset int64, length int64) string {
   data := make([]byte, length)
   n, err := reader.ReadAt(data, offset)
   if err != nil && err != io.EOF {
      return ""
   }
   return string(data[:n])
}

`
}

func (self SignatureParser) PrototypeName() string {
	return "Signature"
}

func (self *SignatureParser) Compile(struct_name string, field_name string) string {
	return fmt.Sprintf(`

func (self *%[1]s) %[2]s() *Signature {
  value := ParseSignature(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset, %[3]v)
  return &Signature{value: value, signature: "%v"}
}
`, struct_name, field_name, len(self.Value), self.Value)
}

func (self SignatureParser) GoType() string {
	return "Signature"
}

func (self SignatureParser) Size(value string) string {
	return fmt.Sprintf("%v", len(self.Value))
}
