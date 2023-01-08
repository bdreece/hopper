package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/bdreece/hopper/pkg/config"
	"github.com/bdreece/hopper/pkg/graphql/resolvers"
)

func NewGraphQLHandler(cfg *config.Config) http.Handler {
	schema := resolvers.NewExecutableSchema(resolvers.Config{
		Resolvers: resolvers.NewResolver(cfg),
	})

	return handler.NewDefaultServer(schema)
}
