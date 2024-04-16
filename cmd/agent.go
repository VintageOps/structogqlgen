package cmd

import (
	"fmt"
	"github.com/VintageOps/structogqlgen/pkg/conversion"
	"github.com/VintageOps/structogqlgen/pkg/load"
	"github.com/urfave/cli/v2"
	"log"
	"log/slog"
	"os"
	"strings"
)

type cmdOptions struct {
	fNameContStruct string
	printOpts       conversion.PrettyPrintOptions
}

func Execute() {
	var opts cmdOptions
	app := &cli.App{
		Name:                 "structogqlgen",
		Usage:                "Converts Golang structs into GraphQL types that are readily usable with the popular GraphQL framework, gqlgen",
		EnableBashCompletion: true,
		HideHelpCommand:      true,
		Authors: []*cli.Author{
			{Name: "VintageOps"},
		},
		Description: "StructsToGqlGenTypes is a tool that helps to automatically converts Golang structs into GraphQL types that are readily usable with the popular GraphQL framework, gqlgen.\n" +
			"It aims to reduce the boilerplate code required to define GraphQL schemas manually, thus accelerating the development of GraphQL APIs in Go projects.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "src",
				Usage:       "Path to the source File containing the structs to import",
				Destination: &opts.fNameContStruct,
				Required:    true,
				Aliases:     []string{"s"},
			},
			&cli.BoolFlag{
				Name:        "use-tags",
				Usage:       "Use Tags as field name, when tag is available. If selected without specifying the tags to use, then json tag will be used",
				Destination: &opts.printOpts.UseTags,
				Aliases:     []string{"u"},
			},
			&cli.StringFlag{
				Name:        "tags",
				Usage:       "The tags to use as field name. Specifying this implies that --use-tags is selected.",
				Destination: &opts.printOpts.TagsToUse,
				Aliases:     []string{"t"},
			},
			&cli.StringFlag{
				Name:    "required-tags",
				Usage:   "If there is a tag that make a field required, specified that tag using the format key=value. e.g. validate=required",
				Aliases: []string{"r"},
				Action: func(context *cli.Context, required string) error {
					parts := strings.SplitN(required, "=", 2)
					if len(parts) != 2 {
						return fmt.Errorf("invalid format for required-tags, expected key=value")
					}
					opts.printOpts.RequireTags.Key = parts[0]
					opts.printOpts.RequireTags.Val = parts[1]
					return nil
				},
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		return printStructsAsGraphqlTypes(&opts)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func printStructsAsGraphqlTypes(opts *cmdOptions) error {
	slog.Info("Finding Structs in the provided File", "fileName", opts.fNameContStruct)

	structsFound, err := load.FindStructsInPkg(opts.fNameContStruct)
	if err != nil {
		slog.Info("Error getting structs", "fileName", opts.fNameContStruct)
		return err
	}

	gqlGenTypes, err := buildTypeDefinitions(structsFound)
	if err != nil {
		return err
	}

	prettyPrint, err := conversion.GqlPrettyPrint(gqlGenTypes, &opts.printOpts)
	if err != nil {
		return err
	}
	fmt.Println(prettyPrint)

	return nil
}

func buildTypeDefinitions(structsFound []load.StructDiscovered) ([]conversion.GqlTypeDefinition, error) {
	gqlGenTypes := make([]conversion.GqlTypeDefinition, len(structsFound))
	for idx, structType := range structsFound {
		var err error
		gqlGenTypes[idx], err = conversion.BuildGqlgenType(structType)
		if err != nil {
			return nil, err
		}
	}
	return gqlGenTypes, nil
}
