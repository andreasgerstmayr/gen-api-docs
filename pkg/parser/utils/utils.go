package utils

import (
	"fmt"
	"io"
	"strings"
)

var (
	CommentPadding int
	HideCoreTypes  bool

	CoreTypes = []string{"tolerations", "nodeAffinity", "podAffinity", "podAntiAffinity"}
)

func WriteLine(out io.Writer, line string, comment string) {
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
func SortByCategory(a string, b string) bool {
	catA := sortCategory(a)
	catB := sortCategory(b)
	if catA == catB {
		return a < b
	} else {
		return catA < catB
	}
}
