package binparsergen

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// The conversion process is driven by the conversion spec
// configuration file.
type ConversionSpec struct {
	Module              string              `yaml:"Module"`
	Profile             string              `yaml:"Profile"`
	Filename            string              `yaml:"Filename"`
	Structs             []string            `yaml:"Structs"`
	FieldWhiteList      map[string][]string `yaml:"FieldWhiteList"`
	FieldBlackList      map[string][]string `yaml:"FieldBlackList"`
	GenerateDebugString bool                `yaml:"GenerateDebugString"`
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
