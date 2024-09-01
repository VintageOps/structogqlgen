package cmd

import (
	"errors"
	"fmt"
	"github.com/VintageOps/structogqlgen/pkg/conversion"
	"github.com/VintageOps/structogqlgen/pkg/load"
	"github.com/spf13/cobra"
	"log"
)

// cmdOptions contains the options that are set
type cmdOptions struct {
	fNamePathPkg string                        //directory with package or single source file path
	printOpts    conversion.PrettyPrintOptions //Other print option
}

var requiredTagsFlag map[string]string
var opts cmdOptions

var rootCmd = &cobra.Command{
	Use:   "structogqlgen [path]",
	Short: "Converts Golang structs defined on the specified path, which can be either a go package folder or a single go source file, into GraphQL types for gqlgen",
	Long: `StructsToGqlGenTypes is a tool that helps automatically convert Golang structs into GraphQL types
that are readily usable with the popular GraphQL framework, gqlgen. It aims to reduce the boilerplate code
required to define GraphQL schemas manually, thus accelerating the development of GraphQL APIs in Go projects.`,
	Version: "0.2",
	Example: `structogqlgen pkg/examples --use-json-tags`,
}

func init() {
	// Define flags
	rootCmd.PersistentFlags().BoolVarP(&opts.printOpts.UseJsonTags, "use-json-tags", "j", false, "Use JSON Tag as field name when available. If not present, the field name will be used.")
	rootCmd.PersistentFlags().StringVarP(&opts.printOpts.UseCustomTags, "use-custom-tags", "c", "", "Specify a custom tag to use as field name. This takes precedence over JSON tags.")

	var tagFieldToIgnore string
	rootCmd.PersistentFlags().StringVarP(&tagFieldToIgnore, "tags-value-ignored", "i", "-", "Specify a tag value that signals to ignore a field with this tag value."+
		"Automatically set to '-' for JSON tags if not specified.")
	opts.printOpts.TagFieldToIgnore = &tagFieldToIgnore

	// required-tags checks to be done on the PersistentPreRunE
	rootCmd.PersistentFlags().StringToStringVarP(&requiredTagsFlag, "required-tags", "r", nil, "If there is a tag that make a field required, specified that tag using the format `key=value`. e.g. validate=required")

}

func Execute() {

	// Check Argument that should mutually exclude
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Check mutual exclusion
		if len(args) == 0 {
			return errors.New("the path to the package folder or to a single go source file must be provided")
		}

		if tagIgnored, _ := cmd.Flags().GetString("tags-value-ignored"); tagIgnored != "" {
			opts.printOpts.TagFieldToIgnore = &tagIgnored
		}

		// Check on required tag key value
		if requiredTagsFlag != nil && len(requiredTagsFlag) > 0 {
			if len(requiredTagsFlag) > 1 {
				return fmt.Errorf("invalid number of arguments for required-tags, expected only one key=value pair")
			}
			for k, v := range requiredTagsFlag {
				opts.printOpts.RequireTags.Key = k
				opts.printOpts.RequireTags.Val = v
			}
		}

		return nil
	}

	// Define the Run
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			opts.fNamePathPkg = args[0]
		}
		return printStructsAsGraphqlTypes(&opts)
	}

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func printStructsAsGraphqlTypes(opts *cmdOptions) error {
	var err error
	var structsFound []load.StructDiscovered

	// Build StructsFound
	structsFound, err = load.GetStructsFromPath(opts.fNamePathPkg)
	if err != nil {
		return err
	}

	if structsFound == nil || len(structsFound) == 0 {
		return fmt.Errorf("no structs found in given path %s", opts.fNamePathPkg)
	}

	// Conversion
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
