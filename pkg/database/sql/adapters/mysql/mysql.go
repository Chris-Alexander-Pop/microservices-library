package mysql

import (
	"fmt"

	"github.com/chris-alexander-pop/system-design-library/pkg/database"
	"github.com/chris-alexander-pop/system-design-library/pkg/database/sql"
	"github.com/chris-alexander-pop/system-design-library/pkg/errors"
	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// New creates a new MySQL connection.
func New(cfg sql.Config) (*gorm.DB, error) {
	if cfg.Driver != database.DriverMySQL {
		return nil, errors.New(errors.CodeInvalidArgument, fmt.Sprintf("invalid driver %s for mysql adapter", cfg.Driver), nil)
	}

	tlsParam := "false"

	// Load TLS Config
	tlsConfig, err := database.LoadTLSConfig(database.Config{
		SSLMode:     cfg.SSLMode,
		SSLCert:     cfg.SSLCert,
		SSLKey:      cfg.SSLKey,
		SSLRootCert: cfg.SSLRootCert,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to load tls config")
	}

	if tlsConfig != nil {
		err = mysqldriver.RegisterTLSConfig("custom", tlsConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to register mysql tls config")
		}
		tlsParam = "custom"
	} else if cfg.SSLMode == "require" || cfg.SSLMode == "true" {
		tlsParam = "true"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, tlsParam)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: database.NewGORMLogger().LogMode(gormlogger.Info),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mysql")
	}

	return db, nil
}
