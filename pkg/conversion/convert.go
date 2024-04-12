package conversion

import (
	"bytes"
	"fmt"
	"go/types"
)

// BuildGqlgenType builds a gqlgen Type using a struct definition
func BuildGqlgenType(structDef *types.Struct) (string, error) {

	var gqlType bytes.Buffer

	for i := 0; i < structDef.NumFields(); i++ {
		field := structDef.Field(i)
		if field.Exported() {
			gqlType.WriteString(fmt.Sprintf("%s: %s!\n", field.Name(), field.Type().String()))
		}
	}

	return gqlType.String(), nil
}
