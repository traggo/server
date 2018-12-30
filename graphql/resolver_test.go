package graphql

import (
	"testing"

	"github.com/traggo/server/test"
)

func TestNewResolver_doesNotThrow(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	resolver := NewResolver(db, 4)
	resolver.RootMutation()
	resolver.RootQuery()
}
