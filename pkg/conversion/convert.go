// Package conversion provides tools to systematically convert Go data structures into GraphQL schema elements
// including types and fields, based on the type information extracted from Go source files using functions
// from package github.com/VintageOps/structogqlgen/pkg/load
package conversion

import (
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

// GqlFieldsDefinition represents the definition of a GraphQL field.
type GqlFieldsDefinition struct {
	GqlFieldName         string                // GqlFieldName represents the name of a graphQL field
	GqlFieldType         string                // GqlFieldType is a string representing the type of GraphQL field
	GqlFieldTags         string                // GqlFieldTags represents the tags of a GraphQL field
	GqlFieldIsEmbedded   bool                  // GqlFieldIsEmbedded represents whether a GraphQL field is an embedded field.
	IsCustomScalar       bool                  // IsCustomScalar is True if this field need to define a Scalar which will be type Name
	NestedCustomType     []GqlTypeDefinition   // NestedCustomType represents any custom types that might be needed to be defined for this type.
	GqlGenFieldsEmbedded []GqlFieldsDefinition // GqlGenFieldsEmbedded represents fields for Embedded Structs
}

// gqlTypeIsCustScalar represents indicates whether a graphql type must be represented as a custom scalar type or not.
type gqlTypeIsCustScalar struct {
	gqlType        string
	isCustomScalar bool
}

// MapBasicKindToGqlType maps Go basic types to corresponding GraphQL types and indicates if that types need a custom scalar.
// Mirrors: https://pkg.go.dev/go/types#BasicKind
var MapBasicKindToGqlType = map[types.BasicKind]gqlTypeIsCustScalar{
	types.Bool:           {gqlType: "Boolean"},
	types.Int:            {gqlType: "Int"},
	types.Int8:           {gqlType: "Int"},
	types.Int16:          {gqlType: "Int"},
	types.Int32:          {gqlType: "Int"},
	types.Int64:          {gqlType: "BigInt", isCustomScalar: true}, // Graphql Int represents a signed 32‐bit integer
	types.Uint:           {gqlType: "Int"},
	types.Uint8:          {gqlType: "Int"},
	types.Uint16:         {gqlType: "Int"},
	types.Uint32:         {gqlType: "Int"},
	types.Uint64:         {gqlType: "BigInt", isCustomScalar: true}, // Graphql Int represents a signed 32‐bit integer
	types.Uintptr:        {gqlType: "BigInt", isCustomScalar: true}, // Graphql Int represents a signed 32‐bit integer
	types.Float32:        {gqlType: "Float"},
	types.Float64:        {gqlType: "Float"},
	types.Complex64:      {gqlType: "ComplexNumber", isCustomScalar: true}, // Custom scalar
	types.Complex128:     {gqlType: "ComplexNumber", isCustomScalar: true}, // Custom scalar
	types.String:         {gqlType: "String"},
	types.UnsafePointer:  {gqlType: "UnsafePointer", isCustomScalar: true}, // Custom Scalar
	types.UntypedBool:    {gqlType: "Boolean"},
	types.UntypedInt:     {gqlType: "BigInt", isCustomScalar: true}, // Custom scalar or String
	types.UntypedRune:    {gqlType: "Int"},
	types.UntypedFloat:   {gqlType: "Float"},
	types.UntypedComplex: {gqlType: "ComplexNumber", isCustomScalar: true}, // Custom scalar
	types.UntypedString:  {gqlType: "String"},
	types.UntypedNil:     {gqlType: "UntypedNil", isCustomScalar: true}, // Special handling may be needed
}

// ConvertCustomError represents a custom error type that can be used in Go programs.
type ConvertCustomError string

func (e ConvertCustomError) Error() string {
	return string(e)
}

// InvalidTypeErr represents an error indicating an invalid type.
const (
	InvalidTypeErr = ConvertCustomError("invalid type")
)

// BuildGqlTypes builds an array of GqlTypeDefinitions for a given array of struct definitions.
// It calls BuildGqlgenType for each struct definition and populates the array with the results.
// If any error occurs during the process, it returns the error immediately.
func BuildGqlTypes(structsFound []load.StructDiscovered) ([]GqlTypeDefinition, error) {
	gqlGenTypes := make([]GqlTypeDefinition, len(structsFound))
	for idx, structType := range structsFound {
		var err error
		gqlGenTypes[idx], err = BuildGqlgenType(structType)
		if err != nil {
			return nil, err
		}
	}
	return gqlGenTypes, nil
}

// BuildGqlgenType builds a GqlTypeDefinition for a given struct definition.
// It converts the struct fields into GqlFieldsDefinition, populating the field name and tags.
// It also determines the field type by invoking ConvertType and handles any custom types or scalars.
func BuildGqlgenType(structDef load.StructDiscovered) (GqlTypeDefinition, error) {

	var gqlTypeDef GqlTypeDefinition

	gqlTypeDef.GqlTypeName = structDef.Name.Id()
	gqlTypeDef.GqlFields = make([]GqlFieldsDefinition, structDef.Obj.NumFields())
	for i := 0; i < structDef.Obj.NumFields(); i++ {
		field := structDef.Obj.Field(i)
		tags := structDef.Obj.Tag(i)
		isEmbedded := field.Embedded()
		// Populate Field Name and Tag
		gqlTypeDef.GqlFields[i] = GqlFieldsDefinition{GqlFieldName: field.Name(), GqlFieldTags: tags, GqlFieldIsEmbedded: isEmbedded}
		// Find Field Type and Scalars
		err := ConvertType(field.Type(), &gqlTypeDef.GqlFields[i])
		if err != nil {
			return gqlTypeDef, err
		}
	}

	return gqlTypeDef, nil
}

// ConvertType converts a Go type into a GqlFieldsDefinition by performing type-specific conversions.
// It handles basic types, slices, pointers, maps, named types, and interfaces. .
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
		return fmt.Errorf("%s: %v", InvalidTypeErr, t.String())
	}
}

// convertBasicType converts a Go basic type into a GqlFieldsDefinition by mapping it to a GraphQL type.
func convertBasicType(t *types.Basic, gqlFieldDef *GqlFieldsDefinition) error {

	if t.Kind() == types.Invalid {
		return fmt.Errorf("%v: %s", InvalidTypeErr, t.String())
	}

	if val, ok := MapBasicKindToGqlType[t.Kind()]; ok {
		gqlFieldDef.GqlFieldType = val.gqlType
		gqlFieldDef.IsCustomScalar = val.isCustomScalar
		return nil
	}

	return fmt.Errorf("%v: %s", InvalidTypeErr, t.String())
}

// convertSliceType converts a Go type representing a slice into a GqlFieldsDefinition.
func convertSliceType(t *types.Slice, gqlFieldDef *GqlFieldsDefinition) error {
	var sliceTypeSql GqlFieldsDefinition
	err := ConvertType(t.Elem(), &sliceTypeSql)
	if err != nil {
		return err
	}
	gqlFieldDef.GqlFieldType = fmt.Sprintf("[%s]", sliceTypeSql.GqlFieldType)
	return nil
}

// convertPointerType converts a pointer type into a GqlFieldsDefinition.
func convertPointerType(t *types.Pointer, gqlFieldDef *GqlFieldsDefinition) error {
	var pointerTypeSql GqlFieldsDefinition
	err := ConvertType(t.Elem(), &pointerTypeSql)
	if err != nil {
		return err
	}
	gqlFieldDef.GqlFieldType = pointerTypeSql.GqlFieldType
	return nil
}

// convertMapType converts a Go map type into a GqlFieldsDefinition representing a struct.
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

// convertNamedType converts a named type into a GqlFieldsDefinition.
func convertNamedType(t *types.Named, gqlFieldDef *GqlFieldsDefinition) error {
	if ts, ok := t.Underlying().(*types.Struct); ok {
		gqlFieldDef.GqlFieldType = t.Obj().Id()
		// If the field is embedded, then need to populate
		if gqlFieldDef.GqlFieldIsEmbedded {
			var newStructDiscManual load.StructDiscovered
			newStructDiscManual.Name = t.Obj()
			newStructDiscManual.Obj = ts
			nestStructTypeDef, err := BuildGqlgenType(newStructDiscManual)
			if err != nil {
				return err
			}
			gqlFieldDef.GqlGenFieldsEmbedded = nestStructTypeDef.GqlFields
		}
		return nil
	} else {
		gqlFieldDef.GqlFieldType = t.Obj().Name()
		gqlFieldDef.IsCustomScalar = true
	}
	return nil
}

// convertInterfaceType converts a *types.Interface into a GqlFieldsDefinition.
func convertInterfaceType(t *types.Interface, gqlFieldDef *GqlFieldsDefinition) error {
	if t.Empty() {
		// Empty Interface
		gqlFieldDef.GqlFieldType = "interfaceEmpty"
	} else {
		gqlFieldDef.GqlFieldType = fmt.Sprintf("interface%s", gqlFieldDef.GqlFieldName)
	}
	gqlFieldDef.IsCustomScalar = true
	return nil
}
