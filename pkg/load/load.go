// Package load provides functionalities to parse and analyze Go source files to discover and retrieve struct type definitions.
package load

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
)

// Doc: https://github.com/golang/example/blob/master/gotypes/go-types.md

// StructDiscovered represents a discovered struct.
type StructDiscovered struct {
	Name *types.TypeName
	Obj  *types.Struct
}

// GetStructsFromSourceFile finds all structs defined in a Source File.
func GetStructsFromSourceFile(sourceFilePath string) ([]StructDiscovered, error) {

	var structTypes []StructDiscovered

	// Parse the provided source file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, sourceFilePath, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parsed the file, error was: %v", err)
	}

	// Type checks the parsed AST using types.Config.Check
	// A Config controls various options of the type checker.
	// The defaults work fine except for one setting:
	// we must specify how to deal with imports.
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("mypkg", fset, []*ast.File{file}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to type check the file, error was: %v", err)
	}

	// Get the package's scope, containing package-level declarations
	scope := pkg.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		// Check if this is  a type declaration (defined or alias)
		if typeName, ok := obj.(*types.TypeName); ok {
			// Check if the TypeName's underlying type is a Struct
			if structType, ok := typeName.Type().Underlying().(*types.Struct); ok {
				var newStruct StructDiscovered
				newStruct.Name = typeName
				newStruct.Obj = structType
				structTypes = append(structTypes, newStruct)
			}
		}
	}

	if len(structTypes) == 0 {
		return structTypes, fmt.Errorf("no structs found")
	}

	return structTypes, nil
}
