package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/bdreece/hopper/pkg/config"
)

func NewGraphQLServer(cfg *config.Config) http.Handler {
	schema := NewExecutableSchema(Config{
		Resolvers: NewResolver(cfg),
	})

	return handler.NewDefaultServer(schema)
}
