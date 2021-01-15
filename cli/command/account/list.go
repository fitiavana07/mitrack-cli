package account

import (
	"fmt"

	"github.com/fitiavana07/mitrack/cli"
	"github.com/spf13/cobra"
)

// NewListCommand returns a new `mitrack account ls` command.
func NewListCommand(mitrackCli cli.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List accounts",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runList(mitrackCli)
		},
		Example: `
$ mitrack account ls
`,
	}

	return cmd
}

func runList(mitrackCli cli.Cli) {
	accounts := mitrackCli.AccService().List()
	for _, acc := range accounts {
		fmt.Printf("%s %s - %s (%s)\n", acc.ID.Short(), acc.Type.Initial(), acc.Name, acc.Alias)
	}
	return
}
