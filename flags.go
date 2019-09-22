package binparsergen

import (
	"fmt"
)

type Flags struct {
	BaseParser
	Maskmap map[string]int `json:"maskmap,omitempty"`
	Bitmap  map[string]int `json:"bitmap,omitempty"`
	Target  string         `json:"target,omitempty"`
}

func (self *Flags) Prototype() string {
	return `
type Flags struct {
    Value uint64
    Names  map[string]bool
}

func (self Flags) DebugString() string {
    names := []string{}
    for k, _ := range self.Names {
      names = append(names, k)
    }

    sort.Strings(names)

    return fmt.Sprintf("%d (%s)", self.Value, strings.Join(names, ","))
}

func (self Flags) IsSet(flag string) bool {
    result, _ := self.Names[flag]
    return result
}

`
}

func (self Flags) PrototypeName() string {
	return "Flags"
}

func (self Flags) Compile(struct_name string, field_name string) string {
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
func (self *%[1]s) %[2]s() *Flags {
   value := %[3]s(self.Reader, self.Profile.Off_%[1]s_%[2]s + self.Offset)
   names := make(map[string]bool)

`, struct_name, field_name, parser_func)
	for k, v := range self.Maskmap {
		result += fmt.Sprintf(`
   if value & %v != 0 {
      names["%s"] = true
   }
`, v, k)
	}

	for k, v := range self.Bitmap {
		result += fmt.Sprintf(`
   if value & (1 << %v) != 0 {
      names["%s"] = true
   }
`, v, k)
	}

	result += `
   return &Flags{Value: uint64(value), Names: names}
}

`
	return result
}

func (self Flags) GoType() string {
	return "uint64"
}

func (self Flags) Size(value string) string {
	return "8"
}
