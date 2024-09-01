// Package load provides functionalities to parse and analyze Go source files to discover and retrieve struct type definitions.
package load

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"path/filepath"
	"strings"
)

/* Docs:
   https://github.com/golang/example/blob/master/gotypes/go-types.md
   https://pkg.go.dev/golang.org/x/tools/go/packages#Package
*/

// StructDiscovered represents a discovered struct.
type StructDiscovered struct {
	Name    *types.TypeName
	Obj     *types.Struct
	PkgName string
}

// GetStructsFromPath find all structs defined in a specific directory path
func GetStructsFromPath(path string) ([]StructDiscovered, error) {

	var structTypes []StructDiscovered

	// Determine the absolute path or resolve it relative to the module root
	absPkgPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Error resolving package path: %v", err)
	}

	// Configure the loader to load Go source files and their syntax trees.
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedFiles | packages.NeedSyntax | packages.NeedName,
		Dir:  filepath.Dir(absPkgPath), // Set the Dir to the directory containing the package
	}

	// Load the packages in the specified directory.
	pkgs, err := packages.Load(cfg, absPkgPath)
	if err != nil {
		return nil, err
	}

	// Iterate over the loaded packages.
	for _, pkg := range pkgs {
		if pkg.Errors != nil || len(pkg.Errors) > 0 {
			var allErrors []string
			for _, strErr := range pkg.Errors {
				allErrors = append(allErrors, strErr.Msg)
			}
			return nil, fmt.Errorf("packages on path %s contain the following errors: %v", path, allErrors)
		}
		newStructs := populateStructDiscoveredInPkg(pkg)
		structTypes = append(structTypes, newStructs...)
	}

	return structTypes, nil
}

// populateStructDiscoveredInPkg returns a list of all structs discovered in a given package
func populateStructDiscoveredInPkg(pkg *packages.Package) []StructDiscovered {

	var newStructs []StructDiscovered

	// Nested Function to Traverse AST and find struct types.
	astInspectionFunc := func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if ok {
			// Check if this is a struct
			_, ok := typeSpec.Type.(*ast.StructType)
			if ok {
				// Lookup the type in the package's type information.
				if pkg.TypesInfo != nil {
					obj := pkg.TypesInfo.Defs[typeSpec.Name]
					if obj, ok := obj.(*types.TypeName); ok {
						if structObj, ok := obj.Type().Underlying().(*types.Struct); ok {
							// Create an instance of StructDiscovered and append to the list.
							discovered := StructDiscovered{
								Name:    obj,
								Obj:     structObj,
								PkgName: pkg.Name,
							}
							newStructs = append(newStructs, discovered)

							// Now let's check for anonymous structs within this struct
							newStructs = append(newStructs, findAnonymousStructs(structObj, pkg, obj.Name())...)
						}
					}
				}
			}
		}
		return true
	}

	// Run the inspection
	for _, syntax := range pkg.Syntax {
		ast.Inspect(syntax, astInspectionFunc)
	}

	return newStructs
}

// findAnonymousStructs traverses the fields of a struct and finds any anonymous structs
func findAnonymousStructs(structObj *types.Struct, pkg *packages.Package, parentStructName string) []StructDiscovered {
	var anonymousStructs []StructDiscovered

	// Iterate over the fields of the struct
	for i := 0; i < structObj.NumFields(); i++ {
		field := structObj.Field(i)
		if field.Anonymous() || field.Embedded() {
			continue
		}
		// Check if the field is an anonymous struct
		if anonStruct, ok := field.Type().Underlying().(*types.Struct); ok {
			// Move on only if the anonymous struct is owned by the same package
			if anonStruct.NumFields() > 0 && strings.Compare(anonStruct.Field(0).Pkg().Name(), pkg.Name) == 0 {
				discovered := StructDiscovered{
					Name:    types.NewTypeName(field.Pos(), pkg.Types, fmt.Sprintf("%s.%s", parentStructName, field.Name()), field.Type()),
					Obj:     anonStruct,
					PkgName: pkg.Name,
				}
				anonymousStructs = append(anonymousStructs, discovered)
				// Recursively find any anonymous structs within this struct
				anonymousStructs = append(anonymousStructs, findAnonymousStructs(anonStruct, pkg, discovered.Name.Name())...)
			}
		}
	}

	return anonymousStructs
}
