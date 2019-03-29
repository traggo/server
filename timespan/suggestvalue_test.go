package timespan

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

const date = "2019-06-10T18:30:00Z"

func TestGQL_SuggestTagValue(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	user := db.User(1)
	user.TimeSpan(date, date).Tag("proj", "gotify").Tag("issue", "3")
	user.TimeSpan(date, date).Tag("proj", "traggo").Tag("issue", "3")
	user.TimeSpan(date, date).Tag("proj", "traggo").Tag("issue", "3")
	user.TimeSpan(date, date).Tag("proj", "meh").Tag("issue", "3")
	other := db.User(2)
	other.TimeSpan(date, date).Tag("proj", "secret").Tag("issue", "3")
	resolver := ResolverForTimeSpan{DB: db.DB}

	test.LogDebug()
	tags, err := resolver.SuggestTagValue(fake.User(1), "proj", "")

	require.Nil(t, err)
	expected := []string{"gotify", "traggo", "meh"}
	require.Equal(t, expected, tags)

	tags, err = resolver.SuggestTagValue(fake.User(1), "proj", "uff")

	require.Nil(t, err)
	require.Empty(t, tags)
}
