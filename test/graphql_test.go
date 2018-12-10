package test_test

import (
	"errors"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/traggo/server/test"
)

type fakeTesting struct {
	hasErrors bool
}

func (t *fakeTesting) Errorf(format string, args ...interface{}) {
	t.hasErrors = true
}

func TestGQLTest_expectSuccess(t *testing.T) {
	fields := func(db *gorm.DB) graphql.Fields {
		return graphql.Fields{
			"test": &graphql.Field{
				Name: "test",
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "value", nil
				},
			},
		}
	}

	query := `
query {
   test
}`
	expect := `
{
  "test": "value"
}`
	test.NewGQL(t).Queries(fields).Exec(query).Succeeds(expect)
}

func TestGQLTest_mutation(t *testing.T) {
	state := false
	fields := func(db *gorm.DB) graphql.Fields {
		return graphql.Fields{
			"toggleEnabled": &graphql.Field{
				Name: "test",

				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					state = !state
					return state, nil
				},
			},
		}
	}

	query := `
mutation {
   toggleEnabled 
}`
	expect := `
{
  "toggleEnabled": true
}`
	test.NewGQL(t).Mutations(fields).Exec(query).Succeeds(expect)
	assert.True(t, state)
}

func TestGQLTest_executeBeforeExec(t *testing.T) {
	returnValue := "old"
	fields := func(db *gorm.DB) graphql.Fields {
		return graphql.Fields{
			"test": &graphql.Field{
				Name: "test",
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return returnValue, nil
				},
			},
		}
	}

	query := `
query {
   test
}`

	before := func(db *gorm.DB) {
		returnValue = "new"
	}
	expect := `
{
  "test": "new"
}`

	test.NewGQL(t).Queries(fields).BeforeExec(before).Exec(query).Succeeds(expect)
}

func TestGQLTest_expectError(t *testing.T) {
	fields := func(db *gorm.DB) graphql.Fields {
		return graphql.Fields{
			"test": &graphql.Field{
				Name: "test",
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, errors.New("test error")
				},
			},
		}
	}

	query := `
query {
   test
}`

	test.NewGQL(t).Queries(fields).Exec(query).Errs("test error")
}

func TestGQLTest_expectErrorButWasSuccess(t *testing.T) {
	fields := func(db *gorm.DB) graphql.Fields {
		return graphql.Fields{
			"test": &graphql.Field{
				Name: "test",
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "success", nil
				},
			},
		}
	}

	query := `
query {
   test
}`

	fakeT := &fakeTesting{}

	test.NewGQL(fakeT).Queries(fields).Exec(query).Errs("test error")
	assert.True(t, fakeT.hasErrors)
}

func TestGQLTest_expectSuccessButWasError(t *testing.T) {
	fields := func(db *gorm.DB) graphql.Fields {
		return graphql.Fields{
			"test": &graphql.Field{
				Name: "test",
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, errors.New("error")
				},
			},
		}
	}

	query := `
query {
   test
}`
	expect := `
{
  "test": "new"
}`
	fakeT := &fakeTesting{}

	test.NewGQL(fakeT).Queries(fields).Exec(query).Succeeds(expect)
	assert.True(t, fakeT.hasErrors)
}

func TestGQLTest_expectErrorWithWrongMessage(t *testing.T) {
	fields := func(db *gorm.DB) graphql.Fields {
		return graphql.Fields{
			"test": &graphql.Field{
				Name: "test",
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, errors.New("yes")
				},
			},
		}
	}

	query := `
query {
   test
}`
	fakeT := &fakeTesting{}

	test.NewGQL(fakeT).Queries(fields).Exec(query).Errs("this is not 'yes'")
	assert.True(t, fakeT.hasErrors)
}

func TestGQLTest_expectSuccessWithWrongExpectJSON(t *testing.T) {
	fields := func(db *gorm.DB) graphql.Fields {
		return graphql.Fields{
			"test": &graphql.Field{
				Name: "test",
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "other value", nil
				},
			},
		}
	}

	query := `
query {
   test
}`
	expect := `
{
  "test": "not really"
}`
	fakeT := &fakeTesting{}

	test.NewGQL(fakeT).Queries(fields).Exec(query).Succeeds(expect)
	assert.True(t, fakeT.hasErrors)
}
