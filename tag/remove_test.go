package tag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/schema"
	"github.com/traggo/server/test"
)

func TestGQL_RemoveTag_succeeds_removesTag(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&schema.TagDefinition{Key: "existing tag", Color: "#fff", Type: schema.TypeSingleValue})

	resolver := ResolverForTag{DB: db}
	_, err := resolver.RemoveTag(context.Background(), "existing tag")
	require.Nil(t, err)
	assertTagCount(t, db, 0)
}

func TestGQL_RemoveTag_fails_notExistingTag(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	resolver := ResolverForTag{DB: db}
	_, err := resolver.RemoveTag(context.Background(), "not existing")
	require.EqualError(t, err, "tag with key 'not existing' does not exist")
}
