package binparsergen

import (
	"fmt"
	"strings"
)

/* A profile is a factory object for all its containing types. We
   store all the member offsets for all its structs in this object so
   we can tweak offsets on the fly. This allows us to cater for
   different versions. For example:

   type RegistryProfile struct {
	Off_LARGE_INTEGER_HighPart               int64
	Off_LARGE_INTEGER_LowPart                int64
	Off_LARGE_INTEGER_QuadPart               int64
       ....
   }

   The profile is initialized with the offsets defined in the current
   profile json file but may be overridden if a more up to date file
   is available.

   Profiles also contain methods for all their structs. For example
   this method will be generated to parse a HCELL struct at the
   specified offset.

   func (self *RegistryProfile) HCELL(reader io.ReaderAt, offset int64) *HCELL {
     ....
   }

   This allows structs to be grouped into logical units all controlled
   via the same profile (typically compilation units).

*/
func GenerateProfileCode(
	profile_name string,
	profile map[string]*StructDefinition) string {
	result := fmt.Sprintf("type %s struct {\n", profile_name)
	init := []string{}
	factories := ""
	for struct_name, struct_def := range profile {
		struct_name = NormalizeName(struct_name)
		for field_name, field_def := range struct_def.Fields {
			result += fmt.Sprintf("    Off_%s_%s int64\n",
				struct_name, field_name)
			init = append(init, fmt.Sprintf("%d", field_def.Offset))
		}
		factories += fmt.Sprintf(`
func (self *%s) %s(reader io.ReaderAt, offset int64) *%s {
    return &%s{Reader: reader, Offset: offset, Profile: self}
}
`, profile_name, struct_name, struct_name, struct_name)
	}

	result += fmt.Sprintf(`}

func New%s() *%s {
    // Specific offsets can be tweaked to cater for slight version mismatches.
    self := &%s{%s}
    return self
}
%s
`, profile_name, profile_name, profile_name, strings.Join(init, ","), factories)

	return result
}
