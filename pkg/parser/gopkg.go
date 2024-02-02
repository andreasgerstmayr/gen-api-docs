package parser

import (
	"fmt"
	"reflect"
	"strings"

	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
)

func fieldName(m types.Member) string {
	v := reflect.StructTag(m.Tags).Get("json")
	v = strings.TrimSuffix(v, ",inline")
	v = strings.TrimSuffix(v, ",omitempty")
	return v
}

func fieldEmbedded(m types.Member) bool {
	return strings.Contains(reflect.StructTag(m.Tags).Get("json"), ",inline")
}

func renderPkg(pkg *types.Type, level int, isList bool) error {
	for _, member := range pkg.Members {
		if fieldEmbedded(member) {
			renderPkg(member.Type, level, false)
			continue
		}

		jsonTag := fieldName(member)
		if jsonTag != "" {
			fmt.Printf("%s\n", fieldName(member))
		}
	}
	return nil
}

func GoPackage(pkgName string, typeName string) error {
	p := parser.New()
	err := p.AddDirRecursive(pkgName)
	if err != nil {
		return err
	}

	pkgs, err := p.FindTypes()
	if err != nil {
		return err
	}

	for name, pkg := range pkgs {
		if name == pkgName {
			type_, ok := pkg.Types[typeName]
			if !ok {
				return fmt.Errorf("cannot find type %s", typeName)
			}
			return renderPkg(type_, 0, false)
		}
	}
	return nil
}
