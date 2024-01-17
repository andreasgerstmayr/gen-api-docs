package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

const (
	COMMENT_PADDING = 40
)

func writeLine(out io.Writer, line string, comment string) {
	padding := COMMENT_PADDING - len(line)
	if padding < 0 {
		padding = 0
	}

	if comment != "" {
		fmt.Fprintf(out, "%s%s # %s\n", line, strings.Repeat(" ", padding), comment)
	} else {
		fmt.Fprintf(out, "%s\n", line)
	}
}

func getDefaultValue(props apiextensionsv1.JSONSchemaProps) string {
	switch {
	case props.Default != nil:
		b, _ := props.Default.MarshalJSON()
		return string(b)
	case props.Type == "string":
		return "\"\""
	case props.Type == "boolean":
		return "false"
	case props.Type == "integer":
		return "0"
	case props.XIntOrString:
		return "1Gi"
	case props.XPreserveUnknownFields != nil && *props.XPreserveUnknownFields:
		return "{}"
	default:
		return ""
	}
}

func formatComment(props apiextensionsv1.JSONSchemaProps) string {
	return strings.ReplaceAll(props.Description, "\n", " ")
}

func render(out io.Writer, props apiextensionsv1.JSONSchemaProps, level int, isList bool) {
	switch {
	case props.Type == "object" && len(props.Properties) > 0: // struct
		// sort properties by name
		names := make([]string, 0, len(props.Properties))
		for name := range props.Properties {
			names = append(names, name)
		}
		sort.Strings(names)

		isFirstElem := true
		for _, name := range names {
			member := props.Properties[name]
			var indent string
			if isFirstElem && isList {
				indent = strings.Repeat("  ", level-1) + "- "
			} else {
				indent = strings.Repeat("  ", level)
			}
			comment := formatComment(member)
			value := getDefaultValue(member)

			if value != "" {
				writeLine(out, fmt.Sprintf("%s%s: %s", indent, name, value), comment)
			} else {
				writeLine(out, fmt.Sprintf("%s%s:", indent, name), comment)
				render(out, member, level+1, false)
			}

			isFirstElem = false
		}
	case props.Type == "object" && props.AdditionalProperties != nil: // map
		indent := strings.Repeat("  ", level)
		comment := formatComment(*props.AdditionalProperties.Schema)
		value := getDefaultValue(*props.AdditionalProperties.Schema)
		if value != "" {
			writeLine(out, fmt.Sprintf("%s\"key\": %s", indent, value), comment)
		} else {
			writeLine(out, fmt.Sprintf("%s\"key\":", indent), "")
			render(out, *props.AdditionalProperties.Schema, level+1, false)
		}

	case props.Type == "array":
		indent := strings.Repeat("  ", level-1)
		comment := formatComment(*props.Items.Schema)
		value := getDefaultValue(*props.Items.Schema)
		if value != "" {
			writeLine(out, fmt.Sprintf("%s- %s", indent, value), comment)
		} else {
			render(out, *props.Items.Schema, level, true)
		}
	}
}

func run(in io.Reader, out io.Writer) error {
	crd := apiextensionsv1.CustomResourceDefinition{}
	err := k8syaml.NewYAMLOrJSONDecoder(in, 8192).Decode(&crd)
	if err != nil {
		return err
	}

	spec := crd.Spec
	version := spec.Versions[0]
	schema := *version.Schema.OpenAPIV3Schema
	apiVersionProps := schema.Properties["apiVersion"]
	kindProps := schema.Properties["kind"]

	writeLine(out, fmt.Sprintf("apiVersion: %s/%s", spec.Group, version.Name), formatComment(apiVersionProps))
	writeLine(out, fmt.Sprintf("kind: %s", spec.Names.Kind), formatComment(kindProps))
	writeLine(out, "metadata:", "")
	writeLine(out, "  name: example", "")

	delete(schema.Properties, "kind")
	delete(schema.Properties, "apiVersion")
	delete(schema.Properties, "metadata")
	render(out, schema, 0, false)

	return nil
}

func main() {
	err := run(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
