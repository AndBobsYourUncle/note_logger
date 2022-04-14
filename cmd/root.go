package cmd

import (
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use: "note-logger",
}

func Execute() error {
	rootCommand.SilenceUsage = true

	if err := rootCommand.Execute(); err != nil {
		return err
	}

	return nil
}
