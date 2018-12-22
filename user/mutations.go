package user

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/traggo/server/schema"
	"github.com/traggo/server/user/password"
)

var createPassword = password.CreatePassword

// Mutations for users
func Mutations(db *gorm.DB, passwordStrength int) graphql.Fields {
	return graphql.Fields{
		"createUser": createUser(db, passwordStrength),
		"removeUser": removeUser(db),
		"updateUser": updateUser(db, passwordStrength),
	}
}

func removeUser(db *gorm.DB) *graphql.Field {
	return &graphql.Field{
		Type: schema.UserSchema,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			user := schema.User{ID: p.Args["id"].(int)}
			if db.Find(&user).RecordNotFound() {
				return nil, fmt.Errorf("user with id %d does not exist", user.ID)
			}

			remove := db.Delete(&user)
			return user, remove.Error
		},
	}
}

func createUser(db *gorm.DB, passwordStrength int) *graphql.Field {
	return &graphql.Field{
		Type: schema.UserSchema,
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"pass": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"admin": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			newUser := &schema.User{
				Name:  p.Args["name"].(string),
				Pass:  createPassword(p.Args["pass"].(string), passwordStrength),
				Admin: p.Args["admin"].(bool),
			}

			if !db.Where("name = ?", newUser.Name).Find(&schema.User{}).RecordNotFound() {
				return nil, fmt.Errorf("user with name '%s' does already exist", newUser.Name)
			}

			create := db.Create(&newUser)
			return newUser, create.Error
		},
	}
}

func updateUser(db *gorm.DB, passwordStrength int) *graphql.Field {
	return &graphql.Field{
		Type: schema.UserSchema,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"pass": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"admin": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id := p.Args["id"].(int)

			user := new(schema.User)
			if db.Find(user, id).RecordNotFound() {
				return nil, fmt.Errorf("user with id %d does not exist", id)
			}

			user.Name = p.Args["name"].(string)
			user.Admin = p.Args["admin"].(bool)

			if pass, ok := p.Args["pass"]; ok {
				user.Pass = createPassword(pass.(string), passwordStrength)
			}

			update := db.Save(user)
			return user, update.Error
		},
	}
}
