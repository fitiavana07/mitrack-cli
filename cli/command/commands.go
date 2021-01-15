package command

import (
	"github.com/fitiavana07/mitrack/cli"
	"github.com/fitiavana07/mitrack/cli/command/account"
	"github.com/fitiavana07/mitrack/cli/command/transaction"
	"github.com/spf13/cobra"
)

// AddCommands adds all commands from cli/command to cmd.
func AddCommands(cmd *cobra.Command, mitrackCli cli.Cli) {
	cmd.AddCommand(
		account.NewAccountCommand(mitrackCli),
		transaction.NewTransactionCommand(mitrackCli),
	)
}
