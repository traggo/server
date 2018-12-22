package gql

import (
	"log"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/jinzhu/gorm"
	"github.com/traggo/server/tag"
	"github.com/traggo/server/user"
)

// Handler creates a graphql handler.
func Handler(db *gorm.DB, passwordStrength int) *handler.Handler {
	queryFields := merge(tag.Queries(db), user.Queries(db))
	mutationFields := merge(tag.Mutations(db), user.Mutations(db, passwordStrength))

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queryFields}
	rootMutations := graphql.ObjectConfig{Name: "Mutations", Fields: mutationFields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: graphql.NewObject(rootMutations)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	return handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     false,
		GraphiQL:   true,
		Playground: true,
	})
}

func merge(toMerge ...graphql.Fields) graphql.Fields {
	var result = graphql.Fields{}
	for _, subset := range toMerge {
		for key, value := range subset {
			result[key] = value
		}
	}
	return result
}
