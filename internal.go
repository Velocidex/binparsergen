package binparsergen

var (
	prototypes = make(map[string]string)
)

// How to represent a field in our data structure.
type FieldDefinition struct {
	// A field has an offset within the struct.
	Offset uint64

	// A field may be one of the following parsers. Only one of
	// these parsers is allowed.
	Uint64Parser      *Uint64Parser      `json:"Uint64Parser,omitempty"`
	Uint32Parser      *Uint32Parser      `json:"Uint32Parser,omitempty"`
	Uint16Parser      *Uint16Parser      `json:"Uint16Parser,omitempty"`
	Uint8Parser       *Uint8Parser       `json:"Uint8Parser,omitempty"`
	StructParser      *StructParser      `json:"StructParser,omitempty"`
	ArrayParser       *ArrayParser       `json:"ArrayParser,omitempty"`
	Pointer           *Pointer           `json:"Pointer,omitempty"`
	BitField          *BitField          `json:"BitField,omitempty"`
	Enumeration       *Enumeration       `json:"Enumeration,omitempty"`
	StringParser      *StringParser      `json:"StringParser,omitempty"`
	UTF16StringParser *UTF16StringParser `json:"UTF16StringParser,omitempty"`
}

// Extract the active parser from the field definition.
func (self *FieldDefinition) GetParser() Parser {
	var result Parser = &NullParser{}
	if self.Uint64Parser != nil {
		result = self.Uint64Parser
	} else if self.Uint32Parser != nil {
		result = self.Uint32Parser

	} else if self.Uint16Parser != nil {
		result = self.Uint16Parser

	} else if self.Uint8Parser != nil {
		result = self.Uint8Parser

	} else if self.ArrayParser != nil {
		result = self.ArrayParser

	} else if self.StructParser != nil {
		result = self.StructParser

	} else if self.Pointer != nil {
		result = self.Pointer

	} else if self.BitField != nil {
		result = self.BitField

	} else if self.Enumeration != nil {
		result = self.Enumeration

	} else if self.StringParser != nil {
		result = self.StringParser

	} else if self.UTF16StringParser != nil {
		result = self.UTF16StringParser

	}

	_, pres := prototypes[result.PrototypeName()]
	if !pres {
		prototypes[result.PrototypeName()] = result.Prototype()
	}
	return result
}

// We can consume JSON encoded struct definitions in this format.
type StructDefinition struct {
	Size   uint32
	Fields map[string]*FieldDefinition
}
