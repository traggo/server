package graphql

import (
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/handler"
	"github.com/traggo/server/generated/gqlschema"
	"github.com/traggo/server/logger"
)

// Handler combines graphql handler and playground handler.
func Handler(endpoint string, resolvers gqlschema.ResolverRoot) http.HandlerFunc {
	gqlHandler := handler.GraphQL(gqlschema.NewExecutableSchema(gqlschema.Config{
		Resolvers: resolvers,
	}), handler.RequestMiddleware(logger.GQLLog()))
	playground := handler.Playground("Traggo Playground", endpoint)

	return func(writer http.ResponseWriter, request *http.Request) {
		if acceptHTMLAndNotJSON(request) {
			playground.ServeHTTP(writer, request)
		} else {
			gqlHandler.ServeHTTP(writer, request)
		}
	}
}

func acceptHTMLAndNotJSON(request *http.Request) bool {
	val := request.Header.Get("Accept")
	return strings.Contains(val, "text/html") && !strings.Contains(val, "application/json")
}
