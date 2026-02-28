package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/NR3101/go-ecom-project/graph"
	"github.com/NR3101/go-ecom-project/graph/resolver"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

// createGraphQLHandler initializes the GraphQL server with the necessary resolvers and transports.
func (s *Server) createGraphQLHandler() *handler.Server {
	// Initialize the resolver with the required services.
	rvr := resolver.NewResolver(
		s.authService,
		s.userService,
		s.productService,
		s.cartService,
		s.orderService,
	)

	// Create the executable schema using the generated GraphQL schema and the resolver.
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: rvr})

	// Create a new GraphQL server instance with the executable schema.
	srv := handler.New(schema)

	// Add various transports to the server to handle different types of GraphQL requests.
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	// Enable introspection and automatic persisted queries for better performance and debugging.
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

// graphqlHandler returns a Gin handler function that serves the GraphQL endpoint using the configured GraphQL server.
func (s *Server) graphqlHandler() gin.HandlerFunc {
	// Create the GraphQL handler using the createGraphQLHandler method.
	srv := s.createGraphQLHandler()

	// Return a Gin handler function that serves the GraphQL endpoint.
	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

// playgroundHandler returns a Gin handler function that serves the GraphQL Playground, allowing developers to interact with the GraphQL API in a user-friendly interface.
func (s *Server) playgroundHandler() gin.HandlerFunc {
	srv := playground.Handler("GraphQL Playground", "/graphql")

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}
