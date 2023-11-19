package binparsergen

import "fmt"

type Enumeration struct {
	BaseParser
	Choices map[int]string `json:"choices,omitempty"`
	Target  string         `json:"target,omitempty"`
}

func (self Enumeration) getParser() Parser {
	switch self.Target {
	case "unsigned long":
		return &Uint32Parser{}
	case "unsigned short":
		return &Uint16Parser{}
	case "unsigned char":
		return &Uint8Parser{}
	}
	return &Uint64Parser{}
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
	result := fmt.Sprintf(`
func (self *%[1]s) %[2]s() *Enumeration {
   value := %[3]s(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
   name := "Unknown"
   switch value {
`, struct_name, field_name, self.getParser().PrototypeName())

	for _, k := range SortedIntKeys(self.Choices) {
		v := self.Choices[k]
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

func (self Enumeration) Dependencies() []Parser {
	return []Parser{self.getParser()}
}
