package tag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/schema"
	"github.com/traggo/server/test"
)

func TestGQL_CreateTag_succeeds_addsTag(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()

	resolver := ResolverForTag{DB: db}
	tag, err := resolver.CreateTag(context.Background(), "new tag", "#fff", gqlmodel.TagDefinitionTypeSinglevalue)

	require.Nil(t, err)
	expected := &gqlmodel.TagDefinition{
		Key:   "new tag",
		Color: "#fff",
		Type:  gqlmodel.TagDefinitionTypeSinglevalue,
	}
	require.Equal(t, expected, tag)
	assertTagExist(t, db, schema.TagDefinition{
		Key:   "new tag",
		Color: "#fff",
		Type:  schema.TypeSingleValue,
	})
	assertTagCount(t, db, 1)
}

func TestGQL_CreateTag_fails_tagAlreadyExists(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&schema.TagDefinition{Key: "existing tag", Color: "#fff", Type: schema.TypeSingleValue})

	resolver := ResolverForTag{DB: db}
	_, err := resolver.CreateTag(context.Background(), "existing tag", "#fff", gqlmodel.TagDefinitionTypeSinglevalue)

	require.EqualError(t, err, "tag with key 'existing tag' does already exist")
	assertTagCount(t, db, 1)
}
