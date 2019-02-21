package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"www.velocidex.com/golang/binparsergen"
)

var (
	app = kingpin.New("vtype",
		"A tool for managing vtype definitions.")

	convert_command_file_arg = app.Arg(
		"file", "The yaml file to inspect",
	).Required().String()
)

func doConvert() {
	spec, err := binparsergen.LoadSpecFile(*convert_command_file_arg)
	kingpin.FatalIfError(err, "Reading")

	profile, err := binparsergen.ConvertSpec(spec)
	kingpin.FatalIfError(err, "Parsing")

	fmt.Println(binparsergen.GenerateCode(spec, profile))
}

func main() {
	app.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.CompactUsageTemplate)
	kingpin.MustParse(app.Parse(os.Args[1:]))
	doConvert()
}
