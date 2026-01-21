// Package grpc provides common gRPC utilities, interceptors, and server configurations.
//
// This package standardizes gRPC server setup, including:
//   - Middleware chains (interceptors) for logging, tracing, recovery, and auth
//   - Health checking
//   - Reflection
//   - Error mapping from internal errors to gRPC codes
//
// Usage:
//
//	srv := grpc.NewServer(grpc.Config{Port: 8080})
//	pb.RegisterMyServiceServer(srv.Server, impl)
//	srv.Start()
package grpc
