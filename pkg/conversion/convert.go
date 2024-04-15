package conversion

import (
	"bytes"
	"fmt"
	"github.com/VintageOps/structogqlgen/pkg/load"
	"go/token"
	"go/types"
)

// GqlTypeDefinition contains the definition of a graphQl Type
type GqlTypeDefinition struct {
	GqlTypeName string                // GqlTypeName is the name of a graphQL type.
	GqlFields   []GqlFieldsDefinition // GqlFields is a slice of GqlFieldsDefinition, which represents the fields of a GraphQL type.
}

type GqlFieldsDefinition struct {
	GqlFieldName     string              // GqlFieldName represents the name of a graphQL field
	GqlFieldType     string              // GqlFieldType is a string representing the type of GraphQL field
	GqlFieldTags     string              // GqlFieldTags represents the tags of a GraphQL field
	IsScalar         bool                // True if this field need to define a Scalar which will be type Name
	NestedCustomType []GqlTypeDefinition // NestedCustomType represents any custom types that might be needed to be defined for this type.
}

// ConvertCustomError represents a custom error type that can be used in Go programs.
type ConvertCustomError string

func (e ConvertCustomError) Error() string {
	return string(e)
}

const (
	invalidTypeErr = ConvertCustomError("invalid type")
)

// BuildGqlgenType builds a gqlgen Type using a struct definition
func BuildGqlgenType(structDef load.StructDiscovered) (GqlTypeDefinition, error) {

	var gqlTypeDef GqlTypeDefinition

	gqlTypeDef.GqlTypeName = structDef.Name.Id()
	gqlTypeDef.GqlFields = make([]GqlFieldsDefinition, structDef.Obj.NumFields())
	for i := 0; i < structDef.Obj.NumFields(); i++ {
		field := structDef.Obj.Field(i)
		tags := structDef.Obj.Tag(i)
		// Populate Field Name and Tag
		gqlTypeDef.GqlFields[i] = GqlFieldsDefinition{GqlFieldName: field.Name(), GqlFieldTags: tags}
		// Find Field Type and Scalars
		err := ConvertType(field.Type(), &gqlTypeDef.GqlFields[i])
		if err != nil {
			return gqlTypeDef, err
		}
	}

	return gqlTypeDef, nil
}

func ConvertType(goType types.Type, gqlFieldDef *GqlFieldsDefinition) error {
	switch t := goType.(type) {
	case *types.Basic:
		return convertBasicType(t, gqlFieldDef)
	case *types.Slice:
		return convertSliceType(t, gqlFieldDef)
	case *types.Pointer:
		return convertPointerType(t, gqlFieldDef)
	case *types.Map:
		return convertMapType(t, gqlFieldDef)
	case *types.Named:
		return convertNamedType(t, gqlFieldDef)
	case *types.Interface:
		return convertInterfaceType(t, gqlFieldDef)
	default:
		return fmt.Errorf("%s: %v", invalidTypeErr, t.String())
	}
}

func convertBasicType(t *types.Basic, gqlFieldDef *GqlFieldsDefinition) error {
	baseType, err := ConvertBaseType(t)
	if err != nil {
		return err
	}
	gqlFieldDef.GqlFieldType = baseType
	return nil
}

func convertSliceType(t *types.Slice, gqlFieldDef *GqlFieldsDefinition) error {
	var sliceTypeSql GqlFieldsDefinition
	err := ConvertType(t.Elem(), &sliceTypeSql)
	if err != nil {
		return err
	}
	gqlFieldDef.GqlFieldType = fmt.Sprintf("[%s]", sliceTypeSql.GqlFieldType)
	return nil
}

func convertPointerType(t *types.Pointer, gqlFieldDef *GqlFieldsDefinition) error {
	var pointerTypeSql GqlFieldsDefinition
	err := ConvertType(t.Elem(), &pointerTypeSql)
	if err != nil {
		return err
	}
	gqlFieldDef.GqlFieldType = pointerTypeSql.GqlFieldType
	return nil
}

func convertMapType(t *types.Map, gqlFieldDef *GqlFieldsDefinition) error {
	newStructFieldsName := []string{"key", "values"}
	newStructfields := []*types.Var{
		types.NewVar(token.NoPos, nil, newStructFieldsName[0], t.Key()),
		types.NewVar(token.NoPos, nil, newStructFieldsName[1], t.Elem()),
	}
	tags := []string{"", ""}
	structType := types.NewStruct(newStructfields, tags)
	newStructName := fmt.Sprintf("%sMap", gqlFieldDef.GqlFieldName)
	gqlFieldDef.GqlFieldType = newStructName
	newStruct := types.NewNamed(types.NewTypeName(token.NoPos, nil, newStructName, nil), structType, nil)
	var newStructDiscManual load.StructDiscovered
	newStructDiscManual.Name = newStruct.Obj()
	newStructDiscManual.Obj, _ = newStruct.Underlying().(*types.Struct)
	nestStructTypeDef, err := BuildGqlgenType(newStructDiscManual)
	if err != nil {
		return err
	}
	gqlFieldDef.NestedCustomType = append(gqlFieldDef.NestedCustomType, nestStructTypeDef)
	return nil
}

func convertNamedType(t *types.Named, gqlFieldDef *GqlFieldsDefinition) error {
	if _, ok := t.Underlying().(*types.Struct); ok {
		gqlFieldDef.GqlFieldType = t.Obj().Id()
		return nil
	} else {
		gqlFieldDef.GqlFieldType = t.Obj().Name()
		gqlFieldDef.IsScalar = true
	}
	return nil
}

func convertInterfaceType(t *types.Interface, gqlFieldDef *GqlFieldsDefinition) error {
	gqlFieldDef.GqlFieldType = t.String()
	gqlFieldDef.IsScalar = true
	return nil
}

func ConvertBaseType(goBasicType *types.Basic) (string, error) {
	// Mirrors https://pkg.go.dev/go/types#BasicInfo
	mapBasicTypeToGqlType := map[types.BasicInfo]string{
		types.IsBoolean:  "Boolean",
		types.IsInteger:  "Int",
		types.IsUnsigned: "Int",
		types.IsFloat:    "Float",
		// types.IsComplex:  "Float", Not supported for now
		types.IsString: "String",
		// types.IsUntyped: "Unsupported", Not supported for now
	}

	if mapBasicTypeToGqlType[goBasicType.Info()] != "" {
		return mapBasicTypeToGqlType[goBasicType.Info()], nil
	}
	return "", fmt.Errorf("%v: %s", invalidTypeErr, goBasicType.String())
}

func GqlTypePrettyPrint(gqlTypeDefs []GqlTypeDefinition, useTags bool, tagsToUse string) (string, error) {
	var gqlType bytes.Buffer

	// Write the Scalar on top of the string
	scalarOutput := gqlPrettyPrintScalar(gqlTypeDefs)
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

func gqlPrettyPrintScalar(gqlTypeDefs []GqlTypeDefinition) string {
	var gqlScalarType bytes.Buffer

	for _, gqlTypeDef := range gqlTypeDefs {
		for _, field := range gqlTypeDef.GqlFields {
			if field.IsScalar {
				gqlScalarType.WriteString(fmt.Sprintf("scalar %s\n", field.GqlFieldType))
			}
			if len(field.NestedCustomType) != 0 {
				NestedScalar := gqlPrettyPrintScalar(field.NestedCustomType)
				if NestedScalar != "" {
					gqlScalarType.WriteString(NestedScalar)
				}
			}
		}
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
