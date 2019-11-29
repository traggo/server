package tag

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestGQL_Tags(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	left := db.User(5)
	right := db.User(2)
	left.NewTagDefinition("my tag")
	left.NewTagDefinition("my tag 2")
	right.NewTagDefinition("my tag 5")

	resolver := ResolverForTag{DB: db.DB}
	tags, err := resolver.Tags(fake.User(left.User.ID))

	require.Nil(t, err)
	expected := []*gqlmodel.TagDefinition{
		{Key: "my tag", Type: gqlmodel.TagDefinitionTypeSinglevalue},
		{Key: "my tag 2", Type: gqlmodel.TagDefinitionTypeSinglevalue},
	}
	require.Equal(t, expected, tags)
}
