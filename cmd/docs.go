package cmd

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

//

// versionCmd represents the version command
var docgenCmd = &cobra.Command{
	Use:    "docgen",
	Short:  "Generate the docs",
	Long:   `Just Displays the Version of Rome`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating Docs")
		err := doc.GenMarkdownTree(RootCmd, "./docs")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(docgenCmd)
}
