package introspection

import (
	"context"

	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// ListTables retrieves a list of all tables in the database.
func ListTables(ctx context.Context, db *gorm.DB) ([]string, error) {
	// GORM migrator abstraction handles dialect specifics for listing tables
	return db.Migrator().GetTables()
}

// HasTable checks if a table exists.
func HasTable(ctx context.Context, db *gorm.DB, table string) bool {
	return db.Migrator().HasTable(table)
}

// GetIndexes returns index information for a table.
// GORM Migrator provides GetIndexes logic.
func GetIndexes(ctx context.Context, db *gorm.DB, dst interface{}) ([]gorm.Index, error) {
	return db.Migrator().GetIndexes(dst)
}

// ListCollections retrieves a list of collections/tables from ANY supported database type.
// Unified introspection for the unified manager.
func ListCollections(ctx context.Context, db interface{}) ([]string, error) {
	switch v := db.(type) {
	case *gorm.DB:
		return v.Migrator().GetTables()

	case *mongo.Database:
		return v.ListCollectionNames(ctx, bson.D{})

	case *gocql.Session:
		// Cassandra query system_schema
		// Simplified: assumes v.Keyspace is set in session or we query all
		// Typically requires `SELECT table_name FROM system_schema.tables WHERE keyspace_name = ?`
		// Getting keyspace from session struct is hard (private), so might need user to pass it or just fail gracefully.
		// For now, return empty or implement if keyspace known.
		iter := v.Query("SELECT table_name FROM system_schema.tables").Iter()
		var tables []string
		var tableName string
		for iter.Scan(&tableName) {
			tables = append(tables, tableName)
		}
		if err := iter.Close(); err != nil {
			return nil, err
		}
		return tables, nil

	default:
		return nil, nil // Not supported or KV store
	}
}
