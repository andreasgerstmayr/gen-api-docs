package main

import (
	"flag"
	"log"
	"os"

	"github.com/andreasgerstmayr/gen-api-docs/pkg/parser/crd"
	"github.com/andreasgerstmayr/gen-api-docs/pkg/parser/format"
	"github.com/andreasgerstmayr/gen-api-docs/pkg/parser/gopkg"
)

func main() {
	flag.StringVar(&format.Format, "format", "oneline", "available format: oneline, multiline")
	flag.IntVar(&format.CommentPadding, "padding", 40, "comment padding")
	flag.BoolVar(&format.HideCoreTypes, "hideCoreTypes", true, "hide core types")
	pkgName := flag.String("pkg", "", "package name")
	typeName := flag.String("type", "", "type name")
	flag.Parse()

	var prop format.Prop
	var err error
	if *pkgName != "" && *typeName != "" {
		prop, err = gopkg.Parse(*pkgName, *typeName)
	} else {
		prop, err = crd.Parse(os.Stdin)
	}
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	format.Print(os.Stdout, &prop)
}
