// Package rest provides utilities for building RESTful APIs.
//
// This package includes:
//   - Standardized JSON response formatting
//   - Error handling and mapping to HTTP status codes
//   - Middleware for common concerns (CORS, RequestID, Logging)
//   - Router configuration helpers
//
// Usage:
//
//	r := rest.NewRouter()
//	r.Get("/health", rest.HealthHandler)
//	r.Post("/users", CreateUserHandler)
package rest
