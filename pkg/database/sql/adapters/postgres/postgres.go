package postgres

import (
	"fmt"

	"github.com/chris-alexander-pop/system-design-library/pkg/database"
	"github.com/chris-alexander-pop/system-design-library/pkg/database/sql"
	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New creates a new Postgres connection using GORM.
func New(cfg sql.Config) (*gorm.DB, error) {
	if cfg.Driver != database.DriverPostgres {
		return nil, errors.New(errors.CodeInvalidArgument, fmt.Sprintf("invalid driver %s for postgres adapter", cfg.Driver), nil)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: database.NewGORMLogger(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get sql.DB")
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}
