package schema

import "github.com/graphql-go/graphql"

// User holds information about credentials and authorizations.
type User struct {
	ID    int    `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Name  string `gorm:"unique_index"`
	Pass  []byte
	Admin bool
}

// UserSchema is a gql representation of User.
var UserSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":    &graphql.Field{Type: graphql.Int},
			"name":  &graphql.Field{Type: graphql.String},
			"admin": &graphql.Field{Type: graphql.Boolean},
		},
	})
