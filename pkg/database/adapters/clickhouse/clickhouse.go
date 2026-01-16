package clickhouse

import (
	"fmt"

	"github.com/chris-alexander-pop/system-design-library/pkg/database"
	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New creates a new ClickHouse connection
func New(cfg database.Config) (*gorm.DB, error) {
	if cfg.Driver != "clickhouse" {
		return nil, errors.New(errors.CodeInvalidArgument, fmt.Sprintf("invalid driver %s for clickhouse adapter", cfg.Driver), nil)
	}

	secure := "false"
	if cfg.SSLMode == "require" || cfg.SSLMode == "true" {
		secure = "true"
	}

	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s?secure=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, secure)

	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to clickhouse")
	}

	return db, nil
}
