package gopkg

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"

	"github.com/andreasgerstmayr/gen-api-docs/pkg/parser/format"
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

func formatComment(commentLines []string) []string {
	if len(commentLines) == 1 && commentLines[0] == "" {
		return []string{}
	}

	out := []string{}
	for _, line := range commentLines {
		if !strings.HasPrefix(line, "+") {
			out = append(out, line)
		}
	}

	// trim empty lines at end
	for i := len(out) - 1; i > 0; i-- {
		if out[i] == "" {
			out = out[:i]
		} else {
			break
		}
	}
	return out
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

	if name == "lastTransitionTime" {
		print()
	}

	switch {
	case format.HideCoreTypes && slices.Contains(format.CoreTypes, name):
		return "{}"
	case ts == "k8s.io/apimachinery/pkg/api/resource.Quantity":
		return "0Gi"
	case ts == "k8s.io/apimachinery/pkg/apis/meta/v1.Duration":
		return "0h"
	case ts == "k8s.io/apimachinery/pkg/apis/meta/v1.Time":
		return "\"2006-01-02T15:04:05Z\""
	case ts == "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1.JSON":
		return "{}"
	case ts == "map[string]string":
		return "{}"
	case ts == "string":
		return "\"\""
	case ts == "bool":
		return "false"
	case ts == "int" || ts == "int64":
		return "0"
	default:
		return ""
	}
}

func render(type_ *types.Type, prop *format.Prop, parentCommentLines []string) {
	type_ = deref(type_)

	switch type_.Kind {
	case types.Struct:
		for _, member := range type_.Members {
			if fieldEmbedded(member) {
				render(member.Type, prop, member.CommentLines)
				continue
			}
			name := fieldName(member)
			if name == "" {
				continue
			}

			if name == "status" {
				print()
			}

			comment := formatComment(member.CommentLines)
			if len(comment) == 0 {
				comment = formatComment(member.Type.CommentLines)
			}

			p := &format.Prop{
				Key:         name,
				Comment:     comment,
				ScalarValue: getDefaultValue(name, member.Type),
			}
			prop.Properties = append(prop.Properties, p)
			if p.ScalarValue == "" {
				render(member.Type, p, member.CommentLines)
			}
		}

	case types.Map:
		p := &format.Prop{
			Key:         "\"key\"",
			Comment:     formatComment(type_.CommentLines),
			ScalarValue: getDefaultValue("", type_.Elem),
		}

		if len(parentCommentLines) > 0 && strings.Contains(parentCommentLines[0], "Requests describes the minimum amount of compute resources required.") {
			prop.Properties = []*format.Prop{
				{Key: "cpu", ScalarValue: "\"500m\""},
				{Key: "memory", ScalarValue: "\"1Gi\""},
			}
		} else if len(parentCommentLines) > 0 && strings.Contains(parentCommentLines[0], "Limits describes the maximum amount of compute resources allowed.") {
			prop.Properties = []*format.Prop{
				{Key: "cpu", ScalarValue: "\"750m\""},
				{Key: "memory", ScalarValue: "\"2Gi\""},
			}
		} else {
			prop.Properties = []*format.Prop{p}
			if p.ScalarValue == "" {
				render(type_.Elem, p, parentCommentLines)
			}
		}

	case types.Slice:
		if prop.Key == "permissions" {
			print()
		}
		prop.ListItem = &format.Prop{
			Comment:     formatComment(type_.Elem.CommentLines),
			ScalarValue: getDefaultValue("", type_.Elem),
		}
		if prop.ListItem.ScalarValue == "" {
			render(type_.Elem, prop.ListItem, parentCommentLines)
		}
	}
}

func parsePackage(type_ *types.Type) format.Prop {
	prop := format.Prop{}
	render(type_, &prop, []string{})
	return prop
}

func Parse(pkgName string, typeName string) (format.Prop, error) {
	p := parser.New()
	err := p.AddDirRecursive(pkgName)
	if err != nil {
		return format.Prop{}, err
	}

	pkgs, err := p.FindTypes()
	if err != nil {
		return format.Prop{}, err
	}

	for name, pkg := range pkgs {
		if name == pkgName {
			type_, ok := pkg.Types[typeName]
			if !ok {
				return format.Prop{}, fmt.Errorf("cannot find type %s", typeName)
			}
			return parsePackage(type_), nil
		}
	}
	return format.Prop{}, fmt.Errorf("cannot find package %s", pkgName)
}
