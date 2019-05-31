// Package cli contains the supported lazo commands.
package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "lazo",
	Short: "Lazo is a tool for managing Lazo source code",
	Long:  `Lazo is a tool for managing Lazo source code on the Bazo Blockchain`,
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

// Execute executes the lazo command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
