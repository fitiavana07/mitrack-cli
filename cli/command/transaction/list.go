package transaction

import (
	"fmt"
	"os"
	"time"

	"github.com/fitiavana07/mitrack/cli"
	"github.com/fitiavana07/mitrack/pkg/transaction"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// NewListCommand returns a new `mitrack tx ls` command.
func NewListCommand(mitrackCli cli.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List transactions",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runList(mitrackCli)
		},
		Example: `
$ mitrack tx ls
`,
	}

	return cmd
}

func runList(mitrackCli cli.Cli) {
	txs := mitrackCli.TxService().List()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"DATE", "ACCOUNTS", "Debit", "Credit"})

	data := [][]string{}

	for _, tx := range txs {

		date := time.Unix(tx.Timestamp(), 0)
		dateStr := date.Format(time.RFC3339)

		for i, entry := range tx.Entries() {
			acc, err := mitrackCli.AccService().GetByActualID(entry.AccountID())
			if err != nil {
				return
			}
			accountName := acc.Name
			if entry.Operation() == transaction.OpDebit {
				dateCellContent := ""
				if i == 0 {
					dateCellContent = dateStr
				}
				data = append(data, []string{
					dateCellContent,
					accountName,
					fmt.Sprintf("%d", entry.Amount()),
					"",
				})
			} else if entry.Operation() == transaction.OpCredit {
				data = append(data, []string{
					"",
					accountName,
					"",
					fmt.Sprintf("%d", entry.Amount()),
				})
			}
		}

		note := tx.Note()
		data = append(data, []string{"Note", note, "", ""})
		hash := fmt.Sprintf("%x", tx.Hash())
		data = append(data, []string{"Hash", hash, "", ""})

		data = append(data, []string{"", "", "", ""})
	}

	for _, v := range data {
		table.Append(v)
	}

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	table.SetColumnColor(
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)
	table.SetRowLine(true)
	table.Render()
}
