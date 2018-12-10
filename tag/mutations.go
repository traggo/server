package tag

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/traggo/server/schema"
)

// Mutations for tags
func Mutations(db *gorm.DB) graphql.Fields {
	return graphql.Fields{
		"createTag": createTag(db),
		"removeTag": removeTag(db),
	}
}

func removeTag(db *gorm.DB) *graphql.Field {
	return &graphql.Field{
		Type: schema.TagDefinitionSchema,
		Args: graphql.FieldConfigArgument{
			"key": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			definition := schema.TagDefinition{Key: p.Args["key"].(string)}
			if db.Find(&definition).RecordNotFound() {
				return nil, fmt.Errorf("tag with key '%s' does not exist", definition.Key)
			}

			remove := db.Delete(&definition)
			return definition, remove.Error
		},
	}
}

func createTag(db *gorm.DB) *graphql.Field {
	return &graphql.Field{
		Type: schema.TagDefinitionSchema,
		Args: graphql.FieldConfigArgument{
			"key": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"color": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"type": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(schema.TagDefinitionTypeSchema),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			definition := &schema.TagDefinition{
				Key:   p.Args["key"].(string),
				Color: p.Args["color"].(string),
				Type:  p.Args["type"].(schema.TagDefinitionType),
			}

			if !db.Find(definition).RecordNotFound() {
				return nil, fmt.Errorf("tag with key '%s' does already exist", definition.Key)
			}

			create := db.Create(&definition)
			return definition, create.Error
		},
	}
}
