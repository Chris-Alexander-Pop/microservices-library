// Package graph provides a unified interface for graph databases.
//
// Supported backends:
//   - AWS Neptune
//   - Neo4j
//   - Memory (for testing)
//
// Usage:
//
//	g, err := neptune.New(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer g.Close()
//
//	v := &graph.Vertex{ID: "1", Label: "Person", Properties: map[string]interface{}{"name": "Alice"}}
//	err := g.AddVertex(ctx, v)
package graph
