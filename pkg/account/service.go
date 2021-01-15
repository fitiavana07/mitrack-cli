package account

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fitiavana07/mitrack/pkg/encoding"
)

// AccService provides methods for managing accounts.
type AccService interface {
	// ==== CREATE ====

	// Register registers the Account in the accounts database.
	// The provided account must be a valid initialized account
	// (ie: with valid ID timestamp and alias).
	Register(*Account) error

	// ==== READ ====

	// Count returns the total number of accounts in the DB.
	Count() uint64
	// List returns all accounts in the DB.
	List() []*Account
	// Get returns the account given the alias, short ID (prefix), or full ID (hex).
	// The Order of search trials is: full ID, alias, prefix.
	// For more inspiration, look at daemon/container at moby repo.
	Get(prefixOrAlias string) (*Account, error)
	// GetByID returns the Account identifid by full ID in hex.
	GetByID(id string) (*Account, error)
	// GetByActualID returns the Account identifid by an actual id of type ID.
	GetByActualID(id ID) (*Account, error)
	// GetByAlias would basically find the ID using a alias->ID map,
	// then GetByID(id) to get the Account.
	GetByAlias(alias string) (*Account, error)
	// GetByPrefix returns the Account corresponding to a given ID prefix.
	GetByPrefix(prefix string) (*Account, error)

	// ==== UPDATE ====

	// Update updates the given account in the database.
	Update(*Account) error

	// ==== DELETE ====

	// Delete deletes the found account from the database.
	// (uses Get under the hood to get the *Account).
	Delete(prefixOrAlias string) error

	// ==== CLEAN UP ====
	// Cleanup cleans up used resources.
	// It must be called when the AccService is no more used.
	Cleanup() error
}

// NewAccService returns a new AccService.
func NewAccService(accountsDir string) (AccService, error) {
	// file, err := os.OpenFile(baseFilePath, os.O_RDWR|os.O_CREATE, 0644)

	dbinfoFilePath := filepath.Join(accountsDir, dbInfoFileName)
	_, err := os.Stat(dbinfoFilePath)
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(dbinfoFilePath)
		if err != nil {
			return nil, fmt.Errorf("account.service: could not create .dbinfo: %s", err)
		}
		defer file.Close()

		if _, err = file.WriteString("quick:v0.4"); err != nil {
			return nil, fmt.Errorf("account.service: could not write .dbinfo content: %s", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("account.service: could not read .dbinfo: %s", err)
	}

	return &accService{workDir: accountsDir}, nil
}

type accService struct {
	workDir string
}

const dbInfoFileName = ".dbinfo"

func (s *accService) Register(acc *Account) error {
	path := filepath.Join(s.workDir, fmt.Sprintf("%x", acc.ID))
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("account.service: could not create account file: %s", err)
	}
	defer f.Close()

	encoder := encoding.NewEncoderV3()

	toEncode := []interface{}{
		acc.Type,
		acc.ParentID,
		acc.Timestamp,
		acc.Alias,
		acc.Name,
		acc.Description,
	}

	b := bufio.NewWriter(f)
	for _, v := range toEncode {
		if err = encoder.WriteEncoded(b, v); err != nil {
			return fmt.Errorf("account.service: %s", err)
		}
	}
	if err = b.Flush(); err != nil {
		return fmt.Errorf("account.service: error while writing buffered data into file: %s", err)
	}

	return nil
}

func (s *accService) Count() uint64 {
	// TODO
	return 0
}

func (s *accService) List() []*Account {
	accounts := []*Account{}

	dirEntries, err := os.ReadDir(s.workDir)
	if err != nil {
		return accounts
	}

	for _, entry := range dirEntries {
		fileName := entry.Name()
		actualID, err := DecodeID(fileName)
		if err != nil {
			// the fileName is not of an account
			continue
		}
		acc, err := s.GetByActualID(actualID)
		if err != nil {
			// TODO should we really not return any error?
			continue
		}

		accounts = append(accounts, acc)
	}

	return accounts
}

func (s *accService) Get(prefixOrAlias string) (*Account, error) {
	// TODO
	return nil, nil
}

func (s *accService) GetByID(id string) (*Account, error) {
	actualID, err := DecodeID(id)
	if err != nil {
		return nil, err
	}

	return s.GetByActualID(actualID)
}

func (s *accService) GetByActualID(id ID) (*Account, error) {
	hexValue := id.Hex()
	accountFilePath := filepath.Join(s.workDir, hexValue)
	f, err := os.Open(accountFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("account.service: account not found")
	} else if err != nil {
		return nil, fmt.Errorf("account.service: could not open account file: %s", err)
	}
	defer f.Close()

	decoder := encoding.NewDecoderV3()

	a := Account{ID: id}
	toDecode := []interface{}{
		&(a.Type),
		&(a.ParentID),
		&(a.Timestamp),
		&(a.Alias),
		&(a.Name),
		&(a.Description),
	}

	r := bufio.NewReader(f)
	for _, v := range toDecode {
		if err = decoder.ReadDecoded(r, v); err != nil {
			return nil, fmt.Errorf("account.service: invalid account file format: %s", err)
		}
	}

	return &a, nil
}

func (s *accService) GetByAlias(alias string) (*Account, error) {
	// TODO add alias->account index
	accounts := s.List()

	for _, a := range accounts {
		if a.Alias == alias {
			return a, nil
		}
	}

	return nil, fmt.Errorf("account not found")
}

func (s *accService) GetByPrefix(prefix string) (*Account, error) {
	// TODO use this package: github.com/docker/docker/pkg/truncindex.
	return nil, nil
}

func (s *accService) Update(*Account) error {
	// TODO
	return nil
}
func (s *accService) Delete(prefixOrAlias string) error {
	// TODO
	return nil
}

func (s *accService) Cleanup() error {
	return nil
}

// ErrUnimplemented is returned from unimplemented functions.
var ErrUnimplemented = errors.New("unimplemented")
