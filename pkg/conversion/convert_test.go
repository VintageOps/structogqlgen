package conversion

import (
	"testing"

	"github.com/VintageOps/structogqlgen/pkg/load"
	"go/token"
	"go/types"
)

// TestBuildGqlgenType is a unit test for the BuildGqlgenType function.
func TestBuildGqlgenType(t *testing.T) {
	tests := []struct {
		name      string
		structDef func() load.StructDiscovered
		wantErr   bool
	}{
		{
			name: "EmptyStruct",
			structDef: func() load.StructDiscovered {
				return load.StructDiscovered{
					Name: types.NewTypeName(token.NoPos, nil, "EmptyStruct", types.Typ[types.Invalid]),
					Obj:  types.NewStruct(nil, nil),
				}
			},
			wantErr: false,
		},
		{
			name: "StructWithFields",
			structDef: func() load.StructDiscovered {
				pkg := types.NewPackage("some.pkg/path", "path")
				variable := types.NewVar(token.NoPos, pkg, "someField", types.Typ[types.Bool])
				structType := types.NewStruct([]*types.Var{variable}, []string{"tagValue"})
				return load.StructDiscovered{
					Name: types.NewTypeName(token.NoPos, nil, "StructWithFields", types.Typ[types.Invalid]),
					Obj:  structType,
				}
			},
			wantErr: false,
		},
		{
			name: "ErrorHandlingInConvertType",
			structDef: func() load.StructDiscovered {
				pkg := types.NewPackage("some.pkg/path", "path")
				chanType := types.NewChan(types.SendRecv, types.Typ[types.Bool])
				variable := types.NewVar(token.NoPos, pkg, "invalidField", chanType)
				structType := types.NewStruct([]*types.Var{variable}, []string{"tagValue"})
				return load.StructDiscovered{
					Name: types.NewTypeName(token.NoPos, nil, "ErrorHandlingInConvertType", types.Typ[types.Invalid]),
					Obj:  structType,
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildGqlgenType(tt.structDef())
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildGqlgenType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
