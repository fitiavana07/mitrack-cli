package account

import (
	"github.com/fitiavana07/mitrack/cli"
	"github.com/spf13/cobra"
)

// NewAccountCommand returns a cobra command for `account` subcommands.
func NewAccountCommand(mitrackCli cli.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Manage accounts",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(
		NewRegisterCommand(mitrackCli),
		// TODO
		// NewCountCommand(mitrackCli),
		NewListCommand(mitrackCli),
		// NewShowCommand(mitrackCli),
		// NewUpdateCommand(mitrackCli),
		// NewDeleteCommand(mitrackCli),

		// TODO the format used for list
		// 		fmt.Fprintf(writer, "%s %s - %s (%s)\n", acc.ID.Short(), acc.Type.Initial(), acc.Name, acc.Alias)

	)
	return cmd
}
