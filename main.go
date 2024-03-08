package main

import (
	"flag"
	"log"
	"os"

	"github.com/andreasgerstmayr/gen-api-docs/pkg/parser/crd"
	"github.com/andreasgerstmayr/gen-api-docs/pkg/parser/gopkg"
	"github.com/andreasgerstmayr/gen-api-docs/pkg/parser/utils"
)

func main() {
	flag.IntVar(&utils.CommentPadding, "padding", 40, "comment padding")
	flag.BoolVar(&utils.HideCoreTypes, "hideCoreTypes", true, "hide core types")
	pkgName := flag.String("pkg", "", "package name")
	typeName := flag.String("type", "", "type name")
	flag.Parse()

	var err error
	if *pkgName != "" && *typeName != "" {
		err = gopkg.Parse(*pkgName, *typeName, os.Stdout)
	} else {
		err = crd.Parse(os.Stdin, os.Stdout)
	}

	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
