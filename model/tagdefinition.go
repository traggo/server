package model

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
