package conversion

import (
	"bytes"
	"fmt"
	"github.com/fatih/structtag"
)

// PrettyPrintOptions represents the options for pretty-printing. It contains the following fields:
// - UseJsonTags: a bool indicating whether to use JSON tags
// - UseCustomTags: a string indicating the custom tags to use
// - RequireTags: a SpecTagRequire struct that specifies required tags
type PrettyPrintOptions struct {
	UseJsonTags   bool
	UseCustomTags string
	RequireTags   SpecTagRequire
}

// SpecTagRequire defines the structure for specifying required tags.
type SpecTagRequire struct {
	Key string
	Val string
}

// tagToUse returns the tag that should be used for field definitions.
func (opts *PrettyPrintOptions) tagToUse() string {
	if opts.UseCustomTags != "" {
		return opts.UseCustomTags
	}
	if opts.UseJsonTags {
		return "json"
	}
	return ""
}

// GqlPrettyPrint takes a slice of GqlTypeDefinition and PrettyPrintOptions and returns a string representation of the GraphQL type definitions.
func GqlPrettyPrint(gqlTypeDefs []GqlTypeDefinition, opts *PrettyPrintOptions) (string, error) {
	var gqlType bytes.Buffer

	// Write the Scalar on top of the string
	scalarOutput := gqlPrettyPrintScalar(gqlTypeDefs, nil)
	if scalarOutput != "" {
		gqlType.WriteString(scalarOutput)
	}
	gqlType.WriteString("\n")

	// Write the Type Definition
	gqlTypeDefinition, err := gqlPrettyPrintTypes(gqlTypeDefs, opts)
	if err != nil {
		return "", err
	}
	gqlType.WriteString(gqlTypeDefinition)
	return gqlType.String(), nil
}

// gqlPrettyPrintScalar takes a slice of GqlTypeDefinition and a map of setScalar.
// It returns a string representation of the GraphQL scalar type definitions.
func gqlPrettyPrintScalar(gqlTypeDefs []GqlTypeDefinition, setScalar map[string]bool) string {
	var gqlScalarType bytes.Buffer

	// Dealing with recursive Case
	if setScalar == nil {
		setScalar = make(map[string]bool)
	}

	for _, gqlTypeDef := range gqlTypeDefs {
		for _, field := range gqlTypeDef.GqlFields {
			if field.IsCustomScalar {
				setScalar[field.GqlFieldType] = true
			}
			if len(field.NestedCustomType) != 0 {
				_ = gqlPrettyPrintScalar(field.NestedCustomType, setScalar)
			}
		}
	}

	for scalar := range setScalar {
		gqlScalarType.WriteString(fmt.Sprintf("scalar %s\n", scalar))
	}
	return gqlScalarType.String()
}

// gqlPrettyPrintTypes takes a slice of GqlTypeDefinition and PrettyPrintOptions and returns a string representation of the GraphQL type definitions.
func gqlPrettyPrintTypes(gqlTypeDefs []GqlTypeDefinition, opts *PrettyPrintOptions) (string, error) {
	var gqlType bytes.Buffer
	anyTagToUse := opts.tagToUse()

	for _, gqlTypeDef := range gqlTypeDefs {
		var nestedCustomToWrite string
		gqlType.WriteString(fmt.Sprintf("type %s {\n", gqlTypeDef.GqlTypeName))

		for _, field := range gqlTypeDef.GqlFields {
			fieldDef, err := gqlCreateFieldDefinition(field, anyTagToUse, &opts.RequireTags)
			if err != nil {
				return "", err
			}
			gqlType.WriteString(fieldDef)

			if len(field.NestedCustomType) != 0 {
				nestedCustomToWrite, err = gqlPrettyPrintTypes(field.NestedCustomType, opts)
				if err != nil {
					return "", err
				}
			}
		}
		gqlType.WriteString("}\n\n")
		if nestedCustomToWrite != "" {
			gqlType.WriteString(nestedCustomToWrite)
		}
	}
	return gqlType.String(), nil
}

// gqlCreateFieldDefinition takes a GqlFieldsDefinition, a tag string, and a SpecTagRequire
// and returns a string representation of the GraphQL field definition.
func gqlCreateFieldDefinition(field GqlFieldsDefinition, tag string, requiredTags *SpecTagRequire) (string, error) {
	fieldName := field.GqlFieldName
	requiredFieldmark := ""

	tags, err := parseFieldTags(field)
	if err != nil {
		return "", err
	}

	fieldName, err = updateFieldName(fieldName, tags, tag)
	if err != nil {
		return "", err
	}

	requiredFieldmark, err = updateRequiredFieldMark(tags, requiredTags, requiredFieldmark)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("  %s: %s%s\n", fieldName, field.GqlFieldType, requiredFieldmark), nil
}

// parseFieldTags takes a GqlFieldsDefinition and returns the parsed struct tags using structtag.Parse.
func parseFieldTags(field GqlFieldsDefinition) (tags *structtag.Tags, err error) {
	tags, err = structtag.Parse(field.GqlFieldTags)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// updateFieldName update the field name to output based on the provided tag to use
func updateFieldName(fieldName string, tags *structtag.Tags, tag string) (string, error) {
	if tag != "" {
		specifiedTag, err := tags.Get(tag)
		if err != nil {
			if err.Error() != "tag does not exist" {
				return "", err
			}
		} else {
			fieldName = specifiedTag.Name
		}
	}
	return fieldName, nil
}

// updateRequiredFieldMark appends the fields "!" if the field has a tag that was marked as required
func updateRequiredFieldMark(tags *structtag.Tags, requiredTags *SpecTagRequire, requiredFieldmark string) (string, error) {
	if requiredTags.Key != "" && requiredTags.Val != "" {
		tagValue, err := tags.Get(requiredTags.Key)
		if err != nil {
			if fmt.Sprintf("%v", err) != "tag does not exist" {
				return "", err
			}
		}
		if err == nil && tagValue.Name == requiredTags.Val {
			requiredFieldmark = "!"
		}
	}
	return requiredFieldmark, nil
}
