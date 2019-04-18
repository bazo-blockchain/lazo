package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Lazo",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Lazo compiler v1.0")
	},
}
