package tag

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestUpdate_withKey(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	left := db.User(5)
	right := db.User(2)
	left.NewTagDefinition("coolio")
	right.NewTagDefinition("coolio")
	leftTs := left.TimeSpan("2009-06-30T18:30:00Z", "2009-06-30T18:40:00Z")
	leftTs.Tag("coolio", "mama")
	rightTs := right.TimeSpan("2009-06-30T18:30:00Z", "2009-06-30T18:40:00Z")
	rightTs.Tag("coolio", "mama")

	resolver := ResolverForTag{DB: db.DB}
	newTagName := "mega"
	tag, err := resolver.UpdateTag(fake.User(left.User.ID), "coolio", &newTagName, "#abc")
	require.NoError(t, err)
	require.Equal(t, &gqlmodel.TagDefinition{
		Color: "#abc",
		Key:   "mega",
	}, tag)
	left.AssertHasTagDefinition("coolio", false).AssertHasTagDefinition("mega", true)
	right.AssertHasTagDefinition("coolio", true).AssertHasTagDefinition("mega", false)
	leftTs.AssertHasTag("mega", "mama", true).AssertHasTag("coolio", "mama", false)
	rightTs.AssertHasTag("coolio", "mama", true).AssertHasTag("mega", "mama", false)
}

func TestUpdate_lowercases(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	user := db.User(5)
	user.NewTagDefinition("coolio")
	ts := user.TimeSpan("2009-06-30T18:30:00Z", "2009-06-30T18:40:00Z")
	ts.Tag("coolio", "mama")

	resolver := ResolverForTag{DB: db.DB}
	newTagName := "Mega"
	tag, err := resolver.UpdateTag(fake.User(user.User.ID), "coolio", &newTagName, "#abc")
	require.NoError(t, err)
	require.Equal(t, &gqlmodel.TagDefinition{
		Color: "#abc",
		Key:   "mega",
	}, tag)
	user.AssertHasTagDefinition("coolio", false).AssertHasTagDefinition("mega", true)
	ts.AssertHasTag("mega", "mama", true).AssertHasTag("coolio", "mama", false)
}

func TestUpdate_disallow_space(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	user := db.User(5)
	user.NewTagDefinition("coolio")

	resolver := ResolverForTag{DB: db.DB}
	newTagName := "the coolio"
	_, err := resolver.UpdateTag(fake.User(user.User.ID), "coolio", &newTagName, "#abc")
	require.EqualError(t, err, "tag must not contain spaces")
}

func TestUpdate_withoutKey(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	left := db.User(5)
	right := db.User(2)
	left.NewTagDefinition("coolio")
	right.NewTagDefinition("coolio")
	leftTs := left.TimeSpan("2009-06-30T18:30:00Z", "2009-06-30T18:40:00Z")
	leftTs.Tag("coolio", "mama")
	rightTs := right.TimeSpan("2009-06-30T18:30:00Z", "2009-06-30T18:40:00Z")
	rightTs.Tag("coolio", "mama")

	resolver := ResolverForTag{DB: db.DB}
	tag, err := resolver.UpdateTag(fake.User(left.User.ID), "coolio", nil, "#abc")
	require.NoError(t, err)
	assert.Equal(t, &gqlmodel.TagDefinition{
		Color: "#abc",
		Key:   "coolio",
	}, tag)
}

func TestUpdate_dashboardEntryKey(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	left := db.User(5)
	left.NewTagDefinition("coolio")
	dashboard := left.Dashboard("yeah")
	dashboard.Entry("entry")
	entry := dashboard.Dashboard.Entries[0]
	entry.Keys = "abc,coolio,chicken"
	db.Save(&entry)

	newTag := "yes"
	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.UpdateTag(fake.User(left.User.ID), "coolio", &newTag, "#abc")
	require.NoError(t, err)

	db.Find(&entry)
	require.Equal(t, "abc,yes,chicken", entry.Keys)
}

func TestUpdate_noPermissions(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	left := db.User(5)
	right := db.User(2)
	right.NewTagDefinition("coolio")
	rightTs := right.TimeSpan("2009-06-30T18:30:00Z", "2009-06-30T18:40:00Z")
	rightTs.Tag("coolio", "mama")

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.UpdateTag(fake.User(left.User.ID), "coolio", nil, "#abc")
	require.EqualError(t, err, "tag with key 'coolio' does not exist")
	right.AssertHasTagDefinition("coolio", true)
}
