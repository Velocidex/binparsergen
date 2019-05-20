package main

import (
	"flag"
	"fmt"
	"os"

	"www.velocidex.com/golang/binparsergen"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	spec, err := binparsergen.LoadSpecFile(args[0])
	binparsergen.FatalIfError(err, "Reading")

	profile, err := binparsergen.ConvertSpec(spec)
	binparsergen.FatalIfError(err, "Parsing")

	fmt.Println(binparsergen.GenerateCode(spec, profile))

}
