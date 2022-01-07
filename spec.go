package binparsergen

import (
	"io/ioutil"
	"os"

	yaml "github.com/Velocidex/yaml/v2"
)

// The conversion process is driven by the conversion spec
// configuration file.
type ConversionSpec struct {
	Module              string              `json:"Module"`
	Profile             string              `json:"Profile"`
	Filename            string              `json:"Filename"`
	Structs             []string            `json:"Structs"`
	FieldWhiteList      map[string][]string `json:"FieldWhiteList"`
	FieldBlackList      map[string][]string `json:"FieldBlackList"`
	GenerateDebugString bool                `json:"GenerateDebugString"`
}

func LoadSpecFile(filename string) (*ConversionSpec, error) {
	result := &ConversionSpec{}
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
