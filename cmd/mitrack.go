package cmd

import (
	"github.com/fitiavana07/mitrack/cli"
	"github.com/fitiavana07/mitrack/cli/command"
	"github.com/spf13/cobra"
)

// NewMitrackRootCmd creates the root command.
func NewMitrackRootCmd(mitrackCli cli.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mitrack",
		Short: "A CLI-based finance management tool",
	}

	command.AddCommands(cmd, mitrackCli)

	// think of viper if ever needed a conf file
	// (look at the cobra-generated code)

	return cmd
}
