package introspection

import (
	"context"

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
