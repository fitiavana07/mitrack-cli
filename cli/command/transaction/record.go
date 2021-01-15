package transaction

import (
	"github.com/fitiavana07/mitrack/cli"
	"github.com/spf13/cobra"
)

// NewRecordCommand creates a new `mitrack tx record` command.
func NewRecordCommand(mitrackCli cli.Cli) *cobra.Command {
	options := recordOptions{}
	cmd := &cobra.Command{
		Use:     "r",
		Aliases: []string{"rec", "record"},
		Short:   "Record a new transaction",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.note = args[0]
			return runRecord(mitrackCli, options)
		},
		Example: `
$ mitrack tx rec \
	--debit cash-in-wallet=400,cash-at-home=500 \
	--credit checking-account=900 \
	'naka vola sabotsy namehana'
`,
	}

	flags := cmd.Flags()

	flags.StringToInt64VarP(&options.debitsMap, "debit", "d", map[string]int64{}, "debit lines")
	cmd.MarkFlagRequired("debit")

	flags.StringToInt64VarP(&options.creditsMap, "credit", "c", map[string]int64{}, "credit lines")
	cmd.MarkFlagRequired("credit")

	return cmd
}

func runRecord(mitrackCli cli.Cli, options recordOptions) error {
	_, err := mitrackCli.TxService().RecordFromMaps(options.note, options.debitsMap, options.creditsMap)
	return err
}

type recordOptions struct {
	note       string
	debitsMap  map[string]int64
	creditsMap map[string]int64
}
