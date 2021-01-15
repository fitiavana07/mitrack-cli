package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fitiavana07/mitrack/pkg/account"
	"github.com/fitiavana07/mitrack/pkg/transaction"
)

// Cli represents the mitrack command line interface.
type Cli interface {
	AccService() account.AccService
	TxService() transaction.TxService
	Cleanup() error
}

// MitrackCli represents an instance of the mitrack command line interface.
// Instances are created using NewMitrackCli.
type MitrackCli struct {
	accService account.AccService
	txService  transaction.TxService
}

const (
	dataDirName         = "data"
	configDirName       = "config"
	accountsDirName     = "accounts"
	transactionsDirName = "transactions"
)

// NewMitrackCli returns a new MitrackCli.
func NewMitrackCli(workdir string) (Cli, error) {
	accountsDir := filepath.Join(workdir, accountsDirName)
	transactionsDir := filepath.Join(workdir, transactionsDirName)

	for _, dir := range []string{workdir, accountsDir, transactionsDir} {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
			return nil, err
		}
	}

	accService, err := account.NewAccService(accountsDir)
	if err != nil {
		return nil, err
	}

	txService, err := transaction.NewTxService(transactionsDir, accService)
	if err != nil {
		return nil, err
	}

	return &MitrackCli{accService, txService}, nil
}

// AccService returns the account service.
func (c *MitrackCli) AccService() account.AccService {
	return c.accService
}

// TxService returns the transaction service.
func (c *MitrackCli) TxService() transaction.TxService {
	return c.txService
}

// Cleanup clean up used resources (files, etc.).
func (c *MitrackCli) Cleanup() error {
	accServiceCleanupErr := c.AccService().Cleanup()
	txServiceCleanupErr := c.TxService().Cleanup()
	if accServiceCleanupErr != nil || txServiceCleanupErr != nil {
		// return the full even if one was not nil
		return &CleanupError{accServiceCleanupErr, txServiceCleanupErr}
	}
	return nil
}

// CleanupError is the error returned by Cleanup()
type CleanupError struct {
	accServiceCleanupErr error
	txServiceCleanupErr  error
}

func (e *CleanupError) Error() string {
	return fmt.Sprintf("accServiceCleanupErr=%+v; txServiceCleanupErr=%+v", e.accServiceCleanupErr, e.txServiceCleanupErr)
}
