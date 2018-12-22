package user

import (
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/traggo/server/schema"
)

// Queries for users
func Queries(db *gorm.DB) graphql.Fields {
	return graphql.Fields{
		"users": getUsers(db),
	}
}

func getUsers(db *gorm.DB) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(schema.UserSchema),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var users []schema.User
			db.Find(&users)
			return users, nil
		},
	}
}
