package model

// TagDefinition describes a tag.
type TagDefinition struct {
	Key    string
	UserID int `gorm:"type:int REFERENCES users(id) ON DELETE CASCADE"`
	Color  string
	Usages int `gorm:"-"`
}
