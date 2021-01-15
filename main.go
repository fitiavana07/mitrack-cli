package main

import (
	"path/filepath"

	"github.com/fitiavana07/mitrack/cli"
	"github.com/fitiavana07/mitrack/cmd"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const (
	mitrackDirName = ".mitrack"
)

func main() {
	homeDir, err := homedir.Dir()
	cobra.CheckErr(err)

	mitrackWorkdir := filepath.Join(homeDir, mitrackDirName)

	mitrackCli, err := cli.NewMitrackCli(mitrackWorkdir)
	cobra.CheckErr(err)

	defer func() {
		cobra.CheckErr(mitrackCli.Cleanup())
	}()

	err = cmd.NewMitrackRootCmd(mitrackCli).Execute()
	cobra.CheckErr(err)
}
