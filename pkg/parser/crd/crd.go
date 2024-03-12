package crd

import (
	"fmt"
	"io"
	"slices"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"

	"github.com/andreasgerstmayr/gen-api-docs/pkg/parser/format"
)

func getDefaultValue(name string, props apiextensionsv1.JSONSchemaProps) string {
	switch {
	case props.Default != nil:
		b, _ := props.Default.MarshalJSON()
		return string(b)
	case format.HideCoreTypes && slices.Contains(format.CoreTypes, name):
		return "{}"
	case props.XIntOrString:
		return "0Gi"
	case props.Format == "date-time":
		return "\"2006-01-02T15:04:05Z\""
	case props.XPreserveUnknownFields != nil && *props.XPreserveUnknownFields: // apiextensionsv1.JSON
		return "{}"
	case props.Type == "object" && props.AdditionalProperties != nil && props.AdditionalProperties.Schema != nil && props.AdditionalProperties.Schema.Type == "string": // map[string]string
		return "{}"
	case props.Type == "string":
		return "\"\""
	case props.Type == "boolean":
		return "false"
	case props.Type == "integer":
		return "0"
	default:
		return ""
	}
}

func render(schemaProps apiextensionsv1.JSONSchemaProps, prop *format.Prop) {
	switch {
	case schemaProps.Type == "object" && len(schemaProps.Properties) > 0: // struct
		for name, member := range schemaProps.Properties {
			p := &format.Prop{
				Key:         name,
				Comment:     member.Description,
				ScalarValue: getDefaultValue(name, member),
			}
			prop.Properties = append(prop.Properties, p)
			if p.ScalarValue == "" {
				render(member, p)
			}
		}

	case schemaProps.Type == "object" && schemaProps.AdditionalProperties != nil: // map
		p := &format.Prop{
			Key:         "\"key\"",
			Comment:     schemaProps.AdditionalProperties.Schema.Description,
			ScalarValue: getDefaultValue("", *schemaProps.AdditionalProperties.Schema),
		}

		if strings.Contains(schemaProps.Description, "Requests describes the minimum amount of compute resources required.") {
			prop.Properties = []*format.Prop{
				{Key: "cpu", ScalarValue: "\"500m\""},
				{Key: "memory", ScalarValue: "\"1Gi\""},
			}
		} else if strings.Contains(schemaProps.Description, "Limits describes the maximum amount of compute resources allowed.") {
			prop.Properties = []*format.Prop{
				{Key: "cpu", ScalarValue: "\"750m\""},
				{Key: "memory", ScalarValue: "\"2Gi\""},
			}
		} else {
			prop.Properties = []*format.Prop{p}
			if p.ScalarValue == "" {
				render(*schemaProps.AdditionalProperties.Schema, p)
			}
		}

	case schemaProps.Type == "array":
		prop.ListItem = &format.Prop{
			Comment:     schemaProps.Items.Schema.Description,
			ScalarValue: getDefaultValue("", *schemaProps.Items.Schema),
		}
		if prop.ListItem.ScalarValue == "" {
			render(*schemaProps.Items.Schema, prop.ListItem)
		}
	}
}

func Parse(in io.Reader) (format.Prop, error) {
	crd := apiextensionsv1.CustomResourceDefinition{}
	err := k8syaml.NewYAMLOrJSONDecoder(in, 8192).Decode(&crd)
	if err != nil {
		return format.Prop{}, err
	}

	prop := format.Prop{}
	spec := crd.Spec
	version := spec.Versions[0]
	schema := *version.Schema.OpenAPIV3Schema

	apiVersionProps := schema.Properties["apiVersion"]
	prop.Properties = append(prop.Properties, &format.Prop{
		Key:         "apiVersion",
		ScalarValue: fmt.Sprintf("%s/%s", spec.Group, version.Name),
		Comment:     apiVersionProps.Description,
	})
	delete(schema.Properties, "apiVersion")

	kindProps := schema.Properties["kind"]
	prop.Properties = append(prop.Properties, &format.Prop{
		Key:         "kind",
		ScalarValue: spec.Names.Kind,
		Comment:     kindProps.Description,
	})
	delete(schema.Properties, "kind")

	prop.Properties = append(prop.Properties, &format.Prop{
		Key: "metadata",
		Properties: []*format.Prop{
			{Key: "name", ScalarValue: "example"},
		},
	})
	delete(schema.Properties, "metadata")

	render(schema, &prop)
	return prop, nil
}
