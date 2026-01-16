package mssql

import (
	"fmt"

	"github.com/chris-alexander-pop/system-design-library/pkg/database"
	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New creates a new SQL Server connection
func New(cfg database.Config) (*gorm.DB, error) {
	if cfg.Driver != "sqlserver" {
		return nil, errors.New(errors.CodeInvalidArgument, fmt.Sprintf("invalid driver %s for mssql adapter", cfg.Driver), nil)
	}

	// Azure SQL often requires encrypt=true
	encryption := "disable"
	if cfg.SSLMode == "require" || cfg.SSLMode == "true" {
		encryption = "true"
	}

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, encryption)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to sqlserver")
	}

	return db, nil
}
