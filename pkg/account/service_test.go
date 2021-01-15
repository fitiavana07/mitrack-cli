package account

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/fitiavana07/mitrack/pkg/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO refactor the test and service.go structure to have different backends.
// Ex: FileBackend, InMemoryBackend.

func TestAccServiceNew(t *testing.T) {
	t.Run("fresh new directory", func(t *testing.T) {
		newDir := t.TempDir()
		s, err := NewAccService(newDir)
		defer assert.NoError(t, s.Cleanup())

		assert.NoError(t, err)
		assert.NotNil(t, s)

		dbinfoFile := filepath.Join(newDir, ".dbinfo")
		require.FileExists(t, dbinfoFile, ".dbinfo file was not created")

		b, err := os.ReadFile(dbinfoFile)
		require.NoError(t, err)
		assert.Equal(t, "quick:v0.4", string(b))
	})
	t.Run("initialized dir", func(t *testing.T) {
		accountsDir := t.TempDir()
		s1, err := NewAccService(accountsDir)
		require.NoError(t, err)
		require.NoError(t, s1.Cleanup())

		s2, err := NewAccService(accountsDir)
		assert.NoError(t, err)
		assert.NotNil(t, s2)
		assert.NoError(t, s2.Cleanup())
	})
	t.Run("database with 1 account", func(t *testing.T) {
		newDir := t.TempDir()
		s1, err := NewAccService(newDir)
		require.NoError(t, err)

		acc := NewAccount("Savings Account", TypeAsset)
		s1.Register(acc)

		require.NoError(t, s1.Cleanup())

		s2, err := NewAccService(newDir)
		assert.NoError(t, err)
		assert.NotNil(t, s2)
		assert.NoError(t, s2.Cleanup())
	})
	t.Run("unsupported format version", func(t *testing.T) {
		// TODO
	})
}

func TestAccServiceRegister(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		newDir := t.TempDir()

		s, cleanup := createFakeService(t, newDir)
		defer cleanup()

		acc := NewAccount("Checking Account", TypeAsset)
		err := s.Register(acc)

		assert.NoError(t, err, "account register returned an error")

		accountFilePath := filepath.Join(newDir, fmt.Sprintf("%x", acc.ID))
		require.FileExists(t, accountFilePath, "account file was not created")

		// I want to ensure that the data was correctly written into the file.
		// Just like when I check manually that the file contains the right content.
		f, err := os.Open(accountFilePath)
		require.NoError(t, err)
		defer f.Close()

		r := bufio.NewReader(f)

		decoder := encoding.NewDecoderV3()

		accType := new(Type)
		err = decoder.ReadDecoded(r, accType)
		assert.NoError(t, err, "error reading account type")
		assert.Equal(t, acc.Type, *accType, "account type")

		parentID := new(ID)
		err = decoder.ReadDecoded(r, parentID)
		assert.NoError(t, err, "error reading account parentID")
		assert.Equal(t, acc.ParentID, *parentID, "account parentID")

		timestamp := new(int64)
		err = decoder.ReadDecoded(r, timestamp)
		assert.NoError(t, err, "error reading account timestamp")
		assert.Equal(t, acc.Timestamp, *timestamp, "account timestamp")

		alias := new(string)
		err = decoder.ReadDecoded(r, alias)
		assert.NoError(t, err, "error reading account alias")
		assert.Equal(t, acc.Alias, *alias, "account alias")

		name := new(string)
		err = decoder.ReadDecoded(r, name)
		assert.NoError(t, err, "error reading account name")
		assert.Equal(t, acc.Name, *name, "account name")

		description := new(string)
		err = decoder.ReadDecoded(r, description)
		assert.NoError(t, err, "error reading account description")
		assert.Equal(t, acc.Description, *description, "account description")
	})
	t.Run("duplicate alias", func(t *testing.T) {
		// TODO
	})
}

func TestAccServiceCount(t *testing.T) {
	t.Run("0 account", func(t *testing.T) {
		// TODO
	})
	t.Run("1 account", func(t *testing.T) {
		// TODO
	})
	t.Run("2 account", func(t *testing.T) {
		// TODO
	})
}

func TestAccServiceList(t *testing.T) {
	t.Run("0 account", func(t *testing.T) {
		dir := t.TempDir()

		s, cleanup := createFakeService(t, dir)
		defer cleanup()

		assert.Empty(t, s.List())
	})
	t.Run("1 account", func(t *testing.T) {
		dir := t.TempDir()
		acc := NewAccount("Trosa", TypeLiability)

		s1, cleanup := createFakeService(t, dir)
		require.NoError(t, s1.Register(acc))
		cleanup()

		s2, cleanup := createFakeService(t, dir)
		defer cleanup()

		accounts := s2.List()

		assert.Equal(t, 1, len(accounts), "wrong list length")
		assert.Containsf(t, accounts, acc, "account not in list")
	})
	t.Run("2 accounts", func(t *testing.T) {
		dir := t.TempDir()

		s1, cleanup := createFakeService(t, dir)

		acc1 := NewAccount("Trosa", TypeLiability)
		acc2 := NewAccount("Vola nampindramina", TypeAsset)

		require.NoError(t, s1.Register(acc1))
		require.NoError(t, s1.Register(acc2))

		cleanup()

		s2, cleanup := createFakeService(t, dir)
		defer cleanup()

		accounts := s2.List()

		assert.Equal(t, 2, len(accounts))
		assert.Containsf(t, accounts, acc1, "account not in list")
		assert.Containsf(t, accounts, acc2, "account not in list")
	})
}

func TestAccServiceGet(t *testing.T) {
	t.Run("existing full ID", func(t *testing.T) {
	})
	t.Run("existing alias", func(t *testing.T) {})
	t.Run("existing prefix", func(t *testing.T) {})
	t.Run("not existing", func(t *testing.T) {})
}

func TestAccServiceGetByID(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		dir := t.TempDir()

		acc := NewAccount("Jiro sy Rano", TypeExpense)

		s1, cleanupS1 := createFakeService(t, dir)
		s1.Register(acc)
		cleanupS1()

		s2, cleanupS2 := createFakeService(t, dir)
		defer cleanupS2()

		foundAcc, err := s2.GetByID(acc.ID.Hex())

		assert.NoError(t, err)
		assert.Equal(t, acc, foundAcc)

	})
	t.Run("not existing", func(t *testing.T) {})
}

func TestAccServiceGetByActualID(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		dir := t.TempDir()

		acc := NewAccount("Jiro sy Rano", TypeExpense)
		s1, cleanupS1 := createFakeService(t, dir)
		s1.Register(acc)
		cleanupS1()

		s2, cleanupS2 := createFakeService(t, dir)
		defer cleanupS2()

		foundAcc, err := s2.GetByActualID(acc.ID)

		assert.NoError(t, err)
		assert.Equal(t, acc, foundAcc)
	})
	t.Run("not existing", func(t *testing.T) {})
}

func TestAccServiceGetByAlias(t *testing.T) {
	t.Run("existing", func(t *testing.T) {})
	t.Run("not existing", func(t *testing.T) {})
}

func TestAccServiceGetByPrefix(t *testing.T) {
	t.Run("existing", func(t *testing.T) {})
	t.Run("not existing", func(t *testing.T) {})
}

func TestAccServiceUpdate(t *testing.T) {
	t.Run("update name", func(t *testing.T) {})
	t.Run("update name", func(t *testing.T) {})
	t.Run("update name empty", func(t *testing.T) {})
	t.Run("update name duplicated", func(t *testing.T) {})
	t.Run("update alias", func(t *testing.T) {})
	t.Run("update alias empty", func(t *testing.T) {})
	t.Run("update alias duplicated", func(t *testing.T) {})
	t.Run("update description", func(t *testing.T) {})
	t.Run("update type impossible", func(t *testing.T) {})
	t.Run("update parent", func(t *testing.T) {})
	t.Run("update timestamp impossible", func(t *testing.T) {})
	t.Run("not existing", func(t *testing.T) {})
}

func TestAccServiceDelete(t *testing.T) {
	t.Run("delete success on unused account", func(t *testing.T) {})
	t.Run("delete errors on used account", func(t *testing.T) {
		// TODO for example set "used" to true once used in transaction.
	})
}

func createFakeService(t testing.TB, dir string) (service AccService, cleanup func()) {
	service, err := NewAccService(dir)
	require.NoError(t, err)

	cleanup = func() {
		assert.NoError(t, service.Cleanup())
	}
	return
}
