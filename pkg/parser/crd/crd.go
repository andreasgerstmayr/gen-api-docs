package crd

import (
	"fmt"
	"io"
	"slices"
	"sort"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"

	"github.com/andreasgerstmayr/gen-api-docs/pkg/parser/utils"
)

func formatComment(props apiextensionsv1.JSONSchemaProps) string {
	return strings.ReplaceAll(props.Description, "\n", " ")
}

func getDefaultValue(name string, props apiextensionsv1.JSONSchemaProps) string {
	switch {
	case props.Default != nil:
		b, _ := props.Default.MarshalJSON()
		return string(b)
	case utils.HideCoreTypes && slices.Contains(utils.CoreTypes, name):
		return "{}"
	case props.Type == "string":
		return "\"\""
	case props.Type == "boolean":
		return "false"
	case props.Type == "integer":
		return "0"
	case props.XIntOrString:
		return "0Gi"
	case props.XPreserveUnknownFields != nil && *props.XPreserveUnknownFields: // apiextensionsv1.JSON
		return "{}"
	case props.Type == "object" && props.AdditionalProperties != nil && props.AdditionalProperties.Schema != nil && props.AdditionalProperties.Schema.Type == "string": // map[string]string
		return "{}"
	default:
		return ""
	}
}

func render(out io.Writer, props apiextensionsv1.JSONSchemaProps, level int, isList bool) {
	switch {
	case props.Type == "object" && len(props.Properties) > 0: // struct
		// sort properties by name
		names := make([]string, 0, len(props.Properties))
		for name := range props.Properties {
			names = append(names, name)
		}
		sort.SliceStable(names, func(i int, j int) bool { return utils.SortByCategory(names[i], names[j]) })

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
			value := getDefaultValue(name, member)

			if value != "" {
				utils.WriteLine(out, fmt.Sprintf("%s%s: %s", indent, name, value), comment)
			} else {
				utils.WriteLine(out, fmt.Sprintf("%s%s:", indent, name), comment)
				render(out, member, level+1, false)
			}

			isFirstElem = false
		}

	case props.Type == "object" && props.AdditionalProperties != nil: // map
		indent := strings.Repeat("  ", level)
		comment := formatComment(*props.AdditionalProperties.Schema)
		value := getDefaultValue("", *props.AdditionalProperties.Schema)

		if strings.Contains(props.Description, "Requests describes the minimum amount of compute resources required.") {
			utils.WriteLine(out, fmt.Sprintf("%scpu: \"%s\"", indent, "500m"), comment)
			utils.WriteLine(out, fmt.Sprintf("%smemory: \"%s\"", indent, "1Gi"), comment)
		} else if strings.Contains(props.Description, "Limits describes the maximum amount of compute resources allowed.") {
			utils.WriteLine(out, fmt.Sprintf("%scpu: \"%s\"", indent, "750m"), comment)
			utils.WriteLine(out, fmt.Sprintf("%smemory: \"%s\"", indent, "2Gi"), comment)
		} else if value != "" {
			utils.WriteLine(out, fmt.Sprintf("%s\"key\": %s", indent, value), comment)
		} else {
			utils.WriteLine(out, fmt.Sprintf("%s\"key\":", indent), "")
			render(out, *props.AdditionalProperties.Schema, level+1, false)
		}

	case props.Type == "array":
		indent := strings.Repeat("  ", level-1)
		comment := formatComment(*props.Items.Schema)
		value := getDefaultValue("", *props.Items.Schema)
		if value != "" {
			utils.WriteLine(out, fmt.Sprintf("%s- %s", indent, value), comment)
		} else {
			render(out, *props.Items.Schema, level, true)
		}
	}
}

func Parse(in io.Reader, out io.Writer) error {
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

	utils.WriteLine(out, fmt.Sprintf("apiVersion: %s/%s", spec.Group, version.Name), formatComment(apiVersionProps))
	utils.WriteLine(out, fmt.Sprintf("kind: %s", spec.Names.Kind), formatComment(kindProps))
	utils.WriteLine(out, "metadata:", "")
	utils.WriteLine(out, "  name: example", "")

	delete(schema.Properties, "kind")
	delete(schema.Properties, "apiVersion")
	delete(schema.Properties, "metadata")
	render(out, schema, 0, false)

	return nil
}
