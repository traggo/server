package user_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/traggo/server/schema"
	"github.com/traggo/server/test"
	"github.com/traggo/server/user"
)

func Test_getUsers_withoutUsers(t *testing.T) {
	query := `
query {
   users {
      id,
      name,
      admin
   }
}`
	expected := `
{
  "users": []
}`
	test.NewGQL(t).Queries(user.Queries).Exec(query).Succeeds(expected)
}

func Test_getUsers_withUsers(t *testing.T) {
	query := `
query {
   users {
      id,
      name,
      admin
   }
}`
	createUsers := func(db *gorm.DB) {
		db.Create(&schema.User{Name: "jmattheis", Pass: []byte{1, 2}, Admin: true})
		db.Create(&schema.User{Name: "broderpeters", Pass: []byte{1, 2}, Admin: true})
	}

	expected := `
{
  "users": [
    {
      "admin": true,
      "id": 1,
      "name": "jmattheis"
    },
    {
      "admin": true,
      "id": 2,
      "name": "broderpeters"
    }
  ]
}`
	test.NewGQL(t).Queries(user.Queries).BeforeExec(createUsers).Exec(query).Succeeds(expected)
}

func Test_getUsers_cannotQueryPassword(t *testing.T) {
	query := `
query {
   users {
      pass
   }
}`
	createUsers := func(db *gorm.DB) {
		db.Create(&schema.User{Name: "jmattheis", Pass: []byte{1, 2}, Admin: true})
	}

	expectedError := `Cannot query field "pass" on type "User".`
	test.NewGQL(t).Queries(user.Queries).BeforeExec(createUsers).Exec(query).Errs(expectedError)
}
