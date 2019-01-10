package model

// TagDefinition describes a tag.
type TagDefinition struct {
	Key    string
	UserID int
	Color  string
	Type   TagDefinitionType
}

// TagDefinitionType describes a tag type.
type TagDefinitionType string

const (
	// TypeNoValue used for tags without values
	TypeNoValue TagDefinitionType = "novalue"
	// TypeSingleValue used for tags with one value
	TypeSingleValue TagDefinitionType = "singlevalue"
)
