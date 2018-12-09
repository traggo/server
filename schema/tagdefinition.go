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

type TagDefinitionType string

const (
	TypeNoValue     TagDefinitionType = "novalue"
	TypeSingleValue TagDefinitionType = "singlevalue"
)

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

var TagDefinitionSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TagDefinition",
		Fields: graphql.Fields{
			"key":   &graphql.Field{Type: graphql.String},
			"color": &graphql.Field{Type: graphql.String},
			"type":  &graphql.Field{Type: TagDefinitionTypeSchema},
		},
	})
