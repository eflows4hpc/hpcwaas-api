package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func init() {

	rootCmd.AddCommand(&cobra.Command{
		Use:   "gen-doc <output dir> [man|rst|yaml|md]",
		Short: "Generates command line interface documentation",
		Long: `Generates command line interface documentation

This command accepts two argument the first one is the directory where generated documentation will be stored,
the second one is the documentation format.

Supported documentation format are:

* man: for linux man pages
* rst: for reStructuredText (reST)
* yaml: for a YAML output
* md: for Markdown
`,
		Hidden:       true,
		Args:         cobra.ExactArgs(2),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			outDir := args[0]
			docType := args[1]

			switch strings.ToLower(docType) {
			case "man":
				return doc.GenManTree(rootCmd, &doc.GenManHeader{
					Title: "smcx-prov",
				}, outDir)
			case "rst":
				return doc.GenReSTTree(rootCmd, outDir)
			case "yaml":
				return doc.GenYamlTree(rootCmd, outDir)
			case "md":
				return doc.GenMarkdownTree(rootCmd, outDir)
			default:
				return fmt.Errorf("unsupported documentation output")
			}
		},
	})

}
