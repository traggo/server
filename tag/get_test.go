package tag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func TestGQL_Tags(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	db.Create(&model.TagDefinition{Key: "my tag", Color: "#fff", Type: model.TypeSingleValue})
	db.Create(&model.TagDefinition{Key: "my tag 2", Color: "#fff", Type: model.TypeSingleValue})

	resolver := ResolverForTag{DB: db}
	tags, err := resolver.Tags(context.Background())

	require.Nil(t, err)
	expected := []gqlmodel.TagDefinition{
		{Key: "my tag", Color: "#fff", Type: gqlmodel.TagDefinitionTypeSinglevalue},
		{Key: "my tag 2", Color: "#fff", Type: gqlmodel.TagDefinitionTypeSinglevalue},
	}
	require.Equal(t, expected, tags)
}
