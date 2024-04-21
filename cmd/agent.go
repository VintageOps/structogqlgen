package cmd

import (
	"fmt"
	"github.com/VintageOps/structogqlgen/pkg/conversion"
	"github.com/VintageOps/structogqlgen/pkg/load"
	"github.com/urfave/cli/v2"
	"log"
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
				Usage:       "`SRC_PATH` is the required path to the source file containing the structs to import (required)",
				Destination: &opts.fNameContStruct,
				Required:    true,
				Aliases:     []string{"s"},
			},
			&cli.BoolFlag{
				Name:        "use-json-tags",
				Usage:       "Use JSON Tag as field name when available. If this is selected and a field has no Json tag, then the field name will be used.",
				Destination: &opts.printOpts.UseJsonTags,
				Aliases:     []string{"j"},
			},
			&cli.StringFlag{
				Name:        "use-custom-tags",
				Usage:       "Specify a custom tag to use as field name. Specifying this takes precedence over JSON tags. If specifed and a field does not have this tag, the field name will be used",
				Destination: &opts.printOpts.UseCustomTags,
				Aliases:     []string{"c"},
			},
			&cli.StringFlag{
				Name:    "tags-value-ignored",
				Usage:   "Specify a tag value that signal to ignore Field with tag having this value. When using json tags with use-json-tags option, if this not specified, it is automatically set to '-'",
				Aliases: []string{"i"},
				Action: func(context *cli.Context, s string) error {
					// workaround as Destination: opts.printOpts.TagFieldToIgnore does not set the provided value
					opts.printOpts.TagFieldToIgnore = &s
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "required-tags",
				Usage:   "If there is a tag that make a field required, specified that tag using the format `key=value`. e.g. validate=required",
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

// printStructsAsGraphqlTypes prints the GraphQL type definitions corresponding to the structs found in the provided source file.
// - load.GetStructsFromSourceFile function to find all structs defined in the source file.
// - buildTypeDefinitions function to build the GraphQL type definitions for each struct.
// - conversion.GqlPrettyPrint function to pretty print the GraphQL type definitions and prints the result.
func printStructsAsGraphqlTypes(opts *cmdOptions) error {
	structsFound, err := load.GetStructsFromSourceFile(opts.fNameContStruct)
	if err != nil {
		return err
	}

	gqlGenTypes, err := conversion.BuildGqlTypes(structsFound)
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
