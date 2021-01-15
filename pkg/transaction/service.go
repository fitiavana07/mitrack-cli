package transaction

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fitiavana07/mitrack/pkg/account"
	"github.com/fitiavana07/mitrack/pkg/encoding"
)

// TxService provides methods for managing transactions.
type TxService interface {
	// ==== CREATE ====
	// RecordFromMaps records a transaction using the info given in args.
	// debitsMap and creditsMap are alias->amount maps.
	// This method include transaction verification (ex: sum of debits must
	// equal sum of credits).
	// This method returns a pointer to the created transaction.
	RecordFromMaps(note string, debitsMap, creditsMap map[string]int64) (Transaction, error)

	// ==== READ ====
	// Count returns the total number of transactions in the transactions database.
	Count() uint64
	// List returns all transactions in the transactions database.
	List() []Transaction
	// Get returns the transaction given its prefix (short hash) or full hash.
	// The search is in this order: full hash hex, prefix.
	Get(prefix string) (Transaction, error)
	// GetByHash returns a transaction given its hash in  hex.
	GetByHash(hash string) (Transaction, error)
	// GetByPrefix returns a transaction given a prefix.
	GetByPrefix(prefix string) (Transaction, error)

	// ==== UPDATE ====
	// NO UPDATE, IMMUTABLE

	// ==== DELETE ====
	// NO DELETE, IMMUTABLE

	// ==== CLEAN UP ====
	// Cleanup cleans up used resources.
	// It must be called when the TxService is no more used.
	Cleanup() error
}

// NewTxService returns a new TxService.
// It stores transactions file in the given dir.
// It uses the given accService to search for accounts.
func NewTxService(dir string, accService account.AccService) (TxService, error) {
	dbinfoFilePath := filepath.Join(dir, dbInfoFileName)
	_, err := os.Stat(dbinfoFilePath)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(dbinfoFilePath)
		if err != nil {
			return nil, fmt.Errorf("transaction.service: could not create .dbinfo: %s", err)
		}
		defer f.Close()

		if _, err = f.WriteString("quick:v0.4"); err != nil {
			return nil, fmt.Errorf("transaction.service: could not write .dbinfo content: %s", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("transaction.service: could not read .dbinfo: %s", err)
	}

	return &txService{dir: dir, accService: accService}, nil
}

type txService struct {
	dir        string
	accService account.AccService
}

const dbInfoFileName = ".dbinfo"

func (s *txService) RecordFromMaps(note string, debitsMap, creditsMap map[string]int64) (Transaction, error) {
	entriesLen := len(debitsMap) + len(creditsMap)

	type aliasEntry struct {
		op     Operation
		alias  string
		amount int64
	}
	aliasEntries := make([]aliasEntry, 0, entriesLen)

	for alias, amount := range debitsMap {
		aliasEntries = append(aliasEntries, aliasEntry{OpDebit, alias, amount})
	}
	for alias, amount := range creditsMap {
		aliasEntries = append(aliasEntries, aliasEntry{OpCredit, alias, amount})
	}

	// TODO validation: sum debits = sum credits

	entries := make([]Entry, 0, entriesLen)
	for _, ae := range aliasEntries {
		acc, err := s.accService.GetByAlias(ae.alias)
		if err != nil {
			// TODO refactor the account not found error
			return nil, fmt.Errorf("account of alias %q not found: %s", ae.alias, err)
		}

		entries = append(entries, NewEntry(ae.op, acc.ID, ae.amount))
	}

	tx := &transaction{
		timestamp: time.Now().UTC().Unix(),
		note:      note,
		entries:   entries,
	}

	encoder := encoding.NewEncoderV3()

	b := new(bytes.Buffer)

	var err error
	if err = encoder.WriteEncoded(b, tx.timestamp); err != nil {
		return nil, err
	}

	if err = encoder.WriteEncoded(b, uint16(entriesLen)); err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if err = encoder.WriteEncoded(b, entry.Operation()); err != nil {
			return nil, err
		}
		if err = encoder.WriteEncoded(b, entry.AccountID()); err != nil {
			return nil, err
		}
		if err = encoder.WriteEncoded(b, entry.Amount()); err != nil {
			return nil, err
		}
	}

	if err = encoder.WriteEncoded(b, tx.note); err != nil {
		return nil, err
	}

	tx.hash = sha256.Sum256(b.Bytes())

	txFilePath := filepath.Join(s.dir, fmt.Sprintf("%x", tx.hash))
	f, err := os.Create(txFilePath)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction file: %v", err)
	}
	defer f.Close()

	if _, err = b.WriteTo(f); err != nil {
		return nil, fmt.Errorf("error writing transaction data: %v", err)
	}

	return tx, nil
}

func (s *txService) Count() uint64 {
	// TODO
	return 0
}
func (s *txService) List() []Transaction {
	txs := []Transaction{}

	dirEntries, err := os.ReadDir(s.dir)
	if err != nil {
		return txs
	}

	for _, entry := range dirEntries {
		fileName := entry.Name()
		// TODO refactor account.ID to reuse the DecodeID function.

		if fileName == ".dbinfo" {
			// special file
			continue
		}

		tx, err := s.GetByHash(fileName)
		if err != nil {
			continue
		}

		txs = append(txs, tx)
	}

	return txs
}
func (s *txService) Get(prefix string) (Transaction, error) {
	// TODO
	return nil, nil
}

// GetByHash returns a transaction given its hash in  hex.
func (s *txService) GetByHash(hash string) (Transaction, error) {
	b, err := hex.DecodeString(hash)
	if err != nil {
		return nil, fmt.Errorf("transaction.service: invalid hash")
	}

	actualHash := [sha256.Size]byte{}
	copy(actualHash[:], b[:])

	// open the file
	txFilePath := filepath.Join(s.dir, hash)
	f, err := os.Open(txFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("transaction.service: transaction not found")
	} else if err != nil {
		return nil, fmt.Errorf("transaction.service: could not open transaction file: %v", err)
	}
	defer f.Close()

	tx := transaction{hash: actualHash}

	decoder := encoding.NewDecoderV3()

	r := bufio.NewReader(f)

	if err = decoder.ReadDecoded(r, &tx.timestamp); err != nil {
		return nil, fmt.Errorf("transaction.service: invalid transaction file format: %v", err)
	}
	entriesLen := new(uint16)
	if err = decoder.ReadDecoded(r, entriesLen); err != nil {
		return nil, fmt.Errorf("transaction.service: invalid transaction file format: %v", err)
	}
	entries := make([]Entry, 0, *entriesLen)
	for i := 0; i < int(*entriesLen); i++ {
		entry := entry{}
		if err = decoder.ReadDecoded(r, &entry.operation); err != nil {
			return nil, fmt.Errorf("transaction.service: invalid transaction file format: %v", err)
		}
		if err = decoder.ReadDecoded(r, &entry.accountID); err != nil {
			return nil, fmt.Errorf("transaction.service: invalid transaction file format: %v", err)
		}
		if err = decoder.ReadDecoded(r, &entry.amount); err != nil {
			return nil, fmt.Errorf("transaction.service: invalid transaction file format: %v", err)
		}
		entries = append(entries, &entry)
	}
	if err = decoder.ReadDecoded(r, &tx.note); err != nil {
		return nil, fmt.Errorf("transaction.service: invalid transaction file format: %v", err)
	}

	tx.entries = entries

	return &tx, nil
}

// GetByPrefix returns a transaction given a prefix.
func (s *txService) GetByPrefix(prefix string) (Transaction, error) {
	// TODO
	return nil, nil
}

func (s *txService) Cleanup() error {
	// TODO
	return nil
}
