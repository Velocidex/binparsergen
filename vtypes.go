package binparsergen

import (
	"encoding/json"
	"io/ioutil"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

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

	for type_name, definition_list := range types {
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

		for field_name, field_def := range fields {
			if InString(spec.FieldBlackList[type_name], field_name) {
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
	kingpin.FatalIfError(err, "Decoding target offset")

	var params []json.RawMessage
	err = json.Unmarshal(*field_def[1], &params)
	kingpin.FatalIfError(err, "Decoding target params")

	new_field_def := _ParseParams(params, spec)
	new_field_def.Offset = offset

	return new_field_def
}

func _ParseParams(params []json.RawMessage, spec *ConversionSpec) *FieldDefinition {
	new_field_def := &FieldDefinition{}
	base_parser := BaseParser{Profile: spec.Profile}

	var parser_name string
	err := json.Unmarshal(params[0], &parser_name)
	kingpin.FatalIfError(err, "Decoding parser name")

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
		kingpin.FatalIfError(err, "Decoding")

		target_field_def := _ParseParams([]json.RawMessage{
			vtype_array.Target, vtype_array.TargetArgs}, spec)

		new_field_def.Pointer = &Pointer{
			BaseParser: base_parser,
			Target:     target_field_def,
		}

	case "BitField":
		bitfield := &BitField{BaseParser: base_parser}
		err = json.Unmarshal(params[1], &bitfield)
		kingpin.FatalIfError(err, "Decoding")

		new_field_def.BitField = bitfield

	case "Array":
		vtype_array := &VtypeArray{}
		err = json.Unmarshal(params[1], &vtype_array)
		kingpin.FatalIfError(err, "Decoding")

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
