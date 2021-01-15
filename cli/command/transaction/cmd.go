package transaction

import (
	"github.com/fitiavana07/mitrack/cli"
	"github.com/spf13/cobra"
)

// NewTransactionCommand returns a cobra command for `transaction` subcommands.
func NewTransactionCommand(mitrackCli cli.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tx",
		Aliases: []string{"transaction"},
		Short:   "Record and list Transactions",
		Args:    cobra.NoArgs,
	}
	cmd.AddCommand(
		NewRecordCommand(mitrackCli),
		// TODO
		// NewCountCommand(mitrackCli),
		NewListCommand(mitrackCli),
		// NewShowCommand(mitrackCli),
	)
	return cmd
}
