package parser

import (
	"fmt"
	"io"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var (
	CommentPadding int
)

func writeLine(out io.Writer, line string, comment string) {
	padding := CommentPadding - len(line)
	if padding < 0 {
		padding = 0
	}

	if comment != "" {
		fmt.Fprintf(out, "%s%s # %s\n", line, strings.Repeat(" ", padding), comment)
	} else {
		fmt.Fprintf(out, "%s\n", line)
	}
}

func formatComment(props apiextensionsv1.JSONSchemaProps) string {
	return strings.ReplaceAll(props.Description, "\n", " ")
}

func sortCategory(name string) int {
	switch name {
	case "enabled":
		return 0
	case "nodeSelector", "tolerations", "affinity", "resources":
		return 2
	default:
		return 1
	}
}

// sort property names: enabled < all other property names < nodeSelector, tolerations, affinity, resources
func sortByCategory(a string, b string) bool {
	catA := sortCategory(a)
	catB := sortCategory(b)
	if catA == catB {
		return a < b
	} else {
		return catA < catB
	}
}
