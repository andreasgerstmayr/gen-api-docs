package main

import (
	"flag"
	"log"
	"os"

	"github.com/andreasgerstmayr/gen-api-docs/pkg/parser"
)

func main() {
	flag.IntVar(&parser.CommentPadding, "padding", 40, "comment padding")
	pkgName := flag.String("pkg", "", "package name")
	typeName := flag.String("type", "", "type name")
	flag.Parse()

	var err error
	if *pkgName != "" && *typeName != "" {
		err = parser.GoPackage(*pkgName, *typeName)
	} else {
		err = parser.ParseCRD(os.Stdin, os.Stdout)
	}

	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
