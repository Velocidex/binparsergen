package binparsergen

import (
	"fmt"
	"strings"
)

func GeneratePrototypes() string {
	result := ""
	for _, k := range SortedKeys(prototypes) {
		result += prototypes[k]
	}

	return result
}

func GenerateStructCode(name string, profile_name string, definition *StructDefinition) string {
	result := fmt.Sprintf(`
type %s struct {
    Reader io.ReaderAt
    Offset int64
    Profile *%s
}

func New%s(reader io.ReaderAt) *%s {
    self := &%s{Reader: reader}
    return self
}

func (self *%s) Size() int {
    return %d
}
`,
		name, profile_name,
		name, name,
		name,
		name, definition.Size)

	for _, field_name := range SortedKeys(definition.Fields) {
		field_def := definition.Fields[field_name]
		result += field_def.GetParser().Compile(name, field_name)
	}

	return result
}

func GenerateDebugString(name string, profile_name string, definition *StructDefinition) string {
	result := fmt.Sprintf(
		"func (self *%s) DebugString() string {\n    result := fmt.Sprintf("+
			"\"struct %s @ %%#x:\\n\", self.Offset)\n", name, name)
	for _, field_name := range SortedKeys(definition.Fields) {
		field_def := definition.Fields[field_name]
		if field_def.StringParser != nil ||
			field_def.UTF16StringParser != nil {
			result += fmt.Sprintf(
				"    result += fmt.Sprintf(\"%[1]s: %%v\\n\", string(self.%[1]s()))\n",
				field_name)

		} else if field_def.Uint64Parser != nil ||
			field_def.BitField != nil ||
			field_def.Uint16Parser != nil ||
			field_def.Uint32Parser != nil {
			result += fmt.Sprintf(
				"    result += fmt.Sprintf(\"%[1]s: %%#0x\\n\", self.%[1]s())\n",
				field_name)
		} else if field_def.Enumeration != nil {
			result += fmt.Sprintf(
				"    result += fmt.Sprintf(\"%[1]s: %%v\\n\", self.%[1]s().DebugString())\n",
				field_name)
		} else if field_def.StructParser != nil {
			result += fmt.Sprintf(
				"    result += fmt.Sprintf(\"%[1]s: {\\n%%v}\\n\", self.%[1]s().DebugString())\n",
				field_name)
		}
	}

	result += "    return result\n}\n"

	return result
}

// Convert names to something which is exportable by Go.
func NormalizeName(name string) string {
	name = strings.TrimLeft(name, "_")
	return strings.ToUpper(name[:1]) + name[1:]
}
