package conversion

import (
	"bytes"
	"fmt"
	"github.com/VintageOps/structogqlgen/pkg/load"
)

// BuildGqlgenType builds a gqlgen Type using a struct definition
func BuildGqlgenType(structDef load.StructDiscovered) (string, error) {

	var gqlType bytes.Buffer

	gqlType.WriteString(fmt.Sprintf("type %s struct {\n", structDef.Name.Id()))
	for i := 0; i < structDef.Obj.NumFields(); i++ {
		field := structDef.Obj.Field(i)
		tags := structDef.Obj.Tag(i)
		gqlType.WriteString(fmt.Sprintf("  %s: %s! %s\n", field.Name(), field.Type().String(), tags))
	}
	gqlType.WriteString("}\n")

	return gqlType.String(), nil
}
