package binparsergen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func FatalIfError(err error, format string, args ...interface{}) {
	if err != nil {
		fmt.Printf(format, args...)
		os.Exit(1)
	}
}

func ConvertSpec(spec *ConversionSpec) (map[string]*StructDefinition, error) {
	fd, err := os.Open(spec.Filename)
	if err != nil {
		return nil, err
	}

	definitions, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	var types map[string][]*json.RawMessage

	err = json.Unmarshal(definitions, &types)
	if err != nil {
		return nil, err
	}

	profile := make(map[string]*StructDefinition)

	for _, type_name := range SortedKeys(types) {
		definition_list := types[type_name]
		if !InString(spec.Structs, type_name) {
			continue
		}

		struct_def := &StructDefinition{
			Fields: make(map[string]*FieldDefinition),
		}
		err := json.Unmarshal(*definition_list[0], &struct_def.Size)
		if err != nil {
			return nil, err
		}

		fields := make(map[string][]*json.RawMessage)
		err = json.Unmarshal(*definition_list[1], &fields)
		if err != nil {
			return nil, err
		}

		for _, field_name := range SortedKeys(fields) {
			field_def := fields[field_name]
			if InString(spec.FieldBlackList[type_name], field_name) {
				continue
			}

			allowed_fields, pres := spec.FieldWhiteList[type_name]
			if pres && !InString(allowed_fields, field_name) {
				continue
			}
			struct_def.Fields[field_name] = ParseFieldDef(field_def, spec)
		}

		profile[type_name] = struct_def
	}

	return profile, nil
}

func ParseFieldDef(field_def []*json.RawMessage, spec *ConversionSpec) *FieldDefinition {
	var offset uint64

	err := json.Unmarshal(*field_def[0], &offset)
	FatalIfError(err, "Decoding target offset")

	var params []json.RawMessage
	err = json.Unmarshal(*field_def[1], &params)
	FatalIfError(err, "Decoding target params")

	new_field_def := _ParseParams(params, spec)
	new_field_def.Offset = offset

	return new_field_def
}

func _ParseParams(params []json.RawMessage, spec *ConversionSpec) *FieldDefinition {
	new_field_def := &FieldDefinition{}
	base_parser := BaseParser{Profile: spec.Profile}

	var parser_name string
	err := json.Unmarshal(params[0], &parser_name)
	FatalIfError(err, "Decoding parser name")

	switch parser_name {
	case "unsigned long long", "long long":
		new_field_def.Uint64Parser = &Uint64Parser{BaseParser: base_parser}

	case "unsigned long", "long":
		new_field_def.Uint32Parser = &Uint32Parser{BaseParser: base_parser}

	case "unsigned short", "short":
		new_field_def.Uint16Parser = &Uint16Parser{BaseParser: base_parser}

	case "unsigned char", "char":
		new_field_def.Uint8Parser = &Uint8Parser{BaseParser: base_parser}

	case "Pointer":
		vtype_array := &VtypeArray{}
		err = json.Unmarshal(params[1], &vtype_array)
		FatalIfError(err, "Decoding")

		target_field_def := _ParseParams([]json.RawMessage{
			vtype_array.Target, vtype_array.TargetArgs}, spec)

		new_field_def.Pointer = &Pointer{
			BaseParser: base_parser,
			Target:     target_field_def,
		}

	case "Enumeration":
		enumeration := &Enumeration{BaseParser: base_parser}
		err = json.Unmarshal(params[1], &enumeration)
		FatalIfError(err, "Decoding")

		new_field_def.Enumeration = enumeration

	case "BitField":
		bitfield := &BitField{BaseParser: base_parser}
		err = json.Unmarshal(params[1], &bitfield)
		FatalIfError(err, "Decoding")

		new_field_def.BitField = bitfield

	case "String":
		string_parser := &StringParser{BaseParser: base_parser}
		if len(params) > 1 {
			err = json.Unmarshal(params[1], &string_parser)
			FatalIfError(err, "Decoding")
		}

		new_field_def.StringParser = string_parser

	case "UnicodeString":
		string_parser := &UTF16StringParser{BaseParser: base_parser}
		if len(params) > 1 {
			err = json.Unmarshal(params[1], &string_parser)
			FatalIfError(err, "Decoding")
		}

		new_field_def.UTF16StringParser = string_parser

	case "Array":
		vtype_array := &VtypeArray{}
		err = json.Unmarshal(params[1], &vtype_array)
		FatalIfError(err, "Decoding")

		target_field_def := _ParseParams([]json.RawMessage{
			vtype_array.Target, vtype_array.TargetArgs}, spec)

		new_field_def.ArrayParser = &ArrayParser{
			BaseParser: base_parser,
			Count:      vtype_array.Count,
			Target:     target_field_def,
		}

	default:
		// This must be a reference to a custom type the user
		// may implement themselves.
		if !InString(spec.Structs, parser_name) {
			//log.Warn("Reference to undefined struct ", parser_name)
		}
		new_field_def.StructParser = &StructParser{
			BaseParser: base_parser,
			Target:     NormalizeName(parser_name),
		}
	}
	return new_field_def
}

type VtypeArray struct {
	Target     json.RawMessage
	TargetArgs json.RawMessage
	Count      int
}
