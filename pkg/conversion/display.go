package conversion

import (
	"bytes"
	"fmt"
	"github.com/fatih/structtag"
)

type PrettyPrintOptions struct {
	UseJsonTags   bool
	UseCustomTags string
	RequireTags   specTagRequire
}

type specTagRequire struct {
	Key string
	Val string
}

func (opts *PrettyPrintOptions) tagToUse() string {
	if opts.UseCustomTags != "" {
		return opts.UseCustomTags
	}
	if opts.UseJsonTags {
		return "json"
	}
	return ""
}

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

func gqlPrettyPrintTypes(gqlTypeDefs []GqlTypeDefinition, opts *PrettyPrintOptions) (string, error) {
	var gqlType bytes.Buffer
	anyTagToUse := opts.tagToUse()

	for _, gqlTypeDef := range gqlTypeDefs {
		var nestedCustomToWrite string
		var err error
		gqlType.WriteString(fmt.Sprintf("type %s {\n", gqlTypeDef.GqlTypeName))
		for _, field := range gqlTypeDef.GqlFields {
			if anyTagToUse != "" {
				tags, err := structtag.Parse(field.GqlFieldTags)
				if err != nil {
					return "", err
				}
				specifiedTag, err := tags.Get(anyTagToUse)
				if err != nil && fmt.Sprintf("%v", err) != "tag does not exist" {
					return "", err
				}
				if fmt.Sprintf("%v", err) == "tag does not exist" {
					gqlType.WriteString(fmt.Sprintf("  %s: %s\n", field.GqlFieldName, field.GqlFieldType))
				} else {
					gqlType.WriteString(fmt.Sprintf("  %s: %s\n", specifiedTag.Name, field.GqlFieldType))
				}
			} else {
				gqlType.WriteString(fmt.Sprintf("  %s: %s\n", field.GqlFieldName, field.GqlFieldType))
			}
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
