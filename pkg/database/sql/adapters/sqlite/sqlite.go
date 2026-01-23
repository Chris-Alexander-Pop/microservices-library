package sqlite

import (
	"fmt"

	"github.com/chris-alexander-pop/system-design-library/pkg/database"
	"github.com/chris-alexander-pop/system-design-library/pkg/database/sql"
	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New creates a new SQLite connection.
func New(cfg sql.Config) (*gorm.DB, error) {
	if cfg.Driver != database.DriverSQLite {
		return nil, errors.New(errors.CodeInvalidArgument, fmt.Sprintf("invalid driver %s for sqlite adapter", cfg.Driver), nil)
	}

	// For sqlite, Name is used as filepath
	filepath := cfg.Name
	if filepath == "" {
		filepath = "gorm.db"
	}

	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to sqlite")
	}

	return db, nil
}
