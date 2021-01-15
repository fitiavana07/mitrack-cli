package account

import (
	"github.com/fitiavana07/mitrack/cli"
	"github.com/fitiavana07/mitrack/pkg/account"
	"github.com/spf13/cobra"
)

// NewRegisterCommand creates a new `mitrack account register` command.
func NewRegisterCommand(mitrackCli cli.Cli) *cobra.Command {
	options := registerOptions{}

	cmd := &cobra.Command{
		Use:   "register --type=TYPE NAME",
		Short: "Register a new account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.accountName = args[0]
			return runRegister(mitrackCli, options)
		},
		Example: `
$ mitrack account register --type=asset 'Checking Account'
`,
	}

	flags := cmd.Flags()

	flags.Var(newAccountTypeValue(&options.accountType), "type", "account type (asset|liability|equity|expense|revenue)")
	cmd.MarkFlagRequired("type")

	return cmd
}

func runRegister(mitrackCli cli.Cli, options registerOptions) error {
	a := account.NewAccount(options.accountName, options.accountType)
	return mitrackCli.AccService().Register(a)
}

type registerOptions struct {
	accountName string
	accountType account.Type
}
