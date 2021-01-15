package account

// AliasIDIndex is basically a key-value store for a fast way to retrieve the ID
// of an account given its alias.
// Docker's daemon GetByName (for containers) uses github.com/hashicorp/go-memdb
// under the hood.
type AliasIDIndex interface {
	Put(alias string)
	Get(alias string) string
	Delete(alias string)
}

// aliasIDIndex is an (example) implementation of the AliasIDIndex interface.
type aliasIDIndex struct {
	m map[string]string
}
