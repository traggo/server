package tag

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/traggo/server/schema"
)

// Queries for tags
func Queries(db *gorm.DB) graphql.Fields {
	return graphql.Fields{
		"tags":       getTags(db),
		"suggestTag": suggestTag(db),
	}
}

func suggestTag(db *gorm.DB) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(schema.TagDefinitionSchema),
		Args: graphql.FieldConfigArgument{
			"query": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			query := p.Args["query"].(string)
			var suggestions []schema.TagDefinition
			find := db.Where("Key LIKE ?", query+"%").Find(&suggestions)
			return suggestions, find.Error
		},
	}
}

func getTags(db *gorm.DB) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(schema.TagDefinitionSchema),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var tags []schema.TagDefinition
			db.Find(&tags)
			return tags, nil
		},
	}
}
