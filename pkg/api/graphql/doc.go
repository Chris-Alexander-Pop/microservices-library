// Package graphql provides utilities for building GraphQL APIs.
//
// This package integrates with standard GraphQL libraries to provide:
//   - Schema loading and validation
//   - Resolver middleware
//   - Complexity limits
//   - Error formatting
//
// Usage:
//
//	srv := graphql.NewServer(schema, resolvers)
//	http.Handle("/query", srv)
package graphql
