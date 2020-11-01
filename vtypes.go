package binparsergen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Velocidex/ordereddict"
)

func FatalIfError(err error, format string, args ...interface{}) {
	if err != nil {
		fmt.Printf(format, args...)
		fmt.Printf(": %v\n", err)
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

		ordered_fields := ordereddict.NewDict()
		err = json.Unmarshal(*definition_list[1], &ordered_fields)
		if err != nil {
			return nil, err
		}

		// Preserve the order of the fields.
		struct_def.fields = ordered_fields.Keys()
		for _, field_name := range struct_def.fields {
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
	var offset int64

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
	case "unsigned long long", "uint64":
		new_field_def.Uint64Parser = &Uint64Parser{BaseParser: base_parser}
	case "long long", "int64":
		new_field_def.Int64Parser = &Int64Parser{BaseParser: base_parser}

	case "unsigned long", "uint32":
		new_field_def.Uint32Parser = &Uint32Parser{BaseParser: base_parser}

	case "long":
		new_field_def.Int32Parser = &Int32Parser{BaseParser: base_parser}

	case "unsigned short":
		new_field_def.Uint16Parser = &Uint16Parser{BaseParser: base_parser}

	case "short":
		new_field_def.Int16Parser = &Int16Parser{BaseParser: base_parser}

	case "unsigned char":
		new_field_def.Uint8Parser = &Uint8Parser{BaseParser: base_parser}

	case "char":
		new_field_def.Int8Parser = &Int8Parser{BaseParser: base_parser}

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

	case "Signature":
		sig := &SignatureParser{BaseParser: base_parser}
		err = json.Unmarshal(params[1], &sig)
		FatalIfError(err, "Decoding")

		new_field_def.SignatureParser = sig

	case "Flags":
		flags := &Flags{BaseParser: base_parser}
		err = json.Unmarshal(params[1], &flags)
		FatalIfError(err, "Decoding")

		new_field_def.Flags = flags

	case "BitField":
		bitfield := &BitField{BaseParser: base_parser}
		err = json.Unmarshal(params[1], &bitfield)
		FatalIfError(err, "Decoding")

		new_field_def.BitField = bitfield

	case "String":
		string_parser := &StringParser{BaseParser: base_parser}
		if len(params) > 1 && len(params[1]) > 0 {
			err = json.Unmarshal(params[1], &string_parser)
			FatalIfError(err, fmt.Sprintf("Decoding %v", params))
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
