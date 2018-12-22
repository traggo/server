package user

import (
	"os"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/schema"
	"github.com/traggo/server/test"
)

var (
	ponyPW       = []byte{1}
	unicornPW    = []byte{2}
	fakePassword = func(pw string, strength int) []byte {
		if pw == "pony" {
			return ponyPW
		} else if pw == "unicorn" {
			return unicornPW
		}
		panic("unknown pw")
	}
)

func TestMain(m *testing.M) {
	createPassword = fakePassword
	retCode := m.Run()
	os.Exit(retCode)
}

func userExist(t *testing.T, expected schema.User) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		foundUser := new(schema.User)
		find := db.Find(foundUser, expected.ID)
		assert.Nil(t, find.Error)
		assert.NotNil(t, foundUser)
		assert.Equal(t, expected, *foundUser)
	}
}

func userCount(t *testing.T, expected int) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		count := new(int)
		db.Model(new(schema.User)).Count(count)
		assert.Equal(t, expected, *count)
	}
}

func mutationsWithStrength(db *gorm.DB) graphql.Fields {
	return Mutations(db, 4)
}

func Test_createUser_succeeds(t *testing.T) {
	query := `
mutation {
   createUser(name: "jmattheis", pass: "unicorn", admin:true) {
      name,
      admin,
      id
   }
}`
	expected := `
{
  "createUser": {
    "admin": true,
    "id": 1,
    "name": "jmattheis"
  }
}`
	user := schema.User{
		ID:    1,
		Name:  "jmattheis",
		Pass:  unicornPW,
		Admin: true,
	}

	test.NewGQL(t).Mutations(mutationsWithStrength).DBAssert(userExist(t, user)).Exec(query).Succeeds(expected)

}

func Test_createUser_withExistingUserName_fails(t *testing.T) {
	createUser := func(db *gorm.DB) {
		db.Create(&schema.User{Name: "jmattheis", Pass: []byte{1}, Admin: true})
	}
	query := `
mutation {
   createUser(name: "jmattheis", pass: "unicorn", admin:true) {
      name,
      admin,
      id
   }
}`
	expectedError := "user with name 'jmattheis' does already exist"

	test.NewGQL(t).Mutations(mutationsWithStrength).BeforeExec(createUser).DBAssert(userCount(t, 1)).Exec(query).Errs(expectedError)
}

func Test_removeUser_withExistingUser_succeeds(t *testing.T) {
	createUser := func(db *gorm.DB) {
		db.Create(&schema.User{Name: "jmattheis", Pass: []byte{1}, Admin: true})
	}
	query := `
mutation {
   removeUser(id: 1) {
      name,
      admin,
      id
   }
}`
	expected := `
{
  "removeUser": {
    "admin": true,
    "id": 1,
    "name": "jmattheis"
  }
}`

	test.NewGQL(t).Mutations(mutationsWithStrength).BeforeExec(createUser).DBAssert(userCount(t, 0)).Exec(query).Succeeds(expected)
}

func Test_updateUser_withExistingUser_succeeds(t *testing.T) {
	createUser := func(db *gorm.DB) {
		db.Create(&schema.User{Name: "jmattheis", Pass: unicornPW, Admin: true})
	}
	query := `
mutation {
   updateUser(id: 1, name: "broderpeters", pass:"pony", admin:false) {
      name,
      admin,
      id
   }
}`
	expected := `
{
  "updateUser": {
    "admin": false,
    "id": 1,
    "name": "broderpeters"
  }
}`
	user := schema.User{
		ID:    1,
		Name:  "broderpeters",
		Pass:  ponyPW,
		Admin: false,
	}

	test.NewGQL(t).Mutations(mutationsWithStrength).BeforeExec(createUser).DBAssert(userExist(t, user)).Exec(query).Succeeds(expected)
}

func Test_updateUser_withExistingUser_preservesPassword(t *testing.T) {
	createUser := func(db *gorm.DB) {
		db.Create(&schema.User{Name: "jmattheis", Pass: unicornPW, Admin: true})
	}
	query := `
mutation {
   updateUser(id: 1, name: "broderpeters", admin:false) {
      name,
      admin,
      id
   }
}`
	expected := `
{
  "updateUser": {
    "admin": false,
    "id": 1,
    "name": "broderpeters"
  }
}`
	user := schema.User{
		ID:    1,
		Name:  "broderpeters",
		Pass:  unicornPW,
		Admin: false,
	}

	test.NewGQL(t).Mutations(mutationsWithStrength).BeforeExec(createUser).DBAssert(userExist(t, user)).Exec(query).Succeeds(expected)
}

func Test_removeUser_withoutExistingUser_fails(t *testing.T) {
	query := `
mutation {
   removeUser(id: 1) {
      name,
      admin,
      id
   }
}`

	expectedError := "user with id 1 does not exist"

	test.NewGQL(t).Mutations(mutationsWithStrength).Exec(query).Errs(expectedError)
}

func Test_updateUser_withoutExistingUser_fails(t *testing.T) {
	query := `
mutation {
   updateUser(id: 1, name: "jmattheis2", pass:"123", admin:false) {
      name,
      admin,
      id
   }
}`

	expectedError := "user with id 1 does not exist"

	test.NewGQL(t).Mutations(mutationsWithStrength).DBAssert(userCount(t, 0)).Exec(query).Errs(expectedError)
}
