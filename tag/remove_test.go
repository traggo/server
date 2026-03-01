package tag

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestGQL_RemoveTag_succeeds_removesTag(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	user := db.User(3)
	user.NewTagDefinition("existing")

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.RemoveTag(fake.User(3), "existing")
	require.Nil(t, err)
	user.AssertHasTagDefinition("existing", false)
}

func TestRemove_referencedInDashboardEntry(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	left := db.User(5)
	left.NewTagDefinition("coolio")
	dashboard := left.Dashboard("yeah")
	dashboard.Entry("entry")
	entry := dashboard.Dashboard.Entries[0]
	entry.Keys = "abc,coolio,chicken"
	db.Save(&entry)

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.RemoveTag(fake.User(left.User.ID), "coolio")
	require.EqualError(t, err, "tag 'coolio' is used in dashboard 'yeah' entry 'entry', remove this reference before deleting the tag")
}

func TestGQL_RemoveTag_succeeds_removesTimespans(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	left := db.User(3)
	right := db.User(4)
	left.NewTagDefinition("tag")
	right.NewTagDefinition("tag")
	leftTs := left.TimeSpan("2009-06-30T18:30:00Z", "2009-06-30T18:40:00Z")
	leftTs.Tag("tag", "def")
	rightTs := right.TimeSpan("2009-06-30T18:30:00Z", "2009-06-30T18:40:00Z")
	rightTs.Tag("tag", "def")

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.RemoveTag(fake.User(left.User.ID), "tag")
	require.Nil(t, err)

	assertTagCount(t, db, 1)

	left.AssertHasTagDefinition("tag", false)
	right.AssertHasTagDefinition("tag", true)

	leftTs.AssertExists(true).AssertHasTag("tag", "def", false)
	rightTs.AssertExists(true).AssertHasTag("tag", "def", true)
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
	db.User(3).NewTagDefinition("existing")
	db.User(5)

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.RemoveTag(fake.User(5), "existing")
	require.EqualError(t, err, "tag with key 'existing' does not exist")
}

func TestRemove_crossUserIsolation(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	userA := db.User(5)
	userB := db.User(2)

	userA.NewTagDefinition("project")
	userB.NewTagDefinition("project")

	dashboardB := userB.Dashboard("secretDashboard")
	dashboardB.Entry("secretEntry")
	entryB := dashboardB.Dashboard.Entries[0]
	entryB.Keys = "project,secret"
	db.Save(&entryB)

	resolver := ResolverForTag{DB: db.DB}
	_, err := resolver.RemoveTag(fake.User(userA.User.ID), "project")
	require.NoError(t, err, "removing user's own tag should succeed even if another user has a dashboard entry with that tag")

	userA.AssertHasTagDefinition("project", false)
	userB.AssertHasTagDefinition("project", true)
}
