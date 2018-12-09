package schema

import (
	"github.com/graphql-go/graphql"
)

// TagDefinition describes a tag.
type TagDefinition struct {
	Key   string `gorm:"primary_key;unique_index"`
	Color string
	Type  TagDefinitionType
	Owner uint
}

// TagDefinitionType describes a tag type.
type TagDefinitionType string

const (
	// TypeNoValue used for tags without values
	TypeNoValue TagDefinitionType = "novalue"
	// TypeSingleValue used for tags with one value
	TypeSingleValue TagDefinitionType = "singlevalue"
)

// TagDefinitionTypeSchema is a gql representation of TagDefinitionType
var TagDefinitionTypeSchema = graphql.NewEnum(graphql.EnumConfig{
	Name: "TagDefinitionType",
	Values: graphql.EnumValueConfigMap{
		"novalue": &graphql.EnumValueConfig{
			Value: TypeNoValue,
		},
		"singlevalue": &graphql.EnumValueConfig{
			Value: TypeSingleValue,
		},
	},
})

// TagDefinitionSchema is a gql representation of TagDefinition
var TagDefinitionSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TagDefinition",
		Fields: graphql.Fields{
			"key":   &graphql.Field{Type: graphql.String},
			"color": &graphql.Field{Type: graphql.String},
			"type":  &graphql.Field{Type: TagDefinitionTypeSchema},
		},
	})
