package tag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/traggo/server/generated/gqlmodel"
	"github.com/traggo/server/model"
	"github.com/traggo/server/test"
)

func TestGQL_SuggestTag_matchesTags(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	resolver := ResolverForTag{DB: db}
	db.Create(&model.TagDefinition{Key: "project", Color: "#fff", Type: model.TypeSingleValue})
	db.Create(&model.TagDefinition{Key: "priority", Color: "#fff", Type: model.TypeSingleValue})
	db.Create(&model.TagDefinition{Key: "wood", Color: "#fff", Type: model.TypeSingleValue})

	tags, err := resolver.SuggestTag(context.Background(), "pr")

	require.Nil(t, err)
	expected := []gqlmodel.TagDefinition{
		{Key: "project", Color: "#fff", Type: gqlmodel.TagDefinitionTypeSinglevalue},
		{Key: "priority", Color: "#fff", Type: gqlmodel.TagDefinitionTypeSinglevalue},
	}
	require.Equal(t, expected, tags)
}

func TestGQL_SuggestTag_noMatchingTags(t *testing.T) {
	db := test.InMemoryDB(t)
	defer db.Close()
	resolver := ResolverForTag{DB: db}
	db.Create(&model.TagDefinition{Key: "project", Color: "#fff", Type: model.TypeSingleValue})
	db.Create(&model.TagDefinition{Key: "wood", Color: "#fff", Type: model.TypeSingleValue})

	tags, err := resolver.SuggestTag(context.Background(), "fire")

	require.Nil(t, err)
	var expected []gqlmodel.TagDefinition
	require.Equal(t, expected, tags)
}
