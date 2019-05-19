package binparsergen

import "fmt"

type Enumeration struct {
	BaseParser
	Choices map[int]string `json:"choices,omitempty"`
	Target  string         `json:"target,omitempty"`
}

func (self *Enumeration) Prototype() string {
	return `
type Enumeration struct {
    Value uint64
    Name  string
}

func (self Enumeration) DebugString() string {
    return fmt.Sprintf("%s (%d)", self.Name, self.Value)
}

`
}

func (self Enumeration) PrototypeName() string {
	return "Enumeration"
}

func (self Enumeration) Compile(struct_name string, field_name string) string {
	parser_func := "ParseUint64"
	switch self.Target {
	case "unsigned long":
		parser_func = "ParseUint32"
	case "unsigned short":
		parser_func = "ParseUint16"
	case "unsigned char":
		parser_func = "ParseUint8"
	}
	result := fmt.Sprintf(`
func (self *%[1]s) %[2]s() *Enumeration {
   value := %[3]s(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
   name := "Unknown"
   switch value {
`, struct_name, field_name, parser_func)

	for k, v := range self.Choices {
		result += fmt.Sprintf(`
      case %d:
         name = "%s"
`, k, v)
	}

	result += `}
   return &Enumeration{Value: uint64(value), Name: name}
}

`
	return result
}

func (self Enumeration) GoType() string {
	return "uint64"
}

func (self Enumeration) Size(value string) string {
	return "8"
}
