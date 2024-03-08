package gopkg

import (
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"

	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"

	"github.com/andreasgerstmayr/gen-api-docs/pkg/parser/utils"
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

func formatComment(commentLines []string) string {
	return strings.Join(commentLines, " ")
}

// deref resolves pointers and aliases
func deref(t *types.Type) *types.Type {
	switch t.Kind {
	case types.Pointer:
		return deref(t.Elem)
	case types.Alias:
		return deref(t.Underlying)
	default:
		return t
	}
}

func getDefaultValue(name string, type_ *types.Type) string {
	t := deref(type_)
	ts := t.String()

	switch {
	case utils.HideCoreTypes && slices.Contains(utils.CoreTypes, name):
		return "{}"
	case ts == "string":
		return "\"\""
	case ts == "bool":
		return "false"
	case ts == "int":
		return "0"
	case ts == "k8s.io/apimachinery/pkg/api/resource.Quantity":
		return "0Gi"
	case ts == "k8s.io/apimachinery/pkg/apis/meta/v1.Duration":
		return "0h"
	default:
		return ""
	}
}

func render(out io.Writer, type_ *types.Type, level int, isList bool, parentCommentLines []string) error {
	type_ = deref(type_)

	switch type_.Kind {
	case types.Struct:
		isFirstElem := true
		for _, member := range type_.Members {
			if fieldEmbedded(member) {
				render(out, member.Type, level, false, member.CommentLines)
				continue
			}
			name := fieldName(member)
			if name == "" {
				continue
			}

			var indent string
			if isFirstElem && isList {
				indent = strings.Repeat("  ", level-1) + "- "
			} else {
				indent = strings.Repeat("  ", level)
			}
			comment := formatComment(member.CommentLines)
			value := getDefaultValue(name, member.Type)

			if value != "" {
				utils.WriteLine(out, fmt.Sprintf("%s%s: %s", indent, name, value), comment)
			} else {
				utils.WriteLine(out, fmt.Sprintf("%s%s:", indent, name), comment)
				if name == "requests" {
					fmt.Printf("")
				}
				render(out, member.Type, level+1, false, member.CommentLines)
			}

			isFirstElem = false
		}

	case types.Map:
		indent := strings.Repeat("  ", level)
		comment := formatComment(parentCommentLines)
		value := getDefaultValue("", type_.Elem)

		if len(parentCommentLines) > 0 && strings.Contains(parentCommentLines[0], "Requests describes the minimum amount of compute resources required.") {
			utils.WriteLine(out, fmt.Sprintf("%scpu: \"%s\"", indent, "500m"), comment)
			utils.WriteLine(out, fmt.Sprintf("%smemory: \"%s\"", indent, "1Gi"), comment)
		} else if len(parentCommentLines) > 0 && strings.Contains(parentCommentLines[0], "Limits describes the maximum amount of compute resources allowed.") {
			utils.WriteLine(out, fmt.Sprintf("%scpu: \"%s\"", indent, "750m"), comment)
			utils.WriteLine(out, fmt.Sprintf("%smemory: \"%s\"", indent, "2Gi"), comment)
		} else if value != "" {
			utils.WriteLine(out, fmt.Sprintf("%s\"key\": %s", indent, value), comment)
		} else {
			utils.WriteLine(out, fmt.Sprintf("%s\"key\":", indent), "")
			render(out, type_.Elem, level+1, false, parentCommentLines)
		}
	}

	return nil
}

func Parse(pkgName string, typeName string, out io.Writer) error {
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
			return render(out, type_, 0, false, []string{})
		}
	}
	return nil
}
