package conversion

import (
	"bytes"
	"fmt"
)

func GqlPrettyPrint(gqlTypeDefs []GqlTypeDefinition, useTags bool, tagsToUse string) (string, error) {
	var gqlType bytes.Buffer

	// Write the Scalar on top of the string
	scalarOutput := gqlPrettyPrintScalar(gqlTypeDefs, nil)
	if scalarOutput != "" {
		gqlType.WriteString(scalarOutput)
	}
	gqlType.WriteString("\n")

	// Write the Type Definition
	gqlTypeDefinition, err := gqlPrettyPrintTypes(gqlTypeDefs, useTags, tagsToUse)
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
				fmt.Println(setScalar)
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

func gqlPrettyPrintTypes(gqlTypeDefs []GqlTypeDefinition, useTags bool, tagsToUse string) (string, error) {
	var gqlType bytes.Buffer
	for _, gqlTypeDef := range gqlTypeDefs {
		var nestedCustomToWrite string
		var err error
		gqlType.WriteString(fmt.Sprintf("type %s {\n", gqlTypeDef.GqlTypeName))
		for _, field := range gqlTypeDef.GqlFields {
			if useTags {
				// Need to factor in here to make use of TagsToUse and err if GqlFieldTags is not found
				gqlType.WriteString(fmt.Sprintf("  %s: %s\n", field.GqlFieldTags, field.GqlFieldType))
			} else {
				gqlType.WriteString(fmt.Sprintf("  %s: %s\n", field.GqlFieldName, field.GqlFieldType))
			}
			if len(field.NestedCustomType) != 0 {
				nestedCustomToWrite, err = gqlPrettyPrintTypes(field.NestedCustomType, useTags, tagsToUse)
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
