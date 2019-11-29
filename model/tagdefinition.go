package model

// TagDefinition describes a tag.
type TagDefinition struct {
	Key    string
	UserID int `gorm:"type:int REFERENCES users(id) ON DELETE CASCADE"`
	Color  string
	Type   TagDefinitionType
	Usages int `gorm:"-"`
}

// TagDefinitionType describes a tag type.
type TagDefinitionType string

const (
	// TypeNoValue used for tags without values
	TypeNoValue TagDefinitionType = "novalue"
	// TypeSingleValue used for tags with one value
	TypeSingleValue TagDefinitionType = "singlevalue"
)
