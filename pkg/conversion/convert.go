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

type GqlFieldsDefinition struct {
	GqlFieldName     string              // GqlFieldName represents the name of a graphQL field
	GqlFieldType     string              // GqlFieldType is a string representing the type of GraphQL field
	GqlFieldTags     string              // GqlFieldTags represents the tags of a GraphQL field
	IsCustomScalar   bool                // True if this field need to define a Scalar which will be type Name
	NestedCustomType []GqlTypeDefinition // NestedCustomType represents any custom types that might be needed to be defined for this type.
}

type gqlTypeIsCustScalar struct {
	gqlType        string
	isCustomScalar bool
}

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

	if t.Kind() == types.Invalid {
		return fmt.Errorf("%v: %s", invalidTypeErr, t.String())
	}

	if val, ok := MapBasicKindToGqlType[t.Kind()]; ok {
		gqlFieldDef.GqlFieldType = val.gqlType
		gqlFieldDef.IsCustomScalar = val.isCustomScalar
		return nil
	}

	return fmt.Errorf("%v: %s", invalidTypeErr, t.String())
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
		gqlFieldDef.IsCustomScalar = true
	}
	return nil
}

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
