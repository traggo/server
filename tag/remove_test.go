package tag

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestGQL_RemoveTag_succeeds_removesTag(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(3)

	db.Create(&model.TagDefinition{Key: "existing tag", Color: "#fff", Type: model.TypeSingleValue, UserID: 3})

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.RemoveTag(fake.User(3), "existing tag")
	require.Nil(t, err)
	assertTagCount(t, db, 0)
}

func TestGQL_RemoveTag_fails_notExistingTag(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(3)

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.RemoveTag(fake.User(3), "not existing")
	require.EqualError(t, err, "tag with key 'not existing' does not exist")
}

func TestGQL_RemoveTag_fails_notPermission(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.User(3)
	db.User(5)
	db.Create(&model.TagDefinition{Key: "existing tag", Color: "#fff", Type: model.TypeSingleValue, UserID: 3})

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.RemoveTag(fake.User(5), "existing tag")
	require.EqualError(t, err, "tag with key 'existing tag' does not exist")
}
