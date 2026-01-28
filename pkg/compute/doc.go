/*
Package compute provides compute resource management abstractions.

This package organizes compute functionality into the following subpackages:

  - vm: Virtual machine lifecycle management
  - container: Container runtime and orchestration (Kubernetes, Fargate)
  - serverless: Function-as-a-Service (Lambda, Cloud Functions)

Each subpackage follows the standard adapter pattern with instrumented.go
wrappers and memory adapters for testing.
*/
package compute
