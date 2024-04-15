package cmd

import (
	"fmt"
	"github.com/VintageOps/structogqlgen/pkg/conversion"
	"github.com/VintageOps/structogqlgen/pkg/load"
	"github.com/urfave/cli/v2"
	"log"
	"log/slog"
	"os"
)

func Execute() {
	var fNameContStruct string
	app := &cli.App{
		Name:                 "structogqlgen",
		Usage:                "Converts Golang structs into GraphQL types that are readily usable with the popular GraphQL framework, gqlgen",
		EnableBashCompletion: true,
		HideHelpCommand:      true,
		ArgsUsage:            "<package_to_import>",
		Description: "StructsToGqlGenTypes is a tool that helps to automatically converts Golang structs into GraphQL types that are readily usable with the popular GraphQL framework, gqlgen.\n" +
			"It aims to reduce the boilerplate code required to define GraphQL schemas manually, thus accelerating the development of GraphQL APIs in Go projects.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "src",
				Usage:       "Path for Source File containing the structs to import",
				Destination: &fNameContStruct,
				Required:    true,
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		slog.Info("Finding Structs in the provided File", "fileName", fNameContStruct)
		structsFound, err := load.FindStructsInPkg(fNameContStruct)
		if err != nil {
			slog.Info("Error getting structs", "fileName", fNameContStruct)
			return err
		}
		// For each struct we will build its graphQl Type
		gqlGenTypes := make([]conversion.GqlTypeDefinition, len(structsFound))
		for idx, structType := range structsFound {
			// Build GraphQL type for structType
			gqlGenTypes[idx], err = conversion.BuildGqlgenType(structType)
			if err != nil {
				return err
			}
		}
		prettyPrint, err := conversion.GqlTypePrettyPrint(gqlGenTypes, false, "")
		if err != nil {
			return err
		}
		fmt.Println(prettyPrint)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
