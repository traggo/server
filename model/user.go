package model

// User holds information about credentials and authorizations.
type User struct {
	ID    int    `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Name  string `gorm:"unique"`
	Pass  []byte
	Admin bool
}
