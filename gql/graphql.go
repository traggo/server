package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"github.com/traggo/server/logger"
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
	gqlSchema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create graphql schema")
	}

	return handler.New(&handler.Config{
		Schema:           &gqlSchema,
		Pretty:           false,
		GraphiQL:         true,
		Playground:       true,
		ResultCallbackFn: logger.GQLLog,
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
