package transaction

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/fitiavana07/mitrack/pkg/account"
	"github.com/fitiavana07/mitrack/pkg/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxServiceNew(t *testing.T) {
	t.Run("fresh new directory", func(t *testing.T) {
		accDir := t.TempDir()
		txDir := t.TempDir()

		accService, err := account.NewAccService(accDir)
		require.NoError(t, err)
		defer accService.Cleanup()

		s, err := NewTxService(txDir, accService)
		require.NoError(t, err)
		assert.NotNil(t, s)
		defer s.Cleanup()

		dbinfoFile := filepath.Join(txDir, ".dbinfo")
		require.FileExists(t, dbinfoFile, ".dbinfo file was not created")

		b, err := os.ReadFile(dbinfoFile)
		require.NoError(t, err)
		assert.Equal(t, "quick:v0.4", string(b))
	})
	t.Run("initialized dir", func(t *testing.T) {
	})
	t.Run("unsupported format version", func(t *testing.T) {})
}

func TestTxServiceRecordFromMaps(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		accDir := t.TempDir()
		txDir := t.TempDir()

		accService, cleanup := createTestAccService(t, accDir)
		defer cleanup()

		s, cleanup := createTestTxService(t, txDir, accService)
		defer cleanup()

		accCashInWallet := account.NewAccount("Cash in Wallet", account.TypeAsset)
		accService.Register(accCashInWallet)

		accInitialBalance := account.NewAccount("Initial Balance", account.TypeEquity)
		accService.Register(accInitialBalance)

		note := "Initial Balance of Cash in Wallet"
		amount := int64(46000)
		tx, err := s.RecordFromMaps(
			note,
			map[string]int64{accCashInWallet.Alias: amount},
			map[string]int64{accInitialBalance.Alias: amount},
		)

		assert.NoError(t, err, "transaction recording returned an error")
		require.NotNil(t, tx, "returned transaction is nil")

		var emptyHash [sha256.Size]byte
		assert.NotEqual(t, emptyHash, tx.Hash(), "hash was not set")
		assert.NotEmpty(t, tx.Timestamp(), "timestamp was not set")
		assert.Equal(t, tx.Note(), note, "note was not set")

		assert.Contains(t, tx.Entries(), NewEntry(
			OpDebit,
			accCashInWallet.ID,
			amount,
		), "transaction does not contain the correct debit entry")

		assert.Contains(t, tx.Entries(), NewEntry(
			OpCredit,
			accInitialBalance.ID,
			amount,
		), "transaction does not contain the correct credit entry")

		// TODO there should be some distinction between the transaction creation,
		// and the transaction saving into file.

		txFilePath := filepath.Join(txDir, fmt.Sprintf("%x", tx.Hash()))
		require.FileExists(t, txFilePath, "transaction file was not created")

		// I want to ensure that the data was correctly written into the file.
		// Just like when I check manually that the file contains the right content.
		f, err := os.Open(txFilePath)
		require.NoError(t, err)
		defer f.Close()

		r := bufio.NewReader(f)

		decoder := encoding.NewDecoderV3()

		timestamp := new(int64)
		err = decoder.ReadDecoded(r, timestamp)
		checkNoErrorAndEqual(t, err, tx.Timestamp(), *timestamp, "Timestamp")

		entriesCount := new(uint16)
		err = decoder.ReadDecoded(r, entriesCount)
		checkNoErrorAndEqual(t, err, uint16(len(tx.Entries())), *entriesCount, "EntriesCount")

		for i := 0; i < len(tx.Entries()); i++ {
			entry := tx.Entries()[i]
			op := new(Operation)
			err = decoder.ReadDecoded(r, op)
			checkNoErrorAndEqual(t, err, entry.Operation(), *op, "Operation")

			accountID := new(account.ID)
			err = decoder.ReadDecoded(r, accountID)
			checkNoErrorAndEqual(t, err, entry.AccountID(), *accountID, "AccountID")

			amount := new(int64)
			err = decoder.ReadDecoded(r, amount)
			checkNoErrorAndEqual(t, err, entry.Amount(), *amount, "Amount")
		}

		gotNote := new(string)
		err = decoder.ReadDecoded(r, gotNote)
		checkNoErrorAndEqual(t, err, tx.Note(), *gotNote, "Note")
	})

	t.Run("difference between credits and debits", func(t *testing.T) {})
}

func checkNoErrorAndEqual(t testing.TB, err error, want, got interface{}, name string) {
	t.Helper()
	assert.NoError(t, err, fmt.Sprintf("error reading %s", name))
	assert.Equal(t, want, got, fmt.Sprintf("wrong value for %s", name))
}

func createTestAccService(t testing.TB, dir string) (s account.AccService, cleanup func()) {
	s, err := account.NewAccService(dir)
	require.NoError(t, err)

	cleanup = func() {
		assert.NoError(t, s.Cleanup())
	}
	return
}

func createTestTxService(t testing.TB, dir string, as account.AccService) (s TxService, cleanup func()) {
	s, err := NewTxService(dir, as)
	require.NoError(t, err)

	cleanup = func() {
		assert.NoError(t, s.Cleanup())
	}
	return
}

func TestTxServiceList(t *testing.T) {
	t.Run("1 tx", func(t *testing.T) {
		accDir := t.TempDir()
		txDir := t.TempDir()

		accService, cleanup := createTestAccService(t, accDir)
		defer cleanup()

		s, cleanup := createTestTxService(t, txDir, accService)
		defer cleanup()

		accCashInWallet := account.NewAccount("Cash in Wallet", account.TypeAsset)
		accService.Register(accCashInWallet)

		accInitialBalance := account.NewAccount("Initial Balance", account.TypeEquity)
		accService.Register(accInitialBalance)

		note := "Initial Balance of Cash in Wallet"
		amount := int64(46000)
		tx, err := s.RecordFromMaps(
			note,
			map[string]int64{accCashInWallet.Alias: amount},
			map[string]int64{accInitialBalance.Alias: amount},
		)
		require.NoError(t, err)

		txs := s.List()

		assert.Equal(t, 1, len(txs), "wrong list length")
		assert.Contains(t, txs, tx, "tx not in list")
	})
}

func TestTxServiceGetByHash(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		accDir := t.TempDir()
		txDir := t.TempDir()

		accService, cleanup := createTestAccService(t, accDir)
		defer cleanup()

		s, cleanup := createTestTxService(t, txDir, accService)
		defer cleanup()

		accCashInWallet := account.NewAccount("Cash in Wallet", account.TypeAsset)
		accService.Register(accCashInWallet)

		accInitialBalance := account.NewAccount("Initial Balance", account.TypeEquity)
		accService.Register(accInitialBalance)

		note := "Initial Balance of Cash in Wallet"
		amount := int64(46000)
		tx, err := s.RecordFromMaps(
			note,
			map[string]int64{accCashInWallet.Alias: amount},
			map[string]int64{accInitialBalance.Alias: amount},
		)
		require.NoError(t, err)

		foundTx, err := s.GetByHash(fmt.Sprintf("%x", tx.Hash()))

		assert.NoError(t, err)
		assert.Equal(t, tx, foundTx)
	})
	t.Run("not existing", func(t *testing.T) {
		// should return error
	})
}
