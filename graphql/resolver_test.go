package graphql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func TestNewResolver_doesNotThrow(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	resolver := NewResolver(db.DB, 4, model.Version{Name: "oops", BuildDate: "date", Commit: "aeu"})
	resolver.RootMutation()
	resolver.RootQuery()
	version, e := resolver.RootQuery().Version(context.Background())
	assert.NoError(t, e)
	assert.Equal(t, &gqlmodel.Version{Name: "oops", BuildDate: "date", Commit: "aeu"}, version)
}
