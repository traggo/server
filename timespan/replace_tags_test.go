package timespan

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestReplace_simple(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	user := db.User(1)
	user.NewTagDefinition("old")
	user.NewTagDefinition("new")
	ts := user.TimeSpan("2019-06-11T18:00:00Z", "2019-06-11T18:00:00Z")
	ts.Tag("old", "1")

	resolver := ResolverForTimeSpan{DB: db.DB}

	_, err := resolver.ReplaceTimeSpanTags(fake.User(1),
		gqlmodel.InputTimeSpanTag{Key: "old", Value: "1"},
		gqlmodel.InputTimeSpanTag{Key: "new", Value: "1.0"},
		gqlmodel.InputReplaceOptions{Override: gqlmodel.OverrideModeDiscard})
	require.NoError(t, err)

	ts.AssertHasTag("new", "1.0", true)
	ts.AssertHasTag("old", "1", false)
}

func TestReplace_discard(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	user := db.User(1)
	user.NewTagDefinition("old")
	user.NewTagDefinition("new")
	ts := user.TimeSpan("2019-06-11T18:00:00Z", "2019-06-11T18:00:00Z")
	ts.Tag("old", "1")
	ts.Tag("new", "2.0")

	resolver := ResolverForTimeSpan{DB: db.DB}

	_, err := resolver.ReplaceTimeSpanTags(fake.User(1),
		gqlmodel.InputTimeSpanTag{Key: "old", Value: "1"},
		gqlmodel.InputTimeSpanTag{Key: "new", Value: "1.0"},
		gqlmodel.InputReplaceOptions{Override: gqlmodel.OverrideModeDiscard})
	require.NoError(t, err)

	ts.AssertHasTag("new", "2.0", true)
	ts.AssertHasTag("new", "1.0", false)
	ts.AssertHasTagIgnoreValue("old", false)
}

func TestReplace_override(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	user := db.User(1)
	user.NewTagDefinition("old")
	user.NewTagDefinition("new")
	ts := user.TimeSpan("2019-06-11T18:00:00Z", "2019-06-11T18:00:00Z")
	ts.Tag("old", "1")
	ts.Tag("new", "2.0")

	resolver := ResolverForTimeSpan{DB: db.DB}

	_, err := resolver.ReplaceTimeSpanTags(fake.User(1),
		gqlmodel.InputTimeSpanTag{Key: "old", Value: "1"},
		gqlmodel.InputTimeSpanTag{Key: "new", Value: "1.0"},
		gqlmodel.InputReplaceOptions{Override: gqlmodel.OverrideModeOverride})
	require.NoError(t, err)

	ts.AssertHasTag("new", "2.0", false)
	ts.AssertHasTag("new", "1.0", true)
	ts.AssertHasTagIgnoreValue("old", false)
}

func TestReplace_multiuser(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	left := db.User(1)
	right := db.User(2)

	left.NewTagDefinition("old")
	left.NewTagDefinition("new")
	right.NewTagDefinition("old")
	right.NewTagDefinition("new")
	leftTs := left.TimeSpan("2019-06-11T18:00:00Z", "2019-06-11T18:00:00Z")
	leftTs.Tag("old", "1")
	leftTs.Tag("new", "2.0")
	rightTs := right.TimeSpan("2019-06-11T18:00:00Z", "2019-06-11T18:00:00Z")
	rightTs.Tag("old", "1")
	rightTs.Tag("new", "2.0")

	resolver := ResolverForTimeSpan{DB: db.DB}

	_, err := resolver.ReplaceTimeSpanTags(fake.User(1),
		gqlmodel.InputTimeSpanTag{Key: "old", Value: "1"},
		gqlmodel.InputTimeSpanTag{Key: "new", Value: "1.0"},
		gqlmodel.InputReplaceOptions{Override: gqlmodel.OverrideModeOverride})
	require.NoError(t, err)

	leftTs.AssertHasTag("new", "2.0", false)
	leftTs.AssertHasTag("new", "1.0", true)
	leftTs.AssertHasTagIgnoreValue("old", false)
	rightTs.AssertHasTag("old", "1", true)
	rightTs.AssertHasTag("new", "2.0", true)
}
