package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yoanyombapro1234/FeelguudsPlatform/pkg/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   `version`,
	Short: "Prints podcli version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(version.VERSION)
		return nil
	},
}
