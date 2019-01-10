package tag

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
	"github.com/traggo/server/test/fake"
)

func TestGQL_Tags(t *testing.T) {
	db := test.InMemoryDB(t)
	db.User(5)
	defer db.Close()
	db.Create(&model.TagDefinition{Key: "my tag", Color: "#fff", Type: model.TypeSingleValue, UserID: 5})
	db.Create(&model.TagDefinition{Key: "my tag 2", Color: "#fff", Type: model.TypeSingleValue, UserID: 5})
	db.Create(&model.TagDefinition{Key: "my tag 5", Color: "#fff", Type: model.TypeSingleValue, UserID: 2})

	resolver := ResolverForTag{DB: db.DB}
	tags, err := resolver.Tags(fake.User(5))

	require.Nil(t, err)
	expected := []gqlmodel.TagDefinition{
		{Key: "my tag", Color: "#fff", Type: gqlmodel.TagDefinitionTypeSinglevalue},
		{Key: "my tag 2", Color: "#fff", Type: gqlmodel.TagDefinitionTypeSinglevalue},
	}
	require.Equal(t, expected, tags)
}
