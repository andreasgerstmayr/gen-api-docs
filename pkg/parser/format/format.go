package format

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

var (
	Format         string
	CommentPadding int
	HideCoreTypes  bool

	CoreTypes = []string{"tolerations", "nodeAffinity", "podAffinity", "podAntiAffinity"}
)

type Prop struct {
	Key     string
	Comment []string

	// below types are exclusive
	ScalarValue string
	Properties  []*Prop
	ListItem    *Prop
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

func writeCommentLine(out io.Writer, line string, commentLines []string) {
	padding := CommentPadding - len(line)
	if padding < 0 {
		padding = 0
	}

	comment := strings.Join(commentLines, " ")
	if comment != "" {
		fmt.Fprintf(out, "%s%s # %s\n", line, strings.Repeat(" ", padding), comment)
	} else {
		fmt.Fprintf(out, "%s\n", line)
	}
}

func printOneline(out io.Writer, p *Prop, level int, isList bool) {
	switch {
	// struct
	case len(p.Properties) > 0:
		sort.SliceStable(p.Properties, func(i int, j int) bool { return sortByCategory(p.Properties[i].Key, p.Properties[j].Key) })
		for i, p := range p.Properties {
			var indent string
			if i == 0 && isList {
				indent = strings.Repeat("  ", level-1) + "- "
			} else {
				indent = strings.Repeat("  ", level)
			}

			if p.ScalarValue != "" {
				writeCommentLine(out, fmt.Sprintf("%s%s: %s", indent, p.Key, p.ScalarValue), p.Comment)
			} else {
				writeCommentLine(out, fmt.Sprintf("%s%s:", indent, p.Key), p.Comment)
				printOneline(out, p, level+1, false)
			}
		}

	// list
	case p.ListItem != nil:
		indent := strings.Repeat("  ", level-1)
		if p.ListItem.ScalarValue != "" {
			writeCommentLine(out, fmt.Sprintf("%s- %s", indent, p.ListItem.ScalarValue), p.ListItem.Comment)
		} else {
			printOneline(out, p.ListItem, level, true)
		}
	}
}

func writeMultilineComment(out io.Writer, commentLines []string, indent string) {
	if len(commentLines) > 0 {
		fmt.Fprintf(out, "\n")
	}

	for _, line := range commentLines {
		if len(line) > 0 {
			fmt.Fprintf(out, "%s# %s\n", indent, line)
		} else {
			fmt.Fprintf(out, "%s#\n", indent)
		}
	}
}

func printMultiline(out io.Writer, p *Prop, level int, isList bool) {
	switch {
	// struct
	case len(p.Properties) > 0:
		sort.SliceStable(p.Properties, func(i int, j int) bool { return sortByCategory(p.Properties[i].Key, p.Properties[j].Key) })
		for i, p := range p.Properties {
			var indent string
			if i == 0 && isList {
				indent = strings.Repeat("  ", level-1) + "- "
			} else {
				indent = strings.Repeat("  ", level)
			}

			writeMultilineComment(out, p.Comment, strings.Repeat("  ", level))
			if p.ScalarValue != "" {
				fmt.Fprintf(out, "%s%s: %s\n", indent, p.Key, p.ScalarValue)
			} else {
				fmt.Fprintf(out, "%s%s:\n", indent, p.Key)
				printMultiline(out, p, level+1, false)
			}
		}

	// list
	case p.ListItem != nil:
		indent := strings.Repeat("  ", level-1)
		if p.ListItem.ScalarValue != "" {
			writeMultilineComment(out, p.Comment, indent)
			fmt.Fprintf(out, "%s- %s\n", indent, p.ListItem.ScalarValue)
		} else {
			printMultiline(out, p.ListItem, level, true)
		}
	}
}

func Print(out io.Writer, p *Prop) {
	switch Format {
	case "oneline":
		printOneline(out, p, 0, false)
	case "multiline":
		printMultiline(out, p, 0, false)
	default:
		fmt.Fprintf(out, "invalid format: %s", Format)
	}
}
