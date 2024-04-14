package conversion

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/VintageOps/structogqlgen/pkg/load"
	"go/types"
)

// ConvertCustomError represents a custom error type that can be used in Go programs.
type ConvertCustomError string

func (e ConvertCustomError) Error() string {
	return string(e)
}

// nestedStructErr is a constant of type ConvertCustomError that represents an error
// indicating that a nested struct has occurred.
const (
	nestedStructErr = ConvertCustomError("nested struct")
	invalidTypeErr  = ConvertCustomError("invalid type")
)

// BuildGqlgenType builds a gqlgen Type using a struct definition
func BuildGqlgenType(structDef load.StructDiscovered) (string, error) {

	var gqlType bytes.Buffer

	gqlType.WriteString(fmt.Sprintf("type %s struct {\n", structDef.Name.Id()))
	for i := 0; i < structDef.Obj.NumFields(); i++ {
		field := structDef.Obj.Field(i)
		tags := structDef.Obj.Tag(i)
		gqlFieldType, err := ConvertType(field.Type())
		if err != nil && !errors.Is(err, nestedStructErr) {
			//return "", err
			fmt.Println(err)
		}
		if errors.Is(err, nestedStructErr) {
			// Nested Struct
			gqlType.WriteString(fmt.Sprintf("  %s: %s! %s\n", field.Name(), field.Type().String(), tags))
		} else {
			gqlType.WriteString(fmt.Sprintf("  %s: %s! %s\n", field.Name(), gqlFieldType, tags))
		}
	}
	gqlType.WriteString("}\n")

	return gqlType.String(), nil
}

func ConvertType(goType types.Type) (string, error) {
	switch t := goType.(type) {
	case *types.Basic:
		return ConvertBaseType(t)
	case *types.Slice:
		underlyingType, err := ConvertType(t.Elem())
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[%s]", underlyingType), nil
	case *types.Pointer:
		underlyingType, err := ConvertType(t.Elem())
		if err != nil {
			return "", err
		}
		return underlyingType, nil
	case *types.Map:
		// Need to define new graphql type with key being of type key and value of type value
	case *types.Interface:
		// Need to define a Scalar that must be defined by the requester
	case *types.Named:
		return "", nestedStructErr
	default:
		return "", invalidTypeErr
	}
	return "", invalidTypeErr
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
	return "", fmt.Errorf("unsupported basic type: %s", goBasicType.String())
}
