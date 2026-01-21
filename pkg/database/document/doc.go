// Package document provides a generic interface for document-oriented databases.
//
// Supported backends:
//   - AWS DynamoDB
//   - Azure CosmosDB
//   - GCP Firestore
//
// Usage:
//
//	db, err := dynamodb.New(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer db.Close()
//
//	doc := document.Document{"id": "1", "name": "chris"}
//	err := db.Insert(ctx, "users", doc)
package document
