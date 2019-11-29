package tag

import (
	"github.com/stretchr/testify/assert"
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

func TestGQL_RemoveTag_succeeds_removesTimespans(t *testing.T) {
	test.LogDebug()
	db := test.InMemoryDB(t)
	defer db.Close()
	user := db.User(3)
	db.Create(&model.TagDefinition{Key: "tag", Color: "#fff", Type: model.TypeSingleValue, UserID: 3})
	db.Create(&model.TagDefinition{Key: "tag", Color: "#fff", Type: model.TypeSingleValue, UserID: 4})
	ts := user.TimeSpan("2009-06-30T18:30:00Z", "2009-06-30T18:40:00Z")
	ts.Tag("tag", "def")
	other := db.User(4).TimeSpan("2009-06-30T18:30:00Z", "2009-06-30T18:40:00Z")
	other.Tag("tag", "def")

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.RemoveTag(fake.User(3), "tag")
	require.Nil(t, err)
	assertTagCount(t, db, 1)
	assertTagExist(t, db, model.TagDefinition{Key: "tag", Color: "#fff", Type: model.TypeSingleValue, UserID: 4})

	assert.True(t, db.Where(&model.TimeSpanTag{Key: "tag", TimeSpanID: ts.TimeSpan.ID}).Find(new(model.TimeSpanTag)).RecordNotFound(), "should be removed")
	assert.False(t, db.Where(&model.TimeSpanTag{Key: "tag", TimeSpanID: other.TimeSpan.ID}).Find(new(model.TimeSpanTag)).RecordNotFound(), "should not be removed")
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
