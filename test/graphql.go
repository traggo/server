package test

import (
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

// GQLTest holds information about a gql test.
type GQLTest struct {
	t              assert.TestingT
	queryFields    func(db *gorm.DB) graphql.Fields
	mutationFields func(db *gorm.DB) graphql.Fields
	before         func(db *gorm.DB)
	dbAssert       func(db *gorm.DB)
}

// NewGQL create a new graphql testing instance.
func NewGQL(t assert.TestingT) *GQLTest {
	return &GQLTest{
		t: t,
		queryFields: func(db *gorm.DB) graphql.Fields {
			return graphql.Fields{"dummy": &graphql.Field{
				Type: graphql.String,
			}}
		},
		before: func(db *gorm.DB) {
		},
		dbAssert: func(db *gorm.DB) {
		},
	}
}

// BeforeExec will be executed before execution of the gql query.
// Use this for filling the db.
func (g *GQLTest) BeforeExec(f func(db *gorm.DB)) *GQLTest {
	g.before = f
	return g
}

// DBAssert will be executed after execution of the gql query.
func (g *GQLTest) DBAssert(f func(db *gorm.DB)) *GQLTest {
	g.dbAssert = f
	return g
}

// Queries sets the gql query fields.
func (g *GQLTest) Queries(f func(db *gorm.DB) graphql.Fields) *GQLTest {
	g.queryFields = f
	return g
}

// Mutations sets the gql mutation fields.
func (g *GQLTest) Mutations(f func(db *gorm.DB) graphql.Fields) *GQLTest {
	g.mutationFields = f
	return g
}

// Exec executes the gql query.
func (g *GQLTest) Exec(query string) *GQLResult {
	db := InMemoryDB(g.t)
	defer db.Close()

	g.before(db)
	schemaConfig := graphql.SchemaConfig{}

	schemaConfig.Query = graphql.NewObject(graphql.ObjectConfig{Name: "RootQuery", Fields: g.queryFields(db)})

	if g.mutationFields != nil {
		schemaConfig.Mutation = graphql.NewObject(graphql.ObjectConfig{Name: "Mutations", Fields: g.mutationFields(db)})
	}

	schema, err := graphql.NewSchema(schemaConfig)
	assert.Nil(g.t, err)
	result := graphql.Do(graphql.Params{
		RequestString: query,
		Schema:        schema,
	})
	g.dbAssert(db)
	return &GQLResult{Result: result, GQLTest: g}
}

// GQLResult holds information about a gql test result.
type GQLResult struct {
	*GQLTest
	*graphql.Result
}

// Succeeds asserts the gql result to be successful
// and asserts the result to match the "expectedJSON" string.
func (r *GQLResult) Succeeds(expectedJSON string) {
	if assert.Empty(r.t, r.Errors) {
		asJSON := toJSON(r.t, r.Data)
		if !assert.JSONEq(r.t, expectedJSON, asJSON) {
			fmt.Println(asJSON)
		}
	}
}

// Errs asserts the gql result to be unsuccessful
// and asserts error message to match the "expected" string.
func (r *GQLResult) Errs(expected string) {
	if assert.Len(r.t, r.Errors, 1) {
		assert.EqualError(r.t, r.Errors[0], expected)
	}
}

func toJSON(t assert.TestingT, data interface{}) string {
	bytes, err := json.MarshalIndent(data, "", "  ")
	assert.Nil(t, err)
	return string(bytes)
}
